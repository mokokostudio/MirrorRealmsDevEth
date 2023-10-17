package common

import "github.com/aureontu/MRWebServer/mr_services/mpb"

func MainCmdStr(cmd uint8) string {
	return mpb.MainCmd_Cmd(cmd).String()
}

func SubCmdStr(cmd uint8, subCmd uint32) string {
	switch mpb.MainCmd_Cmd(cmd) {
	case mpb.MainCmd_Error:
		return mpb.SubCmd_Error_Cmd(subCmd).String()
	case mpb.MainCmd_TCPGateway:
		return mpb.SubCmd_TCPGateway_Cmd(subCmd).String()
	case mpb.MainCmd_Team:
		return mpb.SubCmd_Team_Cmd(subCmd).String()
	case mpb.MainCmd_Chat:
		return mpb.SubCmd_Chat_Cmd(subCmd).String()
	}
	return ""
}
