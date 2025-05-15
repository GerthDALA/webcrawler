package crawler

import (
	"context"
	"fmt"
	"time"

	"github.com/gerthdala/webcrawler/internal/domain/crawler"
	result "github.com/gerthdala/webcrawler/pkg/utils/result"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PageRepository implements crawler.PageRepository using PostgreSQL
type PageRepository struct {
	db *gorm.DB
}

// NewPageRepository creates a new PageRepository
func NewPageRepository(db *gorm.DB) *PageRepository {
	return &PageRepository{
		db: db,
	}
}

// Save stores a Page
func (r *PageRepository) Save(ctx context.Context, page *crawler.Page) result.Result[*crawler.Page] {
	tx := r.db.WithContext(ctx)
	model := PageModelFromDomain(page)

	if err := tx.Create(model).Error; err != nil {
		return result.Err[*crawler.Page](fmt.Errorf("failed to save Page: %w", err))
	}

	return result.Ok(page)
}

// FindByID finds a Page by its ID
func (r *PageRepository) FindByID(ctx context.Context, id uuid.UUID) result.Result[*crawler.Page] {
	tx := r.db.WithContext(ctx)
	var model PageModel

	if err := tx.Where("id = ?", id).First(&model).Error; err != nil {
		return result.Err[*crawler.Page](fmt.Errorf("failed to find Page by ID: %w", err))
	}

	return result.Ok(model.ToDomain())
}

// FindByURL finds a Page by its URL
func (r *PageRepository) FindByURL(ctx context.Context, url string) result.Result[*crawler.Page] {
	tx := r.db.WithContext(ctx)
	var model PageModel

	if err := tx.Where("url = ?", url).First(&model).Error; err != nil {
		return result.Err[*crawler.Page](fmt.Errorf("failed to find Page by URL: %w", err))
	}

	return result.Ok(model.ToDomain())
}

// FindRecent finds recently crawled pages
func (r *PageRepository) FindRecent(ctx context.Context, limit int) result.Result[[]crawler.Page] {
	tx := r.db.WithContext(ctx)
	var models []PageModel

	if err := tx.Order("fetched_at DESC").
		Limit(limit).
		Find(&models).Error; err != nil {
		return result.Err[[]crawler.Page](fmt.Errorf("failed to find recent Pages: %w", err))
	}

	pages := make([]crawler.Page, len(models))
	for i, model := range models {
		pages[i] = *model.ToDomain()
	}

	return result.Ok(pages)
}

// CountPages counts the total number of pages
func (r *PageRepository) CountPages(ctx context.Context) result.Result[int] {
	tx := r.db.WithContext(ctx)
	var count int64

	if err := tx.Model(&PageModel{}).Count(&count).Error; err != nil {
		return result.Err[int](fmt.Errorf("failed to count Pages: %w", err))
	}

	return result.Ok(int(count))
}

// Search searches pages by content
func (r *PageRepository) Search(ctx context.Context, query string, limit int) result.Result[[]crawler.Page] {
	tx := r.db.WithContext(ctx)
	var models []PageModel

	// Use PostgreSQL's full-text search
	if err := tx.Where("to_tsvector('english', title || ' ' || plain_text) @@ plainto_tsquery('english', ?)", query).
		Order("fetched_at DESC").
		Limit(limit).
		Find(&models).Error; err != nil {
		return result.Err[[]crawler.Page](fmt.Errorf("failed to search Pages: %w", err))
	}

	pages := make([]crawler.Page, len(models))
	for i, model := range models {
		pages[i] = *model.ToDomain()
	}

	return result.Ok(pages)
}

// DeleteOlderThan deletes pages older than the given duration
func (r *PageRepository) DeleteOlderThan(ctx context.Context, days int) result.Result[int] {
	tx := r.db.WithContext(ctx)
	cutoff := time.Now().AddDate(0, 0, -days)

	resultD := tx.Where("fetched_at < ?", cutoff).Delete(&PageModel{})
	if resultD.Error != nil {
		return result.Err[int](fmt.Errorf("failed to delete old Pages: %w", resultD.Error))
	}

	return result.Ok(int(resultD.RowsAffected))
}