package crawler

import (
	"context"

	result "github.com/gerthdala/webcrawler/pkg/utils/result"
	"github.com/google/uuid"
)

type URLRepository interface {
	Save(ctx context.Context, url *URL) result.Result[*URL]

	FindByID(ctx context.Context, id uuid.UUID) result.Result[*URL]

	FindByNormalizedURL(ctx context.Context, normalizedURL string) result.Result[*URL]

	// FindPending finds URLs with pending status, with limit
	FindPending(ctx context.Context, limit int) result.Result[[]URL]

	// UpdateStatus updates the status of a URL
	UpdateStatus(ctx context.Context, id uuid.UUID, status Status) result.Result[*URL]
	
	// IncrementAttemptCount increments the attempt count of a URL
	IncrementAttemptCount(ctx context.Context, id uuid.UUID) result.Result[*URL]
	
	// FindByDomain finds URLs by domain
	FindByDomain(ctx context.Context, domain string, limit int) result.Result[[]URL]
	
	// CountByStatus counts URLs by status
	CountByStatus(ctx context.Context, status Status) result.Result[int]
	
	// DeleteOlderThan deletes URLs older than the given duration
	DeleteOlderThan(ctx context.Context, days int) result.Result[int]
}

// PageRepository handles Page storage and retrieval
type PageRepository interface {
	// Save stores a Page
	Save(ctx context.Context, page *Page) result.Result[*Page]
	
	// FindByID finds a Page by its ID
	FindByID(ctx context.Context, id uuid.UUID) result.Result[*Page]
	
	// FindByURL finds a Page by its URL
	FindByURL(ctx context.Context, url string) result.Result[*Page]
	
	// FindRecent finds recently crawled pages
	FindRecent(ctx context.Context, limit int) result.Result[[]Page]
	
	// CountPages counts the total number of pages
	CountPages(ctx context.Context) result.Result[int]
	
	// Search searches pages by content
	Search(ctx context.Context, query string, limit int) result.Result[[]Page]
	
	// DeleteOlderThan deletes pages older than the given duration
	DeleteOlderThan(ctx context.Context, days int) result.Result[int]
}

// CrawlJobRepository handles CrawlJob queue operations
type CrawlJobRepository interface {
	// Enqueue adds a job to the queue
	Enqueue(ctx context.Context, job *CrawlJob) result.Result[*CrawlJob]
	
	// Dequeue gets the next job from the queue
	Dequeue(ctx context.Context) result.Result[*CrawlJob]
	
	// Count counts the number of jobs in the queue
	Count(ctx context.Context) result.Result[int]
	
	// Clear clears all jobs from the queue
	Clear(ctx context.Context) result.Result[int]
}