package mysql

import (
	"errors"
	"fmt"
	"time"

	"github.com/flowshield/flowshield/fullnode/app/v1/access/model/mapi"

	"github.com/flowshield/flowshield/fullnode/pkg/util"

	"github.com/flowshield/flowshield/fullnode/app/v1/access/model/mmysql"
	"github.com/flowshield/flowshield/fullnode/app/v1/access/model/mparam"
	"github.com/flowshield/flowshield/fullnode/pkg/logger"
	"github.com/flowshield/flowshield/fullnode/pkg/mysql"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Client struct {
	c *gin.Context
	mysql.DaoMysql
}

func NewClient(c *gin.Context) *Client {
	return &Client{
		DaoMysql: mysql.DaoMysql{TableName: "zta_client"},
		c:        c,
	}
}

func (p *Client) ClientList(param mparam.ClientList) (
	total int64, list []mapi.Client, err error) {
	orm := p.GetOrm().DB
	query := orm.Table(p.TableName).Select(p.TableName + ".*,zta_node.ip as node_ip").
		Joins("left join zta_node on zta_client.peer_id = zta_node.peer_id")
	if len(param.Name) > 0 {
		query = query.Where(fmt.Sprintf(p.TableName+".name like '%%%s%%'", param.Name))
	}
	if len(param.PeerID) > 0 {
		query = query.Where(fmt.Sprintf(p.TableName+".peer_id = '%s'", param.PeerID))
	}
	if param.Working {
		query = query.Where(fmt.Sprintf(p.TableName+".status = %d", mmysql.Success))
		query = query.Where(fmt.Sprintf(p.TableName+".updated_at + 60*60*duration >= %d", time.Now().Unix()))
	}
	if user := util.User(p.c); user != nil {
		query = query.Where(fmt.Sprintf(p.TableName+".user_uuid = '%s'", user.UUID))
	}
	err = query.Model(&list).Count(&total).Error
	if total > 0 {
		offset := param.GetOffset()
		err = query.Limit(param.LimitNum).Offset(offset).
			Order(p.TableName + ".created_at desc").
			Find(&list).Error
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	if err != nil {
		logger.Errorf(p.c, "ClientList err : %v", err)
	}
	return
}

func (p *Client) GetClientByID(id uint64) (info *mmysql.Client, err error) {
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
		logger.Errorf(p.c, "GetClientById err : %v", err)
	}
	return
}

func (p *Client) GetClientByUUID(uuid string) (info *mmysql.Client, err error) {
	orm := p.GetOrm()
	query := orm.Table(p.TableName).Where(fmt.Sprintf("uuid = '%s'", uuid))
	if user := util.User(p.c); user != nil {
		query = query.Where(fmt.Sprintf("user_uuid = '%s'", user.UUID))
	}
	err = query.First(&info).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	if err != nil {
		logger.Errorf(p.c, "GetClientByUUID err : %v", err)
	}
	return
}

func (p *Client) AddClient(data *mmysql.Client) (err error) {
	if user := util.User(p.c); user != nil {
		data.UserUUID = user.UUID
	}
	orm := p.GetOrm()
	sql := orm.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Table(p.TableName).Create(&data)
	})
	err = orm.Exec(sql).Error
	if err != nil {
		logger.Errorf(p.c, "AddClient err : %v", err)
	}
	return
}

func (p *Client) EditClient(data *mmysql.Client) (err error) {
	if user := util.User(p.c); user != nil {
		data.UserUUID = user.UUID
	}
	orm := p.GetOrm()
	sql := orm.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Table(p.TableName).Save(&data)
	})
	err = orm.Exec(sql).Error
	if err != nil {
		logger.Errorf(p.c, "EditClient err : %v", err)
	}
	return
}

func (p *Client) CheckStatusClient(param *mparam.CheckStatus) (list []mmysql.Client, err error) {
	orm := p.GetOrm().DB
	query := orm.Table(p.TableName)
	query = query.Where(fmt.Sprintf("status = %d", param.Status))
	if param.Duration > 0 {
		query = query.Where(fmt.Sprintf("created_at >= %d", time.Now().Add(-(param.Duration * time.Minute)).Unix()))
	}
	err = query.Find(&list).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	if err != nil {
		logger.Errorf(p.c, "CheckStatusClient err : %v", err)
	}
	return
}

func (p *Client) DelClient(uuid string) (err error) {
	orm := p.GetOrm()
	query := orm.Table(p.TableName).Where(fmt.Sprintf("uuid = %s", uuid))
	if user := util.User(p.c); user != nil {
		query = query.Where(fmt.Sprintf("user_uuid = '%s'", user.UUID))
	}
	err = query.Delete(&mmysql.Client{}).Error
	if err != nil {
		logger.Errorf(p.c, "DelClient err : %v", err)
	}
	return
}
