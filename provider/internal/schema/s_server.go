package schema

import "github.com/cloudslit/cloudslit/provider/pkg/util/json"

type ServerType string

const ServerTypeProvider ServerType = "provider"

type ServerInfo struct {
	PeerId   string     `json:"peer_id"`
	Addr     string     `json:"addr"`
	Port     int        `json:"port"`
	MetaData MetaData   `json:"meta_data"`
	Price    int        `json:"price"`
	Type     ServerType `json:"type"`
}

type MetaData struct {
	Ip   string `json:"ip"`
	Loc  string `json:"loc"`
	Colo string `json:"colo"`
}

func (a *ServerInfo) String() string {
	return json.MarshalToString(a)
}
