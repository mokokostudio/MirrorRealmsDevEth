package util

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/aureontu/MRWebServer/mr_services/mpb"
	"github.com/aureontu/MRWebServer/mr_services/mpberr"
	pb "github.com/golang/protobuf/proto"
	"github.com/oldjon/gutil/conv"
	"github.com/oldjon/gutil/encryption"
	gjwt "github.com/oldjon/gutil/jwt"
	gtags "github.com/oldjon/gutil/tags"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// HTTPErrorHandler convert internal error to HTTP error and return to bot
type HTTPErrorHandler struct {
	logger *zap.Logger
}

// NewHTTPErrorHandler create a http error handler
func NewHTTPErrorHandler(logger *zap.Logger) *HTTPErrorHandler {
	hh := HTTPErrorHandler{
		logger: logger,
	}
	return &hh
}

// HTTPErrorHandlerFunc function definition
type HTTPErrorHandlerFunc func(w http.ResponseWriter, r *http.Request) error

func (f HTTPErrorHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	return f(w, r)
}

// Handler HTTPErrorHandler decoration
func (eh *HTTPErrorHandler) Handler(f HTTPErrorHandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := f.ServeHTTP(w, r)
		if err != nil {
			var aErr *AppError
			if errors.As(err, &aErr) {
				// deal with appError
				eh.logger.Info(r.RequestURI, zap.Error(err))
				http.Error(w, aErr.Message, aErr.Code)
				return
			}
			var tErr *TransactionError
			if errors.As(err, &tErr) {
				http.Error(w, tErr.Error(), http.StatusInternalServerError)
				return
			}
			if code := grpc.Code(err); code != codes.Unknown {
				// deal with grpc error
				if code == codes.NotFound {
					eh.logger.Info(r.RequestURI, zap.String("error desc", err.Error()))
					http.Error(w, grpc.ErrorDesc(err), http.StatusNotFound)
				} else if code == codes.AlreadyExists {
					eh.logger.Info(r.RequestURI, zap.String("error desc", err.Error()))
					http.Error(w, grpc.ErrorDesc(err), http.StatusBadRequest)
				} else if code == codes.InvalidArgument {
					eh.logger.Info(r.RequestURI, zap.String("error desc", err.Error()))
					http.Error(w, grpc.ErrorDesc(err), http.StatusBadRequest)
				} else {
					eh.logger.Error(r.RequestURI, zap.Error(err))
					http.Error(w, grpc.ErrorDesc(err), http.StatusInternalServerError)
				}
			} else {
				// grpc unknown error and other type of error
				desc := grpc.ErrorDesc(err)
				if errCode, ok := mpberr.HTTPErrMap[desc]; ok {
					eh.logger.Info(r.RequestURI, zap.String("error desc", desc))
					http.Error(w, desc, errCode)
				} else {
					eh.logger.Error(r.RequestURI, zap.Error(err))
					http.Error(w, desc, http.StatusInternalServerError)
				}
			}
		}
	})
}

// ClaimFromContext get claim from context
func ClaimFromContext(ctx context.Context) (*mpb.JWTClaims, context.Context, error) {
	claimIn, ok := gjwt.FromContext(ctx)
	if !ok {
		errMsg := "Extract jwt token from request context failed, this should not happen"
		return nil, ctx, &AppError{errors.New(errMsg), errMsg, http.StatusInternalServerError} //nolint:goerr113
	}
	claim, ok := claimIn.(*mpb.JWTClaims)
	if !ok {
		errMsg := "jwt token is not correct type, this should not happen"
		return nil, ctx, &AppError{errors.New(errMsg), errMsg, http.StatusInternalServerError} //nolint:goerr113
	}
	tags := gtags.FromContext(ctx)
	tags = tags.Set(TagsKeyUserID, conv.Uin64ToString(claim.UserId))
	tags = tags.Set(TagsKeyAccount, claim.Account)
	tags = tags.Set(TagsKeyRegion, claim.Region)
	ctx = gtags.WithTags(ctx, tags)
	return claim, ctx, nil
}

type AESEncryptionKeyPair struct {
	Index   string // one index refer to Key + IV
	Key     []byte
	IV      []byte
	Retired bool // the key pair is retired , we should not use in standard call , but use it to provide force upgrade information
}

// DefaultAESEncryptionKeyPairs contain a list of keypair currently using
var DefaultAESEncryptionKeyPairs = []*AESEncryptionKeyPair{
	{
		Index:   "Default",
		Retired: true,
	},
}

type HTTPEncryptionOptions struct {
	EnableEncryption      bool
	AESEncryptionKeyPairs []*AESEncryptionKeyPair

	// following is the hack options
	IsPlatformLoginMethodCall bool
}

const HTTPEncryptionHeader = "X-GA"

func getAESEncryptionKeyPair(r *http.Request, options HTTPEncryptionOptions) (*AESEncryptionKeyPair, error) {
	if len(options.AESEncryptionKeyPairs) == 0 {
		return nil, errAESEncryptionKeyPairsIsAnEmptyList
	}

	// select out AESEncryptionKeyPair
	authHeaderValue := r.Header.Get(HTTPEncryptionHeader)
	if authHeaderValue == "" {
		// do not have this header
		// use first one as default
		return options.AESEncryptionKeyPairs[0], nil
	}

	// authHeaderValue follows the format "{keyType} {keyIndex}"
	s := strings.Split(authHeaderValue, " ")
	if len(s) != 2 {
		return nil, errHeaderNotInWellFormat
	}

	keyType := s[0]
	keyIndex := s[1]

	if keyType != "v1" {
		return nil, errHeaderContainNotSupportedKeyType
	}

	// search options.AESEncryptionKeyPairs
	for _, kp := range options.AESEncryptionKeyPairs {
		if keyIndex == kp.Index {
			return kp, nil
		}
	}

	return nil, errHeaderContainNotSupportedKeyIndex
}

// ReadHTTPReq read http request
func ReadHTTPReq(w http.ResponseWriter, r *http.Request, msg pb.Message, options HTTPEncryptionOptions) error {
	var key *AESEncryptionKeyPair

	// get key if we EnableEncryption
	if options.EnableEncryption {
		var err error
		key, err = getAESEncryptionKeyPair(r, options)
		if err != nil {
			return err
		}
	}

	// check key is retired or not
	if options.EnableEncryption {
		if key.Retired {
			if options.IsPlatformLoginMethodCall {
				return &AppError{fmt.Errorf("key[%s] is retired", key.Index), "please upgrade", 801} //nolint:goerr113
			}
			return fmt.Errorf("key[%s] is retired", key.Index) //nolint:goerr113
		}
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return &AppError{err, "invalid request body", http.StatusBadRequest}
	}
	r.Body.Close()

	if options.EnableEncryption {
		body, err = handleDecryption(body, *key)
		if err != nil {
			return err
		}
	}

	err = pb.Unmarshal(body, msg)
	if err != nil {
		return &AppError{err, "invalid request body", http.StatusBadRequest} //nolint:goerr113
	}
	return nil
}

// ReadHTTPJSONReq read http request
func ReadHTTPJSONReq(w http.ResponseWriter, r *http.Request, i interface{}, options HTTPEncryptionOptions) error {
	var key *AESEncryptionKeyPair

	// get key if we EnableEncryption
	if options.EnableEncryption {
		var err error
		key, err = getAESEncryptionKeyPair(r, options)
		if err != nil {
			return err
		}
	}

	// check key is retired or not
	if options.EnableEncryption {
		if key.Retired {
			if options.IsPlatformLoginMethodCall {
				return &AppError{fmt.Errorf("key[%s] is retired", key.Index), "please upgrade", 801} //nolint:goerr113
			}
			return fmt.Errorf("key[%s] is retired", key.Index) //nolint:goerr113
		}
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return &AppError{err, "invalid request body", http.StatusBadRequest} //nolint:goerr113
	}
	r.Body.Close()

	if options.EnableEncryption {
		body, err = handleDecryption(body, *key)
		if err != nil {
			return err
		}
	}

	err = json.Unmarshal(body, &i)
	if err != nil {
		return &AppError{err, "invalid request body", http.StatusBadRequest} //nolint:goerr113
	}
	return nil
}

// WriteHTTPRes write to http response
func WriteHTTPRes(w http.ResponseWriter, msg pb.Message) error {
	res, err := pb.Marshal(msg)
	if err != nil {
		return err
	}
	if _, err := w.Write(res); err != nil {
		return err
	}
	return nil
}

// WriteHTTPJSONRes write to http response
func WriteHTTPJSONRes(w http.ResponseWriter, i interface{}) error {
	w.Header().Set("content-type", "application/json; charset=utf-8")
	res, err := json.Marshal(i)
	if err != nil {
		return err
	}
	if _, err := w.Write(res); err != nil {
		return err
	}
	return nil
}

// GetURLParamString will get param value from params "" when not exits
func GetURLParamString(params *url.Values, name string) string {
	if strs, ok := (*params)[name]; ok {
		if len(strs) > 0 {
			return strs[0]
		}
	}
	return ""
}

// GetURLParamUInt32 will get param value from params 0 when not exits
func GetURLParamUInt32(params *url.Values, name string) uint32 {
	if strs, ok := (*params)[name]; ok {
		if len(strs) > 0 {
			ret := conv.StringToUint32(strs[0])
			return ret
		}
	}
	return 0
}

// GetURLParamInt32 will get param value from params 0 when not exits
func GetURLParamInt32(params *url.Values, name string) int32 {
	if strs, ok := (*params)[name]; ok {
		if len(strs) > 0 {
			ret := conv.StringToInt32(strs[0])
			return ret
		}
	}
	return 0
}

// GetURLParamUInt64 will get param value from params 0 when not exits
func GetURLParamUInt64(params *url.Values, name string) uint64 {
	if strs, ok := (*params)[name]; ok {
		if len(strs) > 0 {
			ret := conv.StringToUint64(strs[0])
			return ret
		}
	}
	return 0
}

// GetURLParamFloat32 will get param value from params 0 when not exits
func GetURLParamFloat32(params *url.Values, name string) float32 {
	if strs, ok := (*params)[name]; ok {
		if len(strs) > 0 {
			return conv.StringToFloat32(strs[0])
		}
	}
	return 0
}

// GetURLParamBool will get param value from params false when not exits
func GetURLParamBool(params *url.Values, name string) bool {
	if strs, ok := (*params)[name]; ok {
		if len(strs) > 0 {
			ret := conv.StringToBool(strs[0])
			return ret
		}
	}
	return false
}

func handleDecryption(body []byte, key AESEncryptionKeyPair) ([]byte, error) {
	decryptedBody, err := encryption.AESDecrypt(body, key.Key, key.IV)
	if err != nil {
		return nil, err
	}
	return decryptedBody, nil
}

func HttpGet(ctx context.Context, reqUrl string) ([]byte, error) {
	var err error
	req := &http.Request{}
	req.URL, err = url.Parse(reqUrl)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = res.Body.Close()
	}()
	return io.ReadAll(res.Body)
}

func HttpGetWithHeader(ctx context.Context, reqUrl string, headers map[string]string) ([]byte, error) {
	var err error
	req := &http.Request{}
	req.URL, err = url.Parse(reqUrl)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	req.Header = http.Header{}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = res.Body.Close()
	}()
	return io.ReadAll(res.Body)
}

func HttpsGet(ctx context.Context, reqUrl string) ([]byte, error) {
	var err error
	req := &http.Request{}
	req.URL, err = url.Parse(reqUrl)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = res.Body.Close()
	}()
	return io.ReadAll(res.Body)
}

func HttpPost(ctx context.Context, reqUrl string, headers map[string]string, body []byte) ([]byte, error) {
	var err error
	//req ,err:= http.NewRequest("POST", reqUrl,strings.NewReader(string(body)))
	req := &http.Request{}
	req.Method = "POST"
	req.URL, err = url.Parse(reqUrl)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	req.Header = http.Header{}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	req.Body = io.NopCloser(bytes.NewBuffer(body))
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = res.Body.Close()
	}()
	return io.ReadAll(res.Body)
}
