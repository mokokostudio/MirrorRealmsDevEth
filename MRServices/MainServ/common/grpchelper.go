package common

import (
	"context"

	"github.com/aureontu/MRWebServer/mr_services/mpb"
	gxgrpc "github.com/oldjon/gx/modules/grpc"
	"go.uber.org/zap"
)

type GRPCClientGetter interface {
	Name() string
	Logger() *zap.Logger
	ConnMgr() *gxgrpc.ConnManager
}

func GetAccountServiceClient(ctx context.Context, svc GRPCClientGetter) (mpb.AccountServiceClient, error) {
	conn, err := svc.ConnMgr().Dial(ctx, "accountservice")
	if err != nil {
		svc.Logger().Error("connect to accountservice failed", zap.String("service", svc.Name()))
		return nil, err
	}
	return mpb.NewAccountServiceClient(conn), nil
}

func GetAPIProxyGRPCClient(ctx context.Context, svc GRPCClientGetter) (mpb.APIProxyGRPCClient, error) {
	conn, err := svc.ConnMgr().Dial(ctx, "apiproxygrpc")
	if err != nil {
		svc.Logger().Error("connect to apiproxygrpc failed", zap.String("service", svc.Name()))
		return nil, err
	}
	return mpb.NewAPIProxyGRPCClient(conn), nil
}

func GetNFTServiceClient(ctx context.Context, svc GRPCClientGetter) (mpb.NFTServiceClient, error) {
	conn, err := svc.ConnMgr().Dial(ctx, "nftservice")
	if err != nil {
		svc.Logger().Error("connect to nftservice failed", zap.String("service", svc.Name()))
		return nil, err
	}
	return mpb.NewNFTServiceClient(conn), nil
}

func GetGMServiceClient(ctx context.Context, svc GRPCClientGetter) (mpb.GMServiceClient, error) {
	conn, err := svc.ConnMgr().Dial(ctx, "gmservice")
	if err != nil {
		svc.Logger().Error("connect to gmservice failed", zap.String("service", svc.Name()))
		return nil, err
	}
	return mpb.NewGMServiceClient(conn), nil
}
