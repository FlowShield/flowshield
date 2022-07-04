package w3s

import (
	"context"
	"encoding/json"
	"os"

	"github.com/web3-storage/go-w3s-client"
	"github.com/cloudslit/cloudslit/fullnode/pkg/confer"
)

var client w3s.Client

func Init(cfg *confer.Web3) (err error) {
	client, err = w3s.NewClient(w3s.WithToken(cfg.W3S.Token))
	if err != nil {
		return
	}
	return
}

func Put(ctx context.Context, data interface{}) (cid string, err error) {
	file, err := dataToFile(data)
	defer os.Remove(file.Name())
	if err != nil {
		return
	}
	cidObj, err := client.Put(ctx, file)
	if err != nil {
		return
	}
	return cidObj.String(), nil
}

func dataToFile(data interface{}) (file *os.File, err error) {
	jsonByes, err := json.Marshal(data)
	if err != nil {
		return
	}
	file, err = os.CreateTemp("", "data")
	if err != nil {
		return
	}
	// TODO 对数据进行加密
	err = os.WriteFile(file.Name(), jsonByes, 0644)
	if err != nil {
		return
	}
	return
}
