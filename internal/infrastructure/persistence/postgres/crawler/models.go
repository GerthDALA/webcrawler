package crawler

import (
	"fmt"
	"time"

	"github.com/gerthdala/webcrawler/internal/domain/crawler"
	result "github.com/gerthdala/webcrawler/pkg/utils/result"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConfig struct {
	Host         string
	Port         int
	User         string
	Password     string
	Database     string
	SSLMode      string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
}

func NewDB(config DBConfig) result.Result[*gorm.DB] {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.Database, config.SSLMode,
	)

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return result.Error[*gorm.DB](fmt.Errorf("failed to connect to database: %w", err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		return result.Error[*gorm.DB](fmt.Errorf("failed to get database connection: %w", err))
	}
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(config.MaxLifetime)

	return result.Ok(db)
}

// URLModel is the database model for URL
type URLModel struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key"`
	URL           string    `gorm:"index;not null"`
	NormalizedURL string    `gorm:"uniqueIndex;not null"`
	Depth         int       `gorm:"not null"`
	Status        string    `gorm:"index;not null"`
	ParentURL     string    `gorm:"index"`
	AttemptCount  int       `gorm:"not null;default:0"`
	LastAttempt   time.Time
	CreatedAt     time.Time `gorm:"index;not null"`
	UpdatedAt     time.Time `gorm:"not null"`
}

// TableName returns the table name for the URL model
func (URLModel) TableName() string {
	return "urls"
}

// ToDomain converts URLModel to domain URL
func (m *URLModel) ToDomain() *crawler.URL {
	return &crawler.URL{
		ID:            m.ID,
		URL:           m.URL,
		NormalizedURL: m.NormalizedURL,
		Depth:         m.Depth,
		Status:        crawler.Status(m.Status),
		ParentURL:     m.ParentURL,
		AttemptCount:  m.AttemptCount,
		LastAttempt:   m.LastAttempt,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

// FromDomain converts domain URL to URLModel
func URLModelFromDomain(url *crawler.URL) *URLModel {
	return &URLModel{
		ID:            url.ID,
		URL:           url.URL,
		NormalizedURL: url.NormalizedURL,
		Depth:         url.Depth,
		Status:        string(url.Status),
		ParentURL:     url.ParentURL,
		AttemptCount:  url.AttemptCount,
		LastAttempt:   url.LastAttempt,
		CreatedAt:     url.CreatedAt,
		UpdatedAt:     url.UpdatedAt,
	}
}

// PageModel is the database model for Page
type PageModel struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key"`
	URL         string    `gorm:"uniqueIndex;not null"`
	StatusCode  int       `gorm:"not null"`
	Title       string
	HTML        string            `gorm:"type:text"`
	PlainText   string            `gorm:"type:text"`
	Headers     map[string]string `gorm:"type:jsonb"`
	Links       []string          `gorm:"type:jsonb"`
	ContentType string
	FetchedAt   time.Time `gorm:"index;not null"`
	ParsedAt    time.Time
}

// TableName returns the table name for the Page model
func (PageModel) TableName() string {
	return "pages"
}

// ToDomain converts PageModel to domain Page
func (m *PageModel) ToDomain() *crawler.Page {
	return &crawler.Page{
		ID:          m.ID,
		URL:         m.URL,
		StatusCode:  m.StatusCode,
		Title:       m.Title,
		HTML:        m.HTML,
		PlainText:   m.PlainText,
		Headers:     m.Headers,
		Links:       m.Links,
		ContentType: m.ContentType,
		FetchedAt:   m.FetchedAt,
		ParsedAt:    m.ParsedAt,
	}
}

// FromDomain converts domain Page to PageModel
func PageModelFromDomain(page *crawler.Page) *PageModel {
	return &PageModel{
		ID:          page.ID,
		URL:         page.URL,
		StatusCode:  page.StatusCode,
		Title:       page.Title,
		HTML:        page.HTML,
		PlainText:   page.PlainText,
		Headers:     page.Headers,
		Links:       page.Links,
		ContentType: page.ContentType,
		FetchedAt:   page.FetchedAt,
		ParsedAt:    page.ParsedAt,
	}
}

// CrawlJobModel is the database model for CrawlJob
type CrawlJobModel struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	URLID     uuid.UUID `gorm:"type:uuid;index;not null"`
	URL       string    `gorm:"not null"`
	Depth     int       `gorm:"not null"`
	Priority  int       `gorm:"index;not null"`
	Status    string    `gorm:"index;not null;default:'pending'"`
	CreatedAt time.Time `gorm:"index;not null"`
	StartedAt *time.Time
	Error     string
}

// TableName returns the table name for the CrawlJob model
func (CrawlJobModel) TableName() string {
	return "crawl_jobs"
}

// ToDomain converts CrawlJobModel to domain CrawlJob
func (m *CrawlJobModel) ToDomain() *crawler.CrawlJob {
	urlEntity := &crawler.URL{
		ID:    m.URLID,
		URL:   m.URL,
		Depth: m.Depth,
	}
	return &crawler.CrawlJob{
		URL:       urlEntity,
		CreatedAt: m.CreatedAt,
		Priority:  m.Priority,
	}
}

// FromDomain converts domain CrawlJob to CrawlJobModel
func CrawlJobModelFromDomain(job *crawler.CrawlJob) *CrawlJobModel {
	return &CrawlJobModel{
		URLID:     job.URL.ID,
		URL:       job.URL.URL,
		Depth:     job.URL.Depth,
		Priority:  job.Priority,
		Status:    "pending",
		CreatedAt: job.CreatedAt,
	}
}
