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

func (r *URLRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status crawler.Status) result.Result[*crawler.URL] {
	tx := r.db.WithContext(ctx)

	if err := tx.Model(&URLModel{}).
		Where("id = ?", id).
		Update("status", string(status)).
		Error; err != nil {
		return result.Err[*crawler.URL](fmt.Errorf("failed to update URL status: %w", err))
	}

	return r.FindByID(ctx, id)
}

func (r *URLRepository) IncrementAttemptCount(ctx context.Context, id uuid.UUID) result.Result[*crawler.URL] {
	tx := r.db.WithContext(ctx)

	if err := tx.Model(&URLModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"attempt_count": gorm.Expr("attempt_count + ?", 1),
			"last_attempt": time.Now(),
			"update_at": time.Now(),
		}).Error; err != nil {
			return result.Err[*crawler.URL](fmt.Errorf("failed to increment attempt count: %w", err))
		}

		return r.FindByID(ctx, id)
}

func (r *URLRepository) FindByDomain(ctx context.Context, domain string, limit int) result.Result[[]crawler.URL] {
	tx := r.db.WithContext(ctx)
	var models []URLModel

	if err := tx.Where("url LIKE ?", "%"+domain+"%").
		Order("created_at DESC").
		Limit(limit).
		Find(&models).Error; err != nil {
		return result.Err[[]crawler.URL](fmt.Errorf("failed to find URLs by domain: %w", err))
	}

	urls := make([]crawler.URL, len(models))
	for i, model := range models {
		urls[i] = *model.ToDomain()
	}

	return result.Ok(urls)
}
// CountByStatus counts URLs by status
func (r *URLRepository) CountByStatus(ctx context.Context, status crawler.Status) result.Result[int] {
	tx := r.db.WithContext(ctx)
	var count int64

	if err := tx.Model(&URLModel{}).
		Where("status = ?", string(status)).
		Count(&count).Error; err != nil {
		return result.Err[int](fmt.Errorf("failed to count URLs by status: %w", err))
	}

	return result.Ok(int(count))
}

// DeleteOlderThan deletes URLs older than the given duration
func (r *URLRepository) DeleteOlderThan(ctx context.Context, days int) result.Result[int] {
	tx := r.db.WithContext(ctx)
	cutoff := time.Now().AddDate(0, 0, -days)

	dbResult := tx.Where("created_at < ?", cutoff).Delete(&URLModel{})
	if dbResult.Error != nil {
		return result.Err[int](fmt.Errorf("failed to delete old URLs: %w", dbResult.Error))
	}

	return result.Ok(int(dbResult.RowsAffected))
}