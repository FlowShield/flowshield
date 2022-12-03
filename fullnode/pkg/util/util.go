package util

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"net"
	"regexp"
	"strings"

	"github.com/go-resty/resty/v2"

	"github.com/flowshield/flowshield/fullnode/app/v1/user/model/mmysql"

	"github.com/gin-gonic/gin"
)

func IsIP(ipv4 string) bool {
	if ip := net.ParseIP(ipv4); ip == nil {
		return false
	}
	return true
}

func IsCIDR(cidr string) bool {
	p := regexp.MustCompile(`^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\/([0-9]|[1-2][0-9]|3[0-2])$|^s*((([0-9A-Fa-f]{1,4}:){7}(:|([0-9A-Fa-f]{1,4})))|(([0-9A-Fa-f]{1,4}:){6}:([0-9A-Fa-f]{1,4})?)|(([0-9A-Fa-f]{1,4}:){5}(((:[0-9A-Fa-f]{1,4}){0,1}):([0-9A-Fa-f]{1,4})?))|(([0-9A-Fa-f]{1,4}:){4}(((:[0-9A-Fa-f]{1,4}){0,2}):([0-9A-Fa-f]{1,4})?))|(([0-9A-Fa-f]{1,4}:){3}(((:[0-9A-Fa-f]{1,4}){0,3}):([0-9A-Fa-f]{1,4})?))|(([0-9A-Fa-f]{1,4}:){2}(((:[0-9A-Fa-f]{1,4}){0,4}):([0-9A-Fa-f]{1,4})?))|(([0-9A-Fa-f]{1,4}:){1}(((:[0-9A-Fa-f]{1,4}){0,5}):([0-9A-Fa-f]{1,4})?))|(:(:|((:[0-9A-Fa-f]{1,4}){1,7}))))(%.+)?s*/([0-9]|[1-9][0-9]|1[0-1][0-9]|12[0-8])$`)
	return p.MatchString(cidr)
}

func GetCookieFromGin(ctx *gin.Context, key string) (value string) {
	value, _ = ctx.Cookie(key)
	return
}

func User(c *gin.Context) (user *mmysql.User) {
	if c == nil {
		return
	}
	if userBytes, ok := c.Get("user"); ok {
		if err := json.Unmarshal(userBytes.([]byte), &user); err == nil {
			return
		}
		return nil
	}
	return nil
}

func NewMd5(str ...string) string {
	h := md5.New()
	for _, v := range str {
		h.Write([]byte(v))
	}
	return hex.EncodeToString(h.Sum(nil))
}

func GetCftrace() (*CfTrace, error) {
	res, err := resty.New().R().Get("https://www.cloudflare.com/cdn-cgi/trace")
	result := new(CfTrace)
	sb := strings.Split(res.String(), "\n")
	for _, item := range sb {
		is := strings.Split(item, "=")
		if is[0] == "ip" {
			result.Ip = is[1]
		}
		if is[0] == "loc" {
			result.Loc = is[1]
		}
		if is[0] == "colo" {
			result.Colo = is[1]
		}
	}
	return result, err
}

type CfTrace struct {
	Ip   string `json:"ip"`
	Loc  string `json:"loc"`
	Colo string `json:"colo"`
}

func Base64Encode(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}
