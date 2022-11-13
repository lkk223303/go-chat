package handlers

import (
	"chatty/models"
	"chatty/repository"
	"chatty/repository/dbrepo"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
	"go.mongodb.org/mongo-driver/mongo"
)

var Repo *Repository

type Repository struct {
	DB  repository.DatabaseRepo
	Bot *linebot.Client
}

// for main.go to init
func NewRepo(client *mongo.Client, bot *linebot.Client) *Repository {
	return &Repository{
		DB:  dbrepo.NewMongoRepo(client),
		Bot: bot,
	}
}

// init handler repo
func NewHandler(r *Repository) {
	Repo = r
}

func (m *Repository) Home(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func (m *Repository) TestMongo(c *gin.Context) {
	m.DB.Testing()

	c.JSON(http.StatusOK, gin.H{"db test": "ok"})
}

// line incoming message
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
				if _, err = m.Bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(msg.Text+" 笑死")).Do(); err != nil {
					log.Println("resp message error ", err)
				}

				msgEvent.UserID = event.Source.UserID
				msgEvent.Message = msg.Text
				msgEvent.TimeStamp = event.Timestamp
				m.DB.InsertMessage(msgEvent)
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
	sendTo := event.UserID
	msg := event.Message
	if _, err := m.Bot.PushMessage(sendTo, linebot.NewTextMessage(msg)).Do(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "sending line message error", "content": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": "ok", "content": "message successfully pushed"})
}

// Query message list of the user from MongoDB
func (m *Repository) Messages(c *gin.Context) {
	userId := c.Param("userId")

	list, err := m.DB.GetMessagesFromUser(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	log.Println(list)
	c.JSON(http.StatusOK, gin.H{"success": list})
}
