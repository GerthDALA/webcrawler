package crawler

import (
	"context"
	"net/url"
	"regexp"
	"strings"
	"sync"

	result "github.com/gerthdala/webcrawler/pkg/utils/result"
)

// URLFilter implements the URLFilterService interface
type URLFilter struct {
	allowedDomains    []string
	allowedExtensions []string
	disallowedPaths   []string
	allowedContentTypes []string
	maxURLLength      int
	mu                sync.RWMutex
}

// URLFilterConfig configuration for the URL filter
type URLFilterConfig struct {
	AllowedDomains    []string
	AllowedExtensions []string
	DisallowedPaths   []string
	AllowedContentTypes []string
	MaxURLLength      int
}

// NewURLFilter creates a new URLFilter
func NewURLFilter(config URLFilterConfig) *URLFilter {
	return &URLFilter{
		allowedDomains:    config.AllowedDomains,
		allowedExtensions: config.AllowedExtensions,
		disallowedPaths:   config.DisallowedPaths,
		allowedContentTypes: config.AllowedContentTypes,
		maxURLLength:      config.MaxURLLength,
	}
}

func (f *URLFilter) ShouldCrawl(ctx context.Context, urlStr string, depth int) result.Result[bool] {
	f.mu.RLock()
	defer f.mu.RUnlock()

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return result.Err[bool](err)
	}
	
	if f.maxURLLength > 0 && len(urlStr) > f.maxURLLength {
		return result.Ok[bool](false)
	}

	if len(f.allowedDomains) > 0 {
		domainAllowed := false
		for _, domain := range f.allowedDomains {
			if strings.HasSuffix(parsedURL.Hostname(), domain) {
				domainAllowed = true
				break
			}
		}
		if !domainAllowed {
			return result.Ok[bool](false)
		}
	}

	if len(f.allowedExtensions) > 0 {
		ext := getExtension(parsedURL.Path)
		if ext != "" {
			extAllowed := false
			for _, allowedExt := range f.allowedExtensions {
				if strings.EqualFold(ext, allowedExt) {
					extAllowed = true
					break
				}
			}
			if !extAllowed {
				return result.Ok[bool](false)
			}
		}
	}

	for _, pathPattern := range f.disallowedPaths {
		matched, err := regexp.MatchString(pathPattern, parsedURL.Path)
		if err != nil {
			continue
		}
		if matched {
			return result.Ok[bool](false)
		}
	}

	return result.Ok[bool](true)
}

func (f *URLFilter) IsAllowedDomain(ctx context.Context, urlStr string ) result.Result[bool] {
	f.mu.RLock()
	defer f.mu.RUnlock()

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return result.Err[bool](err)
	}

	if len(f.allowedDomains) == 0 {
		return result.Ok[bool](true) // No restrictions on domains
	}

	for _, domain := range f.allowedDomains {
		if strings.HasSuffix(parsedURL.Hostname(), domain) {
			return result.Ok[bool](true)
		}
	}

	return result.Ok[bool](false)
}

// IsAllowedContentType checks if a content type is allowed
func (f *URLFilter) IsAllowedContentType(ctx context.Context, contentType string) result.Result[bool] {
	f.mu.RLock()
	defer f.mu.RUnlock()

	// If no allowed content types specified, allow all
	if len(f.allowedContentTypes) == 0 {
		return result.Ok(true)
	}
	mainType := contentType
	if idx := strings.Index(contentType, ";"); idx > 0 {
		mainType = contentType[:idx]
	}

	mainType = strings.TrimSpace(mainType)
	for _, allowedType := range f.allowedContentTypes {
		if strings.HasPrefix(mainType, allowedType) {
			return result.Ok[bool](true)
		}
	}
	return result.Ok[bool](false)
}

func getExtension(path string) string {
	idx := strings.LastIndex(path, ".")
	if idx < 0 || idx == len(path)-1 {
		return ""
	}
	return path[idx+1:]
}