package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// StoryRating - User rating for a story (1-5 stars)
type StoryRating struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null;uniqueIndex:idx_story_rating_user_story"`
	StoryID   uuid.UUID `json:"story_id" gorm:"type:uuid;not null;uniqueIndex:idx_story_rating_user_story;index"`
	Rating    int       `json:"rating" gorm:"not null"` // 1-5
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	User  User  `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Story Story `json:"-" gorm:"foreignKey:StoryID"`
}

func (StoryRating) TableName() string {
	return "story_ratings"
}

func (r *StoryRating) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}
