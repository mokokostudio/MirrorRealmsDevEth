package httpgateway

import (
	"context"
	"net/http"
	"time"

	com "github.com/aureontu/MRWebServer/mr_services/common"
	"github.com/aureontu/MRWebServer/mr_services/mpb"
	"github.com/aureontu/MRWebServer/mr_services/util"
	"github.com/golang-jwt/jwt/v5"
	pb "github.com/golang/protobuf/proto"
	"github.com/oldjon/gutil/conv"
	"github.com/oldjon/gutil/env"
	gjwt "github.com/oldjon/gutil/jwt"
	gxgrpc "github.com/oldjon/gx/modules/grpc"
	gxhttp "github.com/oldjon/gx/modules/http"
	"github.com/oldjon/gx/service"
	"github.com/pkg/errors"
	etcd "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

var httpAESEncryptionKeyPairs = append(util.DefaultAESEncryptionKeyPairs, &util.AESEncryptionKeyPair{
	Index:   conv.Uint32ToString(1),
	Key:     []byte("ad#oiwUbn^asd!q1"),
	IV:      []byte("HUI@as0908(&^@!cs"),
	Retired: false,
})

type HTTPGateway struct {
	name       string
	logger     *zap.Logger
	config     env.ModuleConfig
	mux        *http.ServeMux
	etcdClient *etcd.Client
	connMgr    *gxgrpc.ConnManager

	kvm              *service.KVMgr
	signingKey       []byte
	signingMethod    jwt.SigningMethod
	signingDuration  time.Duration
	protocolEncode   string
	isSandbox        bool
	enableEncryption bool
	// registrationLimiter *util.RedisLimiter

	// HTTPClient Client
	metrics *metrics

	// wfClient *MultiLangWordFilter // TODO

	// res          *gatewayRes //TODO
	centerRegion string
}

// NewHTTPGateway create a HTTPGateway entity
func NewHTTPGateway(driver service.ModuleDriver) (gxhttp.GXHttpHandler, error) {
	mux := http.NewServeMux()
	host := driver.Host()
	gateway := HTTPGateway{
		name:            driver.ModuleName(),
		logger:          driver.Logger(),
		config:          driver.ModuleConfig(),
		mux:             mux,
		etcdClient:      host.EtcdSession().Client(),
		kvm:             host.KVManager(),
		signingMethod:   jwt.SigningMethodHS256,
		signingDuration: 24 * time.Hour,
		metrics:         newGatewayMetrics(driver),
	}

	gateway.protocolEncode = gateway.config.GetString("protocol_code")
	gateway.centerRegion = gateway.config.GetString("center_region")
	gateway.isSandbox = gateway.config.GetBool("is_sandbox")
	gateway.enableEncryption = gateway.config.GetBool("enable_encryption")

	jm := gjwt.New(gjwt.Options{
		KeyGetter: func(token *jwt.Token) (interface{}, error) {
			return gateway.signingKey, nil
		},
		NewClaimsFunc: func() jwt.Claims {
			return &mpb.JWTClaims{}
		},
		SigningMethod: gateway.signingMethod,
	})
	eh := util.NewHTTPErrorHandler(driver.Logger())
	dialer := gxgrpc.Dialer{
		HostName:   driver.Host().Name(),
		EtcdClient: gateway.etcdClient,
		Logger:     gateway.logger,
		Tracer:     driver.Tracer(),
		EnableTLS:  gateway.config.GetBool("enable_tls"),
		CAFile:     gateway.config.GetString("ca_file"),
		CertFile:   gateway.config.GetString("cert_file"),
		KeyFile:    gateway.config.GetString("key_file"),
	}
	_ = jm
	gateway.connMgr = gxgrpc.NewConnManager(&dialer)

	// var err error
	// gateway.HTTPClient, err = wire.NewHTTPClient(host, wire.HTTPClientOptions{})
	// if err != nil {
	// 	return nil, err
	// }
	//
	// gateway.cache, err = wire.NewCacheClient(host, wire.CacheClientOptions{})
	// if err != nil {
	// 	return nil, err
	// }
	//
	// gateway.wfClient, err = util.GetMultiLangWordFilter(gateway.logger)
	// if err != nil {
	// 	return nil, err
	// }
	// gateway.res, err = newGatewayRes(gateway.logger, gateway.metrics)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// if gateway.config.GetBool("registration_limit.open") {
	// 	bot, err := wire.NewRedisClient(host, wire.RedisClientOptions{ConfigKey: "registration_limit.gredis"})
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	gateway.registrationLimiter, err = util.NewRedisLimiter(bot, gateway.logger,
	// 		gateway.config.GetInt64("registration_limit.duration"),
	// 		gateway.config.GetInt64("registration_limit.cnt_per_dur"))
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	mux.Handle("/HelloWorld", eh.Handler(gateway.helloWorld))
	mux.Handle("/LoginByPassword", eh.Handler(gateway.loginByPassword))
	mux.Handle("/GenerateNonce", eh.Handler(gateway.generateNonce))
	//mux.Handle("/RegisterAccount", eh.Handler(gateway.RegisterAccount))
	mux.Handle("/WebLoginByWallet", eh.Handler(gateway.WebLoginByWallet))
	mux.Handle("/SendEmailBindCode", jm.Handler(eh.Handler(gateway.SendEmailBindCode)))
	mux.Handle("/WebBindEmail", jm.Handler(eh.Handler(gateway.webBindEmail)))
	mux.Handle("/ChangePassword", jm.Handler(eh.Handler(gateway.changePassword)))
	mux.Handle("/SendEmailResetPasswordCode", eh.Handler(gateway.sendEmailResetPasswordCode))
	mux.Handle("/CheckEmailResetPasswordCode", eh.Handler(gateway.checkEmailResetPasswordCode))
	mux.Handle("/ResetPasswordByEmail", eh.Handler(gateway.resetPasswordByEmail))
	mux.Handle("/ResetPasswordByEmailAndVCode", eh.Handler(gateway.resetPasswordByEmailAndVCode))
	mux.Handle("/GetAccountInfo", jm.Handler(eh.Handler(gateway.getAccountInfo)))
	mux.Handle("/GetAptosResources", jm.Handler(eh.Handler(gateway.getAptosResources)))
	mux.Handle("/GetAptosNFTs", jm.Handler(eh.Handler(gateway.getAptosNFTs)))
	mux.Handle("/GetAptosNFTMetadatas", jm.Handler(eh.Handler(gateway.getAptosNFTMetaDatas)))
	mux.Handle("/GetAptosNFTsV2", jm.Handler(eh.Handler(gateway.getAptosNFTsV2)))
	mux.Handle("/TestGetAptosNFTsV2", eh.Handler(gateway.testGetAptosNFTsV2))

	//mux.Handle("/GetUser", jm.Handler(eh.Handler(gateway.GetUser)))
	//mux.Handle("/GetItems", jm.Handler(eh.Handler(gateway.GetItems)))
	//mux.Handle("/ExchangeItems", jm.Handler(eh.Handler(gateway.ExchangeItems))) //TODO Delete in live
	// mail
	//mux.Handle("GetMailList", jm.Handler(eh.Handler(gateway.GetMailList)))
	//mux.Handle("ReadMails", jm.Handler(eh.Handler(gateway.ReadMails)))
	//mux.Handle("DelMails", jm.Handler(eh.Handler(gateway.DelMails)))
	//mux.Handle("GetMailsAwards", jm.Handler(eh.Handler(gateway.GetMailsAwards)))

	return &gateway, nil
}

func (hg *HTTPGateway) helloWorld(w http.ResponseWriter, r *http.Request) error {
	ip := getRemoteIPAddress(r)
	_, _ = w.Write([]byte("Hello World: " + ip))
	return nil
}

func (hg *HTTPGateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger := hg.logger.With(zap.String("path", r.URL.Path))
	logger.Info("handling http")
	defer logger.Info("finish http")
	hg.mux.ServeHTTP(w, r)
}

func (hg *HTTPGateway) Serve(ctx context.Context) error {
	signingKey, err := hg.kvm.GetOrGenerate(ctx, com.JWTGatewayTokenKey, 32)
	if err != nil {
		return errors.WithStack(err)
	}
	hg.signingKey = signingKey

	<-ctx.Done()
	return ctx.Err()
}

func (hg *HTTPGateway) Logger() *zap.Logger {
	return hg.logger
}

func (hg *HTTPGateway) ConnMgr() *gxgrpc.ConnManager {
	return hg.connMgr
}

func (hg *HTTPGateway) Name() string {
	return hg.name
}

func (hg *HTTPGateway) readHTTPReq(w http.ResponseWriter, r *http.Request, msg pb.Message) error {
	var err error

	_, isLogin := msg.(*mpb.CReqWebLoginByWallet)
	//other login way
	options := util.HTTPEncryptionOptions{
		EnableEncryption:          hg.enableEncryption,
		AESEncryptionKeyPairs:     httpAESEncryptionKeyPairs,
		IsPlatformLoginMethodCall: isLogin,
	}

	if hg.protocolEncode == "json" {
		err = util.ReadHTTPJSONReq(w, r, msg, options)
	} else {
		err = util.ReadHTTPReq(w, r, msg, options)
	}
	if err != nil {
		hg.metrics.incReadHTTPFail(r.URL.Path, err)
	}
	return err
}

func (hg *HTTPGateway) writeHTTPRes(w http.ResponseWriter, msg pb.Message) error {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "content-type")
	if hg.protocolEncode == "json" {
		return util.WriteHTTPJSONRes(w, msg)
	}
	return util.WriteHTTPRes(w, msg)
}
