package services

import (
	"nekozanedex/internal/models"
	"nekozanedex/internal/repositories"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

type GenreService interface {
	GetAllGenres() ([]models.Genre, error)
	GetGenreByID(id uuid.UUID) (*models.Genre, error)
	CreateGenre(name string, description *string) (*models.Genre, error)
	UpdateGenre(id uuid.UUID, name string, description *string) (*models.Genre, error)
	DeleteGenre(id uuid.UUID) error
}

type genreService struct {
	genreRepo repositories.GenreRepository
}

func NewGenreService(genreRepo repositories.GenreRepository) GenreService {
	return &genreService{genreRepo: genreRepo}
}

// GetAllGenres - Lấy tất cả thể loại
func (s *genreService) GetAllGenres() ([]models.Genre, error) {
	return s.genreRepo.GetAllGenres()
}

// GetGenreByID - Lấy thể loại theo ID
func (s *genreService) GetGenreByID(id uuid.UUID) (*models.Genre, error) {
	return s.genreRepo.FindGenreByID(id)
}

// CreateGenre - Tạo thể loại mới
func (s *genreService) CreateGenre(name string, description *string) (*models.Genre, error) {
	genre := &models.Genre{
		Name:        name,
		Slug:        slug.Make(name),
		Description: description,
	}

	if err := s.genreRepo.CreateGenre(genre); err != nil {
		return nil, err
	}

	return genre, nil
}

// UpdateGenre - Cập nhật thể loại
func (s *genreService) UpdateGenre(id uuid.UUID, name string, description *string) (*models.Genre, error) {
	genre, err := s.genreRepo.FindGenreByID(id)
	if err != nil {
		return nil, err
	}

	genre.Name = name
	genre.Slug = slug.Make(name)
	if description != nil {
		genre.Description = description
	}

	if err := s.genreRepo.UpdateGenre(genre); err != nil {
		return nil, err
	}

	return genre, nil
}

// DeleteGenre - Xóa thể loại
func (s *genreService) DeleteGenre(id uuid.UUID) error {
	return s.genreRepo.DeleteGenre(id)
}
