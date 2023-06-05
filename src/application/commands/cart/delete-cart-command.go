package commands

import "go.mongodb.org/mongo-driver/bson/primitive"

type DeleteCartCommand struct {
	ID primitive.ObjectID `json:"id"`
}
