package crawler

import (
	"net/url"
	"time"

	result "github.com/gerthdala/webcrawler/pkg/utils/result"

	"github.com/google/uuid"
)

// Status represents the status of a URL in the crawling process
type Status string

const (
	StatusPending  Status = "pending"
	StatusFetching Status = "fetching"
	StatusFetched  Status = "fetched"
	StatusFailed   Status = "failed"
)

// URL represents a URL to be crawled
type URL struct {
	ID            uuid.UUID
	URL           string
	NormalizedURL string
	Depth         int
	Status        Status
	ParentURL     string
	AttemptCount  int
	LastAttempt   time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// NewURL creates a new URL entity
func NewURL(rawURL string, depth int, parentURL string) result.Result[*URL] {
	// Validate URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return result.Err[*URL](err)
	}

	// Normalize URL
	normalizedURL := normalizeURL(parsedURL)

	return result.Ok(&URL{
		ID:            uuid.New(),
		URL:           rawURL,
		NormalizedURL: normalizedURL,
		Depth:         depth,
		Status:        StatusPending,
		ParentURL:     parentURL,
		AttemptCount:  0,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	})
}

// normalizeURL standardizes a URL by removing fragments, default ports, etc.
func normalizeURL(u *url.URL) string {
	// Make a copy to avoid modifying the original
	normalized := *u

	// Remove fragments
	normalized.Fragment = ""

	// Use https by default if the scheme is empty
	if normalized.Scheme == "" {
		normalized.Scheme = "https"
	}

	// Remove default ports
	if (normalized.Scheme == "http" && normalized.Port() == "80") ||
		(normalized.Scheme == "https" && normalized.Port() == "443") {
		normalized.Host = normalized.Hostname()
	}

	// Ensure paths end with a trailing slash if they don't have an extension
	// This helps with canonical URLs
	if normalized.Path == "" {
		normalized.Path = "/"
	}

	return normalized.String()
}

// Page represents a fetched web page
type Page struct {
	ID          uuid.UUID
	URL         string
	StatusCode  int
	Title       string
	HTML        string
	PlainText   string
	Headers     map[string]string
	Links       []string
	ContentType string
	FetchedAt   time.Time
	ParsedAt    time.Time
}

// NewPage creates a new Page entity
func NewPage(url string, statusCode int, html string, headers map[string]string) *Page {
	return &Page{
		ID:         uuid.New(),
		URL:        url,
		StatusCode: statusCode,
		HTML:       html,
		Headers:    headers,
		Links:      []string{},
		FetchedAt:  time.Now(),
	}
}

// AddLinks adds extracted links to the page
func (p *Page) AddLinks(links []string) {
	p.Links = links
	p.ParsedAt = time.Now()
}

// SetTitle sets the page title
func (p *Page) SetTitle(title string) {
	p.Title = title
	p.ParsedAt = time.Now()
}

// SetPlainText sets the page plain text content
func (p *Page) SetPlainText(text string) {
	p.PlainText = text
	p.ParsedAt = time.Now()
}

// CrawlJob represents a job to crawl a URL
type CrawlJob struct {
	URL       *URL
	CreatedAt time.Time
	Priority  int
}

// NewCrawlJob creates a new CrawlJob
func NewCrawlJob(url *URL, priority int) *CrawlJob {
	return &CrawlJob{
		URL:       url,
		CreatedAt: time.Now(),
		Priority:  priority,
	}
}
