package gmservice

import (
	"context"
	"time"

	com "github.com/aureontu/MRWebServer/mr_services/common"
	"github.com/aureontu/MRWebServer/mr_services/mpb"
	"github.com/aureontu/MRWebServer/mr_services/util"
	"github.com/golang-jwt/jwt/v5"
	"github.com/oldjon/gutil/env"
	gprotocol "github.com/oldjon/gutil/protocol"
	gxgrpc "github.com/oldjon/gx/modules/grpc"
	"github.com/oldjon/gx/service"
	"github.com/pkg/errors"
	etcd "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type GMService struct {
	mpb.UnimplementedGMServiceServer
	name            string
	logger          *zap.Logger
	config          env.ModuleConfig
	etcdClient      *etcd.Client
	host            service.Host
	connMgr         *gxgrpc.ConnManager
	signingMethod   jwt.SigningMethod
	signingDuration time.Duration
	signingKey      []byte
	rm              *GMResourceMgr
	kvm             *service.KVMgr
	serverEnv       uint32
	sm              *util.ServiceMetrics
	tcpMsgCoder     gprotocol.FrameCoder
}

// NewGMService create a gmservice entity
func NewGMService(driver service.ModuleDriver) (gxgrpc.GRPCServer, error) {
	svc := &GMService{
		name:            driver.ModuleName(),
		logger:          driver.Logger(),
		config:          driver.ModuleConfig(),
		etcdClient:      driver.Host().EtcdSession().Client(),
		host:            driver.Host(),
		kvm:             driver.Host().KVManager(),
		sm:              util.NewServiceMetrics(driver),
		signingMethod:   jwt.SigningMethodHS256,
		signingDuration: com.Dur1Day,
	}

	dialer := gxgrpc.Dialer{
		HostName:   driver.Host().Name(),
		Tracer:     driver.Tracer(),
		EtcdClient: svc.etcdClient,
		Logger:     svc.logger,
		EnableTLS:  svc.config.GetBool("enable_tls"),
		CAFile:     svc.config.GetString("ca_file"),
		CertFile:   svc.config.GetString("cert_file"),
		KeyFile:    svc.config.GetString("key_file"),
	}
	svc.connMgr = gxgrpc.NewConnManager(&dialer)

	var err error
	svc.rm, err = newGMResourceMgr(svc.logger, svc.sm)
	if err != nil {
		return nil, err
	}

	svc.serverEnv = uint32(svc.config.GetInt64("server_env"))
	svc.tcpMsgCoder = gprotocol.NewFrameCoder(svc.config.GetString("protocol_code"))

	return svc, err
}

func (svc *GMService) Register(grpcServer *grpc.Server) {
	mpb.RegisterGMServiceServer(grpcServer, svc)
}

func (svc *GMService) Serve(ctx context.Context) error {
	signingKey, err := svc.kvm.GetOrGenerate(ctx, com.GMJWTGatewayTokenKey, 32)
	if err != nil {
		return errors.WithStack(err)
	}
	svc.signingKey = signingKey

	<-ctx.Done()
	return ctx.Err()
}

func (svc *GMService) Logger() *zap.Logger {
	return svc.logger
}

func (svc *GMService) ConnMgr() *gxgrpc.ConnManager {
	return svc.connMgr
}

func (svc *GMService) Name() string {
	return svc.name
}

func (svc *GMService) generateAdminLoginToken(account string) (string, error) {
	var sToken string
	now := time.Now()
	claim := &mpb.JWTAdminClaims{}
	claim.Account = account
	claim.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(now.Add(svc.signingDuration)),
	}
	token := jwt.NewWithClaims(svc.signingMethod, claim)
	sToken, err := token.SignedString(svc.signingKey)
	if err != nil {
		return "", err
	}
	return sToken, nil
}

func (svc *GMService) AdminLoginByPassword(ctx context.Context, req *mpb.ReqAdminLoginByPassword) (*mpb.ResAdminLoginByPassword, error) {
	rsc, err := svc.rm.getAdminRSC(req.Account)
	if err != nil {
		return nil, err
	}
	if rsc.Password != req.Password {
		return nil, err
	}

	token, err := svc.generateAdminLoginToken(req.Account)
	if err != nil {
		return nil, err
	}

	return &mpb.ResAdminLoginByPassword{
		Token: token,
	}, nil
}
