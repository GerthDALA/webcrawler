package crawler

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gerthdala/webcrawler/internal/domain/crawler"
	result "github.com/gerthdala/webcrawler/pkg/utils/result"
)

// HTMLParser implements the ParserService interface
type HTMLParser struct{}

// NewHTMLParser creates a new HTMLParser
func NewHTMLParser() *HTMLParser {
	return &HTMLParser{}
}

func (p *HTMLParser) Parse(ctx context.Context, page *crawler.Page) result.Result[*crawler.Page] {

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(page.HTML))
	if err != nil {
		return result.Err[*crawler.Page](fmt.Errorf("error creating document: %w", err))
	}
	// Extract title
	titleResult := extractTitle(doc)
	if titleResult.IsOk() {
		page.SetTitle(titleResult.Unwrap())
	}

	// Extract text
	textResult := extractText(doc)
	if textResult.IsOk() {
		page.SetPlainText(textResult.Unwrap())
	}

	// Extract links
	linksResult := extractLinks(doc, page.URL)
	if linksResult.IsOk() {
		page.AddLinks(linksResult.Unwrap())
	}
	return result.Ok(page)
}
// ExtractLinks extracts links from HTML content
func (p *HTMLParser) ExtractLinks(ctx context.Context, html, baseURL string) result.Result[[]string] {
	// Create a reader from the HTML string
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return result.Err[[]string](fmt.Errorf("error creating document: %w", err))
	}

	return extractLinks(doc, baseURL)
}

// ExtractTitle extracts the title from HTML content
func (p *HTMLParser) ExtractTitle(ctx context.Context, html string) result.Result[string] {
	// Create a reader from the HTML string
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return result.Err[string](fmt.Errorf("error creating document: %w", err))
	}


	return extractTitle(doc)
}

// ExtractText extracts plain text from HTML content
func (p *HTMLParser) ExtractText(ctx context.Context, html string) result.Result[string] {
	// Create a reader from the HTML string
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return result.Err[string](fmt.Errorf("error creating document: %w", err))
	}

	return extractText(doc)
}

func extractLinks(doc *goquery.Document, baseURL string) result.Result[[]string] {
	// Parse the base URL
	base, err := url.Parse(baseURL)
	if err != nil {
		return result.Err[[]string](fmt.Errorf("error parsing base URL: %w", err))
	}

	//Look for base tag
	baseHref, exists := doc.Find("base[href]").Attr("href")
	if exists {
		baseHrefURL, err := url.Parse(baseHref)
		if err == nil {
			base = base.ResolveReference(baseHrefURL)
		}
	}

	links := make([]string, 0)
	doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists || href == "" || strings.HasPrefix(href, "#") {
			return
		}

		// Parse href
		u, err := url.Parse(href)
		if err != nil {
			return
		}

		resolved := base.ResolveReference(u)
		if resolved.Scheme == "http" || resolved.Scheme == "https" {
			links = append(links, resolved.String())
		}
	})
	return result.Ok(links)
}

// ExtractTitle extracts the title from HTML content
func extractTitle(doc *goquery.Document) result.Result[string] {
	title := doc.Find("title").First().Text()
	title = strings.TrimSpace(title)

	if title == "" {
		title = doc.Find("h1, h2, h3, h4, h5, h6").First().Text()
		title = strings.TrimSpace(title)
	}
	return result.Ok(title)
}

// ExtractText extracts plain text from HTML content
func extractText(doc *goquery.Document) result.Result[string] {
	// Remove script and style elements
	doc.Find("script, style, noscript, iframe, object, embed").Remove()
	text := doc.Text()
	text = normalizeWhitespace(text)
	return result.Ok(text)
}
func normalizeWhitespace(text string) string {
	// Replace line breaks and tabs with spaces
	text = strings.ReplaceAll(text, "\n", "")
	text = strings.ReplaceAll(text, "\r", " ")
	text = strings.ReplaceAll(text, "\t", " ")

	// Replace multiple spaces with a single space
	for strings.Contains(text, "  ") {
		text = strings.ReplaceAll(text, "  ", " ")
	}

	return strings.TrimSpace(text)
}