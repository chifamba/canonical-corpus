package normalizer

import (
	"crypto/sha256"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

var (
	multiSpace    = regexp.MustCompile(`[ \t]+`)
	multiNewline  = regexp.MustCompile(`\n{3,}`)
	controlChars  = regexp.MustCompile(`[\x00-\x08\x0B\x0C\x0E-\x1F\x7F]`)
	windowsCRLF   = regexp.MustCompile(`\r\n`)
	trailingSpace = regexp.MustCompile(`[ \t]+\n`)
)

// Normalize cleans and normalizes text for corpus storage.
func Normalize(text string) string {
	// Normalize Unicode to NFC.
	text = norm.NFC.String(text)
	// Normalize line endings.
	text = windowsCRLF.ReplaceAllString(text, "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	// Remove control characters (preserve \t and \n).
	text = controlChars.ReplaceAllString(text, "")
	// Remove trailing whitespace from lines.
	text = trailingSpace.ReplaceAllString(text, "\n")
	// Collapse multiple spaces/tabs within lines.
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(multiSpace.ReplaceAllString(line, " "), " \t")
	}
	text = strings.Join(lines, "\n")
	// Cap consecutive blank lines at one.
	text = multiNewline.ReplaceAllString(text, "\n\n")
	return strings.TrimSpace(text)
}

// DeduplicatePassages removes duplicate paragraphs based on content fingerprint.
func DeduplicatePassages(text string) string {
	paragraphs := strings.Split(text, "\n\n")
	seen := make(map[string]bool)
	result := make([]string, 0, len(paragraphs))
	for _, p := range paragraphs {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		fp := fingerprint(p)
		if seen[fp] {
			continue
		}
		seen[fp] = true
		result = append(result, p)
	}
	return strings.Join(result, "\n\n")
}

// fingerprint produces a canonical hash of a paragraph for dedup comparison.
func fingerprint(s string) string {
	s = strings.ToLower(s)
	s = multiSpace.ReplaceAllString(s, " ")
	var b strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
			b.WriteRune(r)
		}
	}
	sum := sha256.Sum256([]byte(strings.TrimSpace(b.String())))
	return fmt.Sprintf("%x", sum)
}
