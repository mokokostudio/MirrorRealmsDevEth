package nftservice

import (
	"fmt"
	"os"
	"strings"

	"github.com/aureontu/MRWebServer/mr_services/mpb"
	"github.com/aureontu/MRWebServer/mr_services/mpberr"
	"github.com/aureontu/MRWebServer/mr_services/util"
	gconv "github.com/oldjon/gutil/conv"
	gcsv "github.com/oldjon/gutil/csv"
	gdm "github.com/oldjon/gutil/dirmonitor"
	"go.uber.org/zap"
)

const (
	//csvSuffix   = ".csv"
	baseCSVPath = "./resources/nft/"

	nftMetaDataPath = "metadata/"

	nftCollectionCSV      = "NFTCollection.csv"
	nftMetaDataVersionCSV = "NFTMetaDataVersion.csv"
	nftConfigCSV          = "NFTConfig.csv"
)

type NFTResourceMgr struct {
	logger *zap.Logger
	dm     *gdm.DirMonitor
	mtr    *util.ServiceMetrics

	emailAddrs       []*mpb.EmailAddrRsc
	nftInfoMap       map[mpb.ENFT_NFTType]*mpb.NFTCollectionRsc
	nftInfos         []*mpb.NFTCollectionRsc
	nftCollectionMap map[string]*mpb.NFTCollectionRsc
	nftMetaDatas     map[uint64]string
	nftConfig        *mpb.NFTConfigRsc
}

func newNFTResourceMgr(logger *zap.Logger, sm *util.ServiceMetrics) (*NFTResourceMgr, error) {
	rMgr := &NFTResourceMgr{
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

func (rm *NFTResourceMgr) load() error {
	var err error

	err = rm.dm.BindAndExec(nftConfigCSV, rm.loadNFTConfig)
	if err != nil {
		return err
	}

	err = rm.dm.BindAndExec(nftCollectionCSV, rm.loadNFTCollection)
	if err != nil {
		return err
	}

	err = rm.dm.BindAndExec(nftMetaDataVersionCSV, rm.loadNFMetaDatas)
	if err != nil {
		return err
	}

	return nil
}

func (rm *NFTResourceMgr) watch() error {
	return rm.dm.StartWatch()
}

func (rm *NFTResourceMgr) loadNFTCollection(csvPath string) error {
	datas, err := gcsv.ReadCSV2Array(csvPath)
	if err != nil {
		rm.logger.Error(fmt.Sprintf("load %s failed: %s", csvPath, err.Error()))
		return err
	}
	m := make(map[mpb.ENFT_NFTType]*mpb.NFTCollectionRsc)
	l := make([]*mpb.NFTCollectionRsc, 0, len(datas))
	cm := make(map[string]*mpb.NFTCollectionRsc)
	for _, data := range datas {
		node := &mpb.NFTCollectionRsc{
			NftType: mpb.ENFT_NFTType(gcsv.Str2Int32(data["nfttype"])),
		}
		if len(data["collectionhash"]) > 0 {
			node.CollectionHash = strings.Split(data["collectionhash"], ";")
			for _, vv := range node.CollectionHash {
				cm[vv] = node
			}
		}
		m[node.NftType] = node
		l = append(l, node)
		rm.logger.Debug("loadNFTCollection read:", zap.Any("row", node))
	}

	rm.nftInfoMap = m
	rm.nftInfos = l
	rm.nftCollectionMap = cm
	rm.logger.Debug("loadNFTCollection read finish:", zap.Any("rows", rm.nftInfos))
	return nil
}

func (rm *NFTResourceMgr) getAllNFTCollectionRsc() []*mpb.NFTCollectionRsc {
	return rm.nftInfos
}

func (rm *NFTResourceMgr) getNFTCollectionRscByNFTType(nftType mpb.ENFT_NFTType) (*mpb.NFTCollectionRsc, error) {
	node, ok := rm.nftInfoMap[nftType]
	if !ok {
		rm.logger.Error("getNFTCollectionRscByNFTType failed", zap.Any("nft_type", nftType))
		return nil, mpberr.ErrConfig
	}
	return node, nil
}

func (rm *NFTResourceMgr) getNFTTypeByCollection(collectionHash string) mpb.ENFT_NFTType {
	return rm.nftCollectionMap[collectionHash].GetNftType()
}

func (rm *NFTResourceMgr) loadNFMetaDatas(_ string) error {
	dirPath := baseCSVPath + nftMetaDataPath
	fileNames, err := listDir(dirPath)
	if err != nil {
		return err
	}

	metaDats := make(map[uint64]string)
	for _, v := range fileNames {
		if !strings.HasSuffix(v, ".json") {
			continue
		}
		strs := strings.Split(v, ".")
		index := gconv.StringToUint64(strs[0])
		if index == 0 {
			rm.logger.Error("loadNFMetaDatas failed", zap.Any("nft_type", fileNames))
			return mpberr.ErrConfig
		}
		datas, err := os.ReadFile(dirPath + v)
		if err != nil {
			rm.logger.Error("loadNFMetaDatas failed", zap.Any("nft_type", fileNames))
			return err
		}
		metaDats[index] = string(datas)
	}

	rm.nftMetaDatas = metaDats
	rm.logger.Debug("loadNFMetaDatas read finish:", zap.Any("rows", metaDats))
	return nil
}

func (rm *NFTResourceMgr) getNftMetaData(tokenId uint64) string {
	return rm.nftMetaDatas[tokenId]
}

func listDir(dirPath string) ([]string, error) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(files))
	for _, v := range files {
		names = append(names, v.Name())
	}
	return names, err
}

func (rm *NFTResourceMgr) loadNFTConfig(csvPath string) error {
	datas, err := gcsv.ReadCSV2Array(csvPath)
	if err != nil || len(datas) < 1 {
		rm.logger.Error(fmt.Sprintf("load %s failed: %s", csvPath, err.Error()))
		return err
	}

	c := &mpb.NFTConfigRsc{
		NftGraphiqlQueryLimit:  gcsv.Str2Uint64(datas[0]["nftgraphiqlquerylimit"]),
		NftGraphiqlPageNum:     gcsv.Str2Uint64(datas[0]["nftgraphiqlpagenum"]),
		NftTransferEvent:       datas[0]["nfttransferevent"],
		NftBurnEvent:           datas[0]["nftburnevent"],
		NftWithdrawEvent2:      datas[0]["nftwithdrawevent2"],
		NftDepositEvent2:       datas[0]["nftdepositevent2"],
		NftBurnEvent2:          datas[0]["nftburnevent2"],
		NftMintEvent2:          datas[0]["nftmintevent2"],
		NftTransferOfferEvent2: datas[0]["nfttransferofferevent2"],
		NftTransferClaimEvent2: datas[0]["nfttransferclaimevent2"],
	}

	rm.nftConfig = c
	rm.logger.Debug("loadNFConfig read finish:", zap.Any("rows", c))
	return nil
}

func (rm *NFTResourceMgr) getGraphiQLQueryLimitSecs() uint64 {
	return rm.nftConfig.GetNftGraphiqlQueryLimit()
}

func (rm *NFTResourceMgr) getGraphiQLQueryPageNum() uint64 {
	return rm.nftConfig.GetNftGraphiqlPageNum()
}

func (rm *NFTResourceMgr) getNFTTransferEvent() string {
	return rm.nftConfig.GetNftTransferEvent()
}

func (rm *NFTResourceMgr) getNFTBurnEvent() string {
	return rm.nftConfig.GetNftBurnEvent()
}

func (rm *NFTResourceMgr) getNFTWithdrawEvent2() string {
	return rm.nftConfig.GetNftWithdrawEvent2()
}

func (rm *NFTResourceMgr) getNFTDepositEvent2() string {
	return rm.nftConfig.GetNftDepositEvent2()
}

func (rm *NFTResourceMgr) getNFTBurnEvent2() string {
	return rm.nftConfig.GetNftBurnEvent2()
}

func (rm *NFTResourceMgr) getNFTMintEvent2() string {
	return rm.nftConfig.GetNftMintEvent2()
}

func (rm *NFTResourceMgr) getNFTTransferOfferEvent2() string {
	return rm.nftConfig.GetNftTransferOfferEvent2()
}

func (rm *NFTResourceMgr) getNFTTransferClaimEvent2() string {
	return rm.nftConfig.GetNftTransferClaimEvent2()
}

func (rm *NFTResourceMgr) isValidCollection(collection string) bool {
	_, ok := rm.nftCollectionMap[collection]
	return ok
}
