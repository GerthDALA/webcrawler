package crawler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gerthdala/webcrawler/internal/domain/crawler"
	result "github.com/gerthdala/webcrawler/pkg/utils/result"
)

// HTTPFetcher implements the FetcherService interface
type HTTPFetcher struct {
	client    *http.Client
	userAgent string
	timeout   time.Duration
}

type HTTPFetcherConfig struct {
	UserAgent       string
	Timout          time.Duration
	MaxRedirects    int
	FollowRedirects bool
}

func NewHTTPFetcher(config HTTPFetcherConfig) *HTTPFetcher {
	// Create rederct policy
	redirectPolicy := func(req *http.Request, via []*http.Request) error {
		if config.FollowRedirects {
			return http.ErrUseLastResponse
		}

		if len(via) >= config.MaxRedirects {
			return fmt.Errorf("exceeded max redirects: %d", config.MaxRedirects)
		}
		// copy the original headers to the redirect
		for key, val := range via[0].Header {
			req.Header[key] = val
		}

		return nil
	}

	client := &http.Client{
		Timeout:       config.Timout,
		CheckRedirect: redirectPolicy,
	}

	return &HTTPFetcher{
		client:    client,
		userAgent: config.UserAgent,
		timeout:   config.Timout,
	}
}

// Fetch performs an HTTP GET request to the specified URL
func (f *HTTPFetcher) Fetch(ctx context.Context, url *crawler.URL) result.Result[*crawler.Page] {
	// Create a request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.URL, nil)
	if err != nil {
		return result.Err[*crawler.Page](fmt.Errorf("error creating request: %w", err))
	}

	// Set headers
	req.Header.Set("User-Agent", f.userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Cache-Control", "max-age=0")

	// Execute request
	resp, err := f.client.Do(req)
	if err != nil {
		return result.Err[*crawler.Page](fmt.Errorf("error fetching URL %s: %v", url.URL, err))
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result.Err[*crawler.Page](fmt.Errorf("error reading response body for URL %s: %v", url.URL, err))
	}

	//Extracts headers
	headers := make(map[string]string)
	for name, values := range resp.Header {
		if len(values) > 0 {
			headers[name] = values[0] // Use the first value for simplicity
		}
	}

	// Create a Page
	page := crawler.NewPage(url.URL, resp.StatusCode, string(body), headers)
	page.ContentType = resp.Header.Get("Content-Type")
	
	return result.Ok(page)
}