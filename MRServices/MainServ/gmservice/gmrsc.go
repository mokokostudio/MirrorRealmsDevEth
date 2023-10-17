package gmservice

import (
	"fmt"

	"github.com/aureontu/MRWebServer/mr_services/mpb"
	"github.com/aureontu/MRWebServer/mr_services/mpberr"
	"github.com/aureontu/MRWebServer/mr_services/util"
	gcsv "github.com/oldjon/gutil/csv"
	gdm "github.com/oldjon/gutil/dirmonitor"
	"go.uber.org/zap"
)

const (
	//csvSuffix   = ".csv"
	baseCSVPath = "./resources/gm/"

	adminsCSV = "Admins.csv"
)

type GMResourceMgr struct {
	logger *zap.Logger
	dm     *gdm.DirMonitor
	mtr    *util.ServiceMetrics

	adminMap map[string]*mpb.AdminRsc
}

func newGMResourceMgr(logger *zap.Logger, sm *util.ServiceMetrics) (*GMResourceMgr, error) {
	rMgr := &GMResourceMgr{
		logger: logger,
		mtr:    sm,
	}

	var err error
	rMgr.dm, err = gdm.NewDirMonitor(baseCSVPath)
	if err != nil {
		return nil, err
	}

	err = rMgr.load()
	if err != nil {
		return nil, err
	}

	err = rMgr.watch()
	if err != nil {
		return nil, err
	}

	return rMgr, nil
}

func (rm *GMResourceMgr) load() error {
	var err error

	err = rm.dm.BindAndExec(adminsCSV, rm.loadAdmins)
	if err != nil {
		return err
	}

	return nil
}

func (rm *GMResourceMgr) watch() error {
	return rm.dm.StartWatch()
}

func (rm *GMResourceMgr) loadAdmins(csvPath string) error {
	datas, err := gcsv.ReadCSV2Array(csvPath)
	if err != nil {
		rm.logger.Error(fmt.Sprintf("load %s failed: %s", csvPath, err.Error()))
		return err
	}
	m := make(map[string]*mpb.AdminRsc)
	for _, data := range datas {
		node := &mpb.AdminRsc{
			Account:  data["account"],
			Password: data["password"],
		}

		m[node.Account] = node
		rm.logger.Debug("loadAdmins read:", zap.Any("row", node))
	}

	rm.adminMap = m
	rm.logger.Debug("loadAdmins read finish:", zap.Any("rows", rm.adminMap))
	return nil
}

func (rm *GMResourceMgr) getAdminRSC(account string) (*mpb.AdminRsc, error) {
	accRsc, ok := rm.adminMap[account]
	if !ok {
		return nil, mpberr.ErrAdminAccountOrPasswd
	}
	return accRsc, nil
}
