package schema

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/flowshield/flowshield/client/pkg/certificate"
	"github.com/flowshield/flowshield/client/pkg/errors"
	"github.com/flowshield/flowshield/client/pkg/util/json"
	"github.com/flowshield/flowshield/client/pkg/web3/w3s"
	"net/url"
	"strings"
)

type ClientConfig struct {
	Server  Server `json:"server"`
	Target  Target `json:"target"`
	CaPem   string `json:"ca_pem"`
	CertPem string `json:"cert_pem"`
	KeyPem  string `json:"key_pem"`
}

type Server struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type Target struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func (a *ClientConfig) LoadServerTarget(data []byte) (*ClientConfig, error) {
	var result ClientConfig
	err := json.Unmarshal(data, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	a.Server = result.Server
	a.Target = result.Target
	return a, nil
}

func ParseClientConfig(attrs map[string]interface{}) (*ClientConfig, error) {
	var result ClientConfig
	attrByte, err := json.Marshal(attrs)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	err = json.Unmarshal(attrByte, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if result.Server.Host == "" {
		err := errors.New("server Addr argument is missing")
		return nil, errors.WithStack(err)
	}
	return &result, nil
}

func (a *ClientConfig) String() string {
	return json.MarshalToString(a)
}

// ControlUserDetail Login user information
type ControlUserDetail struct {
	Uuid   string `json:"uuid"`
	Status string `json:"status"`
}

type ControlMachineAuthResult struct {
	ControlCommonResult
	Data string `json:"data"`
}

// GetCode
func (a *ControlMachineAuthResult) GetCode() string {
	purl, _ := url.Parse(a.Data)
	psurl := strings.Split(purl.Path, "/")
	return psurl[len(psurl)-1]
}

// ControlLoginResult Device login
type ControlLoginResult struct {
	ControlCommonResult
	Data string `json:"data"`
}

// ControlClientResult Client list
type ControlClientResult struct {
	ControlCommonResult
	Data ControlClientData
}

// ControlClientData
type ControlClientData struct {
	List     ControlClients  `json:"list"`
	Paginate ControlPaginate `json:"paginate"`
}

type ControlClients []*ControlClient

// ControlClient Controller Client
type ControlClient struct {
	PeerId    string `json:"peer_id"`
	Uuid      string `json:"uuid"`
	UserUuid  string `json:"user_uuid"`
	Name      string `json:"name"`
	ClientCid string `json:"client_cid"`
}

// ToClientOrder 
func (a *ControlClient) ToClientOrder(ctx context.Context, key []byte) (*ClientConfig, error) {
	data, err := w3s.Get(ctx, a.ClientCid, a.Uuid, key)
	if err != nil {
		return nil, err
	}
	// Parse w3s data
	var result ClientConfig
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	capem, err := base64.StdEncoding.DecodeString(result.CaPem)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	certpem, err := base64.StdEncoding.DecodeString(result.CertPem)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	keypem, err := base64.StdEncoding.DecodeString(result.KeyPem)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	result.CaPem = string(capem)
	result.CertPem = string(certpem)
	result.KeyPem = string(keypem)
	// Parse certificate attr
	certConfig, attr, err := certificate.LoadCert([]byte(result.CertPem))
	if err != nil {
		return nil, err
	}
	if certConfig.Type != certificate.TypeClient {
		return nil, fmt.Errorf("Wrong certificate type, expected:%s, get:%s", certificate.TypeClient, certConfig.Type)
	}
	// Load server and target information
	_, err = result.LoadServerTarget(attr)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
