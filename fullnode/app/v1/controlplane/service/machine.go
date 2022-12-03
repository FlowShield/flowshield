package service

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/flowshield/flowshield/fullnode/app/v1/controlplane/model/mparam"

	"github.com/gin-gonic/contrib/sessions"

	"github.com/flowshield/flowshield/fullnode/app/v1/controlplane/dao/redis"

	"github.com/flowshield/flowshield/fullnode/pkg/util"

	"github.com/flowshield/flowshield/fullnode/pconst"
	"github.com/flowshield/flowshield/fullnode/pkg/confer"
	"github.com/gin-gonic/gin"
)

func GetLoginUrl(c *gin.Context, machine string) (code int, loginURL string) {
	// 通过machineID和当前时间戳，计算出唯一的hash，作为登陆的地址path
	hash := util.NewMd5(fmt.Sprintf("%s%d", machine, time.Now().UnixNano()))
	// hash 放入redis缓存
	err := redis.NewMachine(c).SetLoginHash(machine, hash)
	if err != nil {
		code = pconst.CODE_COMMON_SERVER_BUSY
		return
	}
	domain := confer.ConfigAppGetString("domain", "")
	if len(os.Getenv(domain)) > 0 {
		domain = os.Getenv(domain)
	}
	loginURL = fmt.Sprintf("%s/a/%s", domain, hash)
	return
}

func MachineOauth(c *gin.Context, hash string) {
	// 判断当前hash是否存在或者是否在有消息内
	exist, _, err := redis.NewMachine(c).GetLoginHash(hash)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.New("oauth error"))
		return
	}
	if exist {
		session := sessions.Default(c)
		session.Set("machine", hash)
		session.Save()
		c.Redirect(http.StatusSeeOther, "/api/v1/user/login")
	} else {
		// TODO 重定向到404页面
		c.String(http.StatusNotFound, "auth key not exist or expired")
		return
	}
}

func MachineLongPoll(c *gin.Context, param mparam.MachineLongPoll) (data string, code int) {
	result, err := redis.NewMachine(c).SubMachineCookie(param.Category, time.Second*time.Duration(param.Timeout))
	if err != nil {
		code = pconst.CODE_COMMON_SERVER_BUSY
		return
	}
	if len(result) >= 2 {
		data = result[1]
		return
	}
	code = pconst.CODE_COMMON_DATA_NOT_EXIST
	return
}
