package mysql

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/cloudslit/cloudslit/fullnode/app/v1/user/model/mmysql"
	"github.com/cloudslit/cloudslit/fullnode/pkg/logger"
	"github.com/cloudslit/cloudslit/fullnode/pkg/mysql"
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
	orm := p.GetOrm()
	err = orm.Where(mmysql.User{Email: data.Email}).Attrs(mmysql.User{UUID: uuid.NewString(), AvatarUrl: data.AvatarUrl}).FirstOrCreate(&data).Error
	if err != nil {
		logger.Errorf(p.c, "FirstOrCreateUser err : %v", err)
	}
	return
}
