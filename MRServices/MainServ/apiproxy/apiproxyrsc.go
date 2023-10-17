package apiproxy

import (
	"fmt"

	"github.com/aureontu/MRWebServer/mr_services/mpb"
	"github.com/aureontu/MRWebServer/mr_services/mpberr"
	gcsv "github.com/oldjon/gutil/csv"
	gdm "github.com/oldjon/gutil/dirmonitor"
	"go.uber.org/zap"
)

const (
	//csvSuffix   = ".csv"
	baseCSVPath = "./resources/apiproxy/"

	emailAddrCSV   = "EmailAddr.csv"
	moralisInfoCSV = "MoralisInfo.csv"
)

type APIProxyResourceMgr struct {
	logger *zap.Logger
	dm     *gdm.DirMonitor
	mtr    *metrics

	emailAddrs  []*mpb.EmailAddrRsc
	moralisInfo *mpb.MoralisInfoRsc
}

func newAPIProxyResourceMgr(logger *zap.Logger, mtr *metrics) (*APIProxyResourceMgr, error) {
	rMgr := &APIProxyResourceMgr{
		logger: logger,
		mtr:    mtr,
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

func (rm *APIProxyResourceMgr) load() error {
	var err error
	err = rm.dm.BindAndExec(emailAddrCSV, rm.loadEmailAddr)
	if err != nil {
		return err
	}

	err = rm.dm.BindAndExec(moralisInfoCSV, rm.loadMoralisInfo)
	if err != nil {
		return err
	}

	return nil
}

func (rm *APIProxyResourceMgr) watch() error {
	return rm.dm.StartWatch()
}

func (rm *APIProxyResourceMgr) loadEmailAddr(csvPath string) error {
	datas, err := gcsv.ReadCSV2Array(csvPath)
	if err != nil {
		rm.logger.Error(fmt.Sprintf("load %s failed: %s", csvPath, err.Error()))
		return err
	}
	l := make([]*mpb.EmailAddrRsc, 0, len(datas))
	for _, data := range datas {
		node := &mpb.EmailAddrRsc{
			Addr:   data["address"],
			Passwd: data["passwd"],
			Host:   data["host"],
			Port:   data["port"],
		}
		l = append(l, node)
		rm.logger.Debug("loadEmailAddr read:", zap.Any("row", node))
	}

	rm.emailAddrs = l
	rm.logger.Debug("loadEmailAddr read finish:", zap.Any("rows", rm.emailAddrs))

	return nil
}

func (rm *APIProxyResourceMgr) getEmailAddrs() []*mpb.EmailAddrRsc {
	return rm.emailAddrs
}

func (rm *APIProxyResourceMgr) loadMoralisInfo(csvPath string) error {
	datas, err := gcsv.ReadCSV2Array(csvPath)
	if err != nil {
		rm.logger.Error(fmt.Sprintf("load %s failed: %s", csvPath, err.Error()))
		return err
	}
	if len(datas) == 0 {
		err = mpberr.ErrConfig
		rm.logger.Error(fmt.Sprintf("load %s failed: %s", csvPath, err.Error()))
		return err
	}

	rm.moralisInfo = &mpb.MoralisInfoRsc{ApiKey: datas[0]["apikey"]}
	rm.logger.Debug("loadMoralisInfo read finish:", zap.Any("rows", rm.moralisInfo))

	return nil
}

func (rm *APIProxyResourceMgr) getMoralisAPIKey() (string, error) {
	if rm.moralisInfo == nil || rm.moralisInfo.ApiKey == "" {
		return "", mpberr.ErrConfig
	}
	return rm.moralisInfo.ApiKey, nil
}
