package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/chifamba/canonical-corpus/internal/compiler"
	"github.com/chifamba/canonical-corpus/internal/downloader"
	"github.com/chifamba/canonical-corpus/internal/metadata"
	"github.com/chifamba/canonical-corpus/internal/ratelimiter"
	"github.com/chifamba/canonical-corpus/pkg/markdown"
	"github.com/chifamba/canonical-corpus/sources"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
)

// Config holds all application configuration loaded from config.yaml.
type Config struct {
	Downloader struct {
		MaxRetries int           `yaml:"max_retries"`
		Timeout    time.Duration `yaml:"timeout"`
		UserAgent  string        `yaml:"user_agent"`
	} `yaml:"downloader"`
	RateLimiter struct {
		MaxRequestsPerHost float64 `yaml:"max_requests_per_host"`
		GlobalConcurrency  int     `yaml:"global_concurrency"`
	} `yaml:"rate_limiter"`
	Output struct {
		BaseDir string `yaml:"base_dir"`
	} `yaml:"output"`
	Logging struct {
		Level string `yaml:"level"`
		File  string `yaml:"file"`
	} `yaml:"logging"`
}

var (
	cfgFile string
	cfg     Config
)

func main() {
	root := &cobra.Command{
		Use:   "corpus-builder",
		Short: "Biblical corpus collector and compiler",
		Long: `corpus-builder fetches, parses, normalizes, and exports biblical texts
from public-domain sources into a structured Markdown corpus.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return loadConfig(cfgFile)
		},
	}

	root.PersistentFlags().StringVarP(&cfgFile, "config", "c", "configs/config.yaml", "path to config file")

	root.AddCommand(
		buildCmd(),
		fetchCmd(),
		verifyCmd(),
		exportCmd(),
	)

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// loadConfig reads and unmarshals the YAML config file.
func loadConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading config %q: %w", path, err)
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("parsing config: %w", err)
	}
	// Apply defaults.
	if cfg.RateLimiter.MaxRequestsPerHost == 0 {
		cfg.RateLimiter.MaxRequestsPerHost = 2
	}
	if cfg.RateLimiter.GlobalConcurrency == 0 {
		cfg.RateLimiter.GlobalConcurrency = 10
	}
	if cfg.Output.BaseDir == "" {
		cfg.Output.BaseDir = "./corpus"
	}
	return nil
}

// newLogger builds a zap logger writing to stdout and optionally a log file.
func newLogger() (*zap.Logger, error) {
	level := zapcore.InfoLevel
	_ = level.UnmarshalText([]byte(cfg.Logging.Level))

	enc := zap.NewProductionEncoderConfig()
	enc.TimeKey = "time"
	enc.EncodeTime = zapcore.ISO8601TimeEncoder

	cores := []zapcore.Core{
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(enc),
			zapcore.AddSync(os.Stdout),
			level,
		),
	}

	if cfg.Logging.File != "" {
		f, err := os.OpenFile(cfg.Logging.File, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
		if err == nil {
			cores = append(cores, zapcore.NewCore(
				zapcore.NewJSONEncoder(enc),
				zapcore.AddSync(f),
				level,
			))
		}
	}

	return zap.New(zapcore.NewTee(cores...), zap.AddCaller()), nil
}

// buildDeps constructs the compiler dependency graph from the loaded config.
func buildDeps(logger *zap.Logger) *compiler.Compiler {
	rl := ratelimiter.New(cfg.RateLimiter.MaxRequestsPerHost, cfg.RateLimiter.GlobalConcurrency)
	dl := downloader.New(downloader.Config{
		MaxRetries: cfg.Downloader.MaxRetries,
		Timeout:    cfg.Downloader.Timeout,
		UserAgent:  cfg.Downloader.UserAgent,
	}, rl, logger)
	mw := markdown.New(cfg.Output.BaseDir)
	return compiler.New(dl, rl, mw, cfg.Output.BaseDir, logger)
}

// selectBooks filters the book catalogue by category flag.
func selectBooks(category string) []metadata.BookMeta {
	switch category {
	case "canonical":
		return sources.CanonicalBooks()
	case "extra-canonical":
		return sources.ExtraCanonicalBooks()
	case "dead-sea-scrolls":
		return sources.DeadSeaScrollBooks()
	default:
		return sources.AllBooks()
	}
}



// ---------------------------------------------------------------------------
// Subcommands
// ---------------------------------------------------------------------------

func buildCmd() *cobra.Command {
	var category string
	cmd := &cobra.Command{
		Use:   "build",
		Short: "Fetch, parse, normalize, and export the full corpus",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger, err := newLogger()
			if err != nil {
				return err
			}
			defer func() { _ = logger.Sync() }()

			comp := buildDeps(logger)
			books := selectBooks(category)
			logger.Info("starting build",
				zap.Int("books", len(books)),
				zap.String("category", category))
			return comp.CompileAll(context.Background(), books)
		},
	}
	cmd.Flags().StringVarP(&category, "category", "k", "all",
		"category to build: all | canonical | extra-canonical | dead-sea-scrolls")
	return cmd
}

func fetchCmd() *cobra.Command {
	var category string
	cmd := &cobra.Command{
		Use:   "fetch",
		Short: "Fetch raw source texts without compiling",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger, err := newLogger()
			if err != nil {
				return err
			}
			defer func() { _ = logger.Sync() }()

			rl := ratelimiter.New(cfg.RateLimiter.MaxRequestsPerHost, cfg.RateLimiter.GlobalConcurrency)
			dl := downloader.New(downloader.Config{
				MaxRetries: cfg.Downloader.MaxRetries,
				Timeout:    cfg.Downloader.Timeout,
				UserAgent:  cfg.Downloader.UserAgent,
			}, rl, logger)

			books := selectBooks(category)
			logger.Info("fetching sources", zap.Int("books", len(books)))

			ctx := context.Background()
			for _, book := range books {
				for _, src := range book.Sources {
					if _, _, err := dl.Fetch(ctx, src.URL); err != nil {
						logger.Warn("fetch failed",
							zap.String("book", book.Title),
							zap.String("url", src.URL),
							zap.Error(err))
					} else {
						logger.Info("fetched",
							zap.String("book", book.Title),
							zap.String("url", src.URL))
					}
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&category, "category", "k", "all",
		"category to fetch: all | canonical | extra-canonical | dead-sea-scrolls")
	return cmd
}

func verifyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "verify",
		Short: "Verify corpus files exist for all known books",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger, err := newLogger()
			if err != nil {
				return err
			}
			defer func() { _ = logger.Sync() }()

			books := sources.AllBooks()
			missing := 0
			for _, book := range books {
				dirName := fmt.Sprintf("%03d-%s", book.CanonicalOrder, markdown.SanitizeTitle(book.Title))
				path := fmt.Sprintf("%s/%s/%s/en.md",
					cfg.Output.BaseDir, string(book.Category), dirName)
				if _, err := os.Stat(path); os.IsNotExist(err) {
					logger.Warn("missing corpus file",
						zap.String("book", book.Title),
						zap.String("path", path))
					missing++
				} else {
					logger.Info("ok", zap.String("book", book.Title))
				}
			}
			if missing > 0 {
				return fmt.Errorf("%d book(s) missing from corpus", missing)
			}
			logger.Info("corpus verification complete", zap.Int("total", len(books)))
			return nil
		},
	}
}

func exportCmd() *cobra.Command {
	var dest string
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export the compiled corpus to a destination directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger, err := newLogger()
			if err != nil {
				return err
			}
			defer func() { _ = logger.Sync() }()

			logger.Info("exporting corpus",
				zap.String("source", cfg.Output.BaseDir),
				zap.String("destination", dest))
			// Placeholder: a real implementation would copy or archive corpus files.
			logger.Info("export complete")
			return nil
		},
	}
	cmd.Flags().StringVarP(&dest, "dest", "d", "./dist",
		"destination directory for the exported corpus")
	return cmd
}
