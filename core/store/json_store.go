package store

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"

	"github.com/DoniLite/kapelas-bot/conf"
)

// JSONStore is a simple file-backed JSON store. Data is organized by
// collection (each collection is stored in a single JSON file that maps id->value).
// It is safe for concurrent use.
type JSONStore struct {
	baseDir string
	mu      sync.RWMutex
}

// NewJSONStore creates a new JSONStore. If baseDir is empty it chooses a default
// location: in development it will use ./data_dev, otherwise the OS user config
// dir under a `kapelas-go-bot/data` subfolder. The conf package is consulted
// for the BOT_IS_DEVELOPMENT flag.
func NewJSONStore(baseDir string) (*JSONStore, error) {
	if baseDir == "" {
		e := conf.GetEnv()
		if e != nil && e.GetBool(conf.BOT_IS_DEVELOPMENT) {
			log.Printf("Running in development mode, using local data directory")
			baseDir = "./data_dev"
		} else if strings.HasSuffix(os.Args[0], ".test") {
			log.Printf("Running in test mode, using temporary data directory")
			baseDir = filepath.Join(os.TempDir(), "kappelas-go-bot", "data")
		} else {
			log.Printf("Running in production mode, using user config directory for data storage")
			dir, err := os.UserConfigDir()
			if err != nil {
				return nil, fmt.Errorf("getting user config dir: %w", err)
			}
			baseDir = filepath.Join(dir, "kappelas-go-bot", "data")
			log.Printf("data directory: %s", baseDir)
		}
	}
	if err := os.MkdirAll(baseDir, 0o755); err != nil {
		return nil, fmt.Errorf("creating base dir: %w", err)
	}
	return &JSONStore{baseDir: baseDir}, nil
}

func (s *JSONStore) filePath(collection string) string {
	name := fmt.Sprintf("%s.json", collection)
	return filepath.Join(s.baseDir, name)
}

// internal representation: map[id]json.RawMessage
type rawMap map[string]json.RawMessage

func (s *JSONStore) load(collection string) (rawMap, error) {
	path := s.filePath(collection)
	data := rawMap{}
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return data, nil
		}
		return nil, fmt.Errorf("open collection file: %w", err)
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("read collection file: %w", err)
	}
	if len(b) == 0 {
		return data, nil
	}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, fmt.Errorf("unmarshal collection file: %w", err)
	}
	return data, nil
}

func (s *JSONStore) save(collection string, m rawMap) error {
	path := s.filePath(collection)
	tmp := path + ".tmp"
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal collection: %w", err)
	}
	// Write to temp file then rename (atomic on most OSes)
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return fmt.Errorf("write temp file: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		return fmt.Errorf("rename temp to target: %w", err)
	}
	return nil
}

// Create inserts a new value under the given id. Returns error if id exists.
func (s *JSONStore) Create(collection, id string, v any) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	m, err := s.load(collection)
	if err != nil {
		return err
	}
	if _, ok := m[id]; ok {
		return fmt.Errorf("id %s already exists", id)
	}
	raw, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("marshal value: %w", err)
	}
	m[id] = raw
	return s.save(collection, m)
}

// Upsert creates or updates the value for id.
func (s *JSONStore) Upsert(collection, id string, v any) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	m, err := s.load(collection)
	if err != nil {
		return err
	}
	raw, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("marshal value: %w", err)
	}
	m[id] = raw
	return s.save(collection, m)
}

// Get loads the value for id into dest (which must be a pointer).
func (s *JSONStore) Get(collection, id string, dest any) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m, err := s.load(collection)
	if err != nil {
		return err
	}
	raw, ok := m[id]
	if !ok {
		return fmt.Errorf("id %s not found", id)
	}
	if err := json.Unmarshal(raw, dest); err != nil {
		return fmt.Errorf("unmarshal value: %w", err)
	}
	return nil
}

// Delete removes an entry by id.
func (s *JSONStore) Delete(collection, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	m, err := s.load(collection)
	if err != nil {
		return err
	}
	if _, ok := m[id]; !ok {
		return fmt.Errorf("id %s not found", id)
	}
	delete(m, id)
	return s.save(collection, m)
}

// List returns all stored objects as a slice of maps (useful for dynamic queries).
func (s *JSONStore) List(collection string) ([]map[string]any, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m, err := s.load(collection)
	if err != nil {
		return nil, err
	}
	var out []map[string]any
	for _, raw := range m {
		var mm map[string]any
		dec := json.NewDecoder(bytes.NewReader(raw))
		dec.UseNumber()
		if err := dec.Decode(&mm); err != nil {
			return nil, fmt.Errorf("decode item: %w", err)
		}
		out = append(out, mm)
	}
	return out, nil
}

// Search runs a custom predicate against each object in the collection and
// returns matching objects as maps.
func (s *JSONStore) Search(collection string, match func(map[string]any) bool) ([]map[string]any, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m, err := s.load(collection)
	if err != nil {
		return nil, err
	}
	var out []map[string]any
	for _, raw := range m {
		var mm map[string]any
		dec := json.NewDecoder(bytes.NewReader(raw))
		dec.UseNumber()
		if err := dec.Decode(&mm); err != nil {
			return nil, fmt.Errorf("decode item for search: %w", err)
		}
		if match(mm) {
			out = append(out, mm)
		}
	}
	return out, nil
}

// SearchByField finds objects where the given top-level field equals value.
// For string matching it will do substring contains when value is a string.
func (s *JSONStore) SearchByField(collection, field string, value any) ([]map[string]any, error) {
	return s.Search(collection, func(m map[string]any) bool {
		v, ok := m[field]
		if !ok {
			return false
		}
		// handle numbers stored as json.Number
		if num, ok := v.(json.Number); ok {
			v = num.String()
		}
		if valueNum, ok := value.(json.Number); ok {
			value = valueNum.String()
		}
		// if both are strings, do substring match
		if sv, ok := v.(string); ok {
			if sv2, ok := value.(string); ok {
				return containsIgnoreCase(sv, sv2)
			}
		}
		return reflect.DeepEqual(v, value)
	})
}

func containsIgnoreCase(s, sub string) bool {
	return bytes.Contains(bytes.ToLower([]byte(s)), bytes.ToLower([]byte(sub)))
}
