package models

type NodeInfoList struct {
	NodeInfos []*NodeInfo `json:"nodeInfos"`
}

type NodeInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Address  string `json:"address"`
	Port     int    `json:"port"`
}

type DomainInfoList struct {
	DomainInfos []*DomainInfo `json:"domainInfos"`
}

type DomainInfo struct {
	Domain string `json:"domain"`
	IP     string `json:"ip"`
}
