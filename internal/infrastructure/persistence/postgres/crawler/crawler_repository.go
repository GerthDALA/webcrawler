package crawler

import (
	"context"
	"fmt"
	"time"

	"github.com/gerthdala/webcrawler/internal/domain/crawler"
	result "github.com/gerthdala/webcrawler/pkg/utils/result"
	"gorm.io/gorm"
)

// CrawlJobRepository implements crawler.CrawlJobRepository using PostgreSQL
type CrawlJobRepository struct {
	db *gorm.DB
}

// NewCrawlJobRepository creates a new CrawlJobRepository
func NewCrawlJobRepository(db *gorm.DB) *CrawlJobRepository {
	return &CrawlJobRepository{
		db: db,
	}
}

// Enqueue adds a job to the queue
func (r *CrawlJobRepository) Enqueue(ctx context.Context, job *crawler.CrawlJob) result.Result[*crawler.CrawlJob] {
	tx := r.db.WithContext(ctx)
	model := CrawlJobModelFromDomain(job)

	if err := tx.Create(model).Error; err != nil {
		return result.Err[*crawler.CrawlJob](fmt.Errorf("failed to enqueue job: %w", err))
	}

	return result.Ok(job)
}

// Dequeue gets the next job from the queue
func (r *CrawlJobRepository) Dequeue(ctx context.Context) result.Result[*crawler.CrawlJob] {
	
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return result.Err[*crawler.CrawlJob](fmt.Errorf("failed to begin transaction: %w", tx.Error))
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var model CrawlJobModel
	if err := tx.Set("gorm:query_option", "FOR UPDATE SKIP LOCKED").
		Where("status = ?", "pending").
		Order("priority DESC, created_at ASC").
		First(&model).Error; err != nil {
		tx.Rollback()
		return result.Err[*crawler.CrawlJob](fmt.Errorf("failed to dequeue job: %w", err))
	}

	now := time.Now()
	if err := tx.Model(&CrawlJobModel{}).
		Where("id = ?", model.ID).
		Updates(map[string]interface{}{
			"status":     "processing",
			"started_at": now,
		}).Error; err != nil {
		tx.Rollback()
		return result.Err[*crawler.CrawlJob](fmt.Errorf("failed to update job status: %w", err))
	}

	if err := tx.Commit().Error; err != nil {
		return result.Err[*crawler.CrawlJob](fmt.Errorf("failed to commit transaction: %w", err))
	}

	return result.Ok(model.ToDomain())
}

// Count counts the number of jobs in the queue
func (r *CrawlJobRepository) Count(ctx context.Context) result.Result[int] {
	tx := r.db.WithContext(ctx)
	var count int64

	if err := tx.Model(&CrawlJobModel{}).
		Where("status = ?", "pending").
		Count(&count).Error; err != nil {
		return result.Err[int](fmt.Errorf("failed to count jobs: %w", err))
	}

	return result.Ok(int(count))
}

// Clear clears all jobs from the queue
func (r *CrawlJobRepository) Clear(ctx context.Context) result.Result[int] {
	tx := r.db.WithContext(ctx)

	resultD := tx.Where("status = ?", "pending").Delete(&CrawlJobModel{})
	if resultD.Error != nil {
		return result.Err[int](fmt.Errorf("failed to clear jobs: %w", resultD.Error))
	}

	return result.Ok(int(resultD.RowsAffected))
}