package accountservice

import (
	"context"
	"time"

	com "github.com/aureontu/MRWebServer/mr_services/common"
	"github.com/aureontu/MRWebServer/mr_services/mpb"
	"github.com/aureontu/MRWebServer/mr_services/mpberr"
	"github.com/aureontu/MRWebServer/mr_services/util"
	"github.com/golang-jwt/jwt/v5"
	gcrypto "github.com/oldjon/gutil/crypto"
	"github.com/oldjon/gutil/env"
	"github.com/oldjon/gutil/gdb"
	gprotocol "github.com/oldjon/gutil/protocol"
	grmux "github.com/oldjon/gutil/redismutex"
	gxgrpc "github.com/oldjon/gx/modules/grpc"
	"github.com/oldjon/gx/service"
	"github.com/pkg/errors"
	etcd "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type AccountService struct {
	mpb.UnimplementedAccountServiceServer
	name            string
	logger          *zap.Logger
	config          env.ModuleConfig
	etcdClient      *etcd.Client
	host            service.Host
	connMgr         *gxgrpc.ConnManager
	signingMethod   jwt.SigningMethod
	signingDuration time.Duration
	signingKey      []byte
	rm              *AccountResourceMgr
	kvm             *service.KVMgr
	serverEnv       uint32
	sm              *util.ServiceMetrics
	dao             *accountDAO
	tcpMsgCoder     gprotocol.FrameCoder
}

// NewAccountService create an accountservice entity
func NewAccountService(driver service.ModuleDriver) (gxgrpc.GRPCServer, error) {
	svc := &AccountService{
		name:            driver.ModuleName(),
		logger:          driver.Logger(),
		config:          driver.ModuleConfig(),
		etcdClient:      driver.Host().EtcdSession().Client(),
		host:            driver.Host(),
		kvm:             driver.Host().KVManager(),
		sm:              util.NewServiceMetrics(driver),
		signingMethod:   jwt.SigningMethodHS256,
		signingDuration: 24 * 30 * time.Hour,
	}

	dialer := gxgrpc.Dialer{
		HostName:   driver.Host().Name(),
		Tracer:     driver.Tracer(),
		EtcdClient: svc.etcdClient,
		Logger:     svc.logger,
		EnableTLS:  svc.config.GetBool("enable_tls"),
		CAFile:     svc.config.GetString("ca_file"),
		CertFile:   svc.config.GetString("cert_file"),
		KeyFile:    svc.config.GetString("key_file"),
	}
	svc.connMgr = gxgrpc.NewConnManager(&dialer)

	var err error
	svc.rm, err = newAccountResourceMgr(svc.logger, svc.sm)
	if err != nil {
		return nil, err
	}

	redisMux, err := grmux.NewRedisMux(svc.config.SubConfig("redis_mutex"), nil, svc.logger, driver.Tracer())
	if err != nil {
		return nil, err
	}

	accRedis, err := gdb.NewRedisClientByConfig(svc.config.SubConfig("acc_redis"),
		svc.config.GetString("db_marshaller"), driver.Tracer())
	if err != nil {
		return nil, err
	}

	tmpRedis, err := gdb.NewRedisClientByConfig(svc.config.SubConfig("tmp_redis"),
		svc.config.GetString("db_marshaller"), driver.Tracer())
	if err != nil {
		return nil, err
	}

	svc.dao = newAccountDAO(svc.logger, redisMux, accRedis, tmpRedis)

	svc.serverEnv = uint32(svc.config.GetInt64("server_env"))
	svc.tcpMsgCoder = gprotocol.NewFrameCoder(svc.config.GetString("protocol_code"))

	return svc, err
}

func (svc *AccountService) Register(grpcServer *grpc.Server) {
	mpb.RegisterAccountServiceServer(grpcServer, svc)
}

func (svc *AccountService) Serve(ctx context.Context) error {

	signingKey, err := svc.kvm.GetOrGenerate(ctx, com.JWTGatewayTokenKey, 32)
	if err != nil {
		return errors.WithStack(err)
	}
	svc.signingKey = signingKey

	<-ctx.Done()
	return ctx.Err()
}

func (svc *AccountService) Logger() *zap.Logger {
	return svc.logger
}

func (svc *AccountService) ConnMgr() *gxgrpc.ConnManager {
	return svc.connMgr
}

func (svc *AccountService) Name() string {
	return svc.name
}

//func (svc *AccountService) RegisterAccount(ctx context.Context, req *mpb.ReqRegisterAccount) (*mpb.AccountInfo, error) {
//	if req.DeviceId == "" {
//		return nil, mpberr.ErrOk
//	}
//
//	accInfo, err := svc.dao.registerAccount(ctx, req.Account, req.Password, req.Device, req.DeviceId, req.Os, req.Region)
//	if err != nil {
//		return nil, err
//	}
//
//	svc.logger.Info("Register Account", zap.String("account", accInfo.Account), zap.Uint64("user_id", accInfo.UserId),
//		zap.String("device", req.Device), zap.String("os", req.Os), zap.String("device_id", req.DeviceId),
//		zap.String("client_version", req.ClientVersion), zap.String("region", req.Region),
//		zap.String("ip", req.RemoteIp))
//
//	return accInfo, nil
//}

func (svc *AccountService) LoginByPassword(ctx context.Context, req *mpb.ReqLoginByPassword) (*mpb.ResLoginByPassword,
	error) {
	if req.Account == "" {
		return nil, mpberr.ErrParam
	}

	dbAcc, err := svc.dao.getAccountInfoByAccount(ctx, req.Account)
	if err != nil {
		return nil, err
	}

	if dbAcc.Password != req.Password {
		return nil, mpberr.ErrPassword
	}

	res := &mpb.ResLoginByPassword{
		Account: svc.DBAccountInfo2AccountInfo(dbAcc),
	}

	token, err := svc.generateLoginToken(dbAcc.UserId, dbAcc.Account, &mpb.Region{
		Region: "", TcpGatewayId: "0"}, "", dbAcc.AptosAccAddr)
	if err != nil {
		return nil, err
	}

	err = svc.dao.saveToken(ctx, dbAcc.UserId, gcrypto.MD5SumStr(token), dbAcc.Account, "", "PC")
	if err != nil {
		return nil, err
	}

	res.Token = token

	// get wallet resource
	res.Resources, err = svc.getAptosResources(ctx, dbAcc.AptosAccAddr)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (svc *AccountService) GetAccountInfo(ctx context.Context, req *mpb.ReqUserId) (*mpb.AccountInfo, error) {
	if req.UserId == 0 {
		return nil, mpberr.ErrParam
	}

	accInfo, err := svc.dao.getAccountInfo(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return svc.DBAccountInfo2AccountInfo(accInfo), nil
}

func (svc *AccountService) GetAccountInfoByAccount(ctx context.Context, req *mpb.ReqGetAccountInfoByAccount) (*mpb.AccountInfo,
	error) {
	if len(req.Account) == 0 {
		return nil, mpberr.ErrParam
	}

	accInfo, err := svc.dao.getAccountInfoByAccount(ctx, req.Account)
	if err != nil {
		return nil, err
	}

	return svc.DBAccountInfo2AccountInfo(accInfo), nil
}

func (svc *AccountService) GenerateNonce(ctx context.Context, _ *mpb.Empty) (*mpb.ResGenerateNonce, error) {
	res := &mpb.ResGenerateNonce{
		Nonce: util.GenerateRandomCode(com.NonceLen),
	}

	err := svc.dao.saveNonce(ctx, res.Nonce)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *AccountService) generateLoginToken(userId uint64, account string, region *mpb.Region,
	clientVersion, walletAddr string) (string, error) {
	var sToken string
	now := time.Now()
	claim := &mpb.JWTClaims{}
	claim.UserId = userId
	claim.Region = region
	claim.Account = account
	claim.ClientVersion = clientVersion
	claim.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(now.Add(svc.signingDuration)),
	}
	claim.WalletAddr = walletAddr
	token := jwt.NewWithClaims(svc.signingMethod, claim)
	sToken, err := token.SignedString(svc.signingKey)
	if err != nil {
		return "", err
	}
	return sToken, nil
}

func (svc *AccountService) WebLoginByWallet(ctx context.Context, req *mpb.ReqWebLoginByWallet) (*mpb.ResWebLoginByWallet,
	error) {
	// check nonce
	ok, err := svc.dao.checkNonce(ctx, req.Nonce)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, mpberr.ErrRepeatedRequest
	}

	dbAcc, err := svc.dao.getAccountByWallet(ctx, req.WalletAddr, req.PubKey)
	if err != nil {
		return nil, err
	}
	res := &mpb.ResWebLoginByWallet{
		Account: &mpb.AccountInfo{
			Account:         dbAcc.Account,
			UserId:          dbAcc.UserId,
			Nickname:        dbAcc.Nickname,
			Icon:            dbAcc.Icon,
			Email:           dbAcc.Email,
			AptosWalletAddr: dbAcc.AptosAccAddr,
		},
	}

	token, err := svc.generateLoginToken(dbAcc.UserId, dbAcc.Account, &mpb.Region{
		Region: "", TcpGatewayId: "0"}, "", dbAcc.AptosAccAddr)
	if err != nil {
		return nil, err
	}

	err = svc.dao.saveToken(ctx, dbAcc.UserId, gcrypto.MD5SumStr(token), dbAcc.Account, "", "PC")
	if err != nil {
		return nil, err
	}

	res.Token = token

	if dbAcc.Account == "" { // acc not bind email
		return res, nil
	}

	// get wallet account resource
	res.Resources, err = svc.getAptosResources(ctx, req.WalletAddr)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *AccountService) GenerateAndSendEmailBindCode(ctx context.Context, req *mpb.ReqGenerateAndSendEmailBindCode) (
	*mpb.Empty, error) {

	ok, err := svc.dao.checkEmailExist(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if ok {
		return nil, mpberr.ErrEmailBound
	}

	ok, err = svc.dao.checkEmailSendLimit(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, mpberr.ErrEmailSendMax
	}

	rpcReq := &mpb.ReqSendEmailBindCode{
		Email: req.Email,
		Code:  util.GenerateRandomCode(com.VCodeLen),
	}

	err = svc.dao.saveEmailBindCode(ctx, req.Email, rpcReq.Code)
	if err != nil {
		return nil, err
	}

	client, err := com.GetAPIProxyGRPCClient(ctx, svc)
	if err != nil {
		return nil, err
	}

	_, err = client.SendEmailBindCode(ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	return &mpb.Empty{}, nil
}

func (svc *AccountService) WebBindEmail(ctx context.Context, req *mpb.ReqWebBindEmail) (*mpb.ResWebBindEmail, error) {
	// check email bind code
	ok, err := svc.dao.checkEmailBindCode(ctx, req.Email, req.Code)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, mpberr.ErrEmailBindCode
	}
	// check email is bound by other else
	ok, err = svc.dao.checkEmailExist(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if ok {
		return nil, mpberr.ErrEmailBound
	}

	dbAcc, err := svc.dao.bindEmail(ctx, req.UserId, req.Email, req.Email)
	if err != nil {
		return nil, err
	}

	res := &mpb.ResWebBindEmail{
		Account: &mpb.AccountInfo{
			Account:         dbAcc.Account,
			UserId:          dbAcc.UserId,
			Nickname:        dbAcc.Nickname,
			Icon:            dbAcc.Icon,
			Email:           dbAcc.Email,
			AptosWalletAddr: dbAcc.AptosAccAddr,
		},
	}

	token, err := svc.generateLoginToken(dbAcc.UserId, dbAcc.Account, &mpb.Region{
		Region: "", TcpGatewayId: "0"}, "", dbAcc.AptosAccAddr)
	if err != nil {
		return nil, err
	}

	err = svc.dao.saveToken(ctx, dbAcc.UserId, gcrypto.MD5SumStr(token), dbAcc.Account, "", "PC")
	if err != nil {
		return nil, err
	}

	res.Token = token

	// get wallet resource
	res.Resources, err = svc.getAptosResources(ctx, dbAcc.AptosAccAddr)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *AccountService) ChangePassword(ctx context.Context, req *mpb.ReqChangePassword) (*mpb.Empty, error) {
	if req.OldPassword == req.NewPassword {
		return nil, mpberr.ErrNewPWSameWithOldPW
	}
	err := svc.dao.changePassword(ctx, req.UserId, req.OldPassword, req.NewPassword)
	if err != nil {
		return nil, err
	}
	return &mpb.Empty{}, nil
}

func (svc *AccountService) SendEmailResetPasswordCode(ctx context.Context, req *mpb.ReqSendEmailResetPasswordCode) (*mpb.Empty, error) {
	// check mail exist
	ok, err := svc.dao.checkEmailExist(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, mpberr.ErrAccountNotExist
	}

	// check send email limit
	ok, err = svc.dao.checkEmailSendLimit(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, mpberr.ErrEmailSendMax
	}

	// generate code and save
	code := util.GenerateRandomCode(com.VCodeLen)
	err = svc.dao.saveEmailResetPasswordValidationCode(ctx, req.Email, code)
	if err != nil {
		return nil, err
	}

	//send email
	client, err := com.GetAPIProxyGRPCClient(ctx, svc)
	if err != nil {
		return nil, err
	}
	_, err = client.SendEmailResetPasswordValidationCode(ctx, &mpb.ReqSendEmailResetPasswordValidationCode{
		Email: req.Email,
		Code:  code,
	})
	if err != nil {
		return nil, err
	}

	return &mpb.Empty{}, nil
}

func (svc *AccountService) CheckEmailResetPasswordCode(ctx context.Context, req *mpb.ReqCheckEmailResetPasswordCode) (*mpb.ResCheckEmailResetPasswordCode, error) {
	// check v code
	ok, err := svc.dao.checkEmailResetPasswordValidationCode(ctx, req.Email, req.Code)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, mpberr.ErrEmailVerificationCode
	}
	// generate nonce and save
	nonce := util.GenerateRandomCode(com.NonceLen)
	err = svc.dao.saveEmailResetPasswordNonce(ctx, req.Email, nonce)
	if err != nil {
		return nil, err
	}

	return &mpb.ResCheckEmailResetPasswordCode{Nonce: nonce}, nil
}

func (svc *AccountService) ResetPasswordByEmail(ctx context.Context, req *mpb.ReqResetPasswordByEmail) (*mpb.Empty, error) {
	// check nonce
	ok, err := svc.dao.checkEmailResetPasswordNonce(ctx, req.Email, req.Nonce)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, mpberr.ErrRepeatedRequest
	}

	userId, err := svc.dao.getUserIdByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	// reset password
	err = svc.dao.resetPassword(ctx, userId, req.Password)
	if err != nil {
		return nil, err
	}
	return &mpb.Empty{}, nil
}

func (svc *AccountService) ResetPasswordByEmailAndVCode(ctx context.Context, req *mpb.ReqResetPasswordByEmailAndVCode) (*mpb.Empty, error) {
	// check v code
	ok, err := svc.dao.checkEmailResetPasswordValidationCode(ctx, req.Email, req.Code)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, mpberr.ErrEmailVerificationCode
	}

	userId, err := svc.dao.getUserIdByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	// reset password
	err = svc.dao.resetPassword(ctx, userId, req.Password)
	if err != nil {
		return nil, err
	}
	return &mpb.Empty{}, nil
}

func (svc *AccountService) GetAptosAccount(ctx context.Context, req *mpb.ReqUserId) (*mpb.ResGetAptosAccount, error) {
	dbAcc, err := svc.dao.getAccountInfo(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &mpb.ResGetAptosAccount{
		AptosAccAddr: dbAcc.AptosAccAddr,
		PubKey:       util.EncodeAptosPubKey(dbAcc.PublicKey),
	}, nil
}

func (svc *AccountService) getAptosResources(ctx context.Context, addr string) (string, error) {
	return "{}", nil
	// get wallet account resource
	client, err := com.GetAPIProxyGRPCClient(ctx, svc)
	if err != nil {
		return "", err
	}
	rpcRes, err := client.GetAptosResources(ctx, &mpb.ReqGetAptosResources{
		AptosAccAddr: addr})
	if err != nil {
		return "", err
	}

	return rpcRes.Resources, nil
}

func (svc *AccountService) BatchGetAccountsByWalletAddrs(ctx context.Context, req *mpb.ReqBatchGetAccountsByWalletAddrs) (*mpb.ResBatchGetAccountsByWalletAddrs, error) {
	if len(req.Addrs) == 0 {
		return &mpb.ResBatchGetAccountsByWalletAddrs{}, nil
	}
	wAccs, err := svc.dao.batchGetWalletAccs(ctx, req.Addrs)
	if err != nil {
		return nil, err
	}
	uids := make([]uint64, 0, len(req.Addrs))
	for _, v := range wAccs {
		if v == nil || v.UserId == 0 {
			continue
		}
		uids = append(uids, v.UserId)
	}
	if len(uids) == 0 {
		return &mpb.ResBatchGetAccountsByWalletAddrs{}, nil
	}
	dbAccs, err := svc.dao.batchGetAccounts(ctx, uids)
	if err != nil {
		return nil, err
	}
	res := &mpb.ResBatchGetAccountsByWalletAddrs{}
	for _, v := range dbAccs {
		res.Accounts = append(res.Accounts, svc.DBAccountInfo2AccountInfo(v))
	}
	return res, nil
}
