package nftservice

import (
	"context"
	"fmt"
	"time"

	com "github.com/aureontu/MRWebServer/mr_services/common"
	"github.com/aureontu/MRWebServer/mr_services/mpb"
	"github.com/aureontu/MRWebServer/mr_services/mpberr"
	"github.com/oldjon/gutil"
	"github.com/oldjon/gutil/gdb"
	grmux "github.com/oldjon/gutil/redismutex"
	"go.uber.org/zap"
)

type nftDAO struct {
	svc    *NFTService
	logger *zap.Logger
	rMux   *grmux.RedisMutex
	nftDB  *gdb.DB
	tranDB *gdb.DB
	tmpDB  *gdb.DB
	rm     *NFTResourceMgr
}

func newNftDAO(svc *NFTService, rMux *grmux.RedisMutex, nftRedis, tmpRedis, tranRedis gdb.RedisClient) *nftDAO {
	return &nftDAO{
		logger: svc.logger,
		rMux:   rMux,
		nftDB:  gdb.NewDB(nftRedis),
		tranDB: gdb.NewDB(tranRedis),
		tmpDB:  gdb.NewDB(tmpRedis),
		rm:     svc.rm,
	}
}

func (dao *nftDAO) checkGraphiLimit(ctx context.Context, addr string) (bool, error) {
	ok, err := dao.tmpDB.SetEXNX(ctx, com.NFTGraphiLimitKey(addr), 0, time.Second*time.Duration(dao.rm.getGraphiQLQueryLimitSecs()))
	if err != nil {
		dao.logger.Error("checkGraphiLimit failed", zap.Error(err))
		return false, mpberr.ErrDB
	}
	return ok, nil
}

func (dao *nftDAO) getGraphiQLQueryStartIndex(ctx context.Context, addr string) (uint64, error) {
	index, err := gdb.ToUint64(dao.nftDB.Get(ctx, com.NFTGraphiStartIndexKey(addr)))
	if err != nil && !dao.nftDB.IsErrNil(err) {
		dao.logger.Error("getGraphiQLQueryStartIndex failed", zap.Error(err))
		return 0, mpberr.ErrDB
	}
	return index, nil
}

func (dao *nftDAO) setGraphiQLQueryStartIndex(ctx context.Context, addr string, cnt int64) error {
	_, err := dao.nftDB.IncrBy(ctx, com.NFTGraphiStartIndexKey(addr), cnt)
	if err != nil {
		dao.logger.Error("setGraphiQLQueryStartIndex failed", zap.Error(err))
		return mpberr.ErrDB
	}
	return nil
}

func (dao *nftDAO) updateNFTs(ctx context.Context, userId uint64, addNFTs map[string]*mpb.AptosNFTNodeV2,
	delNFTs map[string]bool) error {
	key := com.NFTsKey(userId)
	anyMap := make(map[string]any)
	anyList := make([]any, 0, len(addNFTs)*2)
	batchAddKeys := make([]string, 0, len(addNFTs))
	batchAddValues := make([]any, 0, len(addNFTs))
	for k, v := range addNFTs {
		nft := dao.svc.AptosNFTNodeV22DBAptosNFTNodeV2(v)
		anyMap[k] = nft
		batchAddKeys = append(batchAddKeys, k)
		batchAddValues = append(batchAddValues, userId)
		anyList = append(anyList, uint64(nft.TokenId)<<32|uint64(nft.TransactionTimestampInt), nft.TokenDataId)
	}
	delFields := make([]string, 0, len(delNFTs))
	delAnyList := make([]any, 0, len(delNFTs))
	for k := range delNFTs {
		delFields = append(delFields, k)
		delAnyList = append(delAnyList, k)
	}
	fmt.Println(addNFTs)
	err := dao.rMux.Safely(ctx, key, func() error {
		var err error
		if len(anyMap) > 0 {
			_, err = dao.nftDB.ZAdd(ctx, com.NFTsListKey(userId), anyList...)
			if err != nil {
				dao.logger.Error("updateNFTs ZAdd failed", zap.Uint64("user_id", userId),
					zap.Any("add_nfts", addNFTs), zap.Any("del_nfts", delNFTs), zap.Error(err))
				return err
			}
			err = dao.nftDB.HSetObjects(ctx, key, anyMap)
			if err != nil {
				dao.logger.Error("updateNFTs HSetObjects failed", zap.Uint64("user_id", userId),
					zap.Any("add_nfts", addNFTs), zap.Any("del_nfts", delNFTs), zap.Error(err))
				return err
			}
		}
		if len(delFields) > 0 {
			_, err = dao.nftDB.ZRem(ctx, com.NFTsListKey(userId), delAnyList...)
			if err != nil {
				dao.logger.Error("updateNFTs ZRem failed", zap.Uint64("user_id", userId),
					zap.Any("add_nfts", addNFTs), zap.Any("del_nfts", delNFTs), zap.Error(err))
				return err
			}
			_, err = dao.nftDB.HDel(ctx, key, delFields...)
			if err != nil {
				dao.logger.Error("updateNFTs HDel failed", zap.Uint64("user_id", userId),
					zap.Any("add_nfts", addNFTs), zap.Any("del_nfts", delNFTs), zap.Error(err))
				return err
			}
		}
		return nil
	})
	if err != nil {
		dao.logger.Error("updateNFTs safely failed", zap.Uint64("user_id", userId),
			zap.Any("add_nfts", addNFTs), zap.Any("del_nfts", delNFTs), zap.Error(err))
		return mpberr.ErrDB
	}

	err = dao.nftDB.BatchSet(ctx, batchAddKeys, batchAddValues, 0)
	if err != nil {
		dao.logger.Error("updateNFTs BatchSet failed", zap.Uint64("user_id", userId),
			zap.Any("add_nfts", addNFTs), zap.Any("del_nfts", delNFTs), zap.Error(err))
	}
	err = dao.nftDB.BatchDel(ctx, delFields)
	if err != nil {
		dao.logger.Error("updateNFTs BatchDel failed", zap.Uint64("user_id", userId),
			zap.Any("add_nfts", addNFTs), zap.Any("del_nfts", delNFTs), zap.Error(err))
	}

	return nil
}

func (dao *nftDAO) getNFTs(ctx context.Context, userId uint64) ([]*mpb.DBAptosNFTNodeV2, error) {
	fields, err := dao.nftDB.ZRange(ctx, com.NFTsListKey(userId), 0, -1)
	if err != nil {
		dao.logger.Error("getNFTs ZRange failed", zap.Uint64("user_id", userId), zap.Error(err))
		return nil, mpberr.ErrDB
	}
	if len(fields) == 0 {
		return nil, nil
	}
	nfts := make([]*mpb.DBAptosNFTNodeV2, len(fields))
	for i := range nfts {
		nfts[i] = &mpb.DBAptosNFTNodeV2{}
	}
	err = dao.nftDB.HMGetObjects(ctx, com.NFTsKey(userId), fields, nfts)
	if err != nil {
		dao.logger.Error("getNFTs HGetAllObjects failed", zap.Uint64("user_id", userId), zap.Error(err))
		return nil, mpberr.ErrDB
	}
	return nfts, nil
}

func (dao *nftDAO) getAptosNFTOwnerUserId(ctx context.Context, tokenId string) (uint64, error) {
	key := com.NFTUIDKey(tokenId)
	userId, err := gdb.ToUint64(dao.nftDB.Get(ctx, com.NFTUIDKey(tokenId)))
	if dao.nftDB.IsErrNil(err) {
		return 0, mpberr.ErrNFTNoOwner
	}
	if err != nil {
		dao.logger.Error("getAptosNFTOwnerUserId Get failed", zap.String("key", key), zap.Error(err))
		return 0, mpberr.ErrDB
	}
	return userId, nil
}

func (dao *nftDAO) getCollectionGraphiQLStartIndex(ctx context.Context, addr string) (uint64, error) {
	index, err := gdb.ToUint64(dao.tranDB.Get(ctx, com.CollectionGraphiStartIndex(addr)))
	if err != nil && !dao.tranDB.IsErrNil(err) {
		dao.logger.Error("getCollectionGraphiQLStartIndex failed", zap.Error(err))
		return 0, mpberr.ErrDB
	}
	return index, nil
}

func (dao *nftDAO) setCollectionGraphiQLStartIndex(ctx context.Context, addr string, cnt int64) error {
	_, err := dao.tranDB.IncrBy(ctx, com.CollectionGraphiStartIndex(addr), cnt)
	if err != nil {
		dao.logger.Error("setCollectionGraphiQLStartIndex failed", zap.Error(err))
		return mpberr.ErrDB
	}
	return nil
}

func (dao *nftDAO) saveCollectionTransactions(ctx context.Context, collectionId string, tas []*mpb.DBTokenActivitiesV2) error {
	l := len(tas)
	if l == 0 {
		return nil
	}
	var tavs = make([]any, 0, 2*l)
	for _, v := range tas {
		if len(v.Activities) == 0 {
			continue
		}
		tavs = append(tavs, v.Activities[0].TransactionVersion, v.Activities[0].TransactionVersion)
	}
	zKey := com.CollectionTranVersions(collectionId)
	_, err := dao.tranDB.ZAdd(ctx, zKey, tavs...)
	if err != nil {
		dao.logger.Error("saveCollectionTransactions ZAdd failed", zap.String("key", zKey), zap.Error(err))
		return mpberr.ErrDB
	}

	for i := 0; i < (l+com.DBBatchNum100-1)/com.DBBatchNum100; i++ {
		start := i * com.DBBatchNum100
		end := gutil.Min(l, (i+1)*com.DBBatchNum100)
		keys := make([]string, 0, end-start)
		for _, ta := range tas[start:end] {
			if len(ta.Activities) == 0 {
				continue
			}
			keys = append(keys, com.CollectionTranActivity(ta.Activities[0].TransactionVersion))
		}
		err = dao.tranDB.SetObjects(ctx, keys, tas[start:end])
		if err != nil {
			dao.logger.Error("saveCollectionTransactions BatchSet failed", zap.Any("keys", keys), zap.Error(err))
			return mpberr.ErrDB
		}
	}
	return nil
}

func (dao *nftDAO) getCollectionTransactions(ctx context.Context, collectionId string) ([]*mpb.DBTokenActivitiesV2, error) {
	zKey := com.CollectionTranVersions(collectionId)
	vers, err := gdb.ToUint64Slice(dao.tranDB.ZRange(ctx, zKey, 0, -1))
	if err != nil && !dao.tranDB.IsErrNil(err) {
		dao.logger.Error("getCollectionTransactions ZRange failed", zap.String("key", zKey), zap.Error(err))
		return nil, mpberr.ErrDB
	}
	if len(vers) == 0 {
		return nil, nil
	}
	l := len(vers)
	trans := make([]*mpb.DBTokenActivitiesV2, 0, l)
	for i := 0; i < (l+com.DBBatchNum100-1)/com.DBBatchNum100; i++ {
		start := i * com.DBBatchNum100
		end := gutil.Min(l, (i+1)*com.DBBatchNum100)
		keys := make([]string, 0, end-start)
		values := make([]*mpb.DBTokenActivitiesV2, 0, end-start)
		for ii := start; ii < end; ii++ {
			keys = append(keys, com.CollectionTranActivity(vers[ii]))
			values = append(values, &mpb.DBTokenActivitiesV2{})
		}
		err = dao.tranDB.GetObjects(ctx, keys, values)
		if err != nil {
			dao.logger.Error("getCollectionTransactions GetObjects failed",
				zap.Any("keys", keys), zap.Error(err))
			return nil, mpberr.ErrDB
		}
		trans = append(trans, values...)
	}
	return trans, nil
}
