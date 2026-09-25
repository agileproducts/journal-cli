// Package submission stores manuscript snapshots and revision-specific facts.
// Its filesystem store supports macOS and Linux.
package submission

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"syscall"
	"time"
)

type Check struct {
	Name   string `json:"name"`
	Passed bool   `json:"passed"`
}

type Manuscript struct {
	File         string `json:"file"`
	OriginalName string `json:"original_name"`
	Revision     int    `json:"revision"`
}

type Validation struct {
	Revision  int       `json:"revision"`
	CheckedAt time.Time `json:"checked_at"`
	Checks    []Check   `json:"checks"`
}

type Record struct {
	SchemaVersion int          `json:"schema_version"`
	ID            string       `json:"id"`
	CreatedAt     time.Time    `json:"created_at"`
	Manuscript    Manuscript   `json:"manuscript"`
	Validations   []Validation `json:"validations"`
}

type Store struct{ Root string }

func DefaultRoot() (string, error) {
	if root := os.Getenv("JCLI_HOME"); root != "" {
		return filepath.Abs(root)
	}
	if root := os.Getenv("XDG_DATA_HOME"); filepath.IsAbs(root) {
		return filepath.Join(root, "jcli"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("find submission store: %w", err)
	}
	return filepath.Join(home, ".local", "share", "jcli"), nil
}

var idPattern = regexp.MustCompile(`^sub-[0-9a-f]{32}$`)

func (s Store) directory(id string) (string, error) {
	if s.Root == "" {
		return "", fmt.Errorf("submission store root is empty")
	}
	if !idPattern.MatchString(id) {
		return "", fmt.Errorf("invalid submission ID %q (expected sub- followed by 32 hex characters)", id)
	}
	return filepath.Join(s.Root, "submissions", id), nil
}

// Create snapshots a regular file and publishes a complete new submission.
func (s Store) Create(path string) (Record, error) {
	var empty Record
	if s.Root == "" {
		return empty, fmt.Errorf("submission store root is empty")
	}
	info, err := os.Stat(path)
	if err != nil {
		return empty, fmt.Errorf("read manuscript: %w", err)
	}
	if !info.Mode().IsRegular() {
		return empty, fmt.Errorf("manuscript %q must be a regular file", path)
	}
	src, err := os.Open(path)
	if err != nil {
		return empty, fmt.Errorf("open manuscript: %w", err)
	}
	defer src.Close()
	var random [16]byte
	if _, err := rand.Read(random[:]); err != nil {
		return empty, err
	}
	r := Record{
		SchemaVersion: 1, ID: "sub-" + hex.EncodeToString(random[:]), CreatedAt: time.Now().UTC(),
		Manuscript:  Manuscript{File: "manuscript.md", OriginalName: filepath.Base(path), Revision: 1},
		Validations: []Validation{},
	}
	parent := filepath.Join(s.Root, "submissions")
	if err := os.MkdirAll(parent, 0700); err != nil {
		return empty, fmt.Errorf("create store: %w", err)
	}
	stage, err := os.MkdirTemp(parent, ".upload-")
	if err != nil {
		return empty, err
	}
	defer os.RemoveAll(stage)
	if err := copySnapshot(filepath.Join(stage, r.Manuscript.File), src); err != nil {
		return empty, err
	}
	if err := saveRecord(stage, r); err != nil {
		return empty, err
	}
	dir, err := s.directory(r.ID)
	if err != nil {
		return empty, err
	}
	if err := os.Rename(stage, dir); err != nil {
		return empty, fmt.Errorf("publish submission: %w", err)
	}
	if err := syncDirectory(parent); err != nil {
		return empty, err
	}
	return r, nil
}

func copySnapshot(path string, src io.Reader) error {
	dst, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer dst.Close()
	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("copy manuscript: %w", err)
	}
	if err := dst.Sync(); err != nil {
		return err
	}
	return dst.Close()
}

func (s Store) Load(id string) (Record, error) {
	var r Record
	dir, err := s.directory(id)
	if err != nil {
		return r, err
	}
	data, err := os.ReadFile(filepath.Join(dir, "submission.json"))
	if err != nil {
		return r, fmt.Errorf("load submission %s: %w", id, err)
	}
	if err := json.Unmarshal(data, &r); err != nil {
		return Record{}, fmt.Errorf("decode submission %s: %w", id, err)
	}
	if r.SchemaVersion != 1 || r.ID != id || r.CreatedAt.IsZero() || r.Manuscript.File != "manuscript.md" || r.Manuscript.Revision < 1 {
		return Record{}, fmt.Errorf("invalid or unsupported submission record %s", id)
	}
	return r, nil
}

// ReadManuscript returns the record and the stored snapshot, never the source file.
func (s Store) ReadManuscript(id string) (Record, []byte, error) {
	r, err := s.Load(id)
	if err != nil {
		return r, nil, err
	}
	dir, err := s.directory(id)
	if err != nil {
		return r, nil, err
	}
	data, err := os.ReadFile(filepath.Join(dir, r.Manuscript.File))
	if err != nil {
		return r, nil, fmt.Errorf("read stored manuscript %s: %w", id, err)
	}
	return r, data, nil
}

// RecordValidation locks before reloading so concurrent attempts cannot overwrite
// each other. Locks live on a separate inode from the atomically replaced record.
func (s Store) RecordValidation(id string, revision int, checks []Check) error {
	dir, err := s.directory(id)
	if err != nil {
		return err
	}
	lock, err := os.OpenFile(filepath.Join(dir, ".lock"), os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return fmt.Errorf("lock submission %s: %w", id, err)
	}
	defer lock.Close()
	for {
		err = syscall.Flock(int(lock.Fd()), syscall.LOCK_EX)
		if err != syscall.EINTR {
			break
		}
	}
	if err != nil {
		return fmt.Errorf("lock submission %s: %w", id, err)
	}
	defer syscall.Flock(int(lock.Fd()), syscall.LOCK_UN)
	r, err := s.Load(id)
	if err != nil {
		return err
	}
	if r.Manuscript.Revision != revision {
		return fmt.Errorf("manuscript revision changed: checked %d, current %d", revision, r.Manuscript.Revision)
	}
	r.Validations = append(r.Validations, Validation{Revision: revision, CheckedAt: time.Now().UTC(), Checks: checks})
	return saveRecord(dir, r)
}

func saveRecord(dir string, r Record) error {
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	f, err := os.CreateTemp(dir, ".record-")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	defer f.Close()
	if _, err := f.Write(append(data, '\n')); err != nil {
		return err
	}
	if err := f.Sync(); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	if err := os.Rename(f.Name(), filepath.Join(dir, "submission.json")); err != nil {
		return err
	}
	return syncDirectory(dir)
}

func syncDirectory(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return f.Sync()
}
