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

// CreateProviderParams 创建网络服务参数
type CreateProviderParams struct {
	Certificate Certificate `json:"certificate"`
}

// CreateProviderParams 创建客户端参数
type CreateClientParams struct {
	Certificate Certificate `json:"certificate"`
}

// Certificate 证书信息
type Certificate struct {
	Cert string `json:"cert"`
	Key  string `json:"key"`
	Ca   string `json:"ca"`
}

func (a *CreateProviderParams) String() string {
	return json.MarshalToString(a)
}
