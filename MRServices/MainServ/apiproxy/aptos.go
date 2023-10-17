package apiproxy

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/aureontu/MRWebServer/mr_services/mpb"
	"github.com/aureontu/MRWebServer/mr_services/util"
	"go.uber.org/zap"
)

const (
	ApiGetAccount          = "/v1/accounts/%s"
	ApiGetAccountResources = "/v1/accounts/%s/resources"
	ApiGetAccountResource  = "/v1/accounts/%s/resource/%s"
	ApiGraphiQL            = "/v1/graphql"
)

const (
	MoralisApiGetNFTsByWallets = "/wallets/nfts?%s"
)

type AptosManager struct {
	logger      *zap.Logger
	aptosUrl    string
	moralisUrl  string
	graphiqlUrl string
}

func newAptosManager(aptosUrl, moralisUrl, graphiqlUrl string, logger *zap.Logger) *AptosManager {
	return &AptosManager{
		aptosUrl:    aptosUrl,
		moralisUrl:  moralisUrl,
		graphiqlUrl: graphiqlUrl,
		logger:      logger,
	}
}

func (aptos *AptosManager) GetAccount(ctx context.Context, addr string) ([]byte, error) {
	data, err := util.HttpsGet(ctx, fmt.Sprintf(aptos.aptosUrl+ApiGetAccount, addr))
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (aptos *AptosManager) GetAccountResources(ctx context.Context, addr string) ([]byte, error) {
	data, err := util.HttpsGet(ctx, fmt.Sprintf(aptos.aptosUrl+ApiGetAccountResources, addr))
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (aptos *AptosManager) MoralisGetNFTByWallets(ctx context.Context, addresses, collectionsHash []string, apiKey string) (map[string][]*mpb.MoralisNFTData, error) {
	var nfts = make(map[string][]*mpb.MoralisNFTData, 0)
	var cursor string
	for {
		values := url.Values{}
		values.Add("limit", "100")
		if cursor != "" {
			values.Add("cursor", cursor)
		}
		for i, v := range addresses {
			values.Add(fmt.Sprintf("owner_addresses[%d]", i), v)
		}
		for i, v := range collectionsHash {
			values.Add(fmt.Sprintf("collection_whitelist[%d]", i), v)
		}
		headers := make(map[string]string)
		headers["Accept"] = "application/json"
		headers["X-API-Key"] = apiKey
		bys, err := util.HttpGetWithHeader(ctx, aptos.moralisUrl+fmt.Sprintf(MoralisApiGetNFTsByWallets, values.Encode()), headers)
		if err != nil {
			return nil, err
		}

		n := &mpb.MoralisNFTsData{}
		err = json.Unmarshal(bys, n)
		if err != nil {
			return nil, err
		}

		for _, v := range n.Result {
			nfts[v.OwnerAddress] = append(nfts[v.OwnerAddress], v)
		}

		cursor = n.Cursor
		if cursor == "" {
			break
		}
	}
	return nfts, nil
}

func (aptos *AptosManager) GraphiQLGetAccountTransactions(ctx context.Context, addr string, startIndex, pageNum uint64) (*mpb.AptosAccountTransactions, error) {
	headers := make(map[string]string)
	headers["Accept"] = "application/json"
	headers["Content-Type"] = "application/json"
	query := fmt.Sprintf(`{"query":"query MyQuery {\n  account_transactions(\n    where: {account_address: {_eq: \"%s\"}}\n    limit: %d\n    offset: %d\n    order_by: {transaction_version: asc}\n  ) {\n    transaction_version\n    token_activities_v2 {\n      entry_function_id_str\n      event_account_address\n      event_index\n      from_address\n      to_address\n      token_amount\n      token_data_id\n      token_standard\n      transaction_timestamp\n      transaction_version\n      type\n      current_token_data {\n        collection_id\n        description\n        last_transaction_timestamp\n        last_transaction_version\n        token_data_id\n        token_name\n        token_properties\n        token_standard\n        token_uri\n      }\n    }\n  }\n}\n","variables":null,"operationName":"MyQuery"}`, addr, pageNum, startIndex)

	bys, err := util.HttpPost(ctx, aptos.graphiqlUrl+ApiGraphiQL, headers, []byte(query))
	if err != nil {
		aptos.logger.Error("GraphiQLGetAccountTransactions http post failed", zap.Error(err))
		return nil, err
	}
	res := &mpb.AptosAccountTransactions{}

	err = json.Unmarshal(bys, res)
	if err != nil {
		aptos.logger.Error("GraphiQLGetAccountTransactions json unmarshal failed", zap.Error(err))
		return nil, err
	}

	return res, nil
}

func (aptos *AptosManager) GraphiQLGetCollectionTransactions(ctx context.Context, collectionId string, startIndex, pageNum uint64) (*mpb.AptosTransactions, error) {
	headers := make(map[string]string)
	headers["Accept"] = "application/json"
	headers["Content-Type"] = "application/json"
	query := fmt.Sprintf(`{"query":"query MyQuery {\n  token_activities_v2(\n    where: {current_token_data: {collection_id: {_eq: \"%s\"}}}\n    limit: %d\n    offset: %d\n  ) {\n    type\n    from_address\n    to_address\n    current_token_data {\n      collection_id\n      token_name\n      token_data_id\n    token_uri\n  }\n    transaction_timestamp\n    transaction_version\n  }\n}\n","variables":null,"operationName":"MyQuery"}`, collectionId, pageNum, startIndex)

	bys, err := util.HttpPost(ctx, aptos.graphiqlUrl+ApiGraphiQL, headers, []byte(query))
	if err != nil {
		aptos.logger.Error("GraphiQLGetCollectionTransactions http post failed", zap.Error(err))
		return nil, err
	}
	res := &mpb.AptosTransactions{}

	err = json.Unmarshal(bys, res)
	if err != nil {
		aptos.logger.Error("GraphiQLGetCollectionTransactions json unmarshal failed", zap.Error(err))
		return nil, err
	}

	return res, nil
}
