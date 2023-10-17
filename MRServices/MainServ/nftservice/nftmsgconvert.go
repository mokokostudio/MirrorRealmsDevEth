package nftservice

import (
	"strings"
	"time"

	"github.com/aureontu/MRWebServer/mr_services/mpb"
	gconv "github.com/oldjon/gutil/conv"
)

func (svc *NFTService) AptosNFTNodeV2Origin2AptosNFTNodeV2(in *mpb.AptosNFTNodeV2Origin) *mpb.AptosNFTNodeV2 {
	out := &mpb.AptosNFTNodeV2{
		CollectionId:         in.CollectionId,
		TokenDataId:          in.TokenDataId,
		Description:          in.Description,
		TokenName:            in.TokenName,
		TokenId:              svc.parseTokenId(in.TokenName),
		TokenProperties:      &mpb.AptosNFTNodeV2_Properties{},
		TokenStandard:        in.TokenStandard,
		TokenUri:             in.TokenUri,
		TransactionTimestamp: in.TransactionTimestamp,
	}
	if in.TokenProperties != nil {
		out.TokenProperties.Prop1 = in.TokenProperties["Prop1"]
		out.TokenProperties.Prop2 = in.TokenProperties["Prop2"]
		out.TokenProperties.Quality = in.TokenProperties["Quality"]
		out.TokenProperties.WeaponType = in.TokenProperties["Weapon Type"]
		out.TokenProperties.WeaponId = in.TokenProperties["Weapon ID"]
	}
	return out
}

func (svc *NFTService) parseTokenId(tokenName string) uint32 {
	strs := strings.Split(tokenName, "#")
	if len(strs) != 2 {
		return 0
	}
	return gconv.StringToUint32(strs[1])
}

func (svc *NFTService) AptosNFTNodeV22DBAptosNFTNodeV2(in *mpb.AptosNFTNodeV2) *mpb.DBAptosNFTNodeV2 {
	out := &mpb.DBAptosNFTNodeV2{
		CollectionId:         in.CollectionId,
		TokenDataId:          in.TokenDataId,
		Description:          in.Description,
		TokenName:            in.TokenName,
		TokenId:              svc.parseTokenId(in.TokenName),
		TokenProperties:      &mpb.DBAptosNFTNodeV2_Properties{},
		TokenStandard:        in.TokenStandard,
		TokenUri:             in.TokenUri,
		TransactionTimestamp: in.TransactionTimestamp,
	}
	if in.TokenProperties != nil {
		out.TokenProperties.Prop1 = in.TokenProperties.Prop1
		out.TokenProperties.Prop2 = in.TokenProperties.Prop2
		out.TokenProperties.Quality = in.TokenProperties.Quality
		out.TokenProperties.WeaponType = in.TokenProperties.WeaponType
		out.TokenProperties.WeaponId = in.TokenProperties.WeaponId
	}
	t, err := time.Parse("2006-01-02T15:04:05", out.TransactionTimestamp[:19])
	if err == nil {
		out.TransactionTimestampInt = t.Unix()
	}
	return out
}

func (svc *NFTService) DBAptosNFTNodeV22AptosNFTNodeV2(in *mpb.DBAptosNFTNodeV2) *mpb.AptosNFTNodeV2 {
	out := &mpb.AptosNFTNodeV2{
		CollectionId:         in.CollectionId,
		TokenDataId:          in.TokenDataId,
		Description:          in.Description,
		TokenName:            in.TokenName,
		TokenId:              svc.parseTokenId(in.TokenName),
		TokenProperties:      &mpb.AptosNFTNodeV2_Properties{},
		TokenStandard:        in.TokenStandard,
		TokenUri:             in.TokenUri,
		TransactionTimestamp: in.TransactionTimestamp,
	}
	if in.TokenProperties != nil {
		out.TokenProperties.Prop1 = in.TokenProperties.Prop1
		out.TokenProperties.Prop2 = in.TokenProperties.Prop2
		out.TokenProperties.Quality = in.TokenProperties.Quality
		out.TokenProperties.WeaponType = in.TokenProperties.WeaponType
		out.TokenProperties.WeaponId = in.TokenProperties.WeaponId
	}
	return out
}

func (svc *NFTService) AptosTransactionsTokenActivityV22DBTokenActivityV2(in *mpb.AptosTransactions_TokenActivityV2) *mpb.DBTokenActivityV2 {
	out := &mpb.DBTokenActivityV2{
		Type:                    in.Type,
		FromAddress:             in.FromAddress,
		ToAddress:               in.ToAddress,
		CollectionId:            in.GetCurrentTokenData().GetCollectionId(),
		TokenName:               in.GetCurrentTokenData().GetTokenName(),
		TokenDataId:             in.GetCurrentTokenData().GetTokenDataId(),
		TokenUrl:                in.GetCurrentTokenData().GetTokenUrl(),
		TransactionTimestamp:    in.TransactionTimestamp,
		TransactionTimestampInt: 0,
		TransactionVersion:      in.TransactionVersion,
	}
	t, err := time.Parse("2006-01-02T15:04:05", out.TransactionTimestamp[:19])
	if err == nil {
		out.TransactionTimestampInt = t.Unix()
	}
	return out
}
