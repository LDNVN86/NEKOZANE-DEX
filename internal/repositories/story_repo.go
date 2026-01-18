package repositories

import (
	"strings"

	"nekozanedex/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// escapeSearchQuery - Escape SQL LIKE wildcards to prevent injection and slow query attacks
func escapeSearchQuery(query string) string {
	// Escape backslash first (order matters)
	query = strings.ReplaceAll(query, "\\", "\\\\")
	// Escape % and _ which are LIKE wildcards
	query = strings.ReplaceAll(query, "%", "\\%")
	query = strings.ReplaceAll(query, "_", "\\_")
	return query
}

type StoryRepository interface{
	CreateStory(story *models.Story) error
	FindStoryByID(id uuid.UUID) (*models.Story,error)
	FindStoryBySlug(slug string) (*models.Story,error)
	UpdateStory(story *models.Story) error
	UpdateStoryGenres(storyID uuid.UUID, genreIDs []uuid.UUID) error
	DeleteStory(id uuid.UUID) error
	GetAllStories(page, limit int, published bool) ([]models.Story, int64, error)
	GetStoriesByGenre(genreID uuid.UUID, page, limit int) ([]models.Story, int64, error)
	GetStoriesLatest(limit int) ([]models.Story, error)
	GetStoriesHot(limit int) ([]models.Story, error)
	SearchStories(query string, page, limit int) ([]models.Story, int64, error)
	AdvancedSearchStories(filters *SearchFilters) ([]models.Story, int64, error)
	SearchStoriesAdmin(query string, page, limit int) ([]models.Story, int64, error)
	IncrementViewCountStory(id uuid.UUID) error
}

// SearchFilters - Filters for advanced story search
type SearchFilters struct {
	Query      string   // Text search query
	Status     string   // ongoing, completed, hiatus
	Country    string   // JP, CN, KR, VN
	GenreSlugs []string // Genre slugs to filter by
	YearFrom   *int     // Release year from
	YearTo     *int     // Release year to
	SortBy     string   // latest, popular, name, rating
	Page       int
	Limit      int
}

type storyRepository struct {
	db *gorm.DB
}

func NewStoryRepository(db *gorm.DB) StoryRepository{
	return &storyRepository{db:db}
}

//Create Story - Tạo Story
func (r *storyRepository) CreateStory(story *models.Story) error{
	return r.db.Create(story).Error
}

//Find Story By ID - Tìm Story theo ID
func (r *storyRepository) FindStoryByID(id uuid.UUID) (*models.Story, error){
	var story models.Story
	err := r.db.First(&story, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &story, nil
}

//Find Story By Slug - Tìm Story Theo Slug

func (r *storyRepository) FindStoryBySlug(slug string) (*models.Story, error) {
	var story models.Story
	err := r.db.Preload("Genres").Preload("Chapters", func(db *gorm.DB) *gorm.DB {
		return db.Where("is_published = ?", true).Order("chapter_number ASC")
	}).First(&story, "slug = ? AND is_published = ?", slug, true).Error
	if err != nil {
		return nil, err
	}
	return &story, nil
}

//Update Story - Cập Nhật Story
func (r *storyRepository) UpdateStory(story *models.Story) error{
	return r.db.Save(story).Error
}

// UpdateStoryGenres - Cập nhật thể loại của Story
func (r *storyRepository) UpdateStoryGenres(storyID uuid.UUID, genreIDs []uuid.UUID) error {
	// First, get the story
	var story models.Story
	if err := r.db.First(&story, "id = ?", storyID).Error; err != nil {
		return err
	}
	
	// Get genres by IDs
	var genres []models.Genre
	if len(genreIDs) > 0 {
		if err := r.db.Where("id IN ?", genreIDs).Find(&genres).Error; err != nil {
			return err
		}
	}
	
	// Replace association
	return r.db.Model(&story).Association("Genres").Replace(genres)
}

//Delete Story - Xóa Story
func (r *storyRepository) DeleteStory(id uuid.UUID) error{
	return r.db.Delete(&models.Story{}, "id = ?", id).Error
}

//Get All Stories - Lấy Tất Cả Story
func (r *storyRepository) GetAllStories(page, limit int, published bool) ([]models.Story, int64, error) {
	var stories []models.Story
	var total int64
	query := r.db.Model(&models.Story{})
	if published {
		query = query.Where("is_published = ?", true)
	}
	query.Count(&total)
	offset := (page - 1) * limit
	err := query.Preload("Genres").Offset(offset).Limit(limit).Order("updated_at DESC").Find(&stories).Error
	return stories, total, err
}


//Get Stories By Genre - Lấy Story Theo Thể Loại
func (r *storyRepository) GetStoriesByGenre(genreID uuid.UUID, page, limit int) ([]models.Story, int64, error) {
	var stories []models.Story
	var total int64
	subQuery := r.db.Table("story_genres").Select("story_id").Where("genre_id = ?", genreID)
	
	r.db.Model(&models.Story{}).Where("id IN (?) AND is_published = ?", subQuery, true).Count(&total)
	offset := (page - 1) * limit
	err := r.db.Preload("Genres").Where("id IN (?) AND is_published = ?", subQuery, true).
		Offset(offset).Limit(limit).Order("updated_at DESC").Find(&stories).Error
	return stories, total, err
}

//Get Stories Latest - Lấy Story Latest
func (r *storyRepository) GetStoriesLatest(limit int) ([]models.Story, error) {
	var stories []models.Story
	err := r.db.Preload("Genres").Where("is_published = ?", true).
		Order("updated_at DESC").Limit(limit).Find(&stories).Error
	return stories, err
}

//Get Stories Hot - Lấy Story Hot
func (r *storyRepository) GetStoriesHot(limit int) ([]models.Story, error) {
	var stories []models.Story
	err := r.db.Preload("Genres").Where("is_published = ?", true).
		Order("view_count DESC").Limit(limit).Find(&stories).Error
	return stories, err
}

//Search Stories - Tìm kiếm Story
func (r *storyRepository) SearchStories(query string, page, limit int) ([]models.Story, int64, error) {
	var stories []models.Story
	var total int64
	searchQuery := "%" + escapeSearchQuery(query) + "%"
	
	r.db.Model(&models.Story{}).Where("is_published = ? AND (title ILIKE ? OR description ILIKE ?)", 
		true, searchQuery, searchQuery).Count(&total)
	offset := (page - 1) * limit
	err := r.db.Preload("Genres").Where("is_published = ? AND (title ILIKE ? OR description ILIKE ?)", 
		true, searchQuery, searchQuery).Offset(offset).Limit(limit).Find(&stories).Error
	return stories, total, err
}

// AdvancedSearchStories - Search with filters
func (r *storyRepository) AdvancedSearchStories(filters *SearchFilters) ([]models.Story, int64, error) {
	var stories []models.Story
	var total int64

	query := r.db.Model(&models.Story{}).Where("is_published = ?", true)

	// Text search
	if filters.Query != "" {
		searchQuery := "%" + escapeSearchQuery(filters.Query) + "%"
		query = query.Where("title ILIKE ? OR description ILIKE ?", searchQuery, searchQuery)
	}

	// Status filter
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}

	// Country filter
	if filters.Country != "" {
		query = query.Where("country = ?", filters.Country)
	}

	// Year range filter
	if filters.YearFrom != nil {
		query = query.Where("release_year >= ?", *filters.YearFrom)
	}
	if filters.YearTo != nil {
		query = query.Where("release_year <= ?", *filters.YearTo)
	}

	// Genre filter - filter by genre slugs
	if len(filters.GenreSlugs) > 0 {
		query = query.Where(`id IN (
			SELECT sg.story_id FROM story_genres sg
			JOIN genres g ON sg.genre_id = g.id
			WHERE g.slug IN ?
		)`, filters.GenreSlugs)
	}

	// Count total
	query.Count(&total)

	// Sorting
	orderClause := "updated_at DESC" // default: latest
	switch filters.SortBy {
	case "popular":
		orderClause = "view_count DESC"
	case "name":
		orderClause = "title ASC"
	case "rating":
		orderClause = "rating DESC NULLS LAST"
	case "oldest":
		orderClause = "created_at ASC"
	}

	// Paginate and fetch
	offset := (filters.Page - 1) * filters.Limit
	err := query.Preload("Genres").Order(orderClause).Offset(offset).Limit(filters.Limit).Find(&stories).Error

	return stories, total, err
}

//Increment View Count Story - Tăng Lượt Xem Story
func (r *storyRepository) IncrementViewCountStory(id uuid.UUID) error {
	return r.db.Model(&models.Story{}).Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

//SearchStoriesAdmin - Admin search (includes drafts)
func (r *storyRepository) SearchStoriesAdmin(query string, page, limit int) ([]models.Story, int64, error) {
	var stories []models.Story
	var total int64
	searchQuery := "%" + escapeSearchQuery(query) + "%"
	
	r.db.Model(&models.Story{}).Where("title ILIKE ? OR description ILIKE ?", 
		searchQuery, searchQuery).Count(&total)
	offset := (page - 1) * limit
	err := r.db.Preload("Genres").Where("title ILIKE ? OR description ILIKE ?", 
		searchQuery, searchQuery).Offset(offset).Limit(limit).Order("updated_at DESC").Find(&stories).Error
	return stories, total, err
}