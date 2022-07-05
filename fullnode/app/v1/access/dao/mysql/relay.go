package mysql

import (
	"errors"
	"fmt"

	"github.com/cloudslit/cloudslit/fullnode/pkg/util"

	"github.com/cloudslit/cloudslit/fullnode/app/v1/access/model/mmysql"
	"github.com/cloudslit/cloudslit/fullnode/app/v1/access/model/mparam"
	"github.com/cloudslit/cloudslit/fullnode/pkg/logger"
	"github.com/cloudslit/cloudslit/fullnode/pkg/mysql"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Relay struct {
	c *gin.Context
	mysql.DaoMysql
}

func NewRelay(c *gin.Context) *Relay {
	return &Relay{
		DaoMysql: mysql.DaoMysql{TableName: "zta_relay"},
		c:        c,
	}
}

func (p *Relay) RelayList(param mparam.RelayList) (
	total int64, list []mmysql.Relay, err error) {
	orm := p.GetOrm().DB
	query := orm.Table(p.TableName)
	if len(param.Name) > 0 {
		query = query.Where(fmt.Sprintf("name like '%%%s%%'", param.Name))
	}
	if user := util.User(p.c); user != nil {
		query = query.Where(fmt.Sprintf("user_uuid = '%s'", user.UUID))
	}
	err = query.Model(&list).Count(&total).Error
	if total > 0 {
		offset := param.GetOffset()
		err = query.Limit(param.LimitNum).Offset(offset).
			Order("created_at desc").
			Find(&list).Error
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	if err != nil {
		logger.Errorf(p.c, "RelayList err : %v", err)
	}
	return
}

func (p *Relay) GetRelayByID(id uint64) (info mmysql.Relay, err error) {
	orm := p.GetOrm()
	query := orm.Table(p.TableName).Where(fmt.Sprintf("id = %d", id))
	if user := util.User(p.c); user != nil {
		query = query.Where(fmt.Sprintf("user_uuid = '%s'", user.UUID))
	}
	err = query.First(&info).Error
	if err != nil {
		logger.Errorf(p.c, "GetRelayById err : %v", err)
	}
	return
}

func (p *Relay) AddRelay(data *mmysql.Relay) (err error) {
	if user := util.User(p.c); user != nil {
		data.UserUUID = user.UUID
	}
	orm := p.GetOrm()
	sql := orm.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Table(p.TableName).Create(&data)
	})
	err = orm.Exec(sql).Error
	if err != nil {
		logger.Errorf(p.c, "AddRelay err : %v", err)
	}
	return
}

func (p *Relay) EditRelay(data mmysql.Relay) (err error) {
	if user := util.User(p.c); user != nil {
		data.UserUUID = user.UUID
	}
	orm := p.GetOrm()
	sql := orm.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Table(p.TableName).Save(&data)
	})
	err = orm.Exec(sql).Error
	if err != nil {
		logger.Errorf(p.c, "EditRelay err : %v", err)
	}
	return
}

func (p *Relay) DelRelay(uuid string) (err error) {
	orm := p.GetOrm()
	query := orm.Table(p.TableName).Where(fmt.Sprintf("uuid = %s", uuid))
	if user := util.User(p.c); user != nil {
		query = query.Where(fmt.Sprintf("user_uuid = '%s'", user.UUID))
	}
	err = query.Delete(&mmysql.Relay{}).Error
	if err != nil {
		logger.Errorf(p.c, "DelRelay err : %v", err)
	}
	return
}
