package services

import (
	"errors"

	"nekozanedex/internal/models"
	"nekozanedex/internal/repositories"

	"github.com/google/uuid"
)

type StoryRatingService interface {
	RateStory(userID, storyID uuid.UUID, rating int) (*models.StoryRating, error)
	GetMyRating(userID, storyID uuid.UUID) (*int, error)
	DeleteMyRating(userID, storyID uuid.UUID) error
	GetStoryRating(storyID uuid.UUID) (float64, int64, error)
}

type storyRatingService struct {
	ratingRepo repositories.StoryRatingRepository
	storyRepo  repositories.StoryRepository
}

func NewStoryRatingService(
	ratingRepo repositories.StoryRatingRepository,
	storyRepo repositories.StoryRepository,
) StoryRatingService {
	return &storyRatingService{
		ratingRepo: ratingRepo,
		storyRepo:  storyRepo,
	}
}

// RateStory - Rate a story (1-5 stars), creates or updates existing rating
func (s *storyRatingService) RateStory(userID, storyID uuid.UUID, rating int) (*models.StoryRating, error) {
	// Validate rating range
	if rating < 1 || rating > 5 {
		return nil, errors.New("rating phải từ 1 đến 5")
	}

	// Check story exists
	story, err := s.storyRepo.FindStoryByID(storyID)
	if err != nil {
		return nil, errors.New("truyện không tồn tại")
	}

	// Check if user already rated
	existingRating, err := s.ratingRepo.FindByUserAndStory(userID, storyID)
	if err == nil && existingRating != nil {
		// Update existing rating
		existingRating.Rating = rating
		if err := s.ratingRepo.CreateOrUpdate(existingRating); err != nil {
			return nil, err
		}
		// Update story's cached rating
		s.updateStoryCachedRating(storyID, story)
		return existingRating, nil
	}

	// Create new rating
	newRating := &models.StoryRating{
		UserID:  userID,
		StoryID: storyID,
		Rating:  rating,
	}
	if err := s.ratingRepo.CreateOrUpdate(newRating); err != nil {
		return nil, err
	}

	// Update story's cached rating
	s.updateStoryCachedRating(storyID, story)

	return newRating, nil
}

// GetMyRating - Get current user's rating for a story (nil if not rated)
func (s *storyRatingService) GetMyRating(userID, storyID uuid.UUID) (*int, error) {
	rating, err := s.ratingRepo.FindByUserAndStory(userID, storyID)
	if err != nil {
		return nil, nil // Not rated yet
	}
	return &rating.Rating, nil
}

// DeleteMyRating - Remove user's rating from a story
func (s *storyRatingService) DeleteMyRating(userID, storyID uuid.UUID) error {
	story, err := s.storyRepo.FindStoryByID(storyID)
	if err != nil {
		return errors.New("truyện không tồn tại")
	}

	if err := s.ratingRepo.Delete(userID, storyID); err != nil {
		return err
	}

	// Update story's cached rating
	s.updateStoryCachedRating(storyID, story)
	return nil
}

// GetStoryRating - Get average rating and count for a story
func (s *storyRatingService) GetStoryRating(storyID uuid.UUID) (float64, int64, error) {
	return s.ratingRepo.GetAverageRating(storyID)
}

// updateStoryCachedRating - Update cached rating fields on story model
func (s *storyRatingService) updateStoryCachedRating(storyID uuid.UUID, story *models.Story) {
	avg, count, err := s.ratingRepo.GetAverageRating(storyID)
	if err != nil {
		return
	}

	story.Rating = &avg
	story.RatingCount = int(count)
	s.storyRepo.UpdateStory(story)
}
