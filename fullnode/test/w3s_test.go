package test

import (
	"context"
	"fmt"
	"io/fs"
	"io/ioutil"
	"testing"

	"github.com/ipfs/go-cid"
	"github.com/web3-storage/go-w3s-client"
)

var Ws3client, _ = w3s.NewClient(w3s.WithToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJkaWQ6ZXRocjoweDU4MUJkZEVGNTA3MDlmZjIzQzEwN0Q5YUU2NEVlMjc5M0IyMzk3NWMiLCJpc3MiOiJ3ZWIzLXN0b3JhZ2UiLCJpYXQiOjE2NTY2NDc2MDM2MjUsIm5hbWUiOiJjbG91ZHNsaXQifQ.7iUZuCDn1SNn7CxuR_kdAWf9_PfpuJlqPmy7ZdB2x9U"))

func TestW3S(t *testing.T) {
	cidObj, err := cid.Decode("bafybeibxn7xd2fhrrabuz5ndhi4zkscwx3a2cvprckjbslxyh5j67myb7y")
	fmt.Println(cidObj.String(), err)
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
	fmt.Println(string(data))
	if err != nil {
		t.Fatal(err)
	}
}
