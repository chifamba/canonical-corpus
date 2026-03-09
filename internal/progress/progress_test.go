package progress_test

import (
"os"
"path/filepath"
"testing"

"github.com/chifamba/canonical-corpus/internal/progress"
)

func TestLoadFresh(t *testing.T) {
dir := t.TempDir()
s, err := progress.Load(dir)
if err != nil {
t.Fatalf("Load: %v", err)
}
if s.IsComplete("canonical/001-genesis/en-kjv.md") {
t.Fatal("expected key to be incomplete on fresh state")
}
}

func TestMarkCompleteAndReload(t *testing.T) {
dir := t.TempDir()

s, err := progress.Load(dir)
if err != nil {
t.Fatalf("Load: %v", err)
}

key := "canonical/001-genesis/en-kjv.md"
if err := s.MarkComplete(key); err != nil {
t.Fatalf("MarkComplete: %v", err)
}
if !s.IsComplete(key) {
t.Fatal("key should be complete after MarkComplete")
}

// Reload from disk – should still be complete.
s2, err := progress.Load(dir)
if err != nil {
t.Fatalf("Load (reload): %v", err)
}
if !s2.IsComplete(key) {
t.Fatal("key should persist across reload")
}
}

func TestMarkCompleteDoesNotAffectOtherKeys(t *testing.T) {
dir := t.TempDir()
s, _ := progress.Load(dir)

_ = s.MarkComplete("canonical/001-genesis/en-kjv.md")

if s.IsComplete("canonical/002-exodus/en-kjv.md") {
t.Fatal("unrelated key should not be affected")
}
}

func TestStateFileCreated(t *testing.T) {
dir := t.TempDir()
s, _ := progress.Load(dir)
_ = s.MarkComplete("x")

if _, err := os.Stat(filepath.Join(dir, ".progress.json")); err != nil {
t.Fatalf("progress file not created: %v", err)
}
}

func TestInMemoryModeDoesNotPersist(t *testing.T) {
// baseDir="" → in-memory mode: MarkComplete should succeed but not write a file.
s, err := progress.Load("")
if err != nil {
t.Fatalf("Load(''): %v", err)
}
if err := s.MarkComplete("some/key"); err != nil {
t.Fatalf("MarkComplete in in-memory mode: %v", err)
}
if !s.IsComplete("some/key") {
t.Fatal("key should be complete within the same in-memory instance")
}
}
