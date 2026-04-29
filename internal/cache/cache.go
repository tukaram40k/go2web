package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const cacheDir = ".go2web/cache"

type Entry struct {
	URL           string    `json:"url"`
	CachedAt      time.Time `json:"cached_at"`
	RedirectCount int       `json:"redirect_count"`
}

func KeyForURL(rawURL string) string {
	hash := sha256.Sum256([]byte(rawURL))
	return hex.EncodeToString(hash[:])
}

func Load(rawURL string) ([]byte, Entry, bool, error) {
	key := KeyForURL(rawURL)
	bodyPath, metaPath := cachePaths(key)

	body, err := os.ReadFile(bodyPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, Entry{}, false, nil
		}
		return nil, Entry{}, false, err
	}

	entry := Entry{URL: rawURL}
	metaBytes, err := os.ReadFile(metaPath)
	if err == nil {
		if err := json.Unmarshal(metaBytes, &entry); err != nil {
			return body, Entry{URL: rawURL}, true, fmt.Errorf("failed to parse cache metadata: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return body, Entry{URL: rawURL}, true, err
	}

	return body, entry, true, nil
}

func Store(rawURL string, body []byte, redirectCount int) error {
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return err
	}

	key := KeyForURL(rawURL)
	bodyPath, metaPath := cachePaths(key)

	entry := Entry{
		URL:           rawURL,
		CachedAt:      time.Now().UTC(),
		RedirectCount: redirectCount,
	}

	metaBytes, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return err
	}

	if err := writeFileAtomically(bodyPath, body, 0o644); err != nil {
		return err
	}

	if err := writeFileAtomically(metaPath, metaBytes, 0o644); err != nil {
		return err
	}

	return nil
}

func cachePaths(key string) (string, string) {
	fileBase := filepath.Join(cacheDir, key)
	return fileBase + ".body", fileBase + ".meta.json"
}

func writeFileAtomically(path string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(path)
	tmpFile, err := os.CreateTemp(dir, "tmp-*")
	if err != nil {
		return err
	}

	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(data); err != nil {
		tmpFile.Close()
		return err
	}

	if err := tmpFile.Close(); err != nil {
		return err
	}

	if err := os.Chmod(tmpFile.Name(), perm); err != nil {
		return err
	}

	return os.Rename(tmpFile.Name(), path)
}
