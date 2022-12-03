package agent

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/flowshield/flowshield/fullnode/app/v1/user/model/mmysql"

	mysqlUser "github.com/flowshield/flowshield/fullnode/app/v1/user/dao/mysql"
	serviceUser "github.com/flowshield/flowshield/fullnode/app/v1/user/service"
	"github.com/flowshield/flowshield/fullnode/pkg/logger"
)

// SyncUserBindStatus 每30秒轮循未绑定的用户，判断是否已绑定
func SyncUserBindStatus() {
	closeChan := make(chan interface{}) // 信号监控
	go func() {
		for range closeChan {
			go handleUserBindStatus(closeChan)
		}
	}()
	go handleUserBindStatus(closeChan)
}

func handleUserBindStatus(closeChan chan<- interface{}) {
	userDao := mysqlUser.NewUser(nil)
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf(nil, "handleUserBindStatus stacktrace from panic:\n"+string(debug.Stack()), fmt.Errorf("%v", err))
		}
		closeChan <- struct{}{}
	}()
	for {
		select {
		case <-time.After(time.Second * 10):
			userList, err := userDao.GetUserByStatus(mmysql.UnBind)
			if err != nil {
				logger.Errorf(nil, "handleUserBindStatus get user by status error: %v", err)
				continue
			}
			for _, value := range userList {
				time.Sleep(5 * time.Second)
				serviceUser.CheckAndBindUser(&value)
			}
		}
	}
}
