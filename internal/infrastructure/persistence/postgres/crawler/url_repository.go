package crawler

import (
	"context"
	"fmt"

	"github.com/gerthdala/webcrawler/internal/domain/crawler"
	result "github.com/gerthdala/webcrawler/pkg/utils/result"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type URLRepository struct {
	db *gorm.DB
}

func NewURLRepository(db *gorm.DB) *URLRepository {
	return &URLRepository{
		db: db,
	}
}

func (r *URLRepository) Save(ctx context.Context, url *crawler.URL) result.Result[*crawler.URL] {
	tx := r.db.WithContext(ctx)
	model := URLModelFromDomain(url)

	if err := tx.Create(model).Error; err != nil {
		return result.Err[*crawler.URL](fmt.Errorf("failed to save URL: %w", err))
	}

	return result.Ok(url)
}

func (r *URLRepository) FindByID(ctx context.Context, id uuid.UUID) result.Result[*crawler.URL] {
	tx := r.db.WithContext(ctx)
	var model URLModel

	if err := tx.Where("id = ?", id).First(&model).Error; err != nil {
		return result.Err[*crawler.URL](fmt.Errorf("failed to find URL by ID: %w", err))
	}
	return result.Ok(model.ToDomain())
}

func (r *URLRepository) FindByNormalizedURL(ctx context.Context, normalizedURL string) result.Result[*crawler.URL] {
	tx := r.db.WithContext(ctx)
	var model URLModel

	if err := tx.Where("normalized_url = ?", normalizedURL).First(&model).Error; err != nil {
		return result.Err[*crawler.URL](fmt.Errorf("failed to find URL by normalized URL: %w", err))
	}
	return result.Ok(model.ToDomain())
}

func (r *URLRepository) FindPending(ctx context.Context, limit int) result.Result[[]crawler.URL] {
	tx := r.db.WithContext(ctx)
	var models []URLModel

	if err := tx.Where("status= ?", string(crawler.StatusPending)).
		Order("created_at ASC").
		Limit(limit).
		Find(&models).Error; err !=nil {
		return result.Err[[]crawler.URL](fmt.Errorf("failed to find pending URLs: %w", err))
	}

	urls := make([]crawler.URL, len(models))
	for i, model := range models {
		urls[i] = *model.ToDomain()
	}

	return result.Ok(urls)
}