package handlers

import (
	"chatty/repository"
	"chatty/repository/dbrepo"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var Repo *Repository

type Repository struct {
	DB repository.DatabaseRepo
}

// for main.go to init
func NewRepo(client *mongo.Client) *Repository {
	return &Repository{
		DB: dbrepo.NewMongoRepo(client),
	}
}

// init handler repo
func NewHandler(r *Repository) {
	Repo = r
}

func (m *Repository) Home(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}
