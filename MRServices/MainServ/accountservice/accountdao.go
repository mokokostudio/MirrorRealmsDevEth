package accountservice

import (
	"context"
	"strconv"
	"time"

	com "github.com/aureontu/MRWebServer/mr_services/common"
	"github.com/aureontu/MRWebServer/mr_services/mpb"
	"github.com/aureontu/MRWebServer/mr_services/mpberr"
	"github.com/oldjon/gutil/gdb"
	grmux "github.com/oldjon/gutil/redismutex"
	"go.uber.org/zap"
)

type accountDAO struct {
	logger *zap.Logger
	rMux   *grmux.RedisMutex
	accDB  *gdb.DB
	tmpDB  *gdb.DB
}

func newAccountDAO(logger *zap.Logger, rMux *grmux.RedisMutex, accRedis gdb.RedisClient, tmpRedis gdb.RedisClient) *accountDAO {
	return &accountDAO{
		logger: logger,
		rMux:   rMux,
		accDB:  gdb.NewDB(accRedis),
		tmpDB:  gdb.NewDB(tmpRedis),
	}
}

func (dao *accountDAO) getAccountInfo(ctx context.Context, userId uint64) (*mpb.DBAccountInfo, error) {
	key := com.AccountKey(userId)
	dbAccount := &mpb.DBAccountInfo{}
	err := dao.accDB.GetObject(ctx, key, dbAccount)
	if dao.accDB.IsErrNil(err) {
		return nil, mpberr.ErrAccountNotExist
	} else if err != nil {
		dao.logger.Error("getAccountInfo GetObject failed", zap.String("key", key), zap.Error(err))
		return nil, mpberr.ErrDB
	}
	return dbAccount, nil
}

func (dao *accountDAO) getUserIdByAccount(ctx context.Context, account string) (uint64, error) {
	key := com.AccountUIDKey(account)
	userId, err := gdb.ToUint64(dao.accDB.Get(ctx, key))
	if dao.accDB.IsErrNil(err) {
		return 0, mpberr.ErrAccountNotExist
	} else if err != nil {
		dao.logger.Error("getUserIdByAccount Get failed", zap.String("key", key), zap.Error(err))
		return 0, mpberr.ErrDB
	}
	return userId, nil
}

func (dao *accountDAO) getAccountInfoByAccount(ctx context.Context, account string) (*mpb.DBAccountInfo, error) {
	key := com.AccountUIDKey(account)
	userId, err := gdb.ToUint64(dao.accDB.Get(ctx, key))
	if dao.accDB.IsErrNil(err) {
		return nil, mpberr.ErrAccountNotExist
	} else if err != nil {
		dao.logger.Error("getUserIdByAccount Get failed", zap.String("key", key), zap.Error(err))
		return nil, mpberr.ErrDB
	}
	return dao.getAccountInfo(ctx, userId)
}

func (dao *accountDAO) saveNonce(ctx context.Context, nonce string) error {
	return dao.tmpDB.SetEX(ctx, com.NonceKey(nonce), 1, com.Dur10Mins)
}

func (dao *accountDAO) checkNonce(ctx context.Context, nonce string) (bool, error) {
	key := com.NonceKey(nonce)
	n, err := dao.tmpDB.Del(ctx, key)
	if err != nil {
		dao.logger.Error("checkNonce Del failed", zap.String("key", key), zap.Error(err))
		return false, mpberr.ErrDB
	}
	return n == 1, nil
}

func (dao *accountDAO) getAccountByWallet(ctx context.Context, aptosAccAddr string, pubKey []byte) (*mpb.DBAccountInfo, error) {
	wa := &mpb.DBWalletAcc{}
	key := com.WalletAccKey(aptosAccAddr)
	err := dao.accDB.GetObject(ctx, key, wa)
	if err != nil && !dao.accDB.IsErrNil(err) {
		dao.logger.Error("getAccountByWallet GetObject failed", zap.String("key", key), zap.Error(err))
		return nil, mpberr.ErrDB
	}
	if wa == nil || wa.UserId == 0 {
		// init wallet acc
		userId, err := gdb.ToUint64(dao.accDB.Incr(ctx, com.UserIdIndexKey()))
		if err != nil {
			dao.logger.Error("getAccountByWallet GetObject failed", zap.String("key", key), zap.Error(err))
			return nil, mpberr.ErrDB
		}
		wa = &mpb.DBWalletAcc{
			UserId: userId,
		}
		err = dao.accDB.SetObject(ctx, key, wa)
		if err != nil {
			dao.logger.Error("getAccountByWallet GetObject failed", zap.String("key", key), zap.Error(err))
			return nil, mpberr.ErrDB
		}
		nickname := "MR" + strconv.Itoa(int(userId))
		dbAcc := &mpb.DBAccountInfo{
			UserId:       userId,
			AptosAccAddr: aptosAccAddr,
			PublicKey:    pubKey,
			Nickname:     nickname,
			Icon:         "0",
			RegisterTime: time.Now().Unix(),
		}
		err = dao.accDB.SetObject(ctx, com.AccountKey(userId), dbAcc)
		if err != nil {
			dao.logger.Error("getAccountByWallet SetObject failed",
				zap.String("key", com.AccountKey(userId)), zap.Error(err))
			return nil, mpberr.ErrDB
		}
		return dbAcc, nil
	}

	return dao.getAccountInfo(ctx, wa.UserId)
}

func (dao *accountDAO) batchGetWalletAccs(ctx context.Context, addrs []string) ([]*mpb.DBWalletAcc, error) {
	keys := make([]string, 0, len(addrs))
	was := make([]*mpb.DBWalletAcc, 0, len(addrs))
	for _, addr := range addrs {
		keys = append(keys, com.WalletAccKey(addr))
		was = append(was, &mpb.DBWalletAcc{})
	}
	err := dao.accDB.GetObjects(ctx, keys, was)
	if err != nil {
		dao.logger.Error("batchGetWalletAccs GetObjects failed",
			zap.Any("addrs", addrs), zap.Error(err))
		return nil, err
	}
	return was, err
}

func (dao *accountDAO) batchGetAccounts(ctx context.Context, userIds []uint64) ([]*mpb.DBAccountInfo, error) {
	keys := make([]string, 0, len(userIds))
	accs := make([]*mpb.DBAccountInfo, 0, len(userIds))
	for _, uid := range userIds {
		keys = append(keys, com.AccountKey(uid))
		accs = append(accs, &mpb.DBAccountInfo{})
	}
	err := dao.accDB.GetObjects(ctx, keys, accs)
	if err != nil {
		dao.logger.Error("batchGetWalletAccs GetObjects failed",
			zap.Any("user_ids", userIds), zap.Error(err))
		return nil, err
	}
	return accs, err
}

func (dao *accountDAO) saveToken(ctx context.Context, userId uint64, token, account, deviceId, device string) error {
	key := com.TokenKey(token)
	err := dao.tmpDB.SetObjectEX(ctx, key, &mpb.DBTokenInfo{
		Account:  account,
		Device:   device,
		DeviceId: deviceId,
	}, com.TokenExpireDuration)
	if err != nil {
		dao.logger.Error("save token failed", zap.String("key", key), zap.Error(err))
		return mpberr.ErrDB
	}
	err = dao.tmpDB.Set(ctx, com.UIDTokenKey(userId), token)
	if err != nil {
		dao.logger.Error("save token failed", zap.String("key", key), zap.Error(err))
		return mpberr.ErrDB
	}
	return nil
}

func (dao *accountDAO) checkEmailSendLimit(ctx context.Context, emailAddr string) (bool, error) {
	now := time.Now()
	date := now.Format("20060102")
	key := com.EmailSendDailyLimitKey(emailAddr, date)
	cnt, err := dao.tmpDB.IncrBy(ctx, key, 1)
	if err != nil {
		dao.logger.Error("checkEmailSendLimit IncrBy failed", zap.String("key", key), zap.Error(err))
		return false, mpberr.ErrDB
	}
	if cnt > com.EmailSendDailyLimit {
		return false, nil
	}
	_, err = dao.tmpDB.Expire(ctx, key, com.Dur1Day)
	if err != nil {
		dao.logger.Error("checkEmailSendLimit Expire failed", zap.String("key", key), zap.Error(err))
		return false, err
	}
	return true, nil
}

func (dao *accountDAO) saveEmailBindCode(ctx context.Context, emailAddr, code string) error {
	key := com.EmailBindCodeKey(emailAddr)
	err := dao.tmpDB.SetEX(ctx, key, code, com.Dur10Mins)
	if err != nil {
		dao.logger.Error("saveEmailBindCode SetEX failed", zap.String("key", key), zap.Error(err))
		return mpberr.ErrDB
	}
	return nil
}

func (dao *accountDAO) checkEmailBindCode(ctx context.Context, emailAddr, code string) (bool, error) {
	key := com.EmailBindCodeKey(emailAddr)
	var ok bool
	err := dao.rMux.Safely(ctx, key, func() error {
		dbCode, err := dao.tmpDB.Get(ctx, key)
		if err != nil && !dao.tmpDB.IsErrNil(err) {
			dao.logger.Error("checkEmailBindCode Get failed", zap.String("key", key), zap.Error(err))
			return mpberr.ErrDB
		}
		if dbCode != code {
			return nil
		}

		ok = true
		_, err = dao.tmpDB.Del(ctx, key)
		if err != nil {
			dao.logger.Error("checkEmailBindCode Del failed", zap.String("key", key), zap.Error(err))
			return mpberr.ErrDB
		}
		return nil
	})
	if err != nil {
		dao.logger.Error("checkEmailBindCode Safely failed", zap.String("key", key), zap.Error(err))
		return false, err
	}
	return ok, nil
}

func (dao *accountDAO) bindEmail(ctx context.Context, userId uint64, acc, email string) (*mpb.DBAccountInfo, error) {
	dbAcc := &mpb.DBAccountInfo{}
	key1 := com.AccountKey(userId)
	err := dao.rMux.Safely(ctx, key1, func() error {
		err := dao.accDB.GetObject(ctx, key1, dbAcc)
		if err != nil {
			dao.logger.Error("bindEmail GetObject failed", zap.String("key", key1), zap.Error(err))
			return mpberr.ErrDB
		}
		if dbAcc.Email != "" {
			return mpberr.ErrAccBoundEmail
		}
		dbAcc.Email = email
		dbAcc.Account = acc
		err = dao.accDB.SetObject(ctx, key1, dbAcc)
		if err != nil {
			dao.logger.Error("bindEmail SetObject failed", zap.String("key", key1), zap.Error(err))
			return mpberr.ErrDB
		}
		return nil
	})
	if err != nil {
		dao.logger.Error("bindEmail Safely failed", zap.String("key", key1), zap.Error(err))
		return nil, err
	}

	key2 := com.EmailAccKey(email)
	err = dao.accDB.Set(ctx, key2, userId)
	if err != nil {
		dao.logger.Error("bindEmail Set failed", zap.String("key", key2), zap.Error(err))
		return nil, mpberr.ErrDB
	}
	key3 := com.AccountUIDKey(dbAcc.Account)
	err = dao.accDB.Set(ctx, key3, userId)
	if err != nil {
		dao.logger.Error("bindEmail Set failed", zap.String("key", key3), zap.Error(err))
		return nil, mpberr.ErrDB
	}
	return dbAcc, nil
}

func (dao *accountDAO) changePassword(ctx context.Context, userId uint64, oldPassword, newPassword string) error {
	key := com.AccountKey(userId)
	err := dao.rMux.Safely(ctx, key, func() error {
		var dbAcc = &mpb.DBAccountInfo{}
		err := dao.accDB.GetObject(ctx, key, dbAcc)
		if dao.accDB.IsErrNil(err) {
			return mpberr.ErrAccountNotExist
		} else if err != nil {
			dao.logger.Error("changePassword GetObject failed", zap.String("key", key), zap.Error(err))
			return mpberr.ErrDB
		}
		if dbAcc.Password != "" && dbAcc.Password != oldPassword {
			return mpberr.ErrPassword
		}
		dbAcc.Password = newPassword
		err = dao.accDB.SetObject(ctx, key, dbAcc)
		if err != nil {
			dao.logger.Error("changePassword SetObject failed", zap.String("key", key), zap.Error(err))
			return mpberr.ErrDB
		}
		return nil
	})
	if err != nil {
		dao.logger.Error("changePassword Safely failed", zap.String("key", key), zap.Error(err))
		return err
	}
	return nil
}

func (dao *accountDAO) checkEmailExist(ctx context.Context, emailAddr string) (bool, error) {
	key := com.EmailAccKey(emailAddr)
	ok, err := dao.accDB.Exists(ctx, key)
	if err != nil {
		dao.logger.Error("checkEmailExist Exists failed", zap.String("key", key), zap.Error(err))
		return false, mpberr.ErrDB
	}
	return ok, nil
}

func (dao *accountDAO) saveEmailResetPasswordValidationCode(ctx context.Context, emailAddr, code string) error {
	key := com.EmailResetPasswordValidationCodeKey(emailAddr)
	err := dao.tmpDB.SetEX(ctx, key, code, com.Dur10Mins)
	if err != nil {
		dao.logger.Error("saveEmailBindCode SetEX failed", zap.String("key", key), zap.Error(err))
		return mpberr.ErrDB
	}
	return nil
}

func (dao *accountDAO) checkEmailResetPasswordValidationCode(ctx context.Context, emailAddr, code string) (bool, error) {
	key := com.EmailResetPasswordValidationCodeKey(emailAddr)
	var ok bool
	err := dao.rMux.Safely(ctx, key, func() error {
		dbCode, err := dao.tmpDB.Get(ctx, key)
		if err != nil && !dao.tmpDB.IsErrNil(err) {
			dao.logger.Error("checkEmailResetPasswordValidationCode Get failed",
				zap.String("key", key), zap.Error(err))
			return mpberr.ErrDB
		}
		if dbCode != code {
			return nil
		}

		ok = true
		_, err = dao.tmpDB.Del(ctx, key)
		if err != nil {
			dao.logger.Error("checkEmailResetPasswordValidationCode Del failed",
				zap.String("key", key), zap.Error(err))
			return mpberr.ErrDB
		}
		return nil
	})
	if err != nil {
		dao.logger.Error("checkEmailResetPasswordValidationCode Safely failed",
			zap.String("key", key), zap.Error(err))
		return false, err
	}
	return ok, nil
}

func (dao *accountDAO) saveEmailResetPasswordNonce(ctx context.Context, emailAddr, nonce string) error {
	key := com.EmailResetPasswordNonceKey(emailAddr)
	err := dao.tmpDB.SetEX(ctx, key, nonce, com.Dur10Mins)
	if err != nil {
		dao.logger.Error("saveEmailBindCode SetEX failed", zap.String("key", key), zap.Error(err))
		return mpberr.ErrDB
	}
	return nil
}

func (dao *accountDAO) checkEmailResetPasswordNonce(ctx context.Context, emailAddr, nonce string) (bool, error) {
	key := com.EmailResetPasswordNonceKey(emailAddr)
	var ok bool
	err := dao.rMux.Safely(ctx, key, func() error {
		dbNonce, err := dao.tmpDB.Get(ctx, key)
		if err != nil && !dao.tmpDB.IsErrNil(err) {
			dao.logger.Error("checkEmailResetPasswordNonce Get failed",
				zap.String("key", key), zap.Error(err))
			return mpberr.ErrDB
		}
		if dbNonce != nonce {
			return nil
		}

		ok = true
		_, err = dao.tmpDB.Del(ctx, key)
		if err != nil {
			dao.logger.Error("checkEmailResetPasswordNonce Del failed",
				zap.String("key", key), zap.Error(err))
			return mpberr.ErrDB
		}
		return nil
	})
	if err != nil {
		dao.logger.Error("checkEmailResetPasswordNonce Safely failed",
			zap.String("key", key), zap.Error(err))
		return false, err
	}
	return ok, nil
}

func (dao *accountDAO) getUserIdByEmail(ctx context.Context, email string) (uint64, error) {
	key := com.EmailAccKey(email)
	userId, err := gdb.ToUint64(dao.accDB.Get(ctx, key))

	if dao.accDB.IsErrNil(err) {
		return 0, mpberr.ErrEmailNotExist
	} else if err != nil {
		dao.logger.Error("getUserIdByEmail SetObject failed", zap.String("key", key), zap.Error(err))
		return 0, mpberr.ErrDB
	}
	return userId, nil
}

func (dao *accountDAO) resetPassword(ctx context.Context, userId uint64, password string) error {
	key := com.AccountKey(userId)
	err := dao.rMux.Safely(ctx, key, func() error {
		var dbAcc = &mpb.DBAccountInfo{}
		err := dao.accDB.GetObject(ctx, key, dbAcc)
		if dao.accDB.IsErrNil(err) {
			return mpberr.ErrAccountNotExist
		} else if err != nil {
			dao.logger.Error("resetPassword GetObject failed", zap.String("key", key), zap.Error(err))
			return mpberr.ErrDB
		}

		dbAcc.Password = password
		err = dao.accDB.SetObject(ctx, key, dbAcc)
		if err != nil {
			dao.logger.Error("resetPassword SetObject failed", zap.String("key", key), zap.Error(err))
			return mpberr.ErrDB
		}
		return nil
	})
	if err != nil {
		dao.logger.Error("resetPassword Safely failed", zap.String("key", key), zap.Error(err))
		return err
	}
	return nil
}
