package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"hybrid-app/backend/internal/domain"
)

type Repository struct {
	client            *mongo.Client
	messageCollection *mongo.Collection
}

func New(ctx context.Context, uri, database string) (*Repository, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}
	db := client.Database(database)
	return &Repository{
		client:            client,
		messageCollection: db.Collection("chat_messages"),
	}, nil
}

func (r *Repository) Close(ctx context.Context) error {
	return r.client.Disconnect(ctx)
}

func (r *Repository) AddMessage(message domain.ChatMessage) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.messageCollection.InsertOne(ctx, message)
	return err
}

func (r *Repository) ListMessages(chatID string) ([]domain.ChatMessage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cur, err := r.messageCollection.Find(ctx, bson.M{"chatId": chatID}, options.Find().SetSort(bson.D{{Key: "createdAt", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var messages []domain.ChatMessage
	for cur.Next(ctx) {
		var message domain.ChatMessage
		if err := cur.Decode(&message); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, cur.Err()
}
