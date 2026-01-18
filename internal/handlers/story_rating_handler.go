package handlers

import (
	"nekozanedex/internal/services"
	"nekozanedex/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type StoryRatingHandler struct {
	ratingService services.StoryRatingService
}

func NewStoryRatingHandler(ratingService services.StoryRatingService) *StoryRatingHandler {
	return &StoryRatingHandler{ratingService: ratingService}
}

// RateStoryRequest - Request body for rating a story
type RateStoryRequest struct {
	Rating int `json:"rating" binding:"required,min=1,max=5"`
}

// RateStory godoc
// @Summary Rate a story
// @Tags Story Ratings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param storyId path string true "Story ID"
// @Param body body RateStoryRequest true "Rating (1-5)"
// @Success 200 {object} response.Response
// @Router /api/ratings/story/{storyId} [post]
func (h *StoryRatingHandler) RateStory(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "Chưa đăng nhập")
		return
	}

	storyIDStr := c.Param("storyId")
	storyID, err := uuid.Parse(storyIDStr)
	if err != nil {
		response.BadRequest(c, "Story ID không hợp lệ")
		return
	}

	var req RateStoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Rating phải từ 1 đến 5")
		return
	}

	rating, err := h.ratingService.RateStory(userID.(uuid.UUID), storyID, req.Rating)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Get updated story rating
	avg, count, _ := h.ratingService.GetStoryRating(storyID)

	response.Oke(c, gin.H{
		"my_rating":    rating.Rating,
		"avg_rating":   avg,
		"rating_count": count,
	})
}

// GetMyRating godoc
// @Summary Get user's rating for a story
// @Tags Story Ratings
// @Security BearerAuth
// @Produce json
// @Param storyId path string true "Story ID"
// @Success 200 {object} response.Response
// @Router /api/ratings/story/{storyId}/my [get]
func (h *StoryRatingHandler) GetMyRating(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "Chưa đăng nhập")
		return
	}

	storyIDStr := c.Param("storyId")
	storyID, err := uuid.Parse(storyIDStr)
	if err != nil {
		response.BadRequest(c, "Story ID không hợp lệ")
		return
	}

	rating, err := h.ratingService.GetMyRating(userID.(uuid.UUID), storyID)
	if err != nil {
		response.InternalServerError(c, "Lỗi khi lấy rating")
		return
	}

	response.Oke(c, gin.H{
		"my_rating": rating, // nil if not rated
	})
}

// DeleteMyRating godoc
// @Summary Delete user's rating for a story
// @Tags Story Ratings
// @Security BearerAuth
// @Produce json
// @Param storyId path string true "Story ID"
// @Success 200 {object} response.Response
// @Router /api/ratings/story/{storyId}/my [delete]
func (h *StoryRatingHandler) DeleteMyRating(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "Chưa đăng nhập")
		return
	}

	storyIDStr := c.Param("storyId")
	storyID, err := uuid.Parse(storyIDStr)
	if err != nil {
		response.BadRequest(c, "Story ID không hợp lệ")
		return
	}

	if err := h.ratingService.DeleteMyRating(userID.(uuid.UUID), storyID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Get updated story rating
	avg, count, _ := h.ratingService.GetStoryRating(storyID)

	response.Oke(c, gin.H{
		"message":      "Đã xóa đánh giá",
		"avg_rating":   avg,
		"rating_count": count,
	})
}

// GetStoryRating godoc
// @Summary Get story's average rating
// @Tags Story Ratings
// @Produce json
// @Param storyId path string true "Story ID"
// @Success 200 {object} response.Response
// @Router /api/ratings/story/{storyId} [get]
func (h *StoryRatingHandler) GetStoryRating(c *gin.Context) {
	storyIDStr := c.Param("storyId")
	storyID, err := uuid.Parse(storyIDStr)
	if err != nil {
		response.BadRequest(c, "Story ID không hợp lệ")
		return
	}

	avg, count, err := h.ratingService.GetStoryRating(storyID)
	if err != nil {
		response.InternalServerError(c, "Lỗi khi lấy rating")
		return
	}

	response.Oke(c, gin.H{
		"avg_rating":   avg,
		"rating_count": count,
	})
}
