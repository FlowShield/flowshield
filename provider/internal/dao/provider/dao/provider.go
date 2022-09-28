package dao

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/cloudslit/cloudslit/provider/internal/dao/provider/model"
	"github.com/cloudslit/cloudslit/provider/pkg/mysql"
)

type Provider struct {
	mysql.DaoMysql
}

func NewProvider() *Provider {
	return &Provider{
		DaoMysql: mysql.DaoMysql{TableName: "pr_provider"},
	}
}

func (p *Provider) ListProvider(param *model.Provider) (model.Providers, error) {
	query := p.GetOrm().DB
	if param.Uuid != "" {
		query = query.Where(fmt.Sprintf("uuid = '%s'", param.Uuid))
	}
	if param.Wallet != "" {
		query = query.Where(fmt.Sprintf("wallet = '%s'", param.Wallet))
	}
	if param.Port > 0 {
		query = query.Where(fmt.Sprintf("port = %d", param.Port))
	}
	if param.ServerCid != "" {
		query = query.Where(fmt.Sprintf("server_cid = '%s'", param.ServerCid))
	}
	query = query.Where(fmt.Sprintf("expired_at = 0 or expired_at >= %d", time.Now().Unix()))
	var list model.Providers
	err := query.Order("updated_at desc").Find(&list).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return list, nil
}

func (p *Provider) GetProviderByUuid(uuid string) (info *model.Provider, err error) {
	orm := p.GetOrm()
	err = orm.Table(p.TableName).Where(fmt.Sprintf("uuid = '%s'", uuid)).First(&info).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return
}

func (p *Provider) AddProvider(data *model.Provider) (err error) {
	orm := p.GetOrm()
	sql := orm.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Table(p.TableName).Create(&data)
	})
	err = orm.Exec(sql).Error
	return
}

func (p *Provider) EditProvider(data *model.Provider) (err error) {
	orm := p.GetOrm()
	sql := orm.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Table(p.TableName).Save(&data)
	})
	err = orm.Exec(sql).Error
	return
}

func (p *Provider) DelProvider(id uint64) (err error) {
	orm := p.GetOrm()
	err = orm.Where(fmt.Sprintf("id = %d", id)).Delete(&model.Provider{}).Error
	return
}
