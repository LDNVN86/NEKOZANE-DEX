package repositories

import (
	"nekozanedex/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StoryRatingRepository interface {
	CreateOrUpdate(rating *models.StoryRating) error
	FindByUserAndStory(userID, storyID uuid.UUID) (*models.StoryRating, error)
	Delete(userID, storyID uuid.UUID) error
	GetAverageRating(storyID uuid.UUID) (float64, int64, error)
}

type storyRatingRepository struct {
	db *gorm.DB
}

func NewStoryRatingRepository(db *gorm.DB) StoryRatingRepository {
	return &storyRatingRepository{db: db}
}

// CreateOrUpdate - Create or update user rating for a story
func (r *storyRatingRepository) CreateOrUpdate(rating *models.StoryRating) error {
	// Use upsert: if exists, update; otherwise create
	return r.db.Save(rating).Error
}

// FindByUserAndStory - Get user's rating for a story
func (r *storyRatingRepository) FindByUserAndStory(userID, storyID uuid.UUID) (*models.StoryRating, error) {
	var rating models.StoryRating
	err := r.db.Where("user_id = ? AND story_id = ?", userID, storyID).First(&rating).Error
	if err != nil {
		return nil, err
	}
	return &rating, nil
}

// Delete - Remove user's rating for a story
func (r *storyRatingRepository) Delete(userID, storyID uuid.UUID) error {
	return r.db.Where("user_id = ? AND story_id = ?", userID, storyID).Delete(&models.StoryRating{}).Error
}

// GetAverageRating - Calculate average rating and count for a story
func (r *storyRatingRepository) GetAverageRating(storyID uuid.UUID) (float64, int64, error) {
	var result struct {
		Avg   float64
		Count int64
	}

	err := r.db.Model(&models.StoryRating{}).
		Select("COALESCE(AVG(rating), 0) as avg, COUNT(*) as count").
		Where("story_id = ?", storyID).
		Scan(&result).Error

	return result.Avg, result.Count, err
}
