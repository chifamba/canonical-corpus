package downloader

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/chifamba/canonical-corpus/internal/ratelimiter"
	"github.com/temoto/robotstxt"
	"go.uber.org/zap"
)

const defaultUserAgent = "CorpusBuilder/1.0 (+https://github.com/chifamba/canonical-corpus)"

// Config holds downloader configuration.
type Config struct {
	MaxRetries       int
	Timeout          time.Duration
	UserAgent        string
	BlacklistedHosts []string
}

// Downloader handles HTTP downloads with robots.txt compliance and rate limiting.
type Downloader struct {
	config      Config
	client      *http.Client
	robots      map[string]*robotstxt.RobotsData
	robotsMu    sync.RWMutex
	rateLimiter *ratelimiter.RateLimiter
	logger      *zap.Logger
}

// New creates a new Downloader.
func New(cfg Config, rl *ratelimiter.RateLimiter, logger *zap.Logger) *Downloader {
	if cfg.UserAgent == "" {
		cfg.UserAgent = defaultUserAgent
	}
	if cfg.MaxRetries <= 0 {
		cfg.MaxRetries = 5
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = 30 * time.Second
	}
	return &Downloader{
		config:      cfg,
		client:      &http.Client{Timeout: cfg.Timeout},
		robots:      make(map[string]*robotstxt.RobotsData),
		rateLimiter: rl,
		logger:      logger,
	}
}

// Fetch downloads a URL, respecting robots.txt and rate limits.
// Returns body bytes, content-type, and any error.
func (d *Downloader) Fetch(ctx context.Context, rawURL string) ([]byte, string, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return nil, "", fmt.Errorf("invalid URL %q: %w", rawURL, err)
	}
	host := parsed.Hostname()

	// Skip URLs from blacklisted hosts.
	for _, b := range d.config.BlacklistedHosts {
		if host == b {
			return nil, "", fmt.Errorf("host %q is blacklisted", host)
		}
	}

	allowed, err := d.checkRobots(ctx, rawURL)
	if err != nil {
		d.logger.Warn("robots.txt check failed, proceeding",
			zap.String("url", rawURL), zap.Error(err))
	} else if !allowed {
		return nil, "", fmt.Errorf("URL %q disallowed by robots.txt", rawURL)
	}

	var (
		body        []byte
		contentType string
		fetchErr    error
	)

	for attempt := 0; attempt <= d.config.MaxRetries; attempt++ {
		if attempt > 0 {
			backoff := time.Duration(1<<uint(attempt-1))*time.Second +
				time.Duration(rand.IntN(1000))*time.Millisecond

			// If last error was 429 with Retry-After, use that instead.
			var httpErr *httpError
			if fetchErr != nil && errors.As(fetchErr, &httpErr) && httpErr.RetryAfter != "" {
				if s, err := strconv.Atoi(httpErr.RetryAfter); err == nil {
					backoff = time.Duration(s) * time.Second
				}
			}

			select {
			case <-time.After(backoff):
			case <-ctx.Done():
				return nil, "", ctx.Err()
			}
		}

		if err := d.rateLimiter.Wait(ctx, host); err != nil {
			return nil, "", fmt.Errorf("rate limiter: %w", err)
		}

		body, contentType, fetchErr = d.doFetch(ctx, rawURL)
		d.rateLimiter.Release()

		if fetchErr == nil {
			return body, contentType, nil
		}

		if !isRetryable(fetchErr) {
			return nil, "", fetchErr
		}
		d.logger.Warn("retrying fetch",
			zap.String("url", rawURL),
			zap.Int("attempt", attempt+1),
			zap.Error(fetchErr))
	}
	return nil, "", fmt.Errorf("max retries exceeded for %q: %w", rawURL, fetchErr)
}

func (d *Downloader) doFetch(ctx context.Context, rawURL string) ([]byte, string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("User-Agent", d.config.UserAgent)

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		retryAfter := resp.Header.Get("Retry-After")
		return nil, "", &httpError{StatusCode: resp.StatusCode, RetryAfter: retryAfter}
	}

	if resp.StatusCode >= 500 {
		return nil, "", &httpError{StatusCode: resp.StatusCode}
	}
	if resp.StatusCode >= 400 {
		return nil, "", fmt.Errorf("HTTP %d for %q", resp.StatusCode, rawURL)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("reading body: %w", err)
	}

	return data, resp.Header.Get("Content-Type"), nil
}

// checkRobots checks if the URL is allowed by robots.txt.
func (d *Downloader) checkRobots(ctx context.Context, rawURL string) (bool, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return false, err
	}
	host := parsed.Host

	d.robotsMu.RLock()
	rd, ok := d.robots[host]
	d.robotsMu.RUnlock()

	if !ok {
		robotsURL := fmt.Sprintf("%s://%s/robots.txt", parsed.Scheme, host)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, robotsURL, nil)
		if err != nil {
			return true, nil
		}
		req.Header.Set("User-Agent", d.config.UserAgent)

		resp, err := d.client.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			if resp != nil && resp.Body != nil {
				resp.Body.Close()
			}
			d.robotsMu.Lock()
			d.robots[host] = nil
			d.robotsMu.Unlock()
			return true, nil
		}
		defer resp.Body.Close()

		rd, err = robotstxt.FromResponse(resp)
		if err != nil {
			d.robotsMu.Lock()
			d.robots[host] = nil
			d.robotsMu.Unlock()
			return true, nil
		}
		d.robotsMu.Lock()
		d.robots[host] = rd
		d.robotsMu.Unlock()
	}

	if rd == nil {
		return true, nil
	}
	group := rd.FindGroup(d.config.UserAgent)
	return group.Test(parsed.Path), nil
}

type httpError struct {
	StatusCode int
	RetryAfter string
}

func (e *httpError) Error() string {
	return fmt.Sprintf("HTTP %d", e.StatusCode)
}

func isRetryable(err error) bool {
	if err == nil {
		return false
	}
	var httpErr *httpError
	if errors.As(err, &httpErr) {
		if httpErr.StatusCode >= 500 || httpErr.StatusCode == http.StatusTooManyRequests {
			return true
		}
	}
	var netErr net.Error
	if errors.As(err, &netErr) {
		return netErr.Timeout()
	}
	if errors.Is(err, io.EOF) {
		return true
	}
	msg := err.Error()
	return strings.Contains(msg, "connection reset")
}
