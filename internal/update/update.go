package update

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	githubRepo   = "icarus-itcs/lazycap"
	releaseURL   = "https://api.github.com/repos/" + githubRepo + "/releases/latest"
	cacheTimeout = 5 * time.Minute
)

// Info contains version update information
type Info struct {
	CurrentVersion  string
	LatestVersion   string
	UpdateAvailable bool
	ReleaseURL      string
	DownloadURL     string
	ReleaseNotes    string
	CheckedAt       time.Time
}

// GitHubRelease represents a GitHub release response
type GitHubRelease struct {
	TagName     string  `json:"tag_name"`
	HTMLURL     string  `json:"html_url"`
	Body        string  `json:"body"`
	Assets      []Asset `json:"assets"`
	PublishedAt string  `json:"published_at"`
}

// Asset represents a release asset
type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

var (
	cachedInfo *Info
	lastCheck  time.Time
)

// Check queries GitHub for the latest release and compares versions
func Check(currentVersion string) (*Info, error) {
	// Return cached result if recent
	if cachedInfo != nil && time.Since(lastCheck) < cacheTimeout {
		cachedInfo.CurrentVersion = currentVersion
		return cachedInfo, nil
	}

	info := &Info{
		CurrentVersion: currentVersion,
		CheckedAt:      time.Now(),
	}

	// Query GitHub API
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(releaseURL)
	if err != nil {
		return info, fmt.Errorf("failed to check for updates: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return info, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return info, fmt.Errorf("failed to parse release info: %w", err)
	}

	info.LatestVersion = strings.TrimPrefix(release.TagName, "v")
	info.ReleaseURL = release.HTMLURL
	info.ReleaseNotes = release.Body
	info.UpdateAvailable = isNewerVersion(currentVersion, info.LatestVersion)

	// Find the right download URL for this platform
	info.DownloadURL = findDownloadURL(release.Assets)

	// Cache the result
	cachedInfo = info
	lastCheck = time.Now()

	return info, nil
}

// isNewerVersion compares two semantic versions
// Returns true if latest is newer than current
func isNewerVersion(current, latest string) bool {
	// Strip 'v' prefix if present
	current = strings.TrimPrefix(current, "v")
	latest = strings.TrimPrefix(latest, "v")

	// Handle "dev" version - always consider updates available
	if current == "dev" || current == "" {
		return latest != ""
	}

	// Split into parts
	currentParts := strings.Split(current, ".")
	latestParts := strings.Split(latest, ".")

	// Compare each part
	for i := 0; i < 3; i++ {
		var c, l int
		if i < len(currentParts) {
			_, _ = fmt.Sscanf(currentParts[i], "%d", &c)
		}
		if i < len(latestParts) {
			_, _ = fmt.Sscanf(latestParts[i], "%d", &l)
		}
		if l > c {
			return true
		}
		if l < c {
			return false
		}
	}

	return false
}

// findDownloadURL finds the appropriate download URL for the current platform
func findDownloadURL(assets []Asset) string {
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	// Build expected filename pattern
	var pattern string
	switch goos {
	case "darwin":
		pattern = fmt.Sprintf("lazycap_%s_%s.tar.gz", goos, goarch)
	case "linux":
		pattern = fmt.Sprintf("lazycap_%s_%s.tar.gz", goos, goarch)
	case "windows":
		pattern = fmt.Sprintf("lazycap_%s_%s.zip", goos, goarch)
	default:
		pattern = fmt.Sprintf("lazycap_%s_%s.tar.gz", goos, goarch)
	}

	// Find matching asset
	for _, asset := range assets {
		if strings.Contains(asset.Name, goos) && strings.Contains(asset.Name, goarch) {
			return asset.BrowserDownloadURL
		}
	}

	// Fallback: try exact pattern match
	for _, asset := range assets {
		if asset.Name == pattern {
			return asset.BrowserDownloadURL
		}
	}

	return ""
}

// SelfUpdate downloads and installs the latest version
func SelfUpdate(info *Info) error {
	if info.DownloadURL == "" {
		return fmt.Errorf("no download available for %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	// Get current executable path
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Resolve symlinks
	execPath, err = filepath.EvalSymlinks(execPath)
	if err != nil {
		return fmt.Errorf("failed to resolve executable path: %w", err)
	}

	// Download the new version
	client := &http.Client{Timeout: 5 * time.Minute}
	resp, err := client.Get(info.DownloadURL)
	if err != nil {
		return fmt.Errorf("failed to download update: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	// Create temp file for the new binary
	tmpDir := os.TempDir()
	tmpFile, err := os.CreateTemp(tmpDir, "lazycap-update-*")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer func() { _ = os.Remove(tmpPath) }()

	// Extract the binary from tar.gz (or handle zip for Windows)
	if strings.HasSuffix(info.DownloadURL, ".tar.gz") {
		if err := extractTarGz(resp.Body, tmpFile); err != nil {
			_ = tmpFile.Close()
			return fmt.Errorf("failed to extract update: %w", err)
		}
	} else {
		// For zip files or direct binaries
		if _, err := io.Copy(tmpFile, resp.Body); err != nil {
			_ = tmpFile.Close()
			return fmt.Errorf("failed to write update: %w", err)
		}
	}
	_ = tmpFile.Close()

	// Make executable
	if err := os.Chmod(tmpPath, 0755); err != nil {
		return fmt.Errorf("failed to set executable permission: %w", err)
	}

	// Backup old binary
	backupPath := execPath + ".bak"
	if err := os.Rename(execPath, backupPath); err != nil {
		return fmt.Errorf("failed to backup current binary: %w", err)
	}

	// Move new binary into place
	if err := copyFile(tmpPath, execPath); err != nil {
		// Restore backup on failure
		_ = os.Rename(backupPath, execPath)
		return fmt.Errorf("failed to install update: %w", err)
	}

	// Make sure new binary is executable
	if err := os.Chmod(execPath, 0755); err != nil {
		// Restore backup on failure
		_ = os.Remove(execPath)
		_ = os.Rename(backupPath, execPath)
		return fmt.Errorf("failed to set permissions: %w", err)
	}

	// Remove backup
	_ = os.Remove(backupPath)

	return nil
}

// extractTarGz extracts the lazycap binary from a tar.gz archive
func extractTarGz(r io.Reader, w io.Writer) error {
	gzr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer func() { _ = gzr.Close() }()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Look for the lazycap binary
		if header.Typeflag == tar.TypeReg &&
			(header.Name == "lazycap" || strings.HasSuffix(header.Name, "/lazycap")) {
			_, err := io.Copy(w, tr)
			return err
		}
	}

	return fmt.Errorf("lazycap binary not found in archive")
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() { _ = in.Close() }()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()

	_, err = io.Copy(out, in)
	return err
}

// VersionString returns a formatted version string
func VersionString(version string) string {
	if version == "" || version == "dev" {
		return "dev"
	}
	return strings.TrimPrefix(version, "v")
}
