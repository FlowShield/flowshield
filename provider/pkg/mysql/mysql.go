package mysql

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cloudslit/cloudslit/provider/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type DaoMysql struct {
	TableName string
}

func NewDaoMysql() *DaoMysql {
	return &DaoMysql{}
}

type Connection struct {
	*gorm.DB
	IsRead bool
}

var (
	mysqlReadPool  Connection
	mysqlWritePool Connection
)

// Init 初始化mysql连接池
func Init() (err error) {
	// initMysqlPool(true) //初始化从库，多个
	cfg := config.C.Mysql
	err = initMysqlPool(&cfg, false) // 初始化写库，一个
	return
}

func initMysqlPool(cfg *config.Mysql, isRead bool) (err error) {
	if isRead {
		mysqlReadPool.DB, err = initDb(cfg, isRead)
		mysqlReadPool.IsRead = isRead
	} else {
		mysqlWritePool.DB, err = initDb(cfg, isRead)
		mysqlWritePool.IsRead = isRead
	}
	if err != nil {
		err = errors.New(fmt.Sprintf("initMysqlPool isread: %v ,error: %v", isRead, err))
		return
	}
	if isRead {
		sqlDB, err := mysqlReadPool.DB.DB()
		if err != nil {
			err = errors.New(fmt.Sprintf("initMysqlPool isread:%v ,error: %v", isRead, err))
			return err
		}
		sqlDB.SetMaxIdleConns(cfg.PoolMinCap)                                      // 空闲链接
		sqlDB.SetMaxOpenConns(cfg.PoolMaxCap)                                      // 最大链接
		sqlDB.SetConnMaxLifetime(time.Duration(cfg.PoolIdleTimeout) * time.Second) // 最大空闲时间
	} else {
		sqlDB, err := mysqlWritePool.DB.DB()
		if err != nil {
			err = errors.New(fmt.Sprintf("initMysqlPool isread:%v ,error: %v", isRead, err))
			return err
		}
		sqlDB.SetMaxIdleConns(cfg.PoolMinCap)                                      // 空闲链接
		sqlDB.SetMaxOpenConns(cfg.PoolMaxCap)                                      // 最大链接
		sqlDB.SetConnMaxLifetime(time.Duration(cfg.PoolIdleTimeout) * time.Second) // 最大空闲时间
	}
	return
}

func initDb(cfg *config.Mysql, isRead bool) (resultDb *gorm.DB, err error) {
	// 判断配置可用性
	if cfg.Host == "" || cfg.DBName == "" {
		err = errors.New("dbConfig is null")
		return
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.Parameters)
	gormCfg := &gorm.Config{
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   cfg.Prefix, // 表名前缀，`User`表为`t_users`
			SingularTable: true,       // 使用单数表名，启用该选项后，`User` 表将是`user`
		},
	}
	if config.C.IsDebugMode() {
		newLogger := glogger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			glogger.Config{
				SlowThreshold: time.Second,  // 慢 SQL 阈值
				LogLevel:      glogger.Info, // Log level
				Colorful:      true,         // 禁用彩色打印
			},
		)
		gormCfg.Logger = newLogger
	}
	resultDb, err = gorm.Open(mysql.Open(dsn), gormCfg)
	return resultDb, err
}

func getMysqlPoolConnection(isRead bool) (conn Connection) {
	if isRead {
		conn = mysqlReadPool
	} else {
		conn = mysqlWritePool
	}
	return
}

func (p *DaoMysql) GetReadOrm() Connection {
	return p.getOrm(true)
}

func (p *DaoMysql) GetWriteOrm() Connection {
	return p.getOrm(false)
}

func (p *DaoMysql) GetOrm() Connection {
	return p.getOrm(false)
}

func (p *DaoMysql) getOrm(isRead bool) Connection {
	return getMysqlPoolConnection(isRead)
}
