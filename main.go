package main

import (
	"chatty/bot"
	"chatty/utils"

	"chatty/driver"
	"chatty/handlers"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile("./config.yaml")
	viper.SetDefault("application.port", 8088)

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("read config error: ", err)
	}
	log.Println("Config initialized")
}
func main() {

	// init line
	lineBot := bot.LineConfig()

	// mongo/redis init connect
	client := driver.GetMongoClient()
	rdsCli := driver.GetrRedisClient()
	// init handler
	repo := handlers.NewRepo(client, lineBot, rdsCli)
	handlers.NewHandler(repo)
	util := utils.NewUtil(client, rdsCli)
	go util.AddMessage()

	// server engine
	engine := gin.Default()
	engine.LoadHTMLGlob("template/html/*")
	engine.Static("/assets", "./template/assets")
	// os.Setenv("GIN_MODE", "release")
	// set routes
	apiRoute(engine)

	engine.Run(viper.GetString("application.port"))
}
