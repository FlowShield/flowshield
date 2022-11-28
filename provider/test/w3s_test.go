package test

import (
	"context"
	"flag"
	"fmt"
	"github.com/cloudslit/cloudslit/provider/internal/config"
	"log"
	"net/http"
	"strconv"
	"testing"
	"time"

	w3sutil "github.com/cloudslit/cloudslit/provider/pkg/web3/w3s"
)

const Token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJkaWQ6ZXRocjoweDU4MUJkZEVGNTA3MDlmZjIzQzEwN0Q5YUU2NEVlMjc5M0IyMzk3NWMiLCJpc3MiOiJ3ZWIzLXN0b3JhZ2UiLCJpYXQiOjE2NTY2NDc2MDM2MjUsIm5hbWUiOiJjbG91ZHNsaXQifQ.7iUZuCDn1SNn7CxuR_kdAWf9_PfpuJlqPmy7ZdB2x9U"

var t1 = flag.String("t", "10", "访w3s超时时间")

var httpClient = http.Client{}

func TestW3s(t *testing.T) {
	flag.Parse()
	ctx := context.Background()
	to, _ := strconv.Atoi(*t1)
	httpClient.Timeout = time.Duration(to) * time.Second
	cfg := &config.Web3{
		W3S: config.W3S{
			Token:      Token,
			Timeout:    to,
			RetryCount: 10,
		},
	}
	err := w3sutil.Init(cfg)
	if err != nil {
		log.Fatalf("Init err:%v", err)
	}
	filename := "filename"
	key := []byte("12345678")
	content := time.Now().String()
	cid, err := w3sutil.Put(ctx, content, filename, key)
	if err != nil {
		log.Fatalf("Put err:%v", err)
	}
	fmt.Printf("https://%v.ipfs.w3s.link/%s\n", cid, filename)
	fmt.Printf("https://api.web3.storage/car/%v\n", cid)
	errCount := 0
	sucCount := 0
	sumStartTime := time.Now().UnixNano()
retry:
	startTime := time.Now().UnixNano()
	data, err := w3sutil.Get(ctx, cid, filename, key)
	endTime := time.Now().UnixNano()
	seconds := float64((endTime - startTime) / 1e9)
	log.Printf("cost time : %.4f s", seconds)
	if err != nil {
		errCount++
		log.Println("Get err:", err.Error())
		goto retry
	}
	log.Println(string(data))
	sumEndTime := time.Now().UnixNano()
	sumSeconds := float64((sumEndTime - sumStartTime) / 1e9)
	totalSeconds := fmt.Sprintf("%.4f s", sumSeconds)
	sucCount++
	time.Sleep(time.Second)
	if sucCount >= 3 {
		log.Println("失败次数:", errCount)
		log.Println("总花费时间:", totalSeconds)
		return
	}
	goto retry
}
