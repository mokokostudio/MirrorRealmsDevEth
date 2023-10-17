package httpgateway

//import (
//	com "github.com/oldjon/mirror/ss_services/common"
//	"net/http"
//
//	"github.com/oldjon/mirror/ss_services/sspb"
//	"github.com/oldjon/mirror/ss_services/util"
//)
//
//func (hg *HTTPGateway) GetItems(w http.ResponseWriter, r *http.Request) error {
//	claim, ctx, err := util.ClaimFromContext(r.Context())
//	if err != nil {
//		return err
//	}
//	rpcReq := &sspb.ReqUserIdRegion{
//		UserId: claim.UserId,
//		Region: claim.Region,
//	}
//	client, err := com.GetItemServiceClient(ctx, hg)
//	if err != nil {
//		return err
//	}
//	res, err := client.GetItems(ctx, rpcReq)
//	if err != nil {
//		return err
//	}
//	return hg.writeHTTPRes(w, res)
//}
//
//func (hg *HTTPGateway) ExchangeItems(w http.ResponseWriter, r *http.Request) error {
//	claim, ctx, err := util.ClaimFromContext(r.Context())
//	if err != nil {
//		return err
//	}
//
//	req := &sspb.CReqExchangeItems{}
//	err = hg.readHTTPReq(w, r, req)
//	if err != nil {
//		return err
//	}
//
//	rpcReq := &sspb.ReqExchangeItems{
//		UserId:   claim.UserId,
//		Region:   claim.Region,
//		AddItems: req.AddItems,
//		DelItems: req.DelItems,
//	}
//
//	client, err := com.GetItemServiceClient(ctx, hg)
//	if err != nil {
//		return err
//	}
//	res, err := client.ExchangeItems(ctx, rpcReq)
//	if err != nil {
//		return err
//	}
//	return hg.writeHTTPRes(w, &sspb.CResExchangeItems{
//		ShowItems: res.ShowItems,
//	})
//}
