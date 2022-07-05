package service

import (
	"errors"

	"github.com/cloudslit/cloudslit/fullnode/app/v1/system/dao/mysql"
	"github.com/cloudslit/cloudslit/fullnode/app/v1/system/model/mmysql"
	"github.com/cloudslit/cloudslit/fullnode/app/v1/system/model/mparam"
	"github.com/cloudslit/cloudslit/fullnode/pconst"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"

	"github.com/gin-gonic/gin"
)

func ListOauth2(c *gin.Context) (code int, data []mmysql.Oauth2) {
	data, err := mysql.NewOauth2(c).ListOauth2()
	if err != nil {
		code = pconst.CODE_COMMON_SERVER_BUSY
		return
	}
	return
}

func AddOauth2(c *gin.Context, param *mmysql.Oauth2) (code int) {
	// 判断服务端哨兵是否存在
	data, err := mysql.NewOauth2(c).GetOauth2ByCompany(param.Company)
	if err != nil {
		return pconst.CODE_COMMON_SERVER_BUSY
	}
	if data.ID > 0 {
		return pconst.CODE_COMMON_DATA_ALREADY_EXIST
	}
	err = mysql.NewOauth2(c).AddOauth2(param)
	if err != nil {
		return pconst.CODE_COMMON_SERVER_BUSY
	}
	return
}

func EditOauth2(c *gin.Context, param *mparam.EditOauth2) (code int) {
	info, err := mysql.NewOauth2(c).GetOauth2ByID(param.ID)
	if err != nil {
		code = pconst.CODE_COMMON_SERVER_BUSY
		return
	}
	if info.ID == 0 {
		code = pconst.CODE_COMMON_DATA_NOT_EXIST
		return
	}
	// 判断当前类型是否重复
	info, err = mysql.NewOauth2(c).GetOauth2ByCompany(param.Company)
	if err != nil {
		code = pconst.CODE_COMMON_SERVER_BUSY
		return
	}
	if info.ID > 0 && int64(info.ID) != param.ID {
		code = pconst.CODE_COMMON_DATA_ALREADY_EXIST
		return
	}
	info.Company = param.Company
	info.ClientId = param.ClientId
	info.ClientSecret = param.ClientSecret
	info.RedirectUrl = param.RedirectUrl
	info.Scopes = param.Scopes
	info.AuthUrl = param.AuthUrl
	info.TokenUrl = param.TokenUrl
	err = mysql.NewOauth2(c).EditOauth2(info)
	if err != nil {
		code = pconst.CODE_COMMON_SERVER_BUSY
		return
	}
	return
}

func DelOauth2(c *gin.Context, id uint64) (code int) {
	err := mysql.NewOauth2(c).DelOauth2(id)
	if err != nil {
		code = pconst.CODE_COMMON_SERVER_BUSY
	}
	return
}

func Oauth2Config(info *mmysql.Oauth2) (config *oauth2.Config, err error) {
	if info == nil {
		err = errors.New("nil info")
		return
	}
	config = &oauth2.Config{
		ClientID:     info.ClientId,
		ClientSecret: info.ClientSecret,
		//Endpoint:     oauth2.Endpoint{},
		RedirectURL: info.RedirectUrl,
		Scopes:      info.Scopes,
	}
	if len(info.AuthUrl) > 0 && len(info.TokenUrl) > 0 {
		config.Endpoint = oauth2.Endpoint{
			AuthURL:  info.AuthUrl,
			TokenURL: info.TokenUrl,
		}
	} else {
		switch info.Company {
		case "github":
			config.Endpoint = github.Endpoint
		case "google":
			config.Endpoint = google.Endpoint
		case "facebook":
			config.Endpoint = facebook.Endpoint
		default:
			return nil, errors.New("wrong company")
		}
	}
	return
}
