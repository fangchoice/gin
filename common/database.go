package common

import (
	"fmt"
	"gin-learn/model"

	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

// 声明DB
var DB *gorm.DB

func InitDB() *gorm.DB {
	// viper config
	driverName := viper.GetString("datasource.driverName")
	host := viper.GetString("datasource.host")
	port := viper.GetString("datasource.port")
	database := viper.GetString("datasource.database")
	username := viper.GetString("datasource.username")
	password := viper.GetString("datasource.password")
	charset := viper.GetString("datasource.charset")

	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		username, password, host, port, database, charset)

	db, err := gorm.Open(driverName, args)
	if err != nil {
		panic("failed to conn database err: " + err.Error())
	}

	// 迁移，将会在数据库中生成相应的表
	db.AutoMigrate(&model.User{})

	// 给DB赋值
	DB = db
	return db
}

func GetDB() *gorm.DB {
	return DB
}
