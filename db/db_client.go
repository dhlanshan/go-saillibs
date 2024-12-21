package db

import (
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

var session = map[string]*gorm.DB{}

// dBaseConfig 数据库整体配置
type dBaseConfig struct {
	IsOutLog        bool          `mapstructure:"isOutLog"`        // 是否打印SQL日志
	MaxIdle         int           `mapstructure:"maxIdle"`         // 设置连接池中空闲连接的最大数量
	MaxOpen         int           `mapstructure:"maxOpen"`         // 设置打开数据库连接的最大数量
	ConnMaxLifetime time.Duration `mapstructure:"connMaxLifetime"` // 设置了连接可复用的最大时间
}

// setDefault 设置默认参数
func (c *dBaseConfig) setDefault() {
	if c.MaxIdle == 0 {
		c.MaxIdle = 10
	}
	if c.MaxOpen == 0 {
		c.MaxOpen = 100
	}
	if c.ConnMaxLifetime == 0 {
		c.ConnMaxLifetime = time.Hour
	}
}

// DataBaseClient 数据库客户端
type dataBaseClient struct {
	Key string // 配置前缀
}

func (c dataBaseClient) initSession(dbName string) *gorm.DB {
	// 获取对应的数据库配置信息
	dbKey := fmt.Sprintf("%s.%s", c.Key, dbName)
	dbCfg := viper.GetStringMap(dbKey)
	if len(dbCfg) == 0 {
		panic(errors.New(fmt.Sprintf("未找到<%s>数据库配置信息", dbName)))
	}

	var dialect gorm.Dialector
	dialectType := dbCfg["dialect"].(string)
	if dialectType == "Mysql" {
		_mysql := &mysqlClient{Key: dbKey}
		dialect = _mysql.getDialect()
	} else if dialectType == "Postgresql" {
		_pgsql := &postgresqlClient{Key: dbKey}
		dialect = _pgsql.getDialect()
	} else {
		panic(errors.New(fmt.Sprintf("<%s>数据库方言类型错误", dbName)))
	}

	// gorm配置，整体数据库配置
	var dbConfig dBaseConfig
	if err := mapstructure.Decode(viper.GetStringMap(c.Key), &dbConfig); err != nil {
		panic(errors.New("数据库整体配置转换失败。"))
	}
	dbConfig.setDefault()

	cfg := &gorm.Config{
		SkipDefaultTransaction:                   false,
		NamingStrategy:                           schema.NamingStrategy{SingularTable: true},
		DisableAutomaticPing:                     false,
		DisableForeignKeyConstraintWhenMigrating: true, // 禁用创建外键约束
	}
	if dbConfig.IsOutLog {
		cfg.Logger = logger.Default.LogMode(logger.Info) // 打印sql语句
	}
	db, err := gorm.Open(dialect, cfg)
	if err != nil {
		panic(errors.New(fmt.Sprintf("<%s>数据库连接错误", dbName)))
	}

	// 设置连接池
	sqlDB, err := db.DB()
	if err != nil {
		panic(errors.New(fmt.Sprintf("连接数据库失败: %s", err.Error())))
	}
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdle)            // 用于设置连接池中空闲连接的最大数量
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpen)            // 设置打开数据库连接的最大数量
	sqlDB.SetConnMaxLifetime(dbConfig.ConnMaxLifetime) // 设置了连接可复用的最大时间

	session[dbName] = db

	return db
}

// GetDbClient 获取数据库客户端
func GetDbClient(dbName string) *gorm.DB {
	dbClient, ok := session[dbName]
	if !ok {
		dbClient = dataBaseClient{"dbClient"}.initSession(dbName)
	}

	return dbClient
}
