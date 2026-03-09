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
	URL        string `yaml:"url"`
	Format     string `yaml:"format"` // html, txt, xml, json, pdf, epub
	Language   string `yaml:"language"`
	Notes      string `yaml:"notes,omitempty"`
	Archive    string `yaml:"archive,omitempty"`
	License    string `yaml:"license,omitempty"`
	Translator string `yaml:"translator,omitempty"`
}

// BookMeta represents metadata for a single book.
type BookMeta struct {
	Title          string      `yaml:"title"`
	CanonicalOrder int         `yaml:"canonical_order"`
	Category       Category    `yaml:"category"`
	Language       string      `yaml:"language"`
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
