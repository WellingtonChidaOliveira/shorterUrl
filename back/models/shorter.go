package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Shorter struct {
	ID        primitive.ObjectID `bson:"_id, omitempty"`
	Shorter   string             `bson:"name, omitempty"`
	OriginUrl string             `bson:"url"`
	MostUse   int                `bson:"most_use, omitempty"`
	Status    int                `bson:"status, omitempty"`
	CreatedAt time.Time          `bson:"create_at, omitempty"`
}
