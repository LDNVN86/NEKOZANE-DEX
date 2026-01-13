package handlers

import (
	"nekozanedex/internal/services"
	"nekozanedex/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GenreHandler struct {
	genreService services.GenreService
}

func NewGenreHandler(genreService services.GenreService) *GenreHandler {
	return &GenreHandler{genreService: genreService}
}

// CreateGenreRequest - DTO cho tạo thể loại
type CreateGenreRequest struct {
	Name        string  `json:"name" binding:"required,min=1,max=50"`
	Description *string `json:"description"`
}

// UpdateGenreRequest - DTO cho cập nhật thể loại
type UpdateGenreRequest struct {
	Name        string  `json:"name" binding:"required,min=1,max=50"`
	Description *string `json:"description"`
}

// GetAllGenres godoc
// @Summary Lấy tất cả thể loại
// @Tags Genres
// @Success 200 {object} response.Response
// @Router /api/genres [get]
func (h *GenreHandler) GetAllGenres(c *gin.Context) {
	genres, err := h.genreService.GetAllGenres()
	if err != nil {
		response.InternalServerError(c, "Không thể lấy thể loại")
		return
	}
	response.Oke(c, genres)
}

// CreateGenre godoc
// @Summary Tạo thể loại mới (Admin)
// @Tags Admin Genres
// @Security BearerAuth
// @Param request body CreateGenreRequest true "Genre data"
// @Success 201 {object} response.Response
// @Router /api/admin/genres [post]
func (h *GenreHandler) CreateGenre(c *gin.Context) {
	var req CreateGenreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Dữ liệu không hợp lệ: "+err.Error())
		return
	}

	genre, err := h.genreService.CreateGenre(req.Name, req.Description)
	if err != nil {
		response.InternalServerError(c, "Không thể tạo thể loại: "+err.Error())
		return
	}

	response.Created(c, genre)
}

// UpdateGenre godoc
// @Summary Cập nhật thể loại (Admin)
// @Tags Admin Genres
// @Security BearerAuth
// @Param id path string true "Genre ID"
// @Param request body UpdateGenreRequest true "Genre data"
// @Success 200 {object} response.Response
// @Router /api/admin/genres/{id} [put]
func (h *GenreHandler) UpdateGenre(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "ID không hợp lệ")
		return
	}

	var req UpdateGenreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Dữ liệu không hợp lệ: "+err.Error())
		return
	}

	genre, err := h.genreService.UpdateGenre(id, req.Name, req.Description)
	if err != nil {
		response.InternalServerError(c, "Không thể cập nhật thể loại: "+err.Error())
		return
	}

	response.Oke(c, genre)
}

// DeleteGenre godoc
// @Summary Xóa thể loại (Admin)
// @Tags Admin Genres
// @Security BearerAuth
// @Param id path string true "Genre ID"
// @Success 200 {object} response.Response
// @Router /api/admin/genres/{id} [delete]
func (h *GenreHandler) DeleteGenre(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "ID không hợp lệ")
		return
	}

	if err := h.genreService.DeleteGenre(id); err != nil {
		response.InternalServerError(c, "Không thể xóa thể loại: "+err.Error())
		return
	}

	response.Oke(c, gin.H{"message": "Đã xóa thể loại"})
}
