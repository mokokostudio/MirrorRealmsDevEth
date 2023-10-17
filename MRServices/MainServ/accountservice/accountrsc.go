package accountservice

import (
	"github.com/aureontu/MRWebServer/mr_services/util"
	"go.uber.org/zap"
)

type AccountResourceMgr struct {
}

func newAccountResourceMgr(logger *zap.Logger, sm *util.ServiceMetrics) (*AccountResourceMgr, error) {
	return &AccountResourceMgr{}, nil
}
