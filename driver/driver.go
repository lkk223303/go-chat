package driver

import (
	"context"
	"log"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// mongoDB driver here
var mgoCli *mongo.Client

func initMongoDB() {
	var err error

	username := viper.GetString("mongo.user")
	password := viper.GetString("mongo.password")
	port := viper.GetString("mongo.port")
	host := viper.GetString("mongo.host")
	source := "mongodb://" + host + ":" + port

	clientOptions := options.Client().ApplyURI(source)
	b := &options.Credential{
		Username: username,
		Password: password,
	}
	clientOptions.Auth = b

	// connect to mongo
	mgoCli, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = mgoCli.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
}
func GetMongoClient() *mongo.Client {
	if mgoCli == nil {
		initMongoDB()
	}

	log.Println("mongo inited")
	return mgoCli
}
