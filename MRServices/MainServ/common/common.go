package common

import (
	"context"
	"fmt"
	"strconv"

	gcrypto "github.com/oldjon/gutil/crypto"
)

func IsDefaultAccount(account string) bool {
	if account[0:1] != "m" { // first char is 'm'
		return false
	}
	_, err := strconv.Atoi(account[1:]) // other chars can convert to num
	return err == nil
}

func GenerateDefaultAccountByUserId(userId uint64) string {
	return fmt.Sprintf("M%d", userId)
}

func GenerateDefaultPassword() string {
	return gcrypto.MD5SumStr(DefaultPassword + PasswordSalt)
}

func TimeoutCtx() (context.Context, func()) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxTimeout) // nolint:govet
	return ctx, cancel
}
