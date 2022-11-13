package main

import (
	"chatty/handlers"

	"github.com/gin-gonic/gin"
)

func apiRoute(router *gin.Engine) {
	router.GET("/", handlers.Repo.Home)
	router.GET("/testmgo", handlers.Repo.TestMongo)
	router.GET("/messages/:userId", handlers.Repo.Messages)
	router.POST("/sendmessage", handlers.Repo.SendMessage)
	router.POST("/callback", handlers.Repo.CallBack)
}
