package mysql

import (
	"errors"
	"fmt"

	"github.com/flowshield/flowshield/fullnode/app/v1/user/model/mmysql"
	"github.com/flowshield/flowshield/fullnode/pkg/logger"
	"github.com/flowshield/flowshield/fullnode/pkg/mysql"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type User struct {
	c *gin.Context
	mysql.DaoMysql
}

func NewUser(c *gin.Context) *User {
	return &User{
		DaoMysql: mysql.DaoMysql{TableName: "zta_user"},
		c:        c,
	}
}

func (p *User) FirstOrCreateUser(data *mmysql.User) (info *mmysql.User, err error) {
	orm := p.GetOrm()
	err = orm.Table(p.TableName).Where(fmt.Sprintf("email = '%s'", data.Email)).First(&info).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
		// 执行插入
		sql := orm.ToSQL(func(tx *gorm.DB) *gorm.DB {
			return tx.Table(p.TableName).Create(&data)
		})
		err = orm.Exec(sql).Error
		if err != nil {
			logger.Errorf(p.c, "FirstOrCreateUser err : %v", err)
			return
		}
		info = &mmysql.User{
			Email:     data.Email,
			AvatarUrl: data.AvatarUrl,
			UUID:      data.UUID,
		}
		return
	}
	if err != nil {
		logger.Errorf(p.c, "FirstOrCreateUser err : %v", err)
		return
	}
	//err = orm.Where(mmysql.User{Email: data.Email}).Attrs(mmysql.User{UUID: uuid.NewString(), AvatarUrl: data.AvatarUrl}).FirstOrCreate(&data).Error
	return
}

func (p *User) GetUser(uuid string) (user *mmysql.User, err error) {
	orm := p.GetOrm()
	err = orm.Table(p.TableName).Where(fmt.Sprintf("uuid = '%s'", uuid)).First(&user).Error
	if err != nil {
		logger.Errorf(p.c, "GetUser err : %v", err)
		return
	}
	return
}

func (p *User) GetUserByStatus(status int) (list []mmysql.User, err error) {
	orm := p.GetOrm().DB
	query := orm.Table(p.TableName).Where(fmt.Sprintf("status = %d", status))
	err = query.Find(&list).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	if err != nil {
		logger.Errorf(p.c, "GetUserByStatus err : %v", err)
	}
	return
}

func (p *User) UpdateUser(data *mmysql.User) (err error) {
	orm := p.GetOrm()
	sql := orm.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Table(p.TableName).Save(&data)
	})
	err = orm.Exec(sql).Error
	if err != nil {
		logger.Errorf(p.c, "UpdateUser err : %v", err)
	}
	return
}
