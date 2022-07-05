package mysql

import (
	"errors"
	"fmt"

	"github.com/cloudslit/cloudslit/fullnode/app/v1/user/model/mmysql"
	"github.com/cloudslit/cloudslit/fullnode/pkg/logger"
	"github.com/cloudslit/cloudslit/fullnode/pkg/mysql"
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

func (p *User) FirstOrCreateUser(data *mmysql.User) (err error) {
	var info *mmysql.User
	orm := p.GetOrm()
	err = orm.Table(p.TableName).Where(fmt.Sprintf("email = '%s'", data.Email)).First(&info).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 执行插入
		sql := orm.ToSQL(func(tx *gorm.DB) *gorm.DB {
			return tx.Table(p.TableName).Create(&data)
		})
		err = orm.Exec(sql).Error
		if err != nil {
			logger.Errorf(p.c, "FirstOrCreateUser err : %v", err)
		}
		return
	}
	//err = orm.Where(mmysql.User{Email: data.Email}).Attrs(mmysql.User{UUID: uuid.NewString(), AvatarUrl: data.AvatarUrl}).FirstOrCreate(&data).Error
	return
}
