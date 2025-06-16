package model

import (
	"time"

	"github.com/google/uuid"
)

type Url struct {
	ID          uuid.UUID `bson:"_id, omitempty"`
	OriginalUrl string    `bson:"original_url"`
	Shorter     string    `bson:"shorter"`
	CreatedAt   time.Time `bson:"created_at"`
}
