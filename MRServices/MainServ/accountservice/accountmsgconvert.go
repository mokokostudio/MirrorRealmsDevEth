package accountservice

import "github.com/aureontu/MRWebServer/mr_services/mpb"

func (svc *AccountService) DBAccountInfo2AccountInfo(in *mpb.DBAccountInfo) *mpb.AccountInfo {
	if in == nil {
		return nil
	}
	return &mpb.AccountInfo{
		Account:         in.Account,
		UserId:          in.UserId,
		Email:           in.Email,
		Nickname:        in.Nickname,
		Icon:            in.Icon,
		AptosWalletAddr: in.AptosAccAddr,
	}
}
