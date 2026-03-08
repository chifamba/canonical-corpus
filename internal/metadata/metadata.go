package metadata

import "time"

// Category represents a corpus category.
type Category string

const (
	CategoryCanonical      Category = "canonical"
	CategoryExtraCanonical Category = "extra-canonical"
	CategoryDeadSeaScrolls Category = "dead-sea-scrolls"
)

// SourceRef represents a source reference for a book.
type SourceRef struct {
	URL         string `yaml:"url"`
	Format      string `yaml:"format"` // html, txt, xml, json, pdf, epub
	Language    string `yaml:"language"`
	Translation string `yaml:"translation,omitempty"` // e.g. ESV, NIV, WLC, LXX, SBLGNT, BSZ, AB
	Notes       string `yaml:"notes,omitempty"`
	Archive     string `yaml:"archive,omitempty"`
	License     string `yaml:"license,omitempty"`
	Translator  string `yaml:"translator,omitempty"`
	// AuthHeader and AuthEnvVar support API-key-authenticated sources.
	// AuthHeader is the HTTP header name (e.g. "Authorization").
	// AuthEnvVar is the environment-variable name that holds the key value.
	// AuthPrefix is an optional prefix prepended to the env-var value (e.g. "Token ").
	// These fields are intentionally excluded from YAML to avoid leaking secrets.
	AuthHeader string `yaml:"-"`
	AuthEnvVar string `yaml:"-"`
	AuthPrefix string `yaml:"-"`
}

// BookMeta represents metadata for a single book.
type BookMeta struct {
	Title          string      `yaml:"title"`
	CanonicalOrder int         `yaml:"canonical_order"`
	Category       Category    `yaml:"category"`
	Language       string      `yaml:"language"`
	Translation    string      `yaml:"translation,omitempty"` // primary translation abbreviation
	Sources        []SourceRef `yaml:"sources"`
	DateCollected  time.Time   `yaml:"date_collected,omitempty"`
	License        string      `yaml:"license"`
	Translator     string      `yaml:"translator,omitempty"`
	Notes          string      `yaml:"notes,omitempty"`
}

// Document represents a fully collected and normalized document.
type Document struct {
	Meta    BookMeta
	Content string
	Verses  []Verse
}

// Verse represents a single verse.
type Verse struct {
	Book    string
	Chapter int
	Verse   int
	Text    string
}
