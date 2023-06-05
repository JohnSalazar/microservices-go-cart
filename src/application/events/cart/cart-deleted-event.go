package events

import "go.mongodb.org/mongo-driver/bson/primitive"

type CartDeletedEvent struct {
	ID primitive.ObjectID `json:"id"`
}
