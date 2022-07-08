package schema

type ServerType string

const ServerTypeProvider ServerType = "provider"

const FullNode ServerType = "fullnode"

type ServerInfo struct {
	PeerId   string     `json:"peer_id"`
	Addr     string     `json:"addr"`
	Port     int        `json:"port"`
	MetaData MetaData   `json:"meta_data"`
	Price    uint       `json:"price"`
	Type     ServerType `json:"type"`
}

type MetaData struct {
	Ip   string `json:"ip"`
	Loc  string `json:"loc"`
	Colo string `json:"colo"`
}

type ClientP2P struct {
	ServerCID string `json:"server_cid"`
	Wallet    string `json:"wallet"`
	UUID      string `json:"uuid"`
	Port      int    `json:"port"`
}
