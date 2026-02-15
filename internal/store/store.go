package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/user/tack/internal/model"
)

const (
	tackDir  = ".tack"
	tackFile = "tacks.json"
)

// Store manages persistence of tacks to a JSON file.
type Store struct {
	path  string
	Tacks map[string][]model.Tack `json:"tacks"`
}

// New creates or opens the tack store at ~/.tack/tacks.json.
func New() (*Store, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("could not find home directory: %w", err)
	}

	dir := filepath.Join(home, tackDir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("could not create tack directory: %w", err)
	}

	p := filepath.Join(dir, tackFile)
	s := &Store{path: p}

	if _, err := os.Stat(p); os.IsNotExist(err) {
		s.Tacks = make(map[string][]model.Tack)
		return s, s.Save()
	}

	data, err := os.ReadFile(p)
	if err != nil {
		return nil, fmt.Errorf("could not read tack file: %w", err)
	}

	if err := json.Unmarshal(data, s); err != nil {
		return nil, fmt.Errorf("could not parse tack file: %w", err)
	}

	return s, nil
}

// Save writes the current tacks to disk.
func (s *Store) Save() error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal tacks: %w", err)
	}
	return os.WriteFile(s.path, data, 0o644)
}

// Add appends a new tack and saves.
func (s *Store) Add(t model.Tack) error {
	s.Tacks[t.Dir] = append(s.Tacks[t.Dir], t)
	return s.Save()
}

// ForDir returns all tacks pinned to the given directory.
func (s *Store) ForDir(dir string) []model.Tack {
	var result []model.Tack
	for _, t := range s.Tacks[dir] {
		if t.Dir == dir {
			result = append(result, t)
		}
	}
	return result
}

// FindByID returns a pointer to the tack with the given ID, or nil.
func (s *Store) FindByID(id string) *model.Tack {
	for dir, tacks := range s.Tacks {
		for i := range tacks {
			if tacks[i].ID == id {
				return &s.Tacks[dir][i]
			}
		}
	}
	return nil
}

// Remove deletes a tack by ID and saves.
func (s *Store) Remove(id string) error {
	for dir, tacks := range s.Tacks {
		for i, t := range tacks {
			if t.ID == id {
				s.Tacks[dir] = append(tacks[:i], tacks[i+1:]...)
				return s.Save()
			}
		}
	}
	return fmt.Errorf("tack %q not found", id)
}

// AllDirs returns a deduplicated list of directories that have tacks.
func (s *Store) AllDirs() []string {
	var dirs []string
	for dir := range s.Tacks {
		dirs = append(dirs, dir)
	}
	return dirs
}
