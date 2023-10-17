package mpberr

import (
	"errors"
	"net/http"

	"github.com/aureontu/MRWebServer/mr_services/mpb"
	gprotocol "github.com/oldjon/gutil/protocol"
)

var (
	// common error
	ErrOk              = errors.New(mpb.ErrCode_ERR_OK.String())
	ErrCmd             = errors.New(mpb.ErrCode_ERR_CMD.String())
	ErrParam           = errors.New(mpb.ErrCode_ERR_PARAM.String())
	ErrDB              = errors.New(mpb.ErrCode_ERR_DB.String())
	ErrConfig          = errors.New(mpb.ErrCode_ERR_CONFIG.String())
	ErrRepeatedRequest = errors.New(mpb.ErrCode_ERR_REPEATED_REQUEST.String())

	// tcpgateway error
	ErrMsgDecode    = errors.New(mpb.ErrCode_ERR_MSG_DECODE.String())
	ErrTokenVerify  = errors.New(mpb.ErrCode_ERR_TOKEN_VERIFY.String())
	ErrWrongGateway = errors.New(mpb.ErrCode_ERR_WRONG_GATEWAY.String())

	// account
	ErrAccountExist          = errors.New(mpb.ErrCode_ERR_ACCOUNT_EXIST.String())
	ErrAccountNotExist       = errors.New(mpb.ErrCode_ERR_ACCOUNT_NOT_EXIST.String())
	ErrPassword              = errors.New(mpb.ErrCode_ERR_PASSWORD.String())
	ErrEmailAddress          = errors.New(mpb.ErrCode_ERR_EMAIL_ADDRESS.String())
	ErrEmailSendMax          = errors.New(mpb.ErrCode_ERR_EMAIL_SEND_MAX.String())
	ErrEmailBindCode         = errors.New(mpb.ErrCode_ERR_EMAIL_BIND_CODE.String())
	ErrEmailBound            = errors.New(mpb.ErrCode_ERR_EMAIL_BOUND.String())
	ErrEmailNotExist         = errors.New(mpb.ErrCode_ERR_EMAIL_NOT_EXIST.String())
	ErrAccBoundEmail         = errors.New(mpb.ErrCode_ERR_ACC_BOUND_EMAIL.String())
	ErrAptosPublicKey        = errors.New(mpb.ErrCode_ERR_APTOS_PUBLIC_KEY.String())
	ErrAptosVerifySignature  = errors.New(mpb.ErrCode_ERR_APTOS_VERIFY_SIGNATURE.String())
	ErrNewPWSameWithOldPW    = errors.New(mpb.ErrCode_ERR_NEW_PASSWD_SAME_WITH_OLD_PASSWD.String())
	ErrEmailVerificationCode = errors.New(mpb.ErrCode_ERR_EMAIL_VERIFICATION_CODE.String())

	//nft
	ErrParseNFTId = errors.New(mpb.ErrCode_ERR_PARSE_NFT_ID.String())
	ErrNFTTokenId = errors.New(mpb.ErrCode_ERR_NFT_TOKEN_ID.String())
	ErrNFTNoOwner = errors.New(mpb.ErrCode_ERR_NFT_NO_OWNER.String())
	ErrNFTHashId  = errors.New(mpb.ErrCode_ERR_NFT_HASH_ID.String())

	//gm
	ErrAdminAccountOrPasswd = errors.New(mpb.ErrCode_ERR_ADMIN_ACCOUNT_OR_PASSWD.String())
)

var HTTPErrMap = map[string]int{
	mpb.ErrCode_ERR_OK.String():                              http.StatusOK,
	mpb.ErrCode_ERR_CMD.String():                             http.StatusBadRequest,
	mpb.ErrCode_ERR_PARAM.String():                           http.StatusBadRequest,
	mpb.ErrCode_ERR_DB.String():                              http.StatusBadRequest,
	mpb.ErrCode_ERR_MSG_DECODE.String():                      http.StatusBadRequest,
	mpb.ErrCode_ERR_TOKEN_VERIFY.String():                    http.StatusBadRequest,
	mpb.ErrCode_ERR_WRONG_GATEWAY.String():                   http.StatusBadRequest,
	mpb.ErrCode_ERR_ACCOUNT_EXIST.String():                   http.StatusBadRequest,
	mpb.ErrCode_ERR_ACCOUNT_NOT_EXIST.String():               http.StatusBadRequest,
	mpb.ErrCode_ERR_PASSWORD.String():                        http.StatusBadRequest,
	mpb.ErrCode_ERR_CONFIG.String():                          http.StatusBadRequest,
	mpb.ErrCode_ERR_EMAIL_ADDRESS.String():                   http.StatusBadRequest,
	mpb.ErrCode_ERR_EMAIL_SEND_MAX.String():                  http.StatusBadRequest,
	mpb.ErrCode_ERR_REPEATED_REQUEST.String():                http.StatusBadRequest,
	mpb.ErrCode_ERR_EMAIL_BIND_CODE.String():                 http.StatusBadRequest,
	mpb.ErrCode_ERR_EMAIL_BOUND.String():                     http.StatusBadRequest,
	mpb.ErrCode_ERR_EMAIL_NOT_EXIST.String():                 http.StatusBadRequest,
	mpb.ErrCode_ERR_ACC_BOUND_EMAIL.String():                 http.StatusBadRequest,
	mpb.ErrCode_ERR_APTOS_PUBLIC_KEY.String():                http.StatusBadRequest,
	mpb.ErrCode_ERR_APTOS_VERIFY_SIGNATURE.String():          http.StatusBadRequest,
	mpb.ErrCode_ERR_NEW_PASSWD_SAME_WITH_OLD_PASSWD.String(): http.StatusBadRequest,
	mpb.ErrCode_ERR_EMAIL_VERIFICATION_CODE.String():         http.StatusBadRequest,
	mpb.ErrCode_ERR_PARSE_NFT_ID.String():                    http.StatusBadRequest,
	mpb.ErrCode_ERR_ADMIN_ACCOUNT_OR_PASSWD.String():         http.StatusBadRequest,
	mpb.ErrCode_ERR_NFT_TOKEN_ID.String():                    http.StatusBadRequest,
	mpb.ErrCode_ERR_NFT_NO_OWNER.String():                    http.StatusBadRequest,
	mpb.ErrCode_ERR_NFT_HASH_ID.String():                     http.StatusBadRequest,
}

func ErrMsg(fc gprotocol.FrameCoder, err error) []byte {
	errMsg := &mpb.ErrorMsg{}
	errCode, ok := mpb.ErrCode_value[err.Error()]
	if ok {
		errMsg.Error = mpb.ErrCode(errCode)
	} else {
		errMsg.Error = mpb.ErrCode_ERR_UNKNOWN
	}
	data, _ := fc.EncodeMsg(uint8(mpb.MainCmd_Error), uint32(mpb.SubCmd_Error_None), errMsg)
	return data
}
