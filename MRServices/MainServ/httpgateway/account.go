package httpgateway

import (
	"net/http"
	"strconv"
	"strings"

	com "github.com/aureontu/MRWebServer/mr_services/common"
	"github.com/aureontu/MRWebServer/mr_services/mpb"
	"github.com/aureontu/MRWebServer/mr_services/mpberr"
	"github.com/aureontu/MRWebServer/mr_services/util"
)

func (hg *HTTPGateway) loginByPassword(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	req := &mpb.CReqLoginByPassword{}
	err := hg.readHTTPReq(w, r, req)
	if err != nil {
		return err
	}

	req.Password = strings.ToLower(req.Password)

	if len(req.Password) != com.PasswordLen {
		return mpberr.ErrPassword
	}

	remoteIP := getRemoteIPAddress(r)

	client, err := com.GetAccountServiceClient(ctx, hg)
	if err != nil {
		return err
	}
	rpcReq := mpb.ReqLoginByPassword{
		Account:  req.Account,
		Password: req.Password,
		RemoteIp: remoteIP,
		Region:   getRegionByIP(remoteIP),
	}
	res, err := client.LoginByPassword(ctx, &rpcReq)
	if err != nil {
		return err
	}
	cres := &mpb.CResLoginByPassword{
		Account:   res.Account,
		Resources: res.Resources,
		Token:     res.Token,
	}
	return hg.writeHTTPRes(w, cres)
}

func (hg *HTTPGateway) WebLoginByWallet(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	req := &mpb.CReqWebLoginByWallet{}
	err := hg.readHTTPReq(w, r, req)
	if err != nil {
		return err
	}
	req.WalletAddr, err = util.FixHashId(req.WalletAddr)
	if err != nil {
		return mpberr.ErrParam
	}
	if req.WalletAddr == "" || req.PubKey == "" {
		return mpberr.ErrParam
	}

	nonce := util.ReadNonceFromAptosFullMsg(req.AptosFullMsg)
	if len(nonce) != com.NonceLen {
		return mpberr.ErrParam
	}

	pubKey, err := util.DecodeAptosPubKey(req.PubKey)
	if err != nil {
		return mpberr.ErrAptosPublicKey
	}

	if !util.VerifySignature(pubKey, req.AptosFullMsg, req.AptosSignature) {
		return mpberr.ErrAptosVerifySignature
	}

	remoteIP := getRemoteIPAddress(r)

	client, err := com.GetAccountServiceClient(ctx, hg)
	if err != nil {
		return err
	}
	rpcReq := mpb.ReqWebLoginByWallet{
		WalletAddr: req.WalletAddr,
		PubKey:     pubKey,
		Nonce:      nonce,
		RemoteIp:   remoteIP,
		Region:     getRegionByIP(remoteIP),
	}
	res, err := client.WebLoginByWallet(ctx, &rpcReq)
	if err != nil {
		return err
	}
	cres := &mpb.CResWebLoginByWallet{
		Account:   res.Account,
		Resources: res.Resources,
		Token:     res.Token,
	}
	return hg.writeHTTPRes(w, cres)
}

func (hg *HTTPGateway) generateNonce(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	client, err := com.GetAccountServiceClient(ctx, hg)
	if err != nil {
		return err
	}
	res, err := client.GenerateNonce(ctx, &mpb.Empty{})
	if err != nil {
		return err
	}
	cres := &mpb.CResGenerateNonce{
		Nonce: res.Nonce,
	}
	return hg.writeHTTPRes(w, cres)
}

func (hg *HTTPGateway) SendEmailBindCode(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	req := &mpb.CReqSendEmailBindCode{}
	err := hg.readHTTPReq(w, r, req)
	if err != nil {
		return err
	}

	if !util.CheckEmailAddr(req.Email) {
		return mpberr.ErrEmailAddress
	}

	client, err := com.GetAccountServiceClient(ctx, hg)
	if err != nil {
		return err
	}
	_, err = client.GenerateAndSendEmailBindCode(ctx, &mpb.ReqGenerateAndSendEmailBindCode{Email: req.Email})
	if err != nil {
		return err
	}

	return hg.writeHTTPRes(w, &mpb.Empty{})
}

func (hg *HTTPGateway) webBindEmail(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	claim, ctx, err := util.ClaimFromContext(r.Context())
	if err != nil {
		return err
	}

	req := &mpb.CReqWebBindEmail{}
	err = hg.readHTTPReq(w, r, req)
	if err != nil {
		return err
	}

	if len(req.Code) != com.VCodeLen {
		return mpberr.ErrEmailBindCode
	}

	_, err = strconv.Atoi(req.Code)
	if err != nil {
		return mpberr.ErrEmailBindCode
	}

	if !util.CheckEmailAddr(req.Email) {
		return mpberr.ErrEmailAddress
	}

	client, err := com.GetAccountServiceClient(ctx, hg)
	if err != nil {
		return err
	}
	res, err := client.WebBindEmail(ctx, &mpb.ReqWebBindEmail{UserId: claim.UserId, Email: req.Email, Code: req.Code})
	if err != nil {
		return err
	}
	cres := &mpb.CResWebBindEmail{
		Account:   res.Account,
		Resources: res.Resources,
		Token:     res.Token,
	}
	return hg.writeHTTPRes(w, cres)
}

func (hg *HTTPGateway) getAccountInfo(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	claim, ctx, err := util.ClaimFromContext(r.Context())
	if err != nil {
		return err
	}

	client, err := com.GetAccountServiceClient(ctx, hg)
	if err != nil {
		return err
	}

	res, err := client.GetAccountInfo(ctx, &mpb.ReqUserId{UserId: claim.UserId})
	if err != nil {
		return err
	}
	cres := &mpb.CResGetAccountInfo{
		Account: &mpb.AccountInfo{
			Account:  res.Account,
			UserId:   res.UserId,
			Nickname: res.Nickname,
			Icon:     res.Icon,
			Email:    res.Email,
		},
	}
	return hg.writeHTTPRes(w, cres)
}

func (hg *HTTPGateway) changePassword(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	claim, ctx, err := util.ClaimFromContext(r.Context())
	if err != nil {
		return err
	}

	req := &mpb.CReqChangePassword{}
	err = hg.readHTTPReq(w, r, req)
	if err != nil {
		return err
	}
	req.NewPassword = strings.ToLower(req.NewPassword)
	req.OldPassword = strings.ToLower(req.OldPassword)
	if len(req.NewPassword) != com.PasswordLen {
		return mpberr.ErrParam
	}

	if req.OldPassword == req.NewPassword {
		return mpberr.ErrNewPWSameWithOldPW
	}

	client, err := com.GetAccountServiceClient(ctx, hg)
	if err != nil {
		return err
	}

	res, err := client.ChangePassword(ctx, &mpb.ReqChangePassword{
		UserId:      claim.UserId,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		return err
	}
	return hg.writeHTTPRes(w, res)
}

func (hg *HTTPGateway) sendEmailResetPasswordCode(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	req := &mpb.CReqSendEmailResetPasswordCode{}
	err := hg.readHTTPReq(w, r, req)
	if err != nil {
		return err
	}
	if util.CheckEmailAddr(req.Email) {
		return mpberr.ErrEmailAddress
	}

	client, err := com.GetAccountServiceClient(ctx, hg)
	if err != nil {
		return err
	}

	_, err = client.SendEmailResetPasswordCode(ctx, &mpb.ReqSendEmailResetPasswordCode{
		Email: req.Email,
	})
	if err != nil {
		return err
	}

	return hg.writeHTTPRes(w, &mpb.Empty{})
}

func (hg *HTTPGateway) checkEmailResetPasswordCode(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	req := &mpb.CReqCheckEmailResetPasswordCode{}
	err := hg.readHTTPReq(w, r, req)
	if err != nil {
		return err
	}

	if len(req.Code) != com.VCodeLen {
		return mpberr.ErrEmailVerificationCode
	}

	if util.CheckEmailAddr(req.Email) {
		return mpberr.ErrEmailAddress
	}

	client, err := com.GetAccountServiceClient(ctx, hg)
	if err != nil {
		return err
	}

	res, err := client.CheckEmailResetPasswordCode(ctx, &mpb.ReqCheckEmailResetPasswordCode{
		Email: req.Email,
		Code:  req.Code,
	})
	if err != nil {
		return err
	}

	return hg.writeHTTPRes(w, &mpb.CResCheckEmailResetPasswordCode{Nonce: res.Nonce})
}

func (hg *HTTPGateway) resetPasswordByEmail(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	req := &mpb.CReqResetPasswordByEmail{}
	err := hg.readHTTPReq(w, r, req)
	if err != nil {
		return err
	}

	if len(req.Nonce) != com.NonceLen {
		return mpberr.ErrParam
	}

	req.Password = strings.ToLower(req.Password)
	if len(req.Password) != com.PasswordLen {
		return mpberr.ErrParam
	}

	if util.CheckEmailAddr(req.Email) {
		return mpberr.ErrEmailAddress
	}

	client, err := com.GetAccountServiceClient(ctx, hg)
	if err != nil {
		return err
	}

	_, err = client.ResetPasswordByEmail(ctx, &mpb.ReqResetPasswordByEmail{
		Email:    req.Email,
		Password: req.Password,
		Nonce:    req.Nonce,
	})
	if err != nil {
		return err
	}

	return hg.writeHTTPRes(w, &mpb.Empty{})
}

func (hg *HTTPGateway) resetPasswordByEmailAndVCode(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	req := &mpb.CReqResetPasswordByEmailAndVCode{}
	err := hg.readHTTPReq(w, r, req)
	if err != nil {
		return err
	}

	if len(req.Code) != com.VCodeLen {
		return mpberr.ErrParam
	}

	req.Password = strings.ToLower(req.Password)
	if len(req.Password) != com.PasswordLen {
		return mpberr.ErrParam
	}

	if util.CheckEmailAddr(req.Email) {
		return mpberr.ErrEmailAddress
	}

	client, err := com.GetAccountServiceClient(ctx, hg)
	if err != nil {
		return err
	}

	_, err = client.ResetPasswordByEmailAndVCode(ctx, &mpb.ReqResetPasswordByEmailAndVCode{
		Email:    req.Email,
		Password: req.Password,
		Code:     req.Code,
	})
	if err != nil {
		return err
	}

	return hg.writeHTTPRes(w, &mpb.Empty{})
}

func (hg *HTTPGateway) getAptosResources(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	claim, ctx, err := util.ClaimFromContext(r.Context())
	if err != nil {
		return err
	}

	clientApp, err := com.GetAPIProxyGRPCClient(ctx, hg)
	if err != nil {
		return err
	}
	res, err := clientApp.GetAptosResources(ctx, &mpb.ReqGetAptosResources{AptosAccAddr: claim.WalletAddr})
	if err != nil {
		return err
	}
	cres := &mpb.CResGetAptosResources{
		Resources: res.Resources,
	}
	return hg.writeHTTPRes(w, cres)
}
