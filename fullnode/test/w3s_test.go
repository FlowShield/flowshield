package test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"testing"

	"github.com/ipfs/go-cid"
	"github.com/web3-storage/go-w3s-client"
	"github.com/wumansgy/goEncrypt"
)

type Json struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var Ws3client, _ = w3s.NewClient(w3s.WithToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJkaWQ6ZXRocjoweDU4MUJkZEVGNTA3MDlmZjIzQzEwN0Q5YUU2NEVlMjc5M0IyMzk3NWMiLCJpc3MiOiJ3ZWIzLXN0b3JhZ2UiLCJpYXQiOjE2NTY2NDc2MDM2MjUsIm5hbWUiOiJjbG91ZHNsaXQifQ.7iUZuCDn1SNn7CxuR_kdAWf9_PfpuJlqPmy7ZdB2x9U"))

func TestW3SPut(t *testing.T) {
	file, err := os.CreateTemp("", "test")
	if err != nil {
		t.Error(err)
		return
	}
	defer os.Remove(file.Name())
	jsonByes, err := json.Marshal(Json{Name: "Tom", Age: 10})
	if err != nil {
		t.Error(err)
		return
	}

	cryptText, err := goEncrypt.DesCbcEncrypt(jsonByes, []byte("asd12345"), nil) //得到密文,可以自己传入初始化向量,如果不传就使用默认的初始化向量,8字节
	if err != nil {
		fmt.Println(err)
	}

	err = os.WriteFile(file.Name(), cryptText, 0644)
	if err != nil {
		t.Error(err)
		return
	}
	// Write a file/directory
	cidObj, err := Ws3client.Put(context.Background(), file)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("https://%v.ipfs.dweb.link\n", cidObj)
}

func TestW3SGet(t *testing.T) {
	cidObj, err := cid.Decode("bafybeihvhtz4h5f3wepzyyjpcz4fx3g6emdu35dyy7n7fytdb3gypehwzm")
	res, err := Ws3client.Get(context.Background(), cidObj)
	if err != nil {
		t.Fatal(err)
	}
	_, fsys, err := res.Files()
	if err != nil {
		t.Fatal(err)
	}
	var data []byte
	err = fs.WalkDir(fsys, "/", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			file, err := fsys.Open(path)
			if err != nil {
				t.Fatal(err)
			}
			data, err = ioutil.ReadAll(file)
		}
		return nil
	})
	newplaintext, err := goEncrypt.DesCbcDecrypt(data, []byte("csd88888"), nil) //解密得到密文,可以自己传入初始化向量,如果不传就使用默认的初始化向量,8字节
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(newplaintext))
}

func TestCipher(t *testing.T) {
	plaintext := []byte("床前明月光，疑是地上霜，举头望明月，学习go语言") //明文
	fmt.Println("明文为：", string(plaintext))
	// 传入明文和自己定义的密钥，密钥为8字节
	cryptText, err := goEncrypt.DesCbcEncrypt(plaintext, []byte("asd12345"), nil) //得到密文,可以自己传入初始化向量,如果不传就使用默认的初始化向量,8字节
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("DES的CBC模式加密后的密文为:", base64.StdEncoding.EncodeToString(cryptText))
	// 传入密文和自己定义的密钥，需要和加密的密钥一样，不一样会报错，8字节 如果解密秘钥错误解密后的明文会为空
	newplaintext, err := goEncrypt.DesCbcDecrypt(cryptText, []byte("asd12345"), nil) //解密得到密文,可以自己传入初始化向量,如果不传就使用默认的初始化向量,8字节
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("DES的CBC模式解密完：", string(newplaintext))
}
