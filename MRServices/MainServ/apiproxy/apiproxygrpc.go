package apiproxy

import (
	"context"

	"github.com/aureontu/MRWebServer/mr_services/mpb"
	"github.com/aureontu/MRWebServer/mr_services/util"
	"github.com/oldjon/gutil/env"
	gxgrpc "github.com/oldjon/gx/modules/grpc"
	"github.com/oldjon/gx/service"
	etcd "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var apiproxygrpc *APIProxyGRPCService

func APIProxyGRPCGetMe() *APIProxyGRPCService { // nolint:unused
	if apiproxygrpc == nil {
		panic("apiproxygrpc not initialize")
	}
	return apiproxygrpc
}

type APIProxyGRPCService struct {
	mpb.UnimplementedAPIProxyGRPCServer
	logger         *zap.Logger
	config         env.ModuleConfig
	etcdClient     *etcd.Client
	cm             *gxgrpc.ConnManager // 转发至其他gateway的消息，需要通过grpc
	host           service.Host
	sm             *util.ServiceMetrics
	rm             *APIProxyResourceMgr
	aptos          *AptosManager
	emailSendIndex uint32
}

// NewAPIProxyGRPCService create a APIProxyGRPCService entity
func NewAPIProxyGRPCService(driver service.ModuleDriver) (gxgrpc.GRPCServer, error) {
	as := &APIProxyGRPCService{
		logger:     driver.Logger(),
		config:     driver.ModuleConfig(),
		etcdClient: driver.Host().EtcdSession().Client(),
		host:       driver.Host(),
		sm:         util.NewServiceMetrics(driver),
	}

	dialer := gxgrpc.Dialer{
		HostName:   driver.Host().Name(),
		EtcdClient: as.etcdClient,
		Logger:     as.logger,
		Tracer:     driver.Tracer(),
		EnableTLS:  as.config.GetBool("enable_tls"),
		CAFile:     as.config.GetString("ca_file"),
		CertFile:   as.config.GetString("cert_file"),
		KeyFile:    as.config.GetString("key_file"),
	}
	as.cm = gxgrpc.NewConnManager(&dialer)
	var err error
	as.rm, err = newAPIProxyResourceMgr(as.logger, nil)
	if err != nil {
		return nil, err
	}
	as.aptos = newAptosManager(as.config.GetString("aptos_api_url"), as.config.GetString("moralis_api_url"),
		as.config.GetString("graphiql_api_url"), as.logger)
	as.logger.Info("apiproxy grpc service start success")
	apiproxygrpc = as
	return as, nil
}

func (svc *APIProxyGRPCService) Register(grpcServer *grpc.Server) {
	mpb.RegisterAPIProxyGRPCServer(grpcServer, svc)
}

func (svc *APIProxyGRPCService) Serve(ctx context.Context) error {
	<-ctx.Done()
	return ctx.Err()
}

func (svc *APIProxyGRPCService) SendEmailBindCode(ctx context.Context, req *mpb.ReqSendEmailBindCode) (*mpb.Empty, error) {
	err := svc.sendEmailBindCode(req.Email, req.Code)
	if err != nil {
		return nil, err
	}
	return &mpb.Empty{}, nil
}

func (svc *APIProxyGRPCService) SendEmailResetPasswordValidationCode(ctx context.Context, req *mpb.ReqSendEmailResetPasswordValidationCode) (*mpb.Empty, error) {
	err := svc.sendEmailResetPasswordValidationCode(req.Email, req.Code)
	if err != nil {
		return nil, err
	}
	return &mpb.Empty{}, nil
}

func (svc *APIProxyGRPCService) GetAptosResources(ctx context.Context, req *mpb.ReqGetAptosResources) (*mpb.ResGetAptosResources, error) { // TODO
	datas, err := svc.aptos.GetAccountResources(ctx, req.AptosAccAddr)
	if err != nil {
		return nil, err
	}
	return &mpb.ResGetAptosResources{
		Resources: string(datas),
	}, nil
}

func (svc *APIProxyGRPCService) MoralisGetNFTByWallets(ctx context.Context, req *mpb.ReqMoralisGetNFTByWallets) (*mpb.ResMoralisGetNFTByWallets, error) {
	apikey, err := svc.rm.getMoralisAPIKey()
	if err != nil {
		return nil, err
	}
	nfts, err := svc.aptos.MoralisGetNFTByWallets(ctx, req.WalletAddresses, req.Collections, apikey)
	if err != nil {
		return nil, err
	}
	res := &mpb.ResMoralisGetNFTByWallets{
		Nfts: make(map[string]*mpb.ResMoralisGetNFTByWallets_NFTList),
	}
	for k, v := range nfts {
		res.Nfts[k] = &mpb.ResMoralisGetNFTByWallets_NFTList{List: v}
	}
	return res, nil
}

func (svc *APIProxyGRPCService) GraphiQLGetAccountTransactions(ctx context.Context, req *mpb.ReqGraphiQLGetAccountTransactions) (*mpb.ResGraphiQLGetAccountTransactions, error) {
	data, err := svc.aptos.GraphiQLGetAccountTransactions(ctx, req.Addr, req.StartIndex, req.PageNum)
	if err != nil {
		svc.logger.Error("GraphiQLGetAccountTransactions failed", zap.Error(err))
		return nil, err
	}
	return &mpb.ResGraphiQLGetAccountTransactions{Transactions: data}, err
}

func (svc *APIProxyGRPCService) GraphiQLGetCollectionTransactions(ctx context.Context, req *mpb.ReqGraphiQLGetCollectionTransactions) (*mpb.ResGraphiQLGetCollectionTransactions, error) {
	data, err := svc.aptos.GraphiQLGetCollectionTransactions(ctx, req.CollectionId, req.StartIndex, req.PageNum)
	if err != nil {
		svc.logger.Error("GraphiQLGetCollectionTransactions failed", zap.Error(err))
		return nil, err
	}
	return &mpb.ResGraphiQLGetCollectionTransactions{Transactions: data}, err
}
