package sources_test

import (
"testing"

"github.com/chifamba/canonical-corpus/sources"
)

// Catalogue size constants — keep in sync with sources.go declarations.
const (
canonicalBookCount      = 66 // Protestant 66-book canon (OT+NT)
extraCanonicalBookCount = 25 // Deuterocanonical / Ethiopian Orthodox
deadSeaScrollCount      = 9  // Major DSS texts
)

func TestAllTranslationsPresent(t *testing.T) {
ids := sources.AllTranslationIDs()
expected := []string{"kjv", "web", "asv", "lxx", "hmt", "bdb"}
found := make(map[string]bool)
for _, id := range ids {
found[id] = true
}
for _, want := range expected {
if !found[want] {
t.Errorf("missing translation ID: %q", want)
}
}
}

func TestAllLanguagesPresent(t *testing.T) {
codes := sources.AllLanguageCodes()
expected := []string{"en", "el", "he", "sn"}
found := make(map[string]bool)
for _, c := range codes {
found[c] = true
}
for _, want := range expected {
if !found[want] {
t.Errorf("missing language code: %q", want)
}
}
}

func TestBooksByTranslation(t *testing.T) {
tests := []struct {
id    string
count int // expected total books across all categories
}{
// KJV covers all three catalogue sections.
{"kjv", canonicalBookCount + extraCanonicalBookCount + deadSeaScrollCount},
// WEB, ASV, LXX, and Shona cover the canonical section only.
{"web", canonicalBookCount},
{"asv", canonicalBookCount},
{"lxx", canonicalBookCount},
// Hebrew Masoretic Text covers the OT only (no NT in Hebrew).
{"hmt", 39},
{"bdb", canonicalBookCount},
}
for _, tt := range tests {
books := sources.BooksByTranslation(tt.id)
if len(books) != tt.count {
t.Errorf("BooksByTranslation(%q): got %d books, want %d", tt.id, len(books), tt.count)
}
}
}

func TestBooksByLanguage(t *testing.T) {
tests := []struct {
lang     string
minBooks int
}{
{"en", canonicalBookCount}, // KJV + WEB + ASV all contribute; at minimum 66
{"el", canonicalBookCount},
{"he", 39}, // OT only
{"sn", canonicalBookCount},
}
for _, tt := range tests {
books := sources.BooksByLanguage(tt.lang)
if len(books) < tt.minBooks {
t.Errorf("BooksByLanguage(%q): got %d books, want at least %d", tt.lang, len(books), tt.minBooks)
}
}
}

func TestKJVBooksHaveTranslationID(t *testing.T) {
books := sources.BooksByTranslation("kjv")
for _, b := range books {
if b.TranslationID != "kjv" {
t.Errorf("book %q: TranslationID=%q, want kjv", b.Title, b.TranslationID)
}
}
}

func TestCanonicalBooksCoversAll66(t *testing.T) {
for _, translationID := range []string{"kjv", "web", "asv", "lxx"} {
books := sources.BooksByTranslation(translationID)
var n int
for _, b := range books {
if b.Category == "canonical" {
n++
}
}
if n != canonicalBookCount {
t.Errorf("translation %q: got %d canonical books, want %d", translationID, n, canonicalBookCount)
}
}
}

func TestHebrewOTOnly(t *testing.T) {
books := sources.BooksByTranslation("hmt")
for _, b := range books {
if b.CanonicalOrder > 39 {
t.Errorf("Hebrew translation has NT book %q (order=%d)", b.Title, b.CanonicalOrder)
}
}
}
