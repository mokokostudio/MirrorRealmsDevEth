package common

import (
	"io/ioutil"
	"strings"

	"github.com/oldjon/gx/service"
)

const (
	// CommonRegion common region for cluster
	CommonRegion = "SG"

	// CommonIPRegion common ipregion for cluster
	CommonIPRegion = "CN"

	// CommonLanguage common language on for cluster
	CommonLanguage = "en"

	// CenterServerCN center server cn mark
	CenterServerCN = "CS_CN"
	// CenterServerSG center server sg mark
	CenterServerSG = "CS_SG"
)

// Regions all region
var Regions = []string{"CN", "SG"}

// RegionServerMapping : region -> center server.
var RegionServerMapping = map[string]string{
	"CN": CenterServerCN,
	"SG": CenterServerSG,
}

// RegionEndMatchTopicMapping: region -> endMatch topic
//var RegionEndMatchTopicMapping = map[string]proto.EMessageQueue_Topic{
//	"VN":     proto.EMessageQueue_Topic_VN_EndMatch,
//	"TH":     proto.EMessageQueue_Topic_TH_EndMatch,
//	"ID":     proto.EMessageQueue_Topic_ID_EndMatch,
//	"TW":     proto.EMessageQueue_Topic_TW_EndMatch,
//	"SG":     proto.EMessageQueue_Topic_SG_EndMatch,
//	"RU":     proto.EMessageQueue_Topic_RU_EndMatch,
//	"EUROPE": proto.EMessageQueue_Topic_EUROPE_EndMatch,
//	"IND":    proto.EMessageQueue_Topic_IND_EndMatch,
//	"ME":     proto.EMessageQueue_Topic_ME_ENdMatch,
//	"BR":     proto.EMessageQueue_Topic_BR_EndMatch,
//	"SAC":    proto.EMessageQueue_Topic_SAC_EndMatch,
//	"US":     proto.EMessageQueue_Topic_US_EndMatch,
//	"NA":     proto.EMessageQueue_Topic_NA_EndMatch,
//	"PK":     proto.EMessageQueue_Topic_PK_EndMatch,
//	"BD":     proto.EMessageQueue_Topic_BD_EndMatch,
//}

// GetRegionFromFileName would get region from file name with format xxxxxx_{region}.csv
// If filename is with other format return CommonRegion
func GetRegionFromFileName(filename string) string {
	//check .csv
	if len(filename) <= 4 {
		return ""
	}
	//remove .csv
	filename = filename[0 : len(filename)-4]
	strArr := strings.Split(filename, "_")
	if len(strArr) < 2 {
		return ""
	}
	return strings.ToUpper(strArr[len(strArr)-1])
}

// IsValidRegion would check if region is in Regions
func IsValidRegion(region string) bool {
	if region == "" {
		return false
	}
	for _, validRegion := range Regions {
		if validRegion == region {
			return true
		}
	}
	return false
}

// GetCenterServerByRegion : get center server node name by region return CS_SG/CS_US/CS_IND
func GetCenterServerByRegion(region string) string {
	if cs, ok := RegionServerMapping[region]; ok {
		return cs
	}
	return ""
}

// IsSameCenterServer will check if the two regions in same center server
func IsSameCenterServer(regionA string, regionB string) bool {
	centerServerA := GetCenterServerByRegion(regionA)
	centerServerB := GetCenterServerByRegion(regionB)
	return centerServerA != "" && centerServerB != "" && centerServerA == centerServerB
}

func GetCenterRegionsFromTags(host service.Host) []string {
	if v, ok := host.Tags()["center_server"]; ok {
		if v != "" {
			centerRegions := strings.Split(strings.Trim(v, " "), ";")
			return centerRegions
		}
	}
	return []string{CenterServerCN, CenterServerSG}
}

// GetAllValidRegionCsvFiles get all csv files with valid region
func GetAllValidRegionCsvFiles(dirPath string) ([]string, error) {
	fileList := make([]string, 0)
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	for _, fi := range files {
		if !IsValidRegion(GetRegionFromFileName(fi.Name())) {
			continue
		}
		if !fi.IsDir() && strings.HasSuffix(fi.Name(), ".csv") {
			fileList = append(fileList, fi.Name())
		}

	}
	return fileList, nil
}

// AllLanguageSlice slice language all region support
var AllLanguageSlice = []string{
	"zh-chs", "en", "zh-cht", "id", "ja", "ko", "pt-br", "ru", "es", "th", "tr", "vn", "ar", "fr", "de", "hi", "my", "ur", "ro", "bn",
}
