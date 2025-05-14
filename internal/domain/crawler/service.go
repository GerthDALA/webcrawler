package crawler

import (
	"context"
	"log"
	"time"

	result "github.com/gerthdala/webcrawler/pkg/utils/result"
)

type FetcherService interface {
	// Fetch fetches a URL and returns the page content
	Fetch(ctx context.Context, url *URL) result.Result[*Page]
}

type ParserService interface {
	// Parse parses HTML content and extracts links, title, and text
	Parse(ctx context.Context, page *Page) result.Result[*Page]

	ExtractLinks(ctx context.Context, html, baseURL string) result.Result[[]string]

	ExtractTitle(ctx context.Context, html string) result.Result[string]

	// ExtractText extracts plain text from HTML content
	ExtractText(ctx context.Context, html string) result.Result[string]
}

type URLFilterService interface {
	ShouldCrawl(ctx context.Context, url string, depth int) result.Result[bool]

	IsAllowedDomain(ctx context.Context, url string) result.Result[bool]

	IsAllowedContentType(cxt context.Context, contentType string) result.Result[bool]
}

// RobotsTxtService handles robots.txt processing
type RobotsTxtService interface {
	// IsAllowed checks if a URL is allowed by robots.txt
	IsAllowed(ctx context.Context, url string, userAgent string) result.Result[bool]

	// FetchRobotsTxt fetches and parses robots.txt for a domain
	FetchRobotsTxt(ctx context.Context, domain string) result.Result[interface{}]

	// GetCrawlDelay gets the crawl delay for a domain
	GetCrawlDelay(ctx context.Context, domain string, userAgent string) result.Result[time.Duration]
}

// CrawlService orchestrates the crawling process
type CrawlService struct {
	urlRepo         URLRepository
	pageRepo        PageRepository
	crawlJobRepo    CrawlJobRepository
	fetcher         FetcherService
	parser          ParserService
	filter          URLFilterService
	robotsTxt       RobotsTxtService
	maxDepth        int
	concurrency     int
	politenessDelay time.Duration
	userAgent       string
}

// CrawlServiceConfig configuration for the crawl service
type CrawlServiceConfig struct {
	MaxDepth       int
	Concurrency    int
	PolitenessDelay time.Duration
	UserAgent      string
}

// NewCrawlService creates a new CrawlService
func NewCrawlService(
	urlRepo URLRepository,
	pageRepo PageRepository,
	crawlJobRepo CrawlJobRepository,
	fetcher FetcherService,
	parser ParserService,
	filter URLFilterService,
	robotsTxt RobotsTxtService,
	config CrawlServiceConfig,
) *CrawlService {
	return &CrawlService{
		urlRepo:        urlRepo,
		pageRepo:       pageRepo,
		crawlJobRepo:   crawlJobRepo,
		fetcher:        fetcher,
		parser:         parser,
		filter:         filter,
		robotsTxt:      robotsTxt,
		maxDepth:       config.MaxDepth,
		concurrency:    config.Concurrency,
		politenessDelay: config.PolitenessDelay,
		userAgent:      config.UserAgent,
	}
}

// AddSeed adds a seed URL to the crawler

func (s *CrawlService) AddSeed(ctx context.Context, rawURL string) result.Result[*URL] {
	// Create URL entity
	urlResult := NewURL(rawURL, 0, "")

	if urlResult.IsErr() {
		return result.Err[*URL](urlResult.Error())
	}

	url := urlResult.Unwrap()

	if existingURL := s.urlRepo.FindByNormalizedURL(ctx, url.NormalizedURL); existingURL.IsOk() {
		//URL already exists, return it
		return existingURL
	}

	saveURL := s.urlRepo.Save(ctx, url)

	if saveURL.IsErr() {
		return saveURL
	}

	// Enqueue
	job := NewCrawlJob(url, 0) // Higher priority for seeds
	s.crawlJobRepo.Enqueue(ctx, job)
	
	return saveURL
}

// ProcessedURL processes a single url
func (s *CrawlService) ProcessedURL(ctx context.Context, url *URL) result.Result[*Page] {
	// Check if the url should crawled
	shouldCrawl := s.filter.ShouldCrawl(ctx, url.URL, url.Depth)
	if shouldCrawl.IsErr() || !shouldCrawl.Unwrap() {
		if shouldCrawl.IsErr() {
			return result.Err[*Page](shouldCrawl.Error())
		}

		return result.ErrMsg[*Page]("URL should not be crawled")
	}

	// Check robots.txt
	isAllowedResult := s.robotsTxt.IsAllowed(ctx, url.URL, s.userAgent)
	if isAllowedResult.IsErr() || !isAllowedResult.Unwrap() {
		if isAllowedResult.IsErr() {
			return result.Err[*Page](isAllowedResult.Error())
		}
		return result.ErrMsg[*Page]("URL not allowed by robots.txt")
	}

	s.urlRepo.UpdateStatus(ctx, url.ID, StatusFetching)
	s.urlRepo.IncrementAttemptCount(ctx, url.ID)


	// Fetch URL status
	fetchResult := s.fetcher.Fetch(ctx, url)
	if fetchResult.IsErr() {
		s.urlRepo.UpdateStatus(ctx, url.ID, StatusFailed)
		return fetchResult
	}

	page := fetchResult.Unwrap()

	if saveResult := s.pageRepo.Save(ctx, page); saveResult.IsErr() {
		return result.Err[*Page](saveResult.Error())
	}

	//Update URL status
	s.urlRepo.UpdateStatus(ctx, url.ID, StatusFetched)

	// Process links if depth is allowed
	if url.Depth < s.maxDepth {
		s.processLinks(ctx, page.Links, url.Depth + 1, url.URL)
	}

	return result.Ok(page)
}

func (s *CrawlService) processLinks(ctx context.Context, links []string, depth int, parentURL string) {
	for _, link := range links {
		urlR := NewURL(link, depth, parentURL)

		if urlR.IsErr() {
			log.Printf("Skipping malformed URL: %s (error: %v)", link, urlR.Error())
			continue
		}

		url := urlR.Unwrap()

		// Check if URL already exists
		existingResult := s.urlRepo.FindByNormalizedURL(ctx, url.NormalizedURL)
		if existingResult.IsOk() {
			continue
		}
		
		// Save URL
		saveResult := s.urlRepo.Save(ctx, url)
		if saveResult.IsErr() {
			continue
		}
		
		// Enqueue job
		job := NewCrawlJob(url, depth) // Priority based on depth
		s.crawlJobRepo.Enqueue(ctx, job)


	}
}