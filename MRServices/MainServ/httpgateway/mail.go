package httpgateway

//import (
//	"net/http"
//
//	com "github.com/oldjon/mirror/ss_services/common"
//	"github.com/oldjon/mirror/ss_services/sspb"
//	"github.com/oldjon/mirror/ss_services/util"
//)
//
//func (hg *HTTPGateway) GetMailList(w http.ResponseWriter, r *http.Request) error {
//	claim, ctx, err := util.ClaimFromContext(r.Context())
//	if err != nil {
//		return err
//	}
//
//	client, err := com.GetMailServiceClient(ctx, hg)
//	if err != nil {
//		return err
//	}
//	res, err := client.GetMailList(ctx, &sspb.ReqUserId{
//		UserId: claim.UserId,
//	})
//	if err != nil {
//		return err
//	}
//	return hg.writeHTTPRes(w, res)
//}
//
//func (hg *HTTPGateway) ReadMails(w http.ResponseWriter, r *http.Request) error {
//	claim, ctx, err := util.ClaimFromContext(r.Context())
//	if err != nil {
//		return err
//	}
//	req := &sspb.CReqReadMails{}
//	err = hg.readHTTPReq(w, r, req)
//	if err != nil {
//		return err
//	}
//
//	client, err := com.GetMailServiceClient(ctx, hg)
//	if err != nil {
//		return err
//	}
//
//	res, err := client.ReadMails(ctx, &sspb.ReqReadMails{
//		UserId:  claim.UserId,
//		Option:  req.Option,
//		MailIds: req.MailIds,
//	})
//	if err != nil {
//		return err
//	}
//	return hg.writeHTTPRes(w, res)
//}
//
//func (hg *HTTPGateway) DelMails(w http.ResponseWriter, r *http.Request) error {
//	claim, ctx, err := util.ClaimFromContext(r.Context())
//	if err != nil {
//		return err
//	}
//	req := &sspb.CReqDelMails{}
//	err = hg.readHTTPReq(w, r, req)
//	if err != nil {
//		return err
//	}
//
//	client, err := com.GetMailServiceClient(ctx, hg)
//	if err != nil {
//		return err
//	}
//	res, err := client.DelMails(ctx, &sspb.ReqDelMails{
//		UserId:  claim.UserId,
//		Option:  req.Option,
//		MailIds: req.MailIds,
//	})
//	if err != nil {
//		return err
//	}
//	return hg.writeHTTPRes(w, res)
//}
//
//func (hg *HTTPGateway) GetMailsAwards(w http.ResponseWriter, r *http.Request) error {
//	claim, ctx, err := util.ClaimFromContext(r.Context())
//	if err != nil {
//		return err
//	}
//	req := &sspb.CReqGetMailsAwards{}
//	err = hg.readHTTPReq(w, r, req)
//	if err != nil {
//		return err
//	}
//
//	client, err := com.GetMailServiceClient(ctx, hg)
//	if err != nil {
//		return err
//	}
//	res, err := client.GetMailsAwards(ctx, &sspb.ReqGetMailsAwards{
//		UserId: claim.UserId,
//		Option: req.Option,
//		MailId: req.MailId,
//	})
//	if err != nil {
//		return err
//	}
//	return hg.writeHTTPRes(w, res)
//}
