package markdown

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chifamba/canonical-corpus/internal/metadata"
	"gopkg.in/yaml.v3"
)

// Writer writes documents as Markdown files with YAML front matter.
type Writer struct {
	baseDir string
}

// New creates a new Writer rooted at baseDir.
func New(baseDir string) *Writer {
	return &Writer{baseDir: baseDir}
}

// Write persists a document as <baseDir>/<category>/<NNN-title>/en.md.
func (w *Writer) Write(doc *metadata.Document) error {
	dirName := fmt.Sprintf("%03d-%s", doc.Meta.CanonicalOrder, SanitizeTitle(doc.Meta.Title))
	dir := filepath.Join(w.baseDir, string(doc.Meta.Category), dirName)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating directory %q: %w", dir, err)
	}

	outPath := filepath.Join(dir, "en.md")
	content := BuildMarkdown(doc)
	if err := os.WriteFile(outPath, []byte(content), 0o644); err != nil {
		return fmt.Errorf("writing file %q: %w", outPath, err)
	}
	return nil
}

// SanitizeTitle converts a book title into a filesystem-safe slug.
func SanitizeTitle(title string) string {
	title = strings.ToLower(title)
	title = strings.ReplaceAll(title, " ", "-")
	var b strings.Builder
	for _, r := range title {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

type frontmatterSource struct {
	URL    string `yaml:"url"`
	Format string `yaml:"format"`
}

type frontmatter struct {
	Title          string              `yaml:"title"`
	CanonicalOrder int                 `yaml:"canonical_order"`
	Category       metadata.Category   `yaml:"category"`
	Language       string              `yaml:"language"`
	Sources        []frontmatterSource `yaml:"sources"`
	DateCollected  string              `yaml:"date_collected,omitempty"`
	License        string              `yaml:"license"`
	Translator     string              `yaml:"translator,omitempty"`
}

// FormatFrontmatter returns the YAML front matter string for a document.
func FormatFrontmatter(meta metadata.BookMeta) string {
	sources := make([]frontmatterSource, 0, len(meta.Sources))
	for _, s := range meta.Sources {
		sources = append(sources, frontmatterSource{URL: s.URL, Format: s.Format})
	}

	fm := frontmatter{
		Title:          meta.Title,
		CanonicalOrder: meta.CanonicalOrder,
		Category:       meta.Category,
		Language:       meta.Language,
		Sources:        sources,
		License:        meta.License,
		Translator:     meta.Translator,
	}
	if !meta.DateCollected.IsZero() {
		fm.DateCollected = meta.DateCollected.Format(time.RFC3339)
	}

	data, err := yaml.Marshal(fm)
	if err != nil {
		return ""
	}
	return string(data)
}

// BuildMarkdown assembles the full Markdown content for a document.
func BuildMarkdown(doc *metadata.Document) string {
	var sb strings.Builder

	sb.WriteString("---\n")
	sb.WriteString(FormatFrontmatter(doc.Meta))
	sb.WriteString("---\n\n")

	sb.WriteString(fmt.Sprintf("# %s\n\n", doc.Meta.Title))

	sb.WriteString("## Metadata\n\n")
	sb.WriteString(fmt.Sprintf("- **Language**: %s\n", doc.Meta.Language))
	if len(doc.Meta.Sources) > 0 {
		sb.WriteString(fmt.Sprintf("- **Original Source**: %s\n", doc.Meta.Sources[0].URL))
	}
	translator := doc.Meta.Translator
	if translator == "" {
		translator = "Unknown"
	}
	sb.WriteString(fmt.Sprintf("- **Translator**: %s\n", translator))

	sourceURLs := make([]string, 0, len(doc.Meta.Sources))
	for _, s := range doc.Meta.Sources {
		sourceURLs = append(sourceURLs, s.URL)
	}
	sb.WriteString(fmt.Sprintf("- **Compilation Sources**: %s\n\n", strings.Join(sourceURLs, ", ")))

	sb.WriteString("## Text\n\n")
	sb.WriteString(doc.Content)
	sb.WriteString("\n\n")

	sb.WriteString("## Source References\n\n")
	for _, s := range doc.Meta.Sources {
		sb.WriteString(fmt.Sprintf("- **URL**: %s\n", s.URL))
		sb.WriteString(fmt.Sprintf("  - **Format**: %s\n", strings.ToUpper(s.Format)))
		if s.Notes != "" {
			sb.WriteString(fmt.Sprintf("  - **Notes**: %s\n", s.Notes))
		}
		if s.License != "" {
			sb.WriteString(fmt.Sprintf("  - **License**: %s\n", s.License))
		}
	}

	return sb.String()
}
