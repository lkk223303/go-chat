package utils

import (
	"chatty/models"
	"chatty/repository"
	"chatty/repository/dbrepo"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
)

type Util struct {
	DB  repository.DatabaseRepo
	Rds *redis.Client
	l   sync.Mutex
}

func NewUtil(client *mongo.Client, r *redis.Client) *Util {
	return &Util{
		DB:  dbrepo.NewMongoRepo(client),
		Rds: r,
	}

}

var eventList []models.EventMessage
var batch int = 10

func (u *Util) AddMessage() {
	t := time.NewTicker(10 * time.Second)
	for {

		var msgEvent models.EventMessage
		// 5秒超時
		value, err := u.Rds.BRPop(5*time.Second, models.MsgCache).Result()
		if err == redis.Nil {
			// 查詢不到資料
			time.Sleep(1 * time.Second)
			continue
		}
		if err != nil {
			time.Sleep(1 * time.Second)
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
			u.l.Lock()
			err = u.DB.InsertMessages(eventList)
			u.l.Unlock()
			if err != nil {
				log.Fatal("Add message error: ", err)
			}
			eventList = nil
		}

		// Or insert messages when time's up
		go func() {
			<-t.C
			if len(eventList) > 0 {
				u.l.Lock()
				err = u.DB.InsertMessages(eventList)
				u.l.Unlock()
				if err != nil {
					log.Fatal("Add message error: ", err)
				}
				eventList = nil
			}
		}()

		time.Sleep(500 * time.Millisecond)
	}
}
