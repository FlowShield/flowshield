package mysql

import (
	"errors"
	"fmt"

	"github.com/flowshield/flowshield/fullnode/pkg/util"

	"github.com/flowshield/flowshield/fullnode/app/v1/access/model/mmysql"
	"github.com/flowshield/flowshield/fullnode/app/v1/access/model/mparam"
	"github.com/flowshield/flowshield/fullnode/pkg/logger"
	"github.com/flowshield/flowshield/fullnode/pkg/mysql"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

type Resource struct {
	c *gin.Context
	mysql.DaoMysql
}

func NewResource(c *gin.Context) *Resource {
	return &Resource{
		DaoMysql: mysql.DaoMysql{TableName: "zta_resource"},
		c:        c,
	}
}

func (p *Resource) ResourceList(param mparam.ResourceList) (
	total int64, list []mmysql.Resource, err error) {
	orm := p.GetOrm().DB
	query := orm.Table(p.TableName)
	if len(param.Name) > 0 {
		query = query.Where(fmt.Sprintf("name like '%%%s%%'", param.Name))
	}
	if len(param.Type) > 0 {
		query = query.Where(fmt.Sprintf("`type` = %s", param.Type))
	}
	//if user := util.User(p.c); user != nil {
	//	query = query.Where(fmt.Sprintf("user_uuid = '%s'", user.UUID))
	//}
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
		logger.Errorf(p.c, "ResourceList err : %v", err)
	}
	return
}

func (p *Resource) GetResourceByIDSli(ids []string) (list []mmysql.Resource, err error) {
	orm := p.GetOrm()
	query := orm.Table(p.TableName)
	var idStr string
	for _, value := range ids {
		idStr += fmt.Sprintf("'%s'", value) + ","
	}
	idStr = idStr[:len(idStr)-1]
	if len(idStr) > 0 {
		query = query.Where(fmt.Sprintf("id in (%s)", idStr))
	}
	//if user := util.User(p.c); user != nil {
	//	query = query.Where(fmt.Sprintf("user_uuid = '%s'", user.UUID))
	//}
	err = query.Find(&list).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	if err != nil {
		logger.Errorf(p.c, "GetResourceByIDSli err : %v", err)
	}
	return
}

func (p *Resource) GetResourceByID(id uint64) (info *mmysql.Resource, err error) {
	orm := p.GetOrm()
	query := orm.Table(p.TableName).Where(fmt.Sprintf("id = %d", id))
	//if user := util.User(p.c); user != nil {
	//	query = query.Where(fmt.Sprintf("user_uuid = '%s'", user.UUID))
	//}
	err = query.First(&info).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	if err != nil {
		logger.Errorf(p.c, "GetResourceById err : %v", err)
	}
	return
}

func (p *Resource) GetResourceByUUID(uuid string) (info *mmysql.Resource, err error) {
	orm := p.GetOrm()
	query := orm.Table(p.TableName).Where(fmt.Sprintf("uuid = '%s'", uuid))
	//if user := util.User(p.c); user != nil {
	//	query = query.Where(fmt.Sprintf("user_uuid = '%s'", user.UUID))
	//}
	err = query.First(&info).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	if err != nil {
		logger.Errorf(p.c, "GetResourceByUUID err : %v", err)
	}
	return
}

func (p *Resource) GetResourceByCID(cid string) (info *mmysql.Resource, err error) {
	orm := p.GetOrm()
	query := orm.Table(p.TableName).Where(fmt.Sprintf("cid = '%s'", cid))
	//if user := util.User(p.c); user != nil {
	//	query = query.Where(fmt.Sprintf("user_uuid = '%s'", user.UUID))
	//}
	err = query.First(&info).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	if err != nil {
		logger.Errorf(p.c, "GetResourceByCID err : %v", err)
	}
	return
}

func (p *Resource) AddResource(data *mmysql.Resource) (err error) {
	orm := p.GetOrm()
	sql := orm.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Table(p.TableName).Create(&data)
	})
	err = orm.Exec(sql).Error
	if err != nil {
		logger.Errorf(p.c, "AddResource err : %v", err)
	}
	return
}

func (p *Resource) EditResource(data *mmysql.Resource) (err error) {
	if user := util.User(p.c); user != nil {
		data.UserUUID = user.UUID
	}
	orm := p.GetOrm()
	sql := orm.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Table(p.TableName).Save(&data)
	})
	err = orm.Exec(sql).Error
	if err != nil {
		logger.Errorf(p.c, "EditResource err : %v", err)
	}
	return
}

func (p *Resource) DelResource(uuid string) (err error) {
	orm := p.GetOrm()
	query := orm.Table(p.TableName).Where(fmt.Sprintf("uuid = '%s'", uuid))
	if user := util.User(p.c); user != nil {
		query = query.Where(fmt.Sprintf("user_uuid = '%s'", user.UUID))
	}
	err = query.Delete(&mmysql.Resource{}).Error
	if err != nil {
		logger.Errorf(p.c, "DelResource err : %v", err)
	}
	return
}
