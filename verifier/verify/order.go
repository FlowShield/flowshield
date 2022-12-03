package verify

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/flowshield/flowshield/verifier/pkg/logger"

	"gorm.io/gorm"

	"github.com/flowshield/flowshield/verifier/pkg/mysql"
)

type OrderMysql struct {
	ID          uint   `gorm:"primarykey"`
	CreatedAt   int64  `gorm:"autoCreateTime"`
	UpdatedAt   int64  `gorm:"autoUpdateTime"`
	UserUUID    string `json:"user_uuid" gorm:"user_uuid"`
	Name        string `json:"name"`
	PeerID      string `json:"peer_id"`
	UUID        string `json:"uuid" gorm:"column:uuid"`
	Port        int    `json:"port"`
	Duration    uint   `json:"duration"` // 使用时间：小时
	Price       uint   `json:"price"`    // 金额
	ResourceCid string `json:"resource_cid"`
	ServerCid   string `json:"server_cid"`
	ClientCid   string `json:"client_cid"`
	Status      uint   `json:"status"` // 0:待支付，1:已支付,待回调，2:已完成
	NodeIP      string `json:"node_ip"`
}

func (OrderMysql) TableName() string {
	return "zta_client"
}

type Order struct {
	PeerID    string    `json:"peer_id"`
	Port      int       `json:"port"`
	OrderID   string    `json:"order_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Healthy   *Healthy  `json:"healthy"`
}

type Healthy struct {
	Health bool  `json:"health"`
	Err    error `json:"err"`
}

func (o *Order) CheckHealthy(ip string) {
	if o.Healthy == nil {
		o.Healthy = &Healthy{Health: false, Err: nil}
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s%d", ip, o.Port))
	if err != nil {
		o.Healthy.Health = false
		o.Healthy.Err = err
	} else {
		if conn == nil {
			o.Healthy.Health = false
			o.Healthy.Err = errors.New("conn is nil")
		} else {
			o.Healthy.Health = true
			o.Healthy.Err = nil
		}
	}
}

func orders() (orders []*OrderMysql, err error) {
	orm := mysql.NewDaoMysql().Orm(nil)
	err = orm.Select("zta_client.*,zta_node.ip as node_ip").
		Joins("left join zta_node on zta_client.peer_id = zta_node.peer_id").
		Where(fmt.Sprintf("zta_client.status = %d", 2)).
		Where(fmt.Sprintf("zta_client.updated_at + 60*60*duration >= %d", time.Now().Unix())).
		Order("zta_client.created_at desc").
		Find(&orders).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	if err != nil {
		logger.Errorf(nil, "Orders err : %v", err)
	}
	return
}
