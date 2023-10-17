package httpgateway

import (
	"github.com/aureontu/MRWebServer/mr_services/util"
	"net/http"

	com "github.com/aureontu/MRWebServer/mr_services/common"
	"github.com/aureontu/MRWebServer/mr_services/mpb"
)

func (hg *HTTPGateway) getAptosNFTs(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	claim, ctx, err := util.ClaimFromContext(r.Context())
	if err != nil {
		return err
	}

	req := &mpb.CReqGetAptosNFTs{}
	err = hg.readHTTPReq(w, r, req)
	if err != nil {
		return err
	}

	clientApp, err := com.GetNFTServiceClient(ctx, hg)
	if err != nil {
		return err
	}
	addr := claim.WalletAddr
	//addr := "0x7cb76b35287055db9dc1a743961dc1191b4261fad92a1b76f129f5ee9a7b5aa5"
	rpcReq := &mpb.ReqGetAptosNFTs{WalletAddr: addr}
	for _, v := range req.NftTypes {
		rpcReq.NftTypes = append(rpcReq.NftTypes, mpb.ENFT_NFTType(v))
	}
	res, err := clientApp.GetAptosNFTs(ctx, rpcReq)
	if err != nil {
		return err
	}
	cres := &mpb.CResGetAptosNFTs{
		Nfts: res.Nfts,
	}
	return hg.writeHTTPRes(w, cres)
}

func (hg *HTTPGateway) getAptosNFTMetaDatas(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	req := &mpb.CReqGetAptosNFTMetadatas{}
	err := hg.readHTTPReq(w, r, req)
	if err != nil {
		return err
	}

	clientApp, err := com.GetNFTServiceClient(ctx, hg)
	if err != nil {
		return err
	}

	//addr := "0x7cb76b35287055db9dc1a743961dc1191b4261fad92a1b76f129f5ee9a7b5aa5"
	rpcReq := &mpb.ReqGetAptosNFTMetadatas{NftIds: req.NftIds}
	res, err := clientApp.GetAptosNFTMetadatas(ctx, rpcReq)
	if err != nil {
		return err
	}
	cres := &mpb.CResGetAptosNFTMetadatas{
		Metadatas: res.Metadatas,
	}
	return hg.writeHTTPRes(w, cres)
}

func (hg *HTTPGateway) getAptosNFTsV2(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	claim, ctx, err := util.ClaimFromContext(r.Context())
	if err != nil {
		return err
	}

	req := &mpb.CReqGetAptosNFTsV2{}
	err = hg.readHTTPReq(w, r, req)
	if err != nil {
		return err
	}

	clientApp, err := com.GetNFTServiceClient(ctx, hg)
	if err != nil {
		return err
	}

	//addr := "0x7cb76b35287055db9dc1a743961dc1191b4261fad92a1b76f129f5ee9a7b5aa5"
	rpcReq := &mpb.ReqGetAptosNFTsV2{
		UserId:     claim.UserId,
		WalletAddr: claim.WalletAddr,
	}
	rpcRes, err := clientApp.GetAptosNFTsV2(ctx, rpcReq)
	if err != nil {
		return err
	}
	cres := &mpb.CResGetAptosNFTsV2{
		Nfts: rpcRes.Nfts,
	}
	return hg.writeHTTPRes(w, cres)
}

func (hg *HTTPGateway) testGetAptosNFTsV2(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	clientApp, err := com.GetNFTServiceClient(ctx, hg)
	if err != nil {
		return err
	}

	addr := "0x2d92cbd4d2da3e63138dae1a057538268fb41cc90a2a94fe0768bcb71cfda863"
	rpcReq := &mpb.ReqGetAptosNFTsV2{UserId: 10000001, WalletAddr: addr}
	rpcRes, err := clientApp.GetAptosNFTsV2(ctx, rpcReq)
	if err != nil {
		return err
	}
	cres := &mpb.CResGetAptosNFTsV2{
		Nfts: rpcRes.Nfts,
	}
	return hg.writeHTTPRes(w, cres)
}
