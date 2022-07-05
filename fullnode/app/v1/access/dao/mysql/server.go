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

type Server struct {
	c *gin.Context
	mysql.DaoMysql
}

func NewServer(c *gin.Context) *Server {
	return &Server{
		DaoMysql: mysql.DaoMysql{TableName: "zta_server"},
		c:        c,
	}
}

func (p *Server) ServerList(param mparam.ServerList) (
	total int64, list []mmysql.Server, err error) {
	orm := p.GetOrm().DB
	query := orm.Table(p.TableName)
	if len(param.Name) > 0 {
		query = query.Where(fmt.Sprintf("name like '%%%s%%'", param.Name))
	}
	if param.ResourceID > 0 {
		query = query.Where(fmt.Sprintf("find_in_set (%d,resource_id)", param.ResourceID))
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
		logger.Errorf(p.c, "ServerList err : %v", err)
	}
	return
}

func (p *Server) GetServerByID(id uint64) (info *mmysql.Server, err error) {
	orm := p.GetOrm()
	query := orm.Table(p.TableName).Where(fmt.Sprintf("id = %d", id))
	if user := util.User(p.c); user != nil {
		query = query.Where(fmt.Sprintf("user_uuid = '%s'", user.UUID))
	}
	err = query.First(&info).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	if err != nil {
		logger.Errorf(p.c, "GetServerById err : %v", err)
	}
	return
}

func (p *Server) GetServerByUUID(uuid string) (info *mmysql.Server, err error) {
	orm := p.GetOrm()
	query := orm.Table(p.TableName).Where(fmt.Sprintf("uuid = %s", uuid))
	if user := util.User(p.c); user != nil {
		query = query.Where(fmt.Sprintf("user_uuid = '%s'", user.UUID))
	}
	err = query.First(&info).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	if err != nil {
		logger.Errorf(p.c, "GetServerByUUID err : %v", err)
	}
	return
}

func (p *Server) AddServer(data *mmysql.Server) (err error) {
	if user := util.User(p.c); user != nil {
		data.UserUUID = user.UUID
	}
	orm := p.GetOrm()
	sql := orm.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Table(p.TableName).Create(&data)
	})
	err = orm.Exec(sql).Error
	if err != nil {
		logger.Errorf(p.c, "AddServer err : %v", err)
	}
	return
}

func (p *Server) EditServer(data *mmysql.Server) (err error) {
	if user := util.User(p.c); user != nil {
		data.UserUUID = user.UUID
	}
	orm := p.GetOrm()
	sql := orm.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Table(p.TableName).Save(&data)
	})
	err = orm.Exec(sql).Error
	if err != nil {
		logger.Errorf(p.c, "EditServer err : %v", err)
	}
	return
}

func (p *Server) DelServer(uuid string) (err error) {
	orm := p.GetOrm()
	query := orm.Table(p.TableName).Where(fmt.Sprintf("uuid = %s", uuid))
	if user := util.User(p.c); user != nil {
		query = query.Where(fmt.Sprintf("user_uuid = '%s'", user.UUID))
	}
	err = query.Delete(&mmysql.Server{}).Error
	if err != nil {
		logger.Errorf(p.c, "DelServer err : %v", err)
	}
	return
}
