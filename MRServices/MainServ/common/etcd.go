package common

type EtcdRes struct {
	Op       int    `json:"Op"`
	Addr     string `json:"Addr"`
	Metadata struct {
		Node             string `json:"node"`
		Region           string `json:"region"`
		PublicListenAddr string `json:"public_listen_addr"`
	} `json:"Metadata"`
}
