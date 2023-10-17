package nftservice

import (
	"context"
	"fmt"
	"sort"

	com "github.com/aureontu/MRWebServer/mr_services/common"
	"github.com/aureontu/MRWebServer/mr_services/mpb"
)

func (svc *NFTService) syncCollectionTransactionsV2(ctx context.Context, collectionId string) ([]*mpb.DBTokenActivitiesV2, error) {
	pageNum := svc.rm.getGraphiQLQueryPageNum()
	transactions := make([]*mpb.AptosTransactions_TokenActivityV2, 0, pageNum)
	startIndex, err := svc.dao.getCollectionGraphiQLStartIndex(ctx, collectionId)
	if err != nil {
		return nil, err
	}
	client, err := com.GetAPIProxyGRPCClient(ctx, svc)
	if err != nil {
		return nil, err
	}
	var cnt uint64
	for {
		rpcRes, err := client.GraphiQLGetCollectionTransactions(ctx, &mpb.ReqGraphiQLGetCollectionTransactions{
			CollectionId: collectionId,
			PageNum:      pageNum,
			StartIndex:   startIndex + cnt,
		})
		if err != nil {
			return nil, err
		}

		if rpcRes.Transactions == nil || rpcRes.Transactions.Data == nil || len(rpcRes.Transactions.Data.TokenActivitiesV2) == 0 {
			break
		}
		transactions = append(transactions, rpcRes.Transactions.Data.TokenActivitiesV2...)
		cnt += uint64(len(transactions))
		if len(rpcRes.Transactions.Data.TokenActivitiesV2) < int(pageNum) {
			break
		}
	}

	fmt.Println(transactions)

	if len(transactions) == 0 {
		return nil, nil
	}

	mintEvent := svc.rm.getNFTMintEvent2()
	transferOfferEvent := svc.rm.getNFTTransferOfferEvent2()
	transferClaimEvent := svc.rm.getNFTTransferClaimEvent2()

	var curVersion uint64
	var curDBTran *mpb.DBTokenActivitiesV2
	var dbTrans = make([]*mpb.DBTokenActivitiesV2, 0)
	for _, tran := range transactions {
		if curVersion != tran.TransactionVersion {
			curVersion = tran.TransactionVersion
			curDBTran = &mpb.DBTokenActivitiesV2{}
			dbTrans = append(dbTrans, curDBTran)
		}
		curDBTran.MintEvent = curDBTran.MintEvent || (tran.Type == mintEvent)
		curDBTran.TransferOfferEvent = curDBTran.TransferOfferEvent || (tran.Type == transferOfferEvent)
		curDBTran.TransferClaimEvent = curDBTran.TransferClaimEvent || (tran.Type == transferClaimEvent)
		curDBTran.Activities = append(curDBTran.Activities,
			svc.AptosTransactionsTokenActivityV22DBTokenActivityV2(tran))
	}

	err = svc.dao.saveCollectionTransactions(ctx, collectionId, dbTrans)
	if err != nil {
		return nil, err
	}

	err = svc.dao.setCollectionGraphiQLStartIndex(ctx, collectionId, int64(cnt))
	if err != nil {
		return nil, err
	}

	return dbTrans, nil
}

func (svc *NFTService) AdminGetAptosNFTsInCollection(ctx context.Context, req *mpb.ReqAdminGetAptosNFTsInCollection) (*mpb.ResAdminGetAptosNFTsInCollection, error) {
	trans, err := svc.dao.getCollectionTransactions(ctx, req.CollectionId)
	if err != nil {
		return nil, err
	}
	newTrans, err := svc.syncCollectionTransactionsV2(ctx, req.CollectionId)
	if err != nil {
		return nil, err
	}
	trans = append(trans, newTrans...)

	depositEvent := svc.rm.getNFTDepositEvent2()
	withdrawEvent := svc.rm.getNFTWithdrawEvent2()
	addrMap := make(map[string]bool)
	tokenMap := make(map[string][4]string)

	for _, transaction := range trans {
		for _, act := range transaction.Activities {
			if act.Type == depositEvent { // get nft
				tokenMap[act.TokenDataId] = [4]string{
					act.ToAddress,
					act.TokenDataId,
					act.TokenName,
					act.TransactionTimestamp,
				}
			} else if act.Type == withdrawEvent { // lose nft
				delete(tokenMap, act.TokenDataId)
			}
		}
	}

	if len(tokenMap) == 0 {
		return &mpb.ResAdminGetAptosNFTsInCollection{}, nil
	}

	// get owners
	for _, v := range tokenMap {
		addrMap[v[0]] = true
	}

	addrs := make([]string, 0, len(addrMap))
	for k := range addrMap {
		addrs = append(addrs, k)
	}

	accClient, err := com.GetAccountServiceClient(ctx, svc)
	if err != nil {
		return nil, err
	}
	accRpcRes, err := accClient.BatchGetAccountsByWalletAddrs(ctx, &mpb.ReqBatchGetAccountsByWalletAddrs{
		Addrs: addrs,
	})
	if err != nil {
		return nil, err
	}

	accMap := make(map[string]*mpb.AccountInfo)
	for _, v := range accRpcRes.Accounts {
		accMap[v.AptosWalletAddr] = v
	}

	res := &mpb.ResAdminGetAptosNFTsInCollection{}
	for tokenDataId, info := range tokenMap {
		res.NftList = append(res.NftList, &mpb.AdminGetAptosNFTsInCollectionNode{
			TokenDataId:          tokenDataId,
			TokenId:              svc.parseTokenId(info[2]),
			TokenName:            info[2],
			TransactionTimestamp: info[3],
			Owner:                accMap[info[0]],
			OwnerAddr:            info[0],
		})
	}

	sort.Slice(res.NftList, func(i, j int) bool {
		return res.NftList[i].TokenId < res.NftList[j].TokenId
	})
	return res, nil
}

func (svc *NFTService) AdminGetCollectionNFTBuyers(ctx context.Context, req *mpb.ReqAdminGetCollectionNFTBuyers) (*mpb.ResAdminGetCollectionNFTBuyers, error) {
	trans, err := svc.dao.getCollectionTransactions(ctx, req.CollectionId)
	if err != nil {
		return nil, err
	}
	newTrans, err := svc.syncCollectionTransactionsV2(ctx, req.CollectionId)
	if err != nil {
		return nil, err
	}
	trans = append(trans, newTrans...)

	depositEvent := svc.rm.getNFTDepositEvent2()
	acts := make([]*mpb.DBTokenActivityV2, 0)
	addrMap := make(map[string]bool)
	for _, tran := range trans {
		if !tran.MintEvent {
			continue
		}
		for _, act := range tran.Activities {
			if act.Type != depositEvent {
				continue
			}
			acts = append(acts, act)
			addrMap[act.ToAddress] = true
		}
	}

	accMap := make(map[string]*mpb.AccountInfo)
	if len(addrMap) > 0 {
		addrs := make([]string, 0, len(addrMap))
		for k := range addrMap {
			addrs = append(addrs, k)
		}
		accClient, err := com.GetAccountServiceClient(ctx, svc)
		if err != nil {
			return nil, err
		}
		accRpcRes, err := accClient.BatchGetAccountsByWalletAddrs(ctx, &mpb.ReqBatchGetAccountsByWalletAddrs{
			Addrs: addrs,
		})
		if err != nil {
			return nil, err
		}
		for _, v := range accRpcRes.Accounts {
			accMap[v.AptosWalletAddr] = v
		}
	}
	res := &mpb.ResAdminGetCollectionNFTBuyers{}
	for _, act := range acts {
		res.NftList = append(res.NftList, &mpb.AdminGetCollectionNFTBuyersNode{
			TokenDataId:          act.TokenDataId,
			TokenId:              svc.parseTokenId(act.TokenName),
			TokenName:            act.TokenName,
			TransactionTimestamp: act.TransactionTimestamp,
			Buyer:                accMap[act.ToAddress],
			BuyerAddr:            act.ToAddress,
		})
	}

	sort.Slice(res.NftList, func(i, j int) bool {
		return res.NftList[i].TokenId < res.NftList[j].TokenId
	})

	return res, nil
}

func (svc *NFTService) AdminGetCollectionNFTOffers(ctx context.Context, req *mpb.ReqAdminGetCollectionNFTOffers) (*mpb.ResAdminGetCollectionNFTOffers, error) {
	trans, err := svc.dao.getCollectionTransactions(ctx, req.CollectionId)
	if err != nil {
		return nil, err
	}
	newTrans, err := svc.syncCollectionTransactionsV2(ctx, req.CollectionId)
	if err != nil {
		return nil, err
	}
	trans = append(trans, newTrans...)

	transferOfferEvent := svc.rm.getNFTTransferOfferEvent2()
	acts := make([]*mpb.DBTokenActivityV2, 0)
	addrMap := make(map[string]bool)
	for _, tran := range trans {
		if !tran.TransferOfferEvent {
			continue
		}
		for _, act := range tran.Activities {
			if act.TransactionTimestampInt < req.StartTime ||
				act.TransactionTimestampInt > req.EndTime {
				continue
			}
			if act.Type != transferOfferEvent ||
				act.ToAddress != req.DstAddr {
				continue
			}
			acts = append(acts, act)
			addrMap[act.FromAddress] = true
		}
	}

	sort.Slice(acts, func(i, j int) bool {
		return acts[i].TransactionTimestampInt < acts[j].TransactionTimestampInt
	})

	accMap := make(map[string]*mpb.AccountInfo)
	if len(addrMap) > 0 {
		addrs := make([]string, 0, len(addrMap))
		for k := range addrMap {
			addrs = append(addrs, k)
		}
		accClient, err := com.GetAccountServiceClient(ctx, svc)
		if err != nil {
			return nil, err
		}
		accRpcRes, err := accClient.BatchGetAccountsByWalletAddrs(ctx, &mpb.ReqBatchGetAccountsByWalletAddrs{
			Addrs: addrs,
		})
		if err != nil {
			return nil, err
		}
		for _, v := range accRpcRes.Accounts {
			accMap[v.AptosWalletAddr] = v
		}
	}
	res := &mpb.ResAdminGetCollectionNFTOffers{}
	for _, act := range acts {
		res.NftList = append(res.NftList, &mpb.AdminGetCollectionNFTOffersNode{
			TokenDataId:          act.TokenDataId,
			TokenId:              svc.parseTokenId(act.TokenName),
			TokenName:            act.TokenName,
			TransactionTimestamp: act.TransactionTimestamp,
			Offer:                accMap[act.FromAddress],
			OfferAddr:            act.FromAddress,
		})
	}

	return res, nil
}
