package main

import (
	"chatty/handlers"

	"github.com/gin-gonic/gin"
)

func apiRoute(router *gin.Engine) {
	router.GET("/", handlers.Repo.Home)

}
