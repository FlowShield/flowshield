package mysql

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/flowshield/flowshield/verifier/pkg/confer"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var gormDB *gorm.DB

type DaoMysql struct {
	db *gorm.DB
}

func NewDaoMysql() *DaoMysql {
	return &DaoMysql{
		db: gormDB,
	}
}

func (dao *DaoMysql) Orm(c *gin.Context) *gorm.DB {
	if c == nil {
		return dao.db
	}
	return dao.db.WithContext(c.Request.Context())
}

func Init(cfg *confer.Mysql) error {
	dbConfig := cfg.Write
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4,utf8&parseTime=True&loc=Local", dbConfig.User,
		dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DBName)
	config := &gorm.Config{
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   cfg.Prefix, // 表名前缀，`User`表为`t_users`
			SingularTable: true,       // 使用单数表名，启用该选项后，`User` 表将是`user`
		},
	}
	if confer.ConfigEnvIsDev() {
		newLogger := glogger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			glogger.Config{
				SlowThreshold: time.Second,  // 慢 SQL 阈值
				LogLevel:      glogger.Info, // Log level
				Colorful:      true,         // 禁用彩色打印
			},
		)
		config.Logger = newLogger
	} else {
		config.Logger = glogger.Default.LogMode(glogger.Silent)
	}
	db, err := gorm.Open(mysql.Open(dsn), config)
	if err != nil {
		return err
	}
	sqlDB, err := db.DB()
	if err != nil {
		err = errors.New(fmt.Sprintf("initMysql,error: %v", err))
		return err
	}
	sqlDB.SetMaxIdleConns(cfg.Pool.PoolMinCap)                                      // 空闲链接
	sqlDB.SetMaxOpenConns(cfg.Pool.PoolMaxCap)                                      // 最大链接
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.Pool.PoolIdleTimeout) * time.Second) // 最大空闲时间
	gormDB = db
	return nil
}
