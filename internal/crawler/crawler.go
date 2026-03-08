package crawler

import (
	"bytes"
	"context"
	"net/url"
	"strings"
	"sync"

	"github.com/chifamba/canonical-corpus/internal/downloader"
	"go.uber.org/zap"
	"golang.org/x/net/html"
)

// Crawler discovers linked document URLs from seed URLs.
type Crawler struct {
	downloader *downloader.Downloader
	logger     *zap.Logger
	visited    map[string]bool
	mu         sync.Mutex
}

// New creates a new Crawler.
func New(dl *downloader.Downloader, logger *zap.Logger) *Crawler {
	return &Crawler{
		downloader: dl,
		logger:     logger,
		visited:    make(map[string]bool),
	}
}

// Discover fetches seedURL and, if depth > 0, recursively follows same-host links.
// Returns the deduplicated set of discovered URLs.
func (c *Crawler) Discover(ctx context.Context, seedURL string, depth int) ([]string, error) {
	if depth < 0 {
		return nil, nil
	}

	c.mu.Lock()
	if c.visited[seedURL] {
		c.mu.Unlock()
		return nil, nil
	}
	c.visited[seedURL] = true
	c.mu.Unlock()

	data, ct, err := c.downloader.Fetch(ctx, seedURL)
	if err != nil {
		return nil, err
	}

	base, err := url.Parse(seedURL)
	if err != nil {
		return nil, err
	}

	links := extractLinks(data, ct, base)
	discovered := []string{seedURL}

	if depth > 0 {
		for _, link := range links {
			sub, err := c.Discover(ctx, link, depth-1)
			if err != nil {
				c.logger.Warn("discover error", zap.String("url", link), zap.Error(err))
				continue
			}
			discovered = append(discovered, sub...)
		}
	}

	return discovered, nil
}

// extractLinks parses an HTML document and returns resolved same-host href links.
func extractLinks(data []byte, contentType string, base *url.URL) []string {
	if !strings.Contains(strings.ToLower(contentType), "html") {
		return nil
	}

	doc, err := html.Parse(bytes.NewReader(data))
	if err != nil {
		return nil
	}

	var links []string
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					if resolved := resolveLink(attr.Val, base); resolved != "" {
						links = append(links, resolved)
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)
	return links
}

// resolveLink resolves href relative to base, filtering out non-HTTP and cross-host links.
func resolveLink(href string, base *url.URL) string {
	if strings.HasPrefix(href, "#") ||
		strings.HasPrefix(href, "javascript:") ||
		strings.HasPrefix(href, "mailto:") ||
		strings.HasPrefix(href, "data:") ||
		strings.HasPrefix(href, "vbscript:") {
		return ""
	}
	ref, err := url.Parse(href)
	if err != nil {
		return ""
	}
	resolved := base.ResolveReference(ref)
	if resolved.Host != base.Host {
		return ""
	}
	resolved.Fragment = ""
	return resolved.String()
}
