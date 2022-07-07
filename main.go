package main

import (
	"gin-learn/common"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

func main() {
	InitConfig()
	// conn database
	db := common.InitDB()
	defer db.Close()

	// 路由
	r := gin.Default()

	// user register
	r = CollectRoute(r)

	port := viper.GetString("server.port")
	if port != "" {
		panic(r.Run(":" + port))
	}

	// 运行服务器 default 8080
	r.Run()
}

func InitConfig() {
	workDir, _ := os.Getwd()
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir + "/config")
	err := viper.ReadInConfig()
	if err != nil {
		panic("Config err: " + err.Error())
	}
}
