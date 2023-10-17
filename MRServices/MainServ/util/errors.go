package util

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

var (
	errAESEncryptionKeyPairsIsAnEmptyList = errors.New("AESEncryptionKeyPairs is an empty list")
	errHeaderNotInWellFormat              = errors.New("header not in well format")
	errHeaderContainNotSupportedKeyType   = errors.New("header contain not supported keyType")
	errHeaderContainNotSupportedKeyIndex  = errors.New("header contain not supported keyIndex")

	errRandomPoolEmpty       = errors.New("random pool empty")
	errRandomPoolTooSmall    = errors.New("random pool to small")
	errRandomTotalWeightZero = errors.New("random total weight zero")
	errRandomWeightFuncNil   = errors.New("random weight func nil")
)

// AppError logic error for current app
type AppError struct {
	// The error display in logs for developer to view
	Err error
	// The message that send to bot
	Message string
	// HTTP error code
	Code int
}

func (ae *AppError) Error() string {
	return ae.Message
}

type TransactionError struct {
	Err     error
	Message string
}

func (te *TransactionError) Error() string {
	return "TransactionError:" + te.Message
}

// IsGrpcTimeout check connect timeout or grpc timeout
func IsGrpcTimeout(err error) bool {
	if errors.Is(err, context.DeadlineExceeded) || grpc.Code(err) == codes.DeadlineExceeded {
		return true
	}
	return false
}
