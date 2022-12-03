package mysql

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/flowshield/flowshield/fullnode/pkg/confer"
	migrate "github.com/rubenv/sql-migrate"
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
func Init(cfg *confer.Mysql) (err error) {
	// initMysqlPool(true) //初始化从库，多个
	err = initMysqlPool(cfg, false) // 初始化写库，一个
	return
}

func initMysqlPool(cfg *confer.Mysql, isRead bool) (err error) {
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
		sqlDB.SetMaxIdleConns(cfg.Pool.PoolMinCap)                       // 空闲链接
		sqlDB.SetMaxOpenConns(cfg.Pool.PoolMaxCap)                       // 最大链接
		sqlDB.SetConnMaxLifetime(cfg.Pool.PoolIdleTimeout * time.Second) // 最大空闲时间
	} else {
		sqlDB, err := mysqlWritePool.DB.DB()
		if err != nil {
			err = errors.New(fmt.Sprintf("initMysqlPool isread:%v ,error: %v", isRead, err))
			return err
		}
		sqlDB.SetMaxIdleConns(cfg.Pool.PoolMinCap)                       // 空闲链接
		sqlDB.SetMaxOpenConns(cfg.Pool.PoolMaxCap)                       // 最大链接
		sqlDB.SetConnMaxLifetime(cfg.Pool.PoolIdleTimeout * time.Second) // 最大空闲时间
	}
	return
}

func initDb(cfg *confer.Mysql, isRead bool) (resultDb *gorm.DB, err error) {
	var dbConfig confer.DBBase
	if isRead && len(cfg.Reads) > 0 {
		rand.Seed(time.Now().UnixNano())
		dbConfig = cfg.Reads[rand.Intn(len(cfg.Reads)-1)]
	} else {
		dbConfig = cfg.Write
	}
	// 判断配置可用性
	if dbConfig.Host == "" || dbConfig.DBName == "" {
		err = errors.New("dbConfig is null")
		return
	}
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
	resultDb, err = gorm.Open(mysql.Open(dsn), config)
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

func SqlMigrate() (err error) {
	// docker 环境下，根据docker file 配置，sql文件在统计db目录下
	migrations := &migrate.FileMigrationSource{
		Dir: "./db",
	}
	Orm := NewDaoMysql().GetOrm()
	sqlDB, err := Orm.DB.DB()
	if err != nil {
		//logger.LogErrorw(nil, logger.LogNameMysql, "sqlMigrate Orm.DB.DB() err ", err)
		return
	}
	_, err = migrate.Exec(sqlDB, "mysql", migrations, migrate.Up)
	return
}
