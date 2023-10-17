package httpgateway

//import (
//	com "github.com/aureontu/MRWebServer/mr_services/common"
//	"net/http"
//
//	"github.com/oldjon/mirror/ss_services/sspb"
//)
//
//func (hg *HTTPGateway) GetUser(w http.ResponseWriter, r *http.Request) error {
//	ctx := r.Context()
//	req := &sspb.CReqGetUser{}
//	err := hg.readHTTPReq(w, r, req)
//	if err != nil {
//		return err
//	}
//
//	client, err := com.GetUserServiceClient(ctx, hg)
//	if err != nil {
//		return err
//	}
//	rpcReq := &sspb.ReqUserId{
//		UserId: req.UserId,
//	}
//	res, err := client.GetUserInfo(ctx, rpcReq)
//	if err != nil {
//		return err
//	}
//	return hg.writeHTTPRes(w, res)
//}
