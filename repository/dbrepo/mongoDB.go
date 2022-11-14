package dbrepo

import (
	"chatty/models"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

var MAXTIMEOUT = 5 * time.Second

// mongo method
func (m *mongoDBRepo) Testing() {
	ctx, cancel := context.WithTimeout(context.Background(), MAXTIMEOUT)
	defer cancel()

	collection := m.Client.Database("testing").Collection("numbers")
	res, err := collection.InsertOne(ctx, bson.D{{"name", "pi"}, {"value", 3.14159}})
	if err != nil {
		log.Fatal("insert err", err)
	}
	id := res.InsertedID
	log.Println(id) //ObjectID("636e7fa74744c16297f598c8")
}

func (m *mongoDBRepo) InsertMessage(msg models.EventMessage) error {
	ctx, cancel := context.WithTimeout(context.Background(), MAXTIMEOUT)
	defer cancel()
	collection := m.Client.Database("line").Collection("Messages")
	res, err := collection.InsertOne(ctx, msg)
	if err != nil {
		return fmt.Errorf("insert message error: %s", err)
	}
	log.Println("insert message success, id: ", res.InsertedID)
	return nil
}

func (m *mongoDBRepo) InsertMessages(msgs []models.EventMessage) error {
	ctx, cancel := context.WithTimeout(context.Background(), MAXTIMEOUT)
	defer cancel()
	collection := m.Client.Database("line").Collection("Messages")
	var u []interface{}
	for _, t := range msgs {
		u = append(u, t)
	}
	res, err := collection.InsertMany(ctx, u)
	if err != nil {
		return fmt.Errorf("insert multiple messages error: %s", err)
	}
	inserts := len(res.InsertedIDs)
	log.Printf("insert %d messages success, ids: %s", inserts, res.InsertedIDs)
	return nil
}

func (m *mongoDBRepo) GetMessagesFromUser(userId string) ([]models.EventMessage, error) {
	var msgEvents []models.EventMessage
	ctx, cancel := context.WithTimeout(context.Background(), MAXTIMEOUT)
	defer cancel()
	collection := m.Client.Database("line").Collection("Messages")
	filter := bson.M{"userid": userId}
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	for cur.Next(ctx) {
		var event models.EventMessage
		err = cur.Decode(&event)
		if err != nil {
			return nil, err
		}
		msgEvents = append(msgEvents, event)
	}

	return msgEvents, nil
}
