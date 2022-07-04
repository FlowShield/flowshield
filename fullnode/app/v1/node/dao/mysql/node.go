package mysql

import (
	"errors"
	"fmt"
	"time"

	"github.com/cloudslit/cloudslit/fullnode/app/v1/node/model/mparam"

	"github.com/gin-gonic/gin"
	"github.com/cloudslit/cloudslit/fullnode/app/v1/node/model/mmysql"
	"github.com/cloudslit/cloudslit/fullnode/pkg/logger"
	"github.com/cloudslit/cloudslit/fullnode/pkg/mysql"
	"gorm.io/gorm"
)

type Node struct {
	c *gin.Context
	mysql.DaoMysql
}

func NewNode(c *gin.Context) *Node {
	return &Node{
		DaoMysql: mysql.DaoMysql{TableName: "zta_node"},
		c:        c,
	}
}

func (p *Node) ListNode(param mparam.ListNode) (total int64, list []mmysql.Node, err error) {
	query := p.GetOrm().DB
	if len(param.PeerId) > 0 {
		query = query.Where(fmt.Sprintf("peer_id = '%s'", param.PeerId))
	}
	if len(param.IP) > 0 {
		query = query.Where(fmt.Sprintf("ip = '%s'", param.PeerId))
	}
	if len(param.Loc) > 0 {
		var loc string
		for _, value := range param.Loc {
			loc += fmt.Sprintf("'%s'", value) + ","
		}
		loc = loc[:len(loc)-1]
		if len(loc) > 0 {
			query = query.Where(fmt.Sprintf("loc in (%s)", loc))
		}
	}
	if len(param.Colo) > 0 {
		query = query.Where(fmt.Sprintf("colo = '%s'", param.Colo))
	}
	if param.GasPrice > 0 {
		query = query.Where(fmt.Sprintf("gas_price = %d", param.GasPrice))
	}
	if len(param.Type) > 0 {
		query = query.Where(fmt.Sprintf("`type` = '%s'", param.Type))
	}
	err = query.Model(&list).Count(&total).Error
	if total > 0 {
		offset := param.GetOffset()
		err = query.Limit(param.LimitNum).Offset(offset).
			Order("updated_at desc").
			Find(&list).Error
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	if err != nil {
		logger.Errorf(p.c, "ListNode err : %v", err)
	}
	return
}

func (p *Node) GetNodeByPeerId(peerId string) (info *mmysql.Node, err error) {
	orm := p.GetOrm()
	err = orm.Table(p.TableName).Where(fmt.Sprintf("peer_id = '%s'", peerId)).First(&info).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	if err != nil {
		logger.Errorf(p.c, "GetNodeByPeerId err : %v", err)
	}
	return
}

func (p *Node) AddNode(data *mmysql.Node) (err error) {
	data.CreatedAt = time.Now().Unix()
	data.UpdatedAt = data.CreatedAt
	query := fmt.Sprintf("%d,%d,'%s','%s','%d','%s','%s','%s','%d','%s'", data.CreatedAt, data.UpdatedAt,
		data.PeerId, data.Addr, data.Port, data.IP, data.Loc, data.Colo, data.GasPrice, data.Type)
	err = p.GetOrm().Exec("INSERT INTO `zta_node` (`created_at`," +
		"`updated_at`,`peer_id`,`addr`,`port`,`ip`,`loc`," +
		"`colo`,`gas_price`,`type`) VALUES " +
		"(" + query + ")").Error
	if err != nil {
		logger.Errorf(p.c, "AddNode err : %v", err)
	}
	return
}

func (p *Node) EditNode(data *mmysql.Node) (err error) {
	sql := fmt.Sprintf("UPDATE `zta_node` SET "+
		"`updated_at`=%d,`addr`='%s',`port`=%d,`ip`='%s',`loc`='%s',`colo`='%s',"+
		"`gas_price`=%d,`type`='%s' WHERE `peer_id`='%s'", time.Now().Unix(), data.Addr, data.Port, data.IP, data.Loc, data.Colo,
		data.GasPrice, data.Type, data.PeerId)
	err = p.GetOrm().Exec(sql).Error
	if err != nil {
		logger.Errorf(p.c, "EditNode err : %v", err)
	}
	return
}

func (p *Node) DelNode(id uint64) (err error) {
	orm := p.GetOrm()
	err = orm.Where(fmt.Sprintf("id = %d", id)).Delete(&mmysql.Node{}).Error
	if err != nil {
		logger.Errorf(p.c, "DelNode err : %v", err)
	}
	return
}
