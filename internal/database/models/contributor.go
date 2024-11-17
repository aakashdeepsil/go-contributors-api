package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Contributor struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username  string             `bson:"username" json:"username"`
	Email     string             `bson:"email" json:"email"`
	Name      string             `bson:"name" json:"name"`
	AvatarURL string             `bson:"avatar_url" json:"avatarUrl"`
	Projects  []string           `bson:"projects" json:"projects"`
	JoinedAt  time.Time          `bson:"joined_at" json:"joinedAt"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updatedAt"`
}

type ContributorInput struct {
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Name      string   `json:"name"`
	AvatarURL string   `json:"avatarUrl"`
	Projects  []string `json:"projects"`
}
