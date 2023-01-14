package database

import (
	"douyin/app/common/config/internal/consts"
	"douyin/app/common/log"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// GetMysqlGormDB 获取mysqlDB实例
func (g *Group) GetMysqlGormDB() (db *gorm.DB, err error) {
	if g.agollo == nil {
		return nil, consts.ErrEmptyConfigClient
	}

	dsn, err := g.GetMysqlDsn()
	if err != nil {
		return nil, err
	}

	logger, err := log.GetGormZapWriter(log.GetGormLoggerConfig())
	if err != nil {
		return nil, err
	}

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger,
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// GetMysqlDsn 返回 mysql DSN
func (g *Group) GetMysqlDsn() (dsn string, err error) {
	v, err := g.agollo.GetViper(consts.MainNamespace)
	if err != nil {
		return "", err
	}

	// 拼接dsn字符串
	dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=Asia%%2FShanghai&tls=true",
		v.GetString("Database.Mysql.Username"),        // 数据库用户名
		v.GetString("Database.Mysql.Password"),        // 数据库密码
		v.GetString("Database.Mysql.Address"),         // 数据库地址
		v.GetString("Database.Mysql.port"),            // 数据库端口
		v.GetString("Database.Mysql.DatabaseName"),    // mysql 的数据库名字
		v.GetString("Database.Mysql.DatabaseCharset"), // mysql 的数据库使用的字符集
	)

	return dsn, nil
}
