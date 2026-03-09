// Package progress provides a simple JSON-backed store for tracking which
// corpus books have already been compiled, enabling resumable builds.
package progress

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

const stateFileName = ".progress.json"

// State persists the set of completed book keys across runs.
// A key is formed as "<category>/<NNN-title>/<lang>-<translationID>".
type State struct {
	mu        sync.Mutex
	path      string // empty string → in-memory only (no persistence)
	Completed map[string]bool `json:"completed"`
}

// Load reads (or creates) the progress state file in baseDir.
// If baseDir is empty the returned State operates in-memory only.
func Load(baseDir string) (*State, error) {
	if baseDir == "" {
		return &State{Completed: make(map[string]bool)}, nil
	}

	path := filepath.Join(baseDir, stateFileName)
	s := &State{
		path:      path,
		Completed: make(map[string]bool),
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		// No existing state – start fresh.
		return s, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading progress file %q: %w", path, err)
	}

	if err := json.Unmarshal(data, s); err != nil {
		return nil, fmt.Errorf("parsing progress file %q: %w", path, err)
	}
	if s.Completed == nil {
		s.Completed = make(map[string]bool)
	}
	return s, nil
}

// IsComplete reports whether key has been previously marked complete.
func (s *State) IsComplete(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Completed[key]
}

// MarkComplete records key as successfully compiled and persists the state.
func (s *State) MarkComplete(key string) error {
	s.mu.Lock()
	s.Completed[key] = true
	s.mu.Unlock()
	return s.save()
}

// save writes the state to disk. Caller must NOT hold mu.
// No-op when the state has no backing file (in-memory mode).
func (s *State) save() error {
	// Snapshot the data under the lock, then do I/O outside it to minimise
	// lock contention and avoid blocking readers during disk writes.
	s.mu.Lock()
	path := s.path
	completed := make(map[string]bool, len(s.Completed))
	for k, v := range s.Completed {
		completed[k] = v
	}
	s.mu.Unlock()

	if path == "" {
		return nil // in-memory only
	}

	snapshot := struct {
		Completed map[string]bool `json:"completed"`
	}{completed}

	data, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling progress state: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("creating progress directory: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("writing progress file %q: %w", path, err)
	}
	return nil
}
