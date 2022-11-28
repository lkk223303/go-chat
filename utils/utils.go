package utils

import (
	"chatty/models"
	"chatty/repository"
	"chatty/repository/dbrepo"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/line/line-bot-sdk-go/linebot"
	"go.mongodb.org/mongo-driver/mongo"
)

var uTool *Util

type Util struct {
	DB  repository.DatabaseRepo
	Bot *linebot.Client
	Rds *redis.Client
}

func InitUtil(client *mongo.Client, bot *linebot.Client, r *redis.Client) *Util {
	go consumeMessage()

	return &Util{
		DB:  dbrepo.NewMongoRepo(client),
		Bot: bot,
		Rds: r,
	}
}

func NewUtil(u *Util) {
	uTool = u
}

var eventList []models.EventMessage
var batch int = 10

// consumeMessage consumes messages from redis to mongo
func consumeMessage() {
	t := time.NewTicker(10 * time.Second)
	for {

		var msgEvent models.EventMessage
		// 5秒超時
		value, err := uTool.Rds.BRPop(5*time.Second, models.MsgCache).Result()
		if err == redis.Nil {
			// 查詢不到資料
			time.Sleep(1 * time.Second)
			continue
		}
		if err != nil {
			time.Sleep(1 * time.Second)
			log.Println("WARN: redis pop error: ", err)
			continue
		}
		msgEvents := strings.Split(value[1], "\t")
		msgEvent.UserID = msgEvents[0]
		msgEvent.Message = msgEvents[1]
		msgEvent.TimeStamp, err = time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", msgEvents[2])
		if err != nil {
			log.Println("parsing time format error, ", err)
		}
		log.Println("Income msgEvent: ", msgEvent)
		eventList = append(eventList, msgEvent)

		// Insert messages when hitting batch
		if len(eventList) == batch {

			err = uTool.DB.InsertMessages(eventList)

			if err != nil {
				log.Fatal("Add message error: ", err)
			}
			eventList = nil
		}

		// Or insert messages when time's up
		go func() {
			<-t.C
			if len(eventList) > 0 && len(eventList) < batch && eventList != nil {

				err = uTool.DB.InsertMessages(eventList)

				if err != nil {
					log.Fatal("Add message error: ", err)
				}
				eventList = nil
			}
		}()
		time.Sleep(500 * time.Millisecond)
	}
}

func GetMessagesByUser(userId string) ([]models.EventMessage, error) {
	msgEventList, err := uTool.DB.GetMessagesbyUser(userId)
	if err != nil {
		return nil, err
	}
	if len(msgEventList) == 0 {
		return nil, fmt.Errorf("the user does not contain any messages")
	}
	return msgEventList, nil
}

func RedisPushMessage(event models.EventMessage) error {
	_, err := uTool.Rds.LPush(models.MsgCache, event).Result()
	if err != nil {
		log.Println("redis message error, ", err)
		return err
	}
	return nil
}

// Reply text message
func ReplyLineMessage(replyToken string, msg string) error {
	if _, err := uTool.Bot.ReplyMessage(replyToken, linebot.NewTextMessage(msg)).Do(); err != nil {
		log.Println("reply message error ", err)
		return err
	}
	return nil
}

func SendLineMessage(event models.EventMessage) error {
	sendTo := event.UserID
	msg := event.Message
	if _, err := uTool.Bot.PushMessage(sendTo, linebot.NewTextMessage(msg)).Do(); err != nil {
		return err
	}
	return nil
}

func SendImageMessage(userId string, urls []string) error {
	sendTo := userId
	if _, err := uTool.Bot.PushMessage(sendTo, linebot.NewImageMessage(urls[0], urls[1])).Do(); err != nil {
		return err
	}
	return nil

}

func SendLocationMessage(userId string) error {
	sendTo := userId
	if _, err := uTool.Bot.PushMessage(sendTo, linebot.NewLocationMessage("臺北市政府警察局中正第一分局",
		"100台北市中正區公園路15號", 25.0451104, 121.5173231)).Do(); err != nil {
		return err
	}
	return nil
}
