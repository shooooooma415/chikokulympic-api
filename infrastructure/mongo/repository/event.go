package repository

import (
	"context"
	"fmt"
	"time"

	"chikokulympic-api/domain/entity"
	repo "chikokulympic-api/domain/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type EventRepo struct {
	eventCollection *mongo.Collection
}

func NewEventRepository(db *mongo.Database) repo.EventRepository {
	return &EventRepo{
		eventCollection: db.Collection("events"),
	}
}

func (er *EventRepo) FindEventByEventID(eventID entity.EventID) (*entity.Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var event entity.Event
	filter := bson.M{"_id": eventID}
	err := er.eventCollection.FindOne(ctx, filter).Decode(&event)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("event not found with ID: %s", eventID)
		}
		return nil, fmt.Errorf("error finding event by ID: %w", err)
	}

	return &event, nil
}

func (er *EventRepo) CreateEvent(event entity.Event) (*entity.Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := er.eventCollection.InsertOne(ctx, event)
	if err != nil {
		return nil, fmt.Errorf("error creating event: %w", err)
	}

	return &event, nil
}

func (er *EventRepo) DeleteEvent(event entity.Event) (*entity.Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": event.EventID}
	_, err := er.eventCollection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error deleting event: %w", err)
	}

	return &event, nil
}

func (er *EventRepo) UpdateEvent(event entity.Event) (*entity.Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": event.EventID}
	update := bson.M{"$set": event}

	_, err := er.eventCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("error updating event: %w", err)
	}

	return &event, nil
}
