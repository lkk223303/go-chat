package utils

import (
	"chatty/models"
	"chatty/repository"
	"chatty/repository/dbrepo"
	"log"
	"time"

	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
)

type Util struct {
	DB  repository.DatabaseRepo
	Rds *redis.Client
}

func NewUtil(client *mongo.Client, r *redis.Client) *Util {
	return &Util{
		DB:  dbrepo.NewMongoRepo(client),
		Rds: r,
	}

}

func (u *Util) AddMessage() {
	for {
		// 設定一個5秒的超時時間
		value, err := u.Rds.BRPop(5*time.Second, models.MsgCache).Result()
		if err == redis.Nil {
			// 查詢不到資料
			time.Sleep(1 * time.Second)
			continue
		}
		if err != nil {
			// 查詢出錯
			time.Sleep(1 * time.Second)
			continue
		}

		log.Println("value get", value, " at time: ", time.Now().Unix())
		time.Sleep(time.Second)

	}
}
