package db

import (
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type pgsqlConfig struct {
	User    string `mapstructure:"user"`    // 用户名
	Passwd  string `mapstructure:"passwd"`  // 密码
	Addr    string `mapstructure:"addr"`    // 数据库地址
	DbName  string `mapstructure:"dbName"`  // 数据库名
	Port    string `mapstructure:"port"`    // 端口
	Dialect string `mapstructure:"dialect"` // 引擎类型Postgresql
}

func (c *pgsqlConfig) checkCfg() {
	nullAttrList := make([]string, 0)

	if c.User == "" { // TODO 校验参数为空的，需要抛出异常
		nullAttrList = append(nullAttrList, "user")
	}
}

// postgresqlClient Postgresql客户端
type postgresqlClient struct {
	Key string // Mysql配置前置key
}

func (c *postgresqlClient) getDialect() gorm.Dialector {
	// 配置
	var config pgsqlConfig
	if err := mapstructure.Decode(viper.GetString(c.Key), &config); err != nil {
		panic(errors.New("postgresql配置转换失败。"))
	}
	config.checkCfg()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		config.Addr, config.User, config.Passwd, config.DbName, config.Port,
	)
	return postgres.Open(dsn)
}
