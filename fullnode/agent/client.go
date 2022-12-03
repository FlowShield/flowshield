package agent

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/flowshield/flowshield/fullnode/app/v1/access/service"

	"github.com/flowshield/flowshield/fullnode/app/v1/access/model/mmysql"

	"github.com/flowshield/flowshield/fullnode/app/v1/access/dao/mysql"
	"github.com/flowshield/flowshield/fullnode/app/v1/access/model/mparam"

	"github.com/flowshield/flowshield/fullnode/pkg/logger"
)

// SyncClientOrderStatus 每30秒轮循10分钟内的订单，判断是否已支付
func SyncClientOrderStatus() {
	closeChan := make(chan interface{}) // 信号监控
	go func() {
		for range closeChan {
			go handleClientOrderStatus(closeChan)
		}
	}()
	go handleClientOrderStatus(closeChan)
}

func handleClientOrderStatus(closeChan chan<- interface{}) {
	clientDao := mysql.NewClient(nil)
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf(nil, "handleClientOrderStatus stacktrace from panic:\n"+string(debug.Stack()), fmt.Errorf("%v", err))
		}
		closeChan <- struct{}{}
	}()
	for {
		select {
		case <-time.After(time.Second * 30):
			clientList, err := clientDao.CheckStatusClient(&mparam.CheckStatus{
				Status:   mmysql.WaitingPaid,
				Duration: 10,
			})
			if err != nil {
				continue
			}
			for _, value := range clientList {
				time.Sleep(5 * time.Second)
				service.NotifyClient(nil, &mparam.NotifyClient{UUID: value.UUID})
			}
		}
	}
}
