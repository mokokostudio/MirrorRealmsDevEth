package apiproxy

import (
	"context"
	"net/http"

	"github.com/aureontu/MRWebServer/mr_services/util"
	pb "github.com/golang/protobuf/proto"
	"github.com/oldjon/gutil/conv"
	"github.com/oldjon/gutil/env"
	gxgrpc "github.com/oldjon/gx/modules/grpc"
	gxhttp "github.com/oldjon/gx/modules/http"
	"github.com/oldjon/gx/service"
	etcd "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

var httpAESEncryptionKeyPairs = append(util.DefaultAESEncryptionKeyPairs, &util.AESEncryptionKeyPair{ // nolint:unused
	Index:   conv.Uint32ToString(1),
	Key:     []byte("ad#oiwUbn^asd!q1"),
	IV:      []byte("HUI@as0908(&^@!cs"),
	Retired: false,
})

var apiproxy *APIProxy

func APIProxyGetMe() *APIProxy { // nolint:unused
	if apiproxy == nil {
		panic("apiproxy not initialize")
	}
	return apiproxy
}

type APIProxy struct {
	logger     *zap.Logger
	config     env.ModuleConfig
	mux        *http.ServeMux
	etcdClient *etcd.Client
	connMgr    *gxgrpc.ConnManager

	kvm              *service.KVMgr
	protocolEncode   string
	isSandbox        bool
	enableEncryption bool
	centerRegion     string
	// HTTPClient Client
	metrics *metrics
}

// NewAPIProxy create an apiproxy entity
func NewAPIProxy(driver service.ModuleDriver) (gxhttp.GXHttpHandler, error) {
	mux := http.NewServeMux()
	host := driver.Host()
	ap := APIProxy{
		logger:     driver.Logger(),
		config:     driver.ModuleConfig(),
		mux:        mux,
		etcdClient: host.EtcdSession().Client(),
		kvm:        host.KVManager(),
		metrics:    newMetrics(driver),
	}

	ap.protocolEncode = ap.config.GetString("protocol_code")
	ap.centerRegion = ap.config.GetString("center_region")
	ap.isSandbox = ap.config.GetBool("is_sandbox")
	ap.enableEncryption = ap.config.GetBool("enable_encryption")

	eh := util.NewHTTPErrorHandler(driver.Logger())
	dialer := gxgrpc.Dialer{
		HostName:   driver.Host().Name(),
		EtcdClient: ap.etcdClient,
		Logger:     ap.logger,
		Tracer:     driver.Tracer(),
		EnableTLS:  ap.config.GetBool("enable_tls"),
		CAFile:     ap.config.GetString("ca_file"),
		CertFile:   ap.config.GetString("cert_file"),
		KeyFile:    ap.config.GetString("key_file"),
	}
	ap.connMgr = gxgrpc.NewConnManager(&dialer)

	mux.Handle("/HelloWorld", eh.Handler(ap.helloWorld))
	apiproxy = &ap
	return &ap, nil
}

func (ap *APIProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger := ap.logger.With(zap.String("path", r.URL.Path))
	logger.Info("handling http")
	defer logger.Info("finish http")
	ap.mux.ServeHTTP(w, r)
}

func (ap *APIProxy) Serve(ctx context.Context) error {
	<-ctx.Done()
	return ctx.Err()
}

func (ap *APIProxy) readHTTPReq(w http.ResponseWriter, r *http.Request, msg pb.Message) error { // nolint:unused
	var err error

	//_, isLogin := msg.(*mpb.CReqLogin)// TODO

	options := util.HTTPEncryptionOptions{
		EnableEncryption:      ap.enableEncryption,
		AESEncryptionKeyPairs: httpAESEncryptionKeyPairs,
		//IsPlatformLoginMethodCall: isLogin,
	}

	if ap.protocolEncode == "json" {
		err = util.ReadHTTPJSONReq(w, r, msg, options)
	} else {
		err = util.ReadHTTPReq(w, r, msg, options)
	}
	if err != nil {
		ap.metrics.incReadHTTPFail(r.URL.Path, err)
	}
	return err
}

func (ap *APIProxy) writeHTTPRes(w http.ResponseWriter, msg pb.Message) error { // nolint:unused
	if ap.protocolEncode == "json" {
		return util.WriteHTTPJSONRes(w, msg)
	}
	return util.WriteHTTPRes(w, msg)
}

func (ap *APIProxy) helloWorld(w http.ResponseWriter, r *http.Request) error {
	//_, err := w.Write([]byte("hello world"))

	return APIProxyGRPCGetMe().sendEmailBindCode("lrunwow@gmail.com", "123456")
}
