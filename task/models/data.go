package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Rate struct {
	Cfrom string  `json:"cfrom" bson:"cfrom" binding:"required"`
	Cto   string  `json:"cto" bson:"cto" binding:"required"`
	Conv  float32 `json:"conv" bson:"conv" binding:"required"`
}

type BasicExchange struct {
	Rates     []Rate             `json:"rates" bson:"rates" binding:"required"`
	CodeId    int                `json:"codeId" bson:"codeId" binding:"required"`
	Id        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" binding:"required"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at" binding:"required"`
	Amount    float32            `json:"amount" bson:"amount" binding:"required"`
	NewAmount float32            `json:"NewAmount" bson:"NewAmount" binding:"required"`
}
