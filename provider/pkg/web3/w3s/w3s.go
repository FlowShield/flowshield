package w3s

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudslit/cloudslit/provider/pkg/errors"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cloudslit/cloudslit/provider/internal/config"

	"github.com/ipfs/go-cid"
	"github.com/web3-storage/go-w3s-client"
	"github.com/wumansgy/goEncrypt"
)

var client w3s.Client
var httpClient http.Client

func Init(cfg *config.Web3) (err error) {
	client, err = w3s.NewClient(w3s.WithToken(cfg.W3S.Token))
	if err != nil {
		return
	}
	httpClient.Timeout = time.Duration(cfg.W3S.Timeout) * time.Second
	return
}

func Put(ctx context.Context, data interface{}, filename string, key []byte) (cid string, err error) {
	file, err := dataToFile(data, filename, key)
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

func Get(ctx context.Context, cidStr string, filename string, key []byte) ([]byte, error) {
	var result []byte
	var err error
	result, err = GetByW3sLink(cidStr, filename, key)
	if err != nil {
		log.Println("Failed to obtain from w3s link, err:", err.Error())
		ctx, cancel := context.WithTimeout(ctx, time.Duration(config.C.Web3.W3S.Timeout)*time.Second)
		defer cancel()
		result, err = GetByW3sClient(ctx, cidStr, key)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}
	return result, nil
}

func GetByW3sLink(cid string, filename string, key []byte) ([]byte, error) {
	url := fmt.Sprintf("https://%s.ipfs.w3s.link/%s", cid, filename)
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request error, url:%s, code:%d, msg:%s", url, resp.StatusCode, resp.Status)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	data, err := goEncrypt.DesCbcDecrypt(body, key[:], nil) // 解密得到密文,可以自己传入初始化向量,如果不传就使用默认的初始化向量,8字节
	if err != nil {
		return nil, err
	}
	return data, nil
}

func GetByW3sClient(ctx context.Context, cidStr string, key []byte) (data []byte, err error) {
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
	data, err = goEncrypt.DesCbcDecrypt(data, key[:], nil) // 解密得到密文,可以自己传入初始化向量,如果不传就使用默认的初始化向量,8字节
	if err != nil {
		return nil, err
	}
	return
}

func dataToFile(data interface{}, filename string, key []byte) (file *os.File, err error) {
	jsonByes, err := json.Marshal(data)
	if err != nil {
		return
	}
	file, err = os.Create(filename)
	if err != nil {
		return
	}
	// 对数据进行加密
	cryptText, err := goEncrypt.DesCbcEncrypt(jsonByes, key[:], nil)
	if err != nil {
		return
	}
	err = os.WriteFile(file.Name(), cryptText, 0o644)
	if err != nil {
		return
	}
	return
}
