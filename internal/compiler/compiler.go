package compiler

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/chifamba/canonical-corpus/internal/downloader"
	"github.com/chifamba/canonical-corpus/internal/metadata"
	"github.com/chifamba/canonical-corpus/internal/normalizer"
	"github.com/chifamba/canonical-corpus/internal/parser"
	"github.com/chifamba/canonical-corpus/internal/ratelimiter"
	"github.com/chifamba/canonical-corpus/pkg/markdown"
	"go.uber.org/zap"
)

// Compiler orchestrates fetching, parsing, normalizing, and writing corpus documents.
type Compiler struct {
	downloader     *downloader.Downloader
	rateLimiter    *ratelimiter.RateLimiter
	markdownWriter *markdown.Writer
	logger         *zap.Logger
	outputDir      string
}

// New creates a new Compiler.
func New(dl *downloader.Downloader, rl *ratelimiter.RateLimiter, mw *markdown.Writer, outputDir string, logger *zap.Logger) *Compiler {
	return &Compiler{
		downloader:     dl,
		rateLimiter:    rl,
		markdownWriter: mw,
		logger:         logger,
		outputDir:      outputDir,
	}
}

// CompileBook fetches all sources for a book, parses, normalizes, and writes output.
func (c *Compiler) CompileBook(ctx context.Context, book metadata.BookMeta) error {
	c.logger.Info("compiling book", zap.String("title", book.Title))

	var combinedContent strings.Builder
	for _, src := range book.Sources {
		data, ct, err := c.downloader.FetchSource(ctx, src)
		if err != nil {
			c.logger.Warn("failed to fetch source",
				zap.String("book", book.Title),
				zap.String("url", src.URL),
				zap.Error(err))
			continue
		}

		parsed, err := parser.Parse(data, ct, src.Format)
		if err != nil {
			c.logger.Warn("failed to parse source",
				zap.String("book", book.Title),
				zap.String("url", src.URL),
				zap.Error(err))
			continue
		}

		normalized := normalizer.Normalize(parsed.Content)
		normalized = normalizer.DeduplicatePassages(normalized)

		if combinedContent.Len() > 0 {
			combinedContent.WriteString("\n\n")
		}
		combinedContent.WriteString(normalized)
	}

	if combinedContent.Len() == 0 {
		return fmt.Errorf("no content collected for %q", book.Title)
	}

	book.DateCollected = time.Now().UTC()

	doc := &metadata.Document{
		Meta:    book,
		Content: combinedContent.String(),
	}

	if err := c.markdownWriter.Write(doc); err != nil {
		return fmt.Errorf("writing markdown for %q: %w", book.Title, err)
	}

	c.logger.Info("compiled book", zap.String("title", book.Title))
	return nil
}

// CompileAll compiles every book in the list, collecting all errors.
func (c *Compiler) CompileAll(ctx context.Context, books []metadata.BookMeta) error {
	var failed int
	for _, book := range books {
		if err := c.CompileBook(ctx, book); err != nil {
			c.logger.Error("compile book failed",
				zap.String("title", book.Title),
				zap.Error(err))
			failed++
		}
	}
	if failed > 0 {
		return fmt.Errorf("%d book(s) failed to compile", failed)
	}
	return nil
}
