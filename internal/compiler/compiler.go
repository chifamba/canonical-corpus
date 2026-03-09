package compiler

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chifamba/canonical-corpus/internal/downloader"
	"github.com/chifamba/canonical-corpus/internal/metadata"
	"github.com/chifamba/canonical-corpus/internal/normalizer"
	"github.com/chifamba/canonical-corpus/internal/parser"
	"github.com/chifamba/canonical-corpus/internal/progress"
	"github.com/chifamba/canonical-corpus/internal/ratelimiter"
	"github.com/chifamba/canonical-corpus/pkg/markdown"
	"github.com/chifamba/canonical-corpus/sources"
	"go.uber.org/zap"
)

// Compiler orchestrates fetching, parsing, normalizing, and writing corpus documents.
type Compiler struct {
	downloader     *downloader.Downloader
	rateLimiter    *ratelimiter.RateLimiter
	markdownWriter *markdown.Writer
	logger         *zap.Logger
	outputDir      string
	// progress tracks which books have already been compiled (for resumable builds).
	progress *progress.State
	// force re-compiles books that are already marked complete.
	force bool
}

// New creates a new Compiler.
// Pass force=true to re-compile books that are already marked as complete.
func New(dl *downloader.Downloader, rl *ratelimiter.RateLimiter, mw *markdown.Writer, outputDir string, logger *zap.Logger, force bool) *Compiler {
	prog, err := progress.Load(outputDir)
	if err != nil {
		logger.Warn("could not load progress state; falling back to in-memory tracking", zap.Error(err))
		prog, _ = progress.Load("") // empty baseDir → in-memory only, never fails
	}
	return &Compiler{
		downloader:     dl,
	rateLimiter:    rl,
	markdownWriter: mw,
		logger:         logger,
		outputDir:      outputDir,
		progress:       prog,
		force:          force,
	}
}

// progressKey returns the unique key used to track completion of a book.
func progressKey(book metadata.BookMeta) string {
	dirName := fmt.Sprintf("%03d-%s", book.CanonicalOrder, markdown.SanitizeTitle(book.Title))
	filename := markdown.BuildFilename(book.Language, book.TranslationID)
	return fmt.Sprintf("%s/%s/%s", string(book.Category), dirName, filename)
}

// ImportLocalBibles recursively scans inputDir for JSON bibles and imports them.
func (c *Compiler) ImportLocalBibles(ctx context.Context, inputDir string) error {
	if _, err := os.Stat(inputDir); os.IsNotExist(err) {
		return nil
	}

	return filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		ext := filepath.Ext(path)
		if info.IsDir() || (ext != ".json" && ext != ".xml") {
			return nil
		}

		c.logger.Info("importing local bible file", zap.String("path", path))

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading local bible %q: %w", path, err)
		}

		var books []parser.ParsedBook
		var lang string

		if filepath.Ext(path) == ".json" {
			books, lang, err = parser.ParseExportedJSON(data)
		} else if filepath.Ext(path) == ".xml" {
			books, lang, err = parser.ParseBibleXML(data)
		}

		if err != nil || len(books) == 0 {
			// Not a valid bible file, skip or log.
			c.logger.Warn("skipping non-bible file", zap.String("path", path), zap.Error(err))
			return nil
		}

		translationID := strings.ToLower(filepath.Base(path[:len(path)-len(filepath.Ext(path))]))

		// Try to fix language if it's empty or too long.
		if len(lang) > 3 || lang == "" {
			// Extract from path if possible, e.g. input/bibles_json/EN-English/...
			parts := strings.Split(path, string(os.PathSeparator))
			for _, p := range parts {
				// Check for EN-English style.
				if strings.Contains(p, "-") {
					code := strings.Split(p, "-")[0]
					if len(code) >= 2 && len(code) <= 3 {
						lang = strings.ToLower(code)
						break
					}
				}
			}
			if lang == "" {
				// Check parent directory name, e.g. input/shona/sna1949.xml
				parent := filepath.Base(filepath.Dir(path))
				if parent == "shona" || parent == "sn" {
					lang = "sn"
				}
			}
		}

		for _, b := range books {
			c.logger.Info("importing book from local file", 
				zap.String("title", b.Title), 
				zap.Int("order", b.Number),
				zap.String("translation", translationID),
				zap.String("lang", lang))
			meta := metadata.BookMeta{
				Title:          b.Title,
				CanonicalOrder: b.Number,
				Category:       metadata.CategoryCanonical,
				Language:       lang,
				TranslationID:  translationID,
				License:        "Public Domain / Freely Licensed",
				DateCollected:  info.ModTime().UTC(),
			}

			// Try to find English title for canonical books to normalize directories.
			if canonical, ok := sources.FindBookByOrder(b.Number); ok {
				meta.Title = canonical.Title
			}

			normalized := normalizer.Normalize(b.Content)
			normalized = normalizer.DeduplicatePassages(normalized)

			doc := &metadata.Document{
				Meta:    meta,
				Content: normalized,
			}

			key := progressKey(meta)
			if err := c.markdownWriter.Write(doc); err != nil {
				return fmt.Errorf("writing imported book %q from %q: %w", b.Title, path, err)
			}
			_ = c.progress.MarkComplete(key)
		}

		return nil
	})
}

// CompileBook fetches all sources for a book, parses, normalizes, and writes output.
// It skips the book if it has already been successfully compiled (unless force=true).
func (c *Compiler) CompileBook(ctx context.Context, book metadata.BookMeta) error {
	key := progressKey(book)

	if !c.force && c.progress.IsComplete(key) {
		c.logger.Info("skipping already-compiled book",
			zap.String("title", book.Title),
			zap.String("translation", book.TranslationID))
		return nil
	}

	c.logger.Info("compiling book",
		zap.String("title", book.Title),
		zap.String("translation", book.TranslationID))

	var combinedContent strings.Builder
	for _, src := range book.Sources {
		data, ct, err := c.downloader.Fetch(ctx, src.URL)
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

	// Persist completion so a future interrupted run can resume.
	if err := c.progress.MarkComplete(key); err != nil {
		c.logger.Warn("could not save progress state", zap.Error(err))
	}

	c.logger.Info("compiled book",
		zap.String("title", book.Title),
		zap.String("translation", book.TranslationID))
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
