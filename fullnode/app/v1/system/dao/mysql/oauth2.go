package mysql

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/cloudslit/cloudslit/fullnode/app/v1/system/model/mmysql"
	"github.com/cloudslit/cloudslit/fullnode/pkg/logger"
	"github.com/cloudslit/cloudslit/fullnode/pkg/mysql"
	"gorm.io/gorm"
)

type Oauth2 struct {
	c *gin.Context
	mysql.DaoMysql
}

func NewOauth2(c *gin.Context) *Oauth2 {
	return &Oauth2{
		DaoMysql: mysql.DaoMysql{TableName: "zta_oauth2"},
		c:        c,
	}
}

func (p *Oauth2) ListOauth2() (list []mmysql.Oauth2, err error) {
	orm := p.GetOrm().DB
	err = orm.Table(p.TableName).Find(&list).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	if err != nil {
		logger.Errorf(p.c, "ListOauth2 err : %v", err)
	}
	return
}

func (p *Oauth2) GetOauth2ByID(id int64) (info *mmysql.Oauth2, err error) {
	orm := p.GetOrm()
	err = orm.Table(p.TableName).Where(fmt.Sprintf("id = %d", id)).First(&info).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	if err != nil {
		logger.Errorf(p.c, "GetOauth2ById err : %v", err)
	}
	return
}

func (p *Oauth2) GetOauth2ByCompany(company string) (info *mmysql.Oauth2, err error) {
	orm := p.GetOrm()
	err = orm.Where(fmt.Sprintf("company = '%s'", company)).First(&info).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	if err != nil {
		logger.Errorf(p.c, "GetOauth2ByCompany err : %v", err)
	}
	return
}

func (p *Oauth2) AddOauth2(data *mmysql.Oauth2) (err error) {
	data.CreatedAt = time.Now().Unix()
	data.UpdatedAt = data.CreatedAt
	scopes, err := json.Marshal(data.Scopes)
	query := fmt.Sprintf("%d,%d,'%s','%s','%s','%s','%s','%s','%s'", data.CreatedAt, data.UpdatedAt,
		data.Company, data.ClientId, data.ClientSecret, data.RedirectUrl, string(scopes), data.AuthUrl, data.TokenUrl)
	err = p.GetOrm().Exec("INSERT INTO `zta_oauth2` (`created_at`,`updated_at`," +
		"`company`,`client_id`,`client_secret`,`redirect_url`,`scopes`,`auth_url`," +
		"`token_url`) VALUES (" + query + ")").Error
	if err != nil {
		logger.Errorf(p.c, "AddOauth2 err : %v", err)
	}
	return
}

func (p *Oauth2) EditOauth2(data *mmysql.Oauth2) (err error) {
	orm := p.GetOrm()
	err = orm.Table(p.TableName).Save(&data).Error
	if err != nil {
		logger.Errorf(p.c, "EditOauth2 err : %v", err)
	}
	return
}

func (p *Oauth2) DelOauth2(id uint64) (err error) {
	orm := p.GetOrm()
	err = orm.Where(fmt.Sprintf("id = %d", id)).Delete(&mmysql.Oauth2{}).Error
	//err = orm.Delete(&mmysql.Oauth2{}, id).Error
	if err != nil {
		logger.Errorf(p.c, "DelOauth2 err : %v", err)
	}
	return
}
