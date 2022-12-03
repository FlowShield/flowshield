package mysql

import (
	"errors"
	"fmt"
	"time"

	"github.com/flowshield/flowshield/fullnode/app/v1/node/model/mparam"

	"github.com/flowshield/flowshield/fullnode/app/v1/node/model/mmysql"
	"github.com/flowshield/flowshield/fullnode/pkg/logger"
	"github.com/flowshield/flowshield/fullnode/pkg/mysql"
	"github.com/gin-gonic/gin"
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
	query := p.GetOrm().DB.Where(fmt.Sprintf("updated_at >= %d", time.Now().Add(-time.Minute*10).Unix()))
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
	if param.Price > 0 {
		query = query.Where(fmt.Sprintf("price = %d", param.Price))
	}
	if len(param.Type) > 0 {
		query = query.Where(fmt.Sprintf("`type` = '%s'", param.Type))
	}
	err = query.Model(&list).Count(&total).Error
	if total > 0 {
		offset := param.GetOffset()
		err = query.Limit(param.LimitNum).Offset(offset).
			Order("`type` asc").
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
	orm := p.GetOrm()
	sql := orm.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Table(p.TableName).Create(&data)
	})
	err = orm.Exec(sql).Error
	if err != nil {
		logger.Errorf(p.c, "AddNode err : %v", err)
	}
	return
}

func (p *Node) EditNode(data *mmysql.Node) (err error) {
	orm := p.GetOrm()
	sql := orm.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Table(p.TableName).Save(&data)
	})
	err = orm.Exec(sql).Error
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
