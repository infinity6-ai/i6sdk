package fszcache

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/infinity6-ai/gox/commonz/urlz"
	"github.com/infinity6-ai/gox/fsz/fsz"
)

// FileRef represents a reference to a specific version of a stored asset.
type FileRef struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	Etag      string `json:"etag"`
	CheckedAt string `json:"checked_at,omitempty"`
}

// Fs manages a versioned cache rooted at a gox/fsz URL.
// Layout under baseUrl:
//
//	<baseUrl>/<safe_name>/
//	├── version.txt          (current released version)
//	├── checkedat.txt        (timestamp of last remote check)
//	└── v<version>/
//	    ├── etag.txt         (remote content ETag)
//	    └── blob.bin         (binary asset content)
type Fs struct {
	baseUrl *urlz.Url
	mu      sync.RWMutex
}

// DefaultBaseUrl returns the default cache root URL: file://$HOME/.cache/i6/sdkgo/fs.
func DefaultBaseUrl() (*urlz.Url, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = os.TempDir()
	}
	p := filepath.Join(homeDir, ".cache", "i6", "sdkgo", "fs")
	u, err := urlz.Parse("file://" + filepath.Clean(p))
	if err != nil {
		return nil, fmt.Errorf("failed to parse default cache url: %w", err)
	}
	return u, nil
}

// NewFs creates a new versioned filesystem store rooted at baseUrl.
// If baseUrl is nil, DefaultBaseUrl() is used.
func NewFs(baseUrl *urlz.Url) (*Fs, error) {
	if baseUrl == nil {
		def, err := DefaultBaseUrl()
		if err != nil {
			return nil, fmt.Errorf("failed to resolve default cache url: %w", err)
		}
		baseUrl = def
	}
	return &Fs{
		baseUrl: baseUrl,
	}, nil
}

// BaseUrl returns the root URL of the cache store.
func (f *Fs) BaseUrl() *urlz.Url {
	return f.baseUrl
}

func (f *Fs) fileUrl(name string) *urlz.Url {
	safeName := strings.ReplaceAll(name, "/", "_")
	safeName = strings.ReplaceAll(safeName, "\\", "_")
	return f.baseUrl.MustJoinPathString(safeName)
}

// ResolveUrl returns the *urlz.Url for a version's blob.bin (defaults to released version).
func (f *Fs) ResolveUrl(ctx context.Context, name string, version ...string) *urlz.Url {
	f.mu.RLock()
	defer f.mu.RUnlock()

	fileUrl := f.fileUrl(name)
	ver := ""
	if len(version) > 0 && version[0] != "" {
		ver = version[0]
	} else {
		verText, _ := f.readText(ctx, fileUrl.MustJoinPathString("version.txt"))
		ver = strings.TrimSpace(verText)
	}
	return fileUrl.MustJoinPathString("v"+ver, "blob.bin")
}

// Resolve returns the filesystem path string for a version's blob.bin if backed by file:// scheme,
// or the string representation of the URL for non-file schemes.
func (f *Fs) Resolve(ctx context.Context, name string, version ...string) string {
	u := f.ResolveUrl(ctx, name, version...)
	if u.Scheme == "file" && u.Path != nil {
		return filepath.Clean(u.Path.String())
	}
	return u.String()
}

// Get resolves a version of a file, or returns nil if not found.
func (f *Fs) Get(ctx context.Context, name string, version ...string) (*FileRef, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	fileUrl := f.fileUrl(name)

	releasedVersion, _ := f.readText(ctx, fileUrl.MustJoinPathString("version.txt"))
	ver := ""
	if len(version) > 0 && version[0] != "" {
		ver = version[0]
	} else {
		ver = strings.TrimSpace(releasedVersion)
	}
	if ver == "" {
		return nil, nil
	}

	verDirUrl := fileUrl.MustJoinPathString("v" + ver)
	blobUrl := verDirUrl.MustJoinPathString("blob.bin")
	stat, err := fsz.Stat(ctx, blobUrl)
	if err != nil || stat == nil {
		return nil, nil
	}

	etag, _ := f.readText(ctx, verDirUrl.MustJoinPathString("etag.txt"))

	checkedAt := ""
	if len(version) == 0 || version[0] == strings.TrimSpace(releasedVersion) {
		checkedAt, _ = f.readText(ctx, fileUrl.MustJoinPathString("checkedat.txt"))
	}

	return &FileRef{
		Name:      name,
		Version:   ver,
		Etag:      strings.TrimSpace(etag),
		CheckedAt: strings.TrimSpace(checkedAt),
	}, nil
}

func isVersionDir(name string) bool {
	return strings.HasPrefix(name, "v") && len(name) > 1 && name[1] >= '0' && name[1] <= '9'
}

// Create creates an unreleased version directory and writes etag.txt.
// If an existing version matching etag is already present, it is reused.
func (f *Fs) Create(ctx context.Context, name string, etag string) (*FileRef, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	fileUrl := f.fileUrl(name)

	// Check if latest existing version directory matches this etag
	var vDirs []string
	_ = fsz.Lister(ctx, fileUrl, fsz.Walker{
		Pager: func(page []*fsz.FileStat, _ string) error {
			for _, item := range page {
				if item.Url != nil {
					baseName := filepath.Base(item.Url.Path.String())
					if isVersionDir(baseName) {
						vDirs = append(vDirs, baseName)
					}
				}
			}
			return nil
		},
	})

	if len(vDirs) > 0 {
		sort.Strings(vDirs)
		latest := vDirs[len(vDirs)-1]
		latestEtag, err := f.readText(ctx, fileUrl.MustJoinPathString(latest, "etag.txt"))
		if err == nil && strings.TrimSpace(latestEtag) == etag {
			ver := strings.TrimPrefix(latest, "v")
			return &FileRef{
				Name:    name,
				Version: ver,
				Etag:    etag,
			}, nil
		}
	}

	version := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	verDirUrl := fileUrl.MustJoinPathString("v" + version)
	etagUrl := verDirUrl.MustJoinPathString("etag.txt")

	if err := fsz.Upload(ctx, etagUrl, nil, strings.NewReader(etag)); err != nil {
		return nil, fmt.Errorf("failed to upload etag to %s: %w", etagUrl, err)
	}

	return &FileRef{
		Name:    name,
		Version: version,
		Etag:    etag,
	}, nil
}

// Release marks a version as the released (current) version, updates checkedat.txt,
// and deletes versions older than the released version.
func (f *Fs) Release(ctx context.Context, name string, version string) (*FileRef, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	fileUrl := f.fileUrl(name)
	verDirUrl := fileUrl.MustJoinPathString("v" + version)
	blobUrl := verDirUrl.MustJoinPathString("blob.bin")

	stat, err := fsz.Stat(ctx, blobUrl)
	if err != nil || stat == nil {
		return nil, fmt.Errorf("fs: version %s of %s does not exist or has no blob.bin", version, name)
	}

	verFileUrl := fileUrl.MustJoinPathString("version.txt")
	if err := fsz.Upload(ctx, verFileUrl, nil, strings.NewReader(version)); err != nil {
		return nil, fmt.Errorf("failed to write version.txt: %w", err)
	}

	checkedAt := f.touchLocked(ctx, name)
	f.cleanFileLocked(ctx, name)

	etag, _ := f.readText(ctx, verDirUrl.MustJoinPathString("etag.txt"))

	return &FileRef{
		Name:      name,
		Version:   version,
		Etag:      strings.TrimSpace(etag),
		CheckedAt: checkedAt,
	}, nil
}

// Touch updates checkedat.txt to now (or custom ISO timestamp).
func (f *Fs) Touch(ctx context.Context, name string, checkedAt ...string) string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.touchLocked(ctx, name, checkedAt...)
}

func (f *Fs) touchLocked(ctx context.Context, name string, checkedAt ...string) string {
	fileUrl := f.fileUrl(name)
	ts := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	if len(checkedAt) > 0 && checkedAt[0] != "" {
		ts = checkedAt[0]
	}

	checkedUrl := fileUrl.MustJoinPathString("checkedat.txt")
	_ = fsz.Upload(ctx, checkedUrl, nil, strings.NewReader(ts))
	return ts
}

// Discard removes an unreleased version directory.
func (f *Fs) Discard(ctx context.Context, name string, version string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	fileUrl := f.fileUrl(name)
	verFileUrl := fileUrl.MustJoinPathString("version.txt")
	if released, err := f.readText(ctx, verFileUrl); err == nil && strings.TrimSpace(released) == version {
		return fmt.Errorf("fs: version %s of %s is the released version", version, name)
	}

	verDirUrl := fileUrl.MustJoinPathString("v" + version)
	return fsz.RmTree(ctx, verDirUrl)
}

// CleanFile removes versions older than the released version for an asset.
func (f *Fs) CleanFile(ctx context.Context, name string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.cleanFileLocked(ctx, name)
}

func (f *Fs) cleanFileLocked(ctx context.Context, name string) {
	fileUrl := f.fileUrl(name)
	releasedVer, err := f.readText(ctx, fileUrl.MustJoinPathString("version.txt"))
	if err != nil || strings.TrimSpace(releasedVer) == "" {
		return
	}
	released := "v" + strings.TrimSpace(releasedVer)

	_ = fsz.Lister(ctx, fileUrl, fsz.Walker{
		Pager: func(page []*fsz.FileStat, _ string) error {
			for _, item := range page {
				if item.Url != nil {
					baseName := filepath.Base(item.Url.Path.String())
					if isVersionDir(baseName) && baseName < released {
						_ = fsz.RmTree(ctx, item.Url)
					}
				}
			}
			return nil
		},
	})
}

// Clean removes older versions across all stored assets.
func (f *Fs) Clean(ctx context.Context) {
	f.mu.Lock()
	defer f.mu.Unlock()

	_ = fsz.Lister(ctx, f.baseUrl, fsz.Walker{
		Pager: func(page []*fsz.FileStat, _ string) error {
			for _, item := range page {
				if item.Url != nil {
					baseName := filepath.Base(item.Url.Path.String())
					f.cleanFileLocked(ctx, baseName)
				}
			}
			return nil
		},
	})
}

// Exists checks if a released version of the asset exists.
func (f *Fs) Exists(ctx context.Context, name string) bool {
	ref, _ := f.Get(ctx, name)
	return ref != nil
}

// Delete removes the entire asset folder.
func (f *Fs) Delete(ctx context.Context, name string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	return fsz.RmTree(ctx, f.fileUrl(name))
}

func (f *Fs) readText(ctx context.Context, u *urlz.Url) (string, error) {
	var content string
	err := fsz.Download(ctx, u, func(found bool, _ http.Header, reader io.Reader) error {
		if !found || reader == nil {
			return nil
		}
		raw, err := io.ReadAll(reader)
		if err != nil {
			return err
		}
		content = string(bytes.TrimSpace(raw))
		return nil
	})
	if err != nil {
		return "", err
	}
	return content, nil
}
