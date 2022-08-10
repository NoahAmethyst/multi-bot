package db_util

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

const (
	baseLink = "{db_user}:{db_password}@tcp({db_host}:{db_port})/{db_name}?charset=utf8mb4&parseTime=true"
)

var db *gorm.DB = nil

func ConnectDb(dbHost string, dbPort string, dbName string, dbUser string, dbPassword string) error {
	currentLink := buildLinkUrl(dbHost, dbPort, dbName, dbUser, dbPassword)
	var err error
	mysqlConfig := mysql.New(mysql.Config{
		DSN:                       currentLink,
		SkipInitializeWithVersion: false,
	})
	db, err = gorm.Open(mysqlConfig, &gorm.Config{})
	if err != nil {
		log.Error().Msgf("failed to connect mysql cause:%s", err.Error())
		return err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(16)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(160)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)
	// 设置全局表名禁用复数
	log.Info().Fields(map[string]interface{}{
		"action": "connect to db",
		"host":   dbHost,
		"port":   dbPort,
		"dbName": dbName,
	}).Send()
	return nil
}

func buildLinkUrl(dbHost string, dbPort string, dbName string, dbUser string, dbPassword string) string {
	currentLink := baseLink
	currentLink = strings.Replace(currentLink, "{db_user}", dbUser, -1)
	currentLink = strings.Replace(currentLink, "{db_password}", dbPassword, -1)
	currentLink = strings.Replace(currentLink, "{db_host}", dbHost, -1)
	currentLink = strings.Replace(currentLink, "{db_port}", dbPort, -1)
	currentLink = strings.Replace(currentLink, "{db_name}", dbName, -1)
	return currentLink
}

func GetDb() *gorm.DB {
	return db
}

type JSON json.RawMessage

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (j *JSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	result := json.RawMessage{}
	err := json.Unmarshal(bytes, &result)
	*j = JSON(result)
	return err
}

// 实现 driver.Valuer 接口，Value 返回 json value
func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.RawMessage(j).MarshalJSON()
}
