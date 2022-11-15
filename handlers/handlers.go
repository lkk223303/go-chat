package handlers

import (
	"chatty/models"
	"chatty/repository"
	"chatty/repository/dbrepo"
	"chatty/utils"
	"strings"

	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/line/line-bot-sdk-go/linebot"
	"go.mongodb.org/mongo-driver/mongo"
)

var Repo *Repository

type Repository struct {
	DB  repository.DatabaseRepo
	Bot *linebot.Client
	Rds *redis.Client
}

// for main.go to init
func NewRepo(client *mongo.Client, bot *linebot.Client, r *redis.Client) *Repository {
	return &Repository{
		DB:  dbrepo.NewMongoRepo(client),
		Bot: bot,
		Rds: r,
	}
}

// init handler repo
func NewHandler(r *Repository) {
	Repo = r
}

func (m *Repository) Home(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"Welcome": "this is a welcome message"})
}

func (m *Repository) TestMongo(c *gin.Context) {
	m.DB.Testing()
	c.JSON(http.StatusOK, gin.H{"db test": "ok"})
}

// Line incoming message
func (m *Repository) CallBack(c *gin.Context) {
	var msgEvent models.EventMessage
	events, err := m.Bot.ParseRequest(c.Request)
	if err != nil {
		log.Println("parse request error ", err)
		if err == linebot.ErrInvalidSignature {
			c.Status(400)
		} else {
			c.Status(500)
		}
		return
	}
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch msg := event.Message.(type) {
			case *linebot.TextMessage:

				if strings.Contains(msg.Text, "你家") {
					err = utils.SendLocationMessage(event.Source.UserID)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": "reply message error", "content": err})
					}
				} else {
					err = utils.ReplyLineMessage(event.ReplyToken, "哈哈，我先去洗澡")
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": "reply message error", "content": err})
					}

					msgEvent.UserID = event.Source.UserID
					msgEvent.Message = msg.Text
					msgEvent.TimeStamp = event.Timestamp
					err = utils.RedisPushMessage(msgEvent)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": "redis push error", "content": err})
						return
					}
				}
			}
		}
	}
}

// Create a API send message back to line
func (m *Repository) SendMessage(c *gin.Context) {
	var event models.EventMessage
	err := c.ShouldBindJSON(&event)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client parameter error", "content": err})
		return
	}
	err = utils.SendLineMessage(event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "sending line message error", "content": err})
		return
	}
	log.Println("message sent to line: ", event.Message)
	c.JSON(http.StatusOK, gin.H{"success": "ok", "content": "message successfully sent"})
}

// Query message list of the user from MongoDB
func (m *Repository) Messages(c *gin.Context) {
	userId := c.Param("userId")

	list, err := utils.GetMessagesFromUser(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": list})
}
