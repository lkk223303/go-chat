package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/olahol/melody.v1"
)

type Message struct {
	Event   string `json:"event"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

func NewMessage(event, name, content string) *Message {
	return &Message{
		Event:   event,
		Name:    name,
		Content: content,
	}
}

func (m *Message) GetByteMessage() []byte {
	result, _ := json.Marshal(m)
	return result
}

func main() {

	// load config

	// mongo init connect
	client := GetMgoCli()

	collection := client.Database("testing").Collection("numbers")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := collection.InsertOne(ctx, bson.D{{"name", "pi"}, {"value", 3.14159}})
	if err != nil {
		log.Fatal("insert err", err)
	}
	id := res.InsertedID
	log.Println(id)

	// init logger

	// server engine

	engine := gin.Default()
	engine.LoadHTMLGlob("template/html/*")
	engine.Static("/assets", "./template/assets")

	// set routes
	// apiRoute(engine)

	m := melody.New()
	engine.GET("/ws", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		m.Broadcast(msg)
	})

	m.HandleConnect(func(session *melody.Session) {
		id := session.Request.URL.Query().Get("id")
		m.Broadcast(NewMessage("other", id, "加入聊天室").GetByteMessage())
	})

	m.HandleClose(func(session *melody.Session, i int, s string) error {
		id := session.Request.URL.Query().Get("id")
		m.Broadcast(NewMessage("other", id, "離開聊天室").GetByteMessage())
		return nil
	})
	engine.Run(":8080")
}

var mgoCli *mongo.Client

func initMongoDB() {
	var err error

	username := "root"
	pass := "123456"
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	b := &options.Credential{
		Username: username,
		Password: pass,
	}
	clientOptions.Auth = b

	// 连接到MongoDB
	mgoCli, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// 检查连接
	err = mgoCli.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
}
func GetMgoCli() *mongo.Client {
	if mgoCli == nil {
		initMongoDB()
	}
	log.Println("mongo inited")
	return mgoCli
}
