package schema

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/cloudslit/cloudslit/provider/pkg/certificate"
	"github.com/cloudslit/cloudslit/provider/pkg/errors"
	"github.com/cloudslit/cloudslit/provider/pkg/util/json"
	"github.com/cloudslit/cloudslit/provider/pkg/web3/w3s"
	"time"
)

type NodeType string

const NodeTypeProvider NodeType = "provider"

type NodeInfo struct {
	PeerId   string   `json:"peer_id"`
	Addr     string   `json:"addr"`
	Port     int      `json:"port"`
	MetaData MetaData `json:"meta_data"`
	Price    int      `json:"price"`
	Type     NodeType `json:"type"`
}

type MetaData struct {
	Ip   string `json:"ip"`
	Loc  string `json:"loc"`
	Colo string `json:"colo"`
}

func (a *NodeInfo) String() string {
	return json.MarshalToString(a)
}

type ClientConfig struct {
	Server Server `json:"server"`
	Target Target `json:"target"`
}

type Server struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type Target struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

const PsMsgTypeNode = "node"
const PsMsgTypeOrder = "order"

type PsMessage struct {
	Type string      `json:"type"` // node, order
	Data interface{} `json:"data"`
}

func (a *PsMessage) String() string {
	return json.MarshalToString(a)
}

// ToNodeOrder 转换为节点订单信息
func (a *PsMessage) ToNodeOrder() (*NodeOrder, error) {
	if a.Type != PsMsgTypeOrder {
		return nil, errors.New("Pubsub Message Type Err,expect: order or node")
	}
	var result NodeOrder
	arr, err := json.Marshal(a.Data)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	err = json.Unmarshal(arr, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &result, nil
}

type NodeOrder struct {
	Uuid      string `json:"uuid"`
	Wallet    string `json:"wallet"`
	ServerCid string `json:"server_cid"`
	Port      int    `json:"port"`
}

func (a *NodeOrder) String() string {
	return json.MarshalToString(a)
}

// ProviderConfig provider配置
type ProviderConfig struct {
	CertPem    string `json:"cert_pem"`
	KeyPem     string `json:"key_pem"`
	CaPem      string `json:"ca_pem"`
	CertConfig *certificate.BasicCertConf
}

func (a *ProviderConfig) String() string {
	return json.MarshalToString(a)
}

// ParserConfig 解析config
func ParserConfig(ctx context.Context, cid string, key []byte) (*ProviderConfig, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	data, err := w3s.Get(ctx, cid, key)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var pc ProviderConfig
	err = json.Unmarshal(data, &pc)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	capem, err := base64.StdEncoding.DecodeString(pc.CaPem)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	certpem, err := base64.StdEncoding.DecodeString(pc.CertPem)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	keypem, err := base64.StdEncoding.DecodeString(pc.KeyPem)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	pc.CaPem = string(capem)
	pc.CertPem = string(certpem)
	pc.KeyPem = string(keypem)
	// 解析证书attr
	certConfig, _, err := certificate.LoadCert([]byte(pc.CertPem))
	if err != nil {
		return nil, err
	}
	if certConfig.Type != certificate.TypeServer {
		return nil, fmt.Errorf("证书类型错误，预期：%s, get:%s", certificate.TypeServer, certConfig.Type)
	}
	pc.CertConfig = certConfig
	return &pc, nil
}
