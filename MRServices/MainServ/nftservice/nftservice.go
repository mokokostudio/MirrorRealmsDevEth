package nftservice

import (
	"context"
	"sort"
	"strings"

	com "github.com/aureontu/MRWebServer/mr_services/common"
	"github.com/aureontu/MRWebServer/mr_services/mpb"
	"github.com/aureontu/MRWebServer/mr_services/mpberr"
	"github.com/aureontu/MRWebServer/mr_services/util"
	gconv "github.com/oldjon/gutil/conv"
	"github.com/oldjon/gutil/env"
	"github.com/oldjon/gutil/gdb"
	gprotocol "github.com/oldjon/gutil/protocol"
	grmux "github.com/oldjon/gutil/redismutex"
	gxgrpc "github.com/oldjon/gx/modules/grpc"
	"github.com/oldjon/gx/service"
	etcd "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type NFTService struct {
	mpb.UnimplementedNFTServiceServer
	name        string
	logger      *zap.Logger
	config      env.ModuleConfig
	etcdClient  *etcd.Client
	host        service.Host
	connMgr     *gxgrpc.ConnManager
	rm          *NFTResourceMgr
	kvm         *service.KVMgr
	serverEnv   uint32
	sm          *util.ServiceMetrics
	dao         *nftDAO
	tcpMsgCoder gprotocol.FrameCoder
}

// NewNFTService create an NFTService entity
func NewNFTService(driver service.ModuleDriver) (gxgrpc.GRPCServer, error) {
	svc := &NFTService{
		name:       driver.ModuleName(),
		logger:     driver.Logger(),
		config:     driver.ModuleConfig(),
		etcdClient: driver.Host().EtcdSession().Client(),
		host:       driver.Host(),
		kvm:        driver.Host().KVManager(),
		sm:         util.NewServiceMetrics(driver),
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
	svc.rm, err = newNFTResourceMgr(svc.logger, svc.sm)
	if err != nil {
		return nil, err
	}

	redisMux, err := grmux.NewRedisMux(svc.config.SubConfig("redis_mutex"), nil, svc.logger, driver.Tracer())
	if err != nil {
		return nil, err
	}

	nftRedis, err := gdb.NewRedisClientByConfig(svc.config.SubConfig("nft_redis"),
		svc.config.GetString("db_marshaller"), driver.Tracer())
	if err != nil {
		return nil, err
	}

	tranRedis, err := gdb.NewRedisClientByConfig(svc.config.SubConfig("tran_redis"),
		svc.config.GetString("db_marshaller"), driver.Tracer())
	if err != nil {
		return nil, err
	}

	tmpRedis, err := gdb.NewRedisClientByConfig(svc.config.SubConfig("tmp_redis"),
		svc.config.GetString("db_marshaller"), driver.Tracer())
	if err != nil {
		return nil, err
	}

	svc.dao = newNftDAO(svc, redisMux, nftRedis, tmpRedis, tranRedis)

	svc.serverEnv = uint32(svc.config.GetInt64("server_env"))
	svc.tcpMsgCoder = gprotocol.NewFrameCoder(svc.config.GetString("protocol_code"))

	return svc, err
}

func (svc *NFTService) Register(grpcServer *grpc.Server) {
	mpb.RegisterNFTServiceServer(grpcServer, svc)
}

func (svc *NFTService) Serve(ctx context.Context) error {
	<-ctx.Done()
	return ctx.Err()
}

func (svc *NFTService) Logger() *zap.Logger {
	return svc.logger
}

func (svc *NFTService) ConnMgr() *gxgrpc.ConnManager {
	return svc.connMgr
}

func (svc *NFTService) Name() string {
	return svc.name
}

func (svc *NFTService) GetAptosNFTs(ctx context.Context, req *mpb.ReqGetAptosNFTs) (*mpb.ResGetAptosNFTs, error) {
	var collections []string
	if len(req.NftTypes) == 0 {
		for _, v := range svc.rm.getAllNFTCollectionRsc() {
			collections = append(collections, v.CollectionHash...)
		}
	} else {
		for _, v := range req.NftTypes {
			nftRes, err := svc.rm.getNFTCollectionRscByNFTType(v)
			if err != nil {
				return nil, err
			}
			collections = append(collections, nftRes.CollectionHash...)
		}
	}

	nfts, err := svc.moralisGetNFTByWallets(ctx, []string{req.WalletAddr}, collections)
	if err != nil {
		return nil, err
	}

	res := &mpb.ResGetAptosNFTs{}
	for _, v := range nfts {
		for _, vv := range v {
			nft := &mpb.AptosNFTNode{
				NftType: uint32(svc.rm.getNFTTypeByCollection(vv.CollectionDataIdHash)),
			}
			// parse id
			nft.NftId, err = svc.parseNFTId(vv.Name)
			if err != nil {
				return nil, err
			}
			nft.Metadata = svc.rm.getNftMetaData(nft.NftId)
			res.Nfts = append(res.Nfts, nft)
		}
	}

	return res, nil
}

func (svc *NFTService) moralisGetNFTByWallets(ctx context.Context, addresses []string, collections []string) (map[string][]*mpb.MoralisNFTData, error) {
	client, err := com.GetAPIProxyGRPCClient(ctx, svc)
	if err != nil {
		return nil, err
	}
	res, err := client.MoralisGetNFTByWallets(ctx, &mpb.ReqMoralisGetNFTByWallets{
		WalletAddresses: addresses,
		Collections:     collections,
	})
	if err != nil {
		svc.logger.Error("moralisGetNFTByWallets failed", zap.Any("addresses", addresses),
			zap.Any("collections", collections),
			zap.Error(err))
		return nil, err
	}
	m := make(map[string][]*mpb.MoralisNFTData)
	for k, v := range res.Nfts {
		m[k] = v.List
	}
	return m, nil
}

func (svc *NFTService) parseNFTId(name string) (uint64, error) {
	strs := strings.Split(name, "#")
	if len(strs) != 2 {
		svc.logger.Error("parseNFTId failed", zap.Any("name", name))
		return 0, mpberr.ErrParseNFTId
	}
	id := gconv.StringToUint64(strs[1])
	if id == 0 {
		svc.logger.Error("parseNFTId failed", zap.Any("name", name))
		return 0, mpberr.ErrParseNFTId
	}
	return id, nil
}

func (svc *NFTService) GetAptosNFTMetadatas(ctx context.Context, req *mpb.ReqGetAptosNFTMetadatas) (res *mpb.ResGetAptosNFTMetadatas, err error) {
	res = &mpb.ResGetAptosNFTMetadatas{}
	for _, v := range req.NftIds {
		metadata := svc.rm.getNftMetaData(v)
		if metadata == "" {
			continue
		}
		res.Metadatas = append(res.Metadatas, &mpb.AptosNFTMetadata{
			NftId:    v,
			MetaData: metadata,
		})
	}
	return res, nil
}

func (svc *NFTService) GetAptosNFTsV2(ctx context.Context, req *mpb.ReqGetAptosNFTsV2) (*mpb.ResGetAptosNFTsV2, error) {
	// check query limit
	ok, err := svc.dao.checkGraphiLimit(ctx, req.WalletAddr)
	if err != nil {
		return nil, err
	}
	if ok {
		// try update by graphiql
		err = svc.updateByGraphiQL(ctx, req.UserId, req.WalletAddr)
		if err != nil {
			return nil, err
		}
	}

	// get from db
	res := &mpb.ResGetAptosNFTsV2{}
	dbNFTs, err := svc.dao.getNFTs(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	res.Nfts = make([]*mpb.AptosNFTNodeV2, 0, len(dbNFTs))
	for _, v := range dbNFTs {
		if v == nil {
			continue
		}
		res.Nfts = append(res.Nfts, svc.DBAptosNFTNodeV22AptosNFTNodeV2(v))
	}
	return res, nil
}

func (svc *NFTService) updateByGraphiQL(ctx context.Context, userId uint64, addr string) error {
	pageNum := svc.rm.getGraphiQLQueryPageNum()
	if pageNum == 0 {
		return nil
	}

	startIndex, err := svc.dao.getGraphiQLQueryStartIndex(ctx, addr)
	if err != nil {
		return err
	}

	// get from graphi
	client, err := com.GetAPIProxyGRPCClient(ctx, svc)
	if err != nil {
		return err
	}
	rpcRes, err := client.GraphiQLGetAccountTransactions(ctx, &mpb.ReqGraphiQLGetAccountTransactions{
		Addr:       addr,
		StartIndex: startIndex,
		PageNum:    pageNum,
	})
	if err != nil {
		return err
	}

	if rpcRes.Transactions.Data == nil || len(rpcRes.Transactions.Data.AccountTransactions) == 0 {
		return nil
	}

	addNFTs, delNFTs := svc.parseNFTTransactions(addr, rpcRes.Transactions.Data.AccountTransactions)
	if len(addNFTs) > 0 || len(delNFTs) > 0 {
		err = svc.dao.updateNFTs(ctx, userId, addNFTs, delNFTs)
		if err != nil {
			return err
		}
	}

	err = svc.dao.setGraphiQLQueryStartIndex(ctx, addr, int64(len(rpcRes.Transactions.Data.AccountTransactions)))
	if err != nil {
		return err
	}

	return nil
}

func (svc *NFTService) parseNFTTransactions(addr string, transactions []*mpb.AptosAccountTransactions_AccountTransaction,
) (addNFTs map[string]*mpb.AptosNFTNodeV2, delNFTs map[string]bool) {
	sort.SliceStable(transactions, func(i, j int) bool {
		return transactions[i].TransactionVersion < transactions[j].TransactionVersion
	})

	addNFTs = make(map[string]*mpb.AptosNFTNodeV2)
	delNFTs = make(map[string]bool)
	transferEvent := svc.rm.getNFTTransferEvent()
	burnEvent := svc.rm.getNFTBurnEvent()
	withDrawEvent2 := svc.rm.getNFTWithdrawEvent2()
	depositEvent2 := svc.rm.getNFTDepositEvent2()
	burnEvent2 := svc.rm.getNFTBurnEvent2()

	for _, t := range transactions {
		for _, act := range t.TokenActivitiesV2 {
			if act.CurrentTokenData == nil ||
				!svc.rm.isValidCollection(act.CurrentTokenData.CollectionId) {
				continue
			}
			switch act.Type {
			case transferEvent:
				if act.ToAddress == addr && act.FromAddress != addr { // got nft
					act.CurrentTokenData.TransactionTimestamp = act.TransactionTimestamp
					addNFTs[act.TokenDataId] = svc.AptosNFTNodeV2Origin2AptosNFTNodeV2(act.CurrentTokenData)
				} else if act.FromAddress == addr && act.ToAddress != addr { // lose nft
					if _, ok := addNFTs[act.TokenDataId]; ok {
						delete(addNFTs, act.TokenDataId)
					} else {
						delNFTs[act.TokenDataId] = true
					}
				}
			case burnEvent, burnEvent2:
				if act.FromAddress == addr {
					if _, ok := addNFTs[act.TokenDataId]; ok {
						delete(addNFTs, act.TokenDataId)
					} else {
						delNFTs[act.TokenDataId] = true
					}
				}
			case withDrawEvent2: //lose nft
				if act.FromAddress == addr && act.ToAddress != addr {
					act.CurrentTokenData.TransactionTimestamp = act.TransactionTimestamp
					addNFTs[act.TokenDataId] = svc.AptosNFTNodeV2Origin2AptosNFTNodeV2(act.CurrentTokenData)
				}
			case depositEvent2: //got nft
				if act.ToAddress == addr && act.FromAddress != addr {
					if _, ok := addNFTs[act.TokenDataId]; ok {
						delete(addNFTs, act.TokenDataId)
					} else {
						delNFTs[act.TokenDataId] = true
					}
				}
			}
		}
	}

	return addNFTs, delNFTs
}

func (svc *NFTService) GetAptosNFTOwner(ctx context.Context, req *mpb.ReqGetAptosNFTOwner) (*mpb.ResGetAptosNFTOwner, error) {
	userId, err := svc.dao.getAptosNFTOwnerUserId(ctx, req.TokenId)
	if err != nil {
		return nil, err
	}

	accInfo, err := svc.getAccountInfo(ctx, userId)
	if err != nil {
		return nil, err
	}

	return &mpb.ResGetAptosNFTOwner{Owner: accInfo}, nil
}

func (svc *NFTService) getAccountInfo(ctx context.Context, userId uint64) (*mpb.AccountInfo, error) {
	client, err := com.GetAccountServiceClient(ctx, svc)
	if err != nil {
		return nil, err
	}

	res, err := client.GetAccountInfo(ctx, &mpb.ReqUserId{UserId: userId})
	if err != nil {
		return nil, err
	}

	return res, nil
}
