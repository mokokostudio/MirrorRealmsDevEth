package gmgateway

import (
	"net/http"
	"strings"
	"time"

	com "github.com/aureontu/MRWebServer/mr_services/common"
	"github.com/aureontu/MRWebServer/mr_services/mpb"
	"github.com/aureontu/MRWebServer/mr_services/mpberr"
	"github.com/aureontu/MRWebServer/mr_services/util"
)

func (gg *GMGateway) adminLoginByPassword(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	req := &mpb.CReqAdminLoginByPassword{}
	err := gg.readHTTPReq(w, r, req)
	if err != nil {
		return err
	}
	if len(req.Account) < com.MinAdminAccountLength || len(req.Account) > com.MaxAdminAccountLength {
		return mpberr.ErrAdminAccountOrPasswd
	}

	req.Password = strings.ToLower(req.Password)

	if len(req.Password) != com.PasswordLen {
		return mpberr.ErrAdminAccountOrPasswd
	}

	client, err := com.GetGMServiceClient(ctx, gg)
	if err != nil {
		return err
	}
	rpcReq := mpb.ReqAdminLoginByPassword{
		Account:  req.Account,
		Password: req.Password,
	}
	res, err := client.AdminLoginByPassword(ctx, &rpcReq)
	if err != nil {
		return err
	}
	cres := &mpb.CResAdminLoginByPassword{
		Token: res.Token,
	}
	return gg.writeHTTPRes(w, cres)
}

func (gg *GMGateway) adminGetAptosNFTOwner(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	req := &mpb.CReqAdminGetAptosNFTOwner{}
	err := gg.readHTTPReq(w, r, req)
	if err != nil {
		return err
	}

	req.TokenId, err = util.FixHashId(req.TokenId)
	if err != nil {
		return err
	}

	client, err := com.GetNFTServiceClient(ctx, gg)
	if err != nil {
		return err
	}

	res, err := client.GetAptosNFTOwner(ctx, &mpb.ReqGetAptosNFTOwner{TokenId: req.TokenId})
	if err != nil {
		return err
	}
	cres := &mpb.CResAdminGetAptosNFTOwner{
		Owner: res.Owner,
	}
	return gg.writeHTTPRes(w, cres)
}

func (gg *GMGateway) adminGetAptosNFTsInCollection(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	req := &mpb.CReqAdminGetAptosNFTsInCollection{}
	err := gg.readHTTPReq(w, r, req)
	if err != nil {
		return err
	}

	req.CollectionId, err = util.FixHashId(req.CollectionId)
	if err != nil {
		return err
	}

	client, err := com.GetNFTServiceClient(ctx, gg)
	if err != nil {
		return err
	}

	res, err := client.AdminGetAptosNFTsInCollection(ctx, &mpb.ReqAdminGetAptosNFTsInCollection{CollectionId: req.CollectionId})
	if err != nil {
		return err
	}
	cres := &mpb.CResAdminGetAptosNFTsInCollection{
		NftList: res.NftList,
	}
	return gg.writeHTTPRes(w, cres)
}

func (gg *GMGateway) adminGetCollectionNFTBuyers(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	req := &mpb.CReqAdminGetCollectionNFTBuyers{}
	err := gg.readHTTPReq(w, r, req)
	if err != nil {
		return err
	}

	req.CollectionId, err = util.FixHashId(req.CollectionId)
	if err != nil {
		return err
	}

	client, err := com.GetNFTServiceClient(ctx, gg)
	if err != nil {
		return err
	}

	res, err := client.AdminGetCollectionNFTBuyers(ctx, &mpb.ReqAdminGetCollectionNFTBuyers{CollectionId: req.CollectionId})
	if err != nil {
		return err
	}
	cres := &mpb.CResAdminGetCollectionNFTBuyers{
		NftList: res.NftList,
	}
	return gg.writeHTTPRes(w, cres)
}

func (gg *GMGateway) adminGetCollectionNFTOffers(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	req := &mpb.CReqAdminGetCollectionNFTOffers{}
	err := gg.readHTTPReq(w, r, req)
	if err != nil {
		return err
	}

	req.CollectionId, err = util.FixHashId(req.CollectionId)
	if err != nil {
		return err
	}
	req.DstAddr, err = util.FixHashId(req.DstAddr)
	if err != nil {
		return err
	}
	if req.StartTime == "" {
		req.StartTime = "2000-01-01 00:00:00"
	}
	if req.EndTime == "" {
		req.EndTime = "2100-12-01 23:59:59"
	}
	startTime, err := time.Parse("2006-01-02 15:04:05", req.StartTime)
	if err != nil {
		return mpberr.ErrParam
	}

	endTime, err := time.Parse("2006-01-02 15:04:05", req.EndTime)
	if err != nil {
		return mpberr.ErrParam
	}

	client, err := com.GetNFTServiceClient(ctx, gg)
	if err != nil {
		return err
	}

	res, err := client.AdminGetCollectionNFTOffers(ctx,
		&mpb.ReqAdminGetCollectionNFTOffers{
			CollectionId: req.CollectionId,
			DstAddr:      req.DstAddr,
			StartTime:    startTime.Unix(),
			EndTime:      endTime.Unix(),
		})
	if err != nil {
		return err
	}
	cres := &mpb.CResAdminGetCollectionNFTOffers{
		NftList: res.NftList,
	}
	return gg.writeHTTPRes(w, cres)
}
