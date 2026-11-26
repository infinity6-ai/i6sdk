package fszcache

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// FileRef represents a reference to a specific version of a stored asset.
type FileRef struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	Etag      string `json:"etag"`
	CheckedAt string `json:"checked_at,omitempty"`
}

// Fs manages the local versioned cache directory matching sdkweb (OPFS) and sdkpy layout.
// Layout:
//
//	<base_dir>/<safe_name>/
//	├── version.txt          (current released version)
//	├── checkedat.txt        (timestamp of last remote check)
//	└── v<version>/
//	    ├── etag.txt         (remote content ETag)
//	    └── blob.bin         (binary asset content)
type Fs struct {
	baseDir string
	mu      sync.RWMutex
}

// DefaultBaseDir returns the default cache root path (~/.cache/i6/sdkgo/fs).
func DefaultBaseDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = os.TempDir()
	}
	return filepath.Join(homeDir, ".cache", "i6", "sdkgo", "fs"), nil
}

// NewFs creates a new versioned filesystem store rooted at baseDir.
// If baseDir is empty, DefaultBaseDir() is used.
func NewFs(baseDir string) (*Fs, error) {
	if baseDir == "" {
		def, err := DefaultBaseDir()
		if err != nil {
			return nil, fmt.Errorf("failed to resolve default base dir: %w", err)
		}
		baseDir = def
	}
	baseDir = filepath.Clean(baseDir)
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache base dir %s: %w", baseDir, err)
	}
	return &Fs{
		baseDir: baseDir,
	}, nil
}

// BaseDir returns the root directory of the cache.
func (f *Fs) BaseDir() string {
	return f.baseDir
}

func (f *Fs) fileDir(name string) string {
	safeName := strings.ReplaceAll(name, "/", "_")
	safeName = strings.ReplaceAll(safeName, "\\", "_")
	return filepath.Join(f.baseDir, safeName)
}

// Resolve returns the absolute path for a version's blob.bin (defaults to released version).
func (f *Fs) Resolve(name string, version ...string) string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	fileDir := f.fileDir(name)
	ver := ""
	if len(version) > 0 && version[0] != "" {
		ver = version[0]
	} else {
		verFile := filepath.Join(fileDir, "version.txt")
		if raw, err := os.ReadFile(verFile); err == nil {
			ver = strings.TrimSpace(string(raw))
		}
	}
	return filepath.Join(fileDir, "v"+ver, "blob.bin")
}

// Get resolves a version of a file, or returns nil if not found.
func (f *Fs) Get(name string, version ...string) (*FileRef, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	fileDir := f.fileDir(name)
	info, err := os.Stat(fileDir)
	if err != nil || !info.IsDir() {
		return nil, nil
	}

	releasedVersion := ""
	verFile := filepath.Join(fileDir, "version.txt")
	if raw, err := os.ReadFile(verFile); err == nil {
		releasedVersion = strings.TrimSpace(string(raw))
	}

	ver := ""
	if len(version) > 0 && version[0] != "" {
		ver = version[0]
	} else {
		ver = releasedVersion
	}
	if ver == "" {
		return nil, nil
	}

	verDir := filepath.Join(fileDir, "v"+ver)
	blobFile := filepath.Join(verDir, "blob.bin")
	blobInfo, err := os.Stat(blobFile)
	if err != nil || blobInfo.IsDir() {
		return nil, nil
	}

	etag := ""
	etagFile := filepath.Join(verDir, "etag.txt")
	if raw, err := os.ReadFile(etagFile); err == nil {
		etag = strings.TrimSpace(string(raw))
	}

	checkedAt := ""
	if len(version) == 0 || version[0] == releasedVersion {
		checkedFile := filepath.Join(fileDir, "checkedat.txt")
		if raw, err := os.ReadFile(checkedFile); err == nil {
			checkedAt = strings.TrimSpace(string(raw))
		}
	}

	return &FileRef{
		Name:      name,
		Version:   ver,
		Etag:      etag,
		CheckedAt: checkedAt,
	}, nil
}

// Create creates an unreleased version directory and writes etag.txt.
// If an existing version matching etag is already present, it is reused.
func (f *Fs) Create(name string, etag string) (*FileRef, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	fileDir := f.fileDir(name)
	if err := os.MkdirAll(fileDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create file directory for %s: %w", name, err)
	}

	// Check if latest existing version directory matches this etag
	entries, err := os.ReadDir(fileDir)
	if err == nil {
		var vDirs []string
		for _, e := range entries {
			if e.IsDir() && strings.HasPrefix(e.Name(), "v") {
				vDirs = append(vDirs, e.Name())
			}
		}
		if len(vDirs) > 0 {
			sort.Strings(vDirs)
			latest := vDirs[len(vDirs)-1]
			latestEtagFile := filepath.Join(fileDir, latest, "etag.txt")
			if raw, err := os.ReadFile(latestEtagFile); err == nil && strings.TrimSpace(string(raw)) == etag {
				ver := strings.TrimPrefix(latest, "v")
				return &FileRef{
					Name:    name,
					Version: ver,
					Etag:    etag,
				}, nil
			}
		}
	}

	version := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	verDir := filepath.Join(fileDir, "v"+version)
	if err := os.MkdirAll(verDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create version directory %s: %w", verDir, err)
	}

	etagFile := filepath.Join(verDir, "etag.txt")
	if err := os.WriteFile(etagFile, []byte(etag), 0644); err != nil {
		return nil, fmt.Errorf("failed to write etag file %s: %w", etagFile, err)
	}

	return &FileRef{
		Name:    name,
		Version: version,
		Etag:    etag,
	}, nil
}

// Release marks a version as the released (current) version, updates checkedat.txt,
// and deletes versions older than the released version.
func (f *Fs) Release(name string, version string) (*FileRef, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	fileDir := f.fileDir(name)
	verDir := filepath.Join(fileDir, "v"+version)
	blobFile := filepath.Join(verDir, "blob.bin")
	if info, err := os.Stat(blobFile); err != nil || info.IsDir() {
		return nil, fmt.Errorf("fs: version %s of %s does not exist or has no blob.bin", version, name)
	}

	verFile := filepath.Join(fileDir, "version.txt")
	if err := os.WriteFile(verFile, []byte(version), 0644); err != nil {
		return nil, fmt.Errorf("failed to write version.txt: %w", err)
	}

	checkedAt := f.touchLocked(name)
	f.cleanFileLocked(name)

	etag := ""
	if raw, err := os.ReadFile(filepath.Join(verDir, "etag.txt")); err == nil {
		etag = strings.TrimSpace(string(raw))
	}

	return &FileRef{
		Name:      name,
		Version:   version,
		Etag:      etag,
		CheckedAt: checkedAt,
	}, nil
}

// Touch updates checkedat.txt to now (or custom ISO timestamp).
func (f *Fs) Touch(name string, checkedAt ...string) string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.touchLocked(name, checkedAt...)
}

func (f *Fs) touchLocked(name string, checkedAt ...string) string {
	fileDir := f.fileDir(name)
	_ = os.MkdirAll(fileDir, 0755)

	ts := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	if len(checkedAt) > 0 && checkedAt[0] != "" {
		ts = checkedAt[0]
	}

	checkedFile := filepath.Join(fileDir, "checkedat.txt")
	_ = os.WriteFile(checkedFile, []byte(ts), 0644)
	return ts
}

// Discard removes an unreleased version directory.
func (f *Fs) Discard(name string, version string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	fileDir := f.fileDir(name)
	verFile := filepath.Join(fileDir, "version.txt")
	if raw, err := os.ReadFile(verFile); err == nil && strings.TrimSpace(string(raw)) == version {
		return fmt.Errorf("fs: version %s of %s is the released version", version, name)
	}

	verDir := filepath.Join(fileDir, "v"+version)
	_ = os.RemoveAll(verDir)
	return nil
}

// CleanFile removes versions older than the released version for an asset.
func (f *Fs) CleanFile(name string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.cleanFileLocked(name)
}

func (f *Fs) cleanFileLocked(name string) {
	fileDir := f.fileDir(name)
	verFile := filepath.Join(fileDir, "version.txt")
	raw, err := os.ReadFile(verFile)
	if err != nil {
		return
	}
	released := "v" + strings.TrimSpace(string(raw))

	entries, err := os.ReadDir(fileDir)
	if err != nil {
		return
	}
	for _, e := range entries {
		if e.IsDir() && strings.HasPrefix(e.Name(), "v") && e.Name() < released {
			_ = os.RemoveAll(filepath.Join(fileDir, e.Name()))
		}
	}
}

// Clean removes older versions across all stored assets.
func (f *Fs) Clean() {
	f.mu.Lock()
	defer f.mu.Unlock()

	entries, err := os.ReadDir(f.baseDir)
	if err != nil {
		return
	}
	for _, e := range entries {
		if e.IsDir() {
			f.cleanFileLocked(e.Name())
		}
	}
}

// Exists checks if a released version of the asset exists.
func (f *Fs) Exists(name string) bool {
	ref, _ := f.Get(name)
	return ref != nil
}

// Delete removes the entire asset folder.
func (f *Fs) Delete(name string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	return os.RemoveAll(f.fileDir(name))
}
