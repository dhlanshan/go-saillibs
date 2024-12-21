package db

import (
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type mysqlConfig struct {
	User     string `mapstructure:"user"`     // 用户名
	Passwd   string `mapstructure:"passwd"`   // 密码
	Addr     string `mapstructure:"addr"`     // 数据库地址
	DbName   string `mapstructure:"dbName"`   // 数据库名
	ChartSet string `mapstructure:"chartSet"` // 编码格式
	TimeOut  string `mapstructure:"timeOut"`  // 连接超时时间,该值必须是带有单位后缀（“ms”、“s”、“m”、“h”）的十进制数，例如“30s”、“0.5m”或“1m30s”
	Dialect  string `mapstructure:"dialect"`  // 引擎类型Mysql
}

func (c *mysqlConfig) checkCfg() {
	nullAttrList := make([]string, 0)

	if c.User == "" { // TODO 校验参数为空的，需要抛出异常
		nullAttrList = append(nullAttrList, "user")
	}
	if c.TimeOut == "" {
		c.TimeOut = "300s"
	}
}

// mysqlClient Mysql客户端
type mysqlClient struct {
	Key string // Mysql配置前置key
}

func (c *mysqlClient) getDialect() gorm.Dialector {
	// 配置
	var cfg mysqlConfig
	if err := mapstructure.Decode(viper.GetStringMap(c.Key), &cfg); err != nil {
		panic(errors.New("mysql配置转换失败。"))
	}
	cfg.checkCfg()

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=true&loc=Local&timeout=%s",
		cfg.User, cfg.Passwd, cfg.Addr, cfg.DbName, cfg.ChartSet, cfg.TimeOut,
	)
	return mysql.Open(dsn)
}
