package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type Models struct {
	LogEntry LogEntry
}

type LogEntry struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string    `bson:"name" json:"name"`
	Data      string    `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func New(mongo *mongo.Client) Models {
	client = mongo

	return Models{
		LogEntry: LogEntry{},
	}
}

func (l *LogEntry) Insert(entry LogEntry) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*12)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	_, err := collection.InsertOne(ctx, LogEntry{
		Name:      entry.Name,
		Data:      entry.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		log.Println("Error inserting into logs", err)
		return err
	}

	return nil
}

func (l *LogEntry) All() ([]*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*12)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	options := options.Find()
	options.SetSort(bson.D{{"created_at", -1}})

	cursor, err := collection.Find(ctx, bson.D{}, options)
	if err != nil {
		log.Println("Error getting all records", err)
		return nil, err
	}

	defer cursor.Close(ctx)
	var logs []*LogEntry
	for cursor.Next(ctx) {
		var item LogEntry
		err := cursor.Decode(&item)
		if err != nil {
			log.Println("Error getting all records", err)
			return nil, err
		}

		logs = append(logs, &item)
	}

	return logs, nil
}
