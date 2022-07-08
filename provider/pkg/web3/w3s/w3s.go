package w3s

import (
	"context"
	"encoding/json"
	"github.com/cloudslit/cloudslit/provider/internal/config"
	"io/fs"
	"io/ioutil"
	"os"

	"github.com/ipfs/go-cid"
	"github.com/web3-storage/go-w3s-client"
)

var client w3s.Client

func Init(cfg *config.Web3) (err error) {
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

func Get(ctx context.Context, cidStr string) (data []byte, err error) {
	cidObj, err := cid.Decode(cidStr)
	if err != nil {
		return nil, err
	}
	res, err := client.Get(ctx, cidObj)
	if err != nil {
		return nil, err
	}
	_, fsys, err := res.Files()
	if err != nil {
		return nil, err
	}
	err = fs.WalkDir(fsys, "/", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			file, err := fsys.Open(path)
			if err != nil {
				return err
			}
			data, err = ioutil.ReadAll(file)
			if err != nil {
				return err
			}
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	return
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
