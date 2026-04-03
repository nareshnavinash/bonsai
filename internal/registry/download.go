package registry

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// Download fetches a GGUF file from HuggingFace with progress reporting.
// Writes to ~/.bonsai/models/{filename}, creating directories as needed.
func Download(ctx context.Context, model *BonsaiModel, progressFn func(downloaded, total int64)) error {
	dir := ModelsDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("cannot create models directory: %w", err)
	}

	destPath := model.LocalPath()
	partialPath := destPath + ".partial"

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, model.DownloadURL(), nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("download request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: HTTP %d", resp.StatusCode)
	}

	total := resp.ContentLength

	// Write to partial file
	out, err := os.Create(partialPath)
	if err != nil {
		return fmt.Errorf("cannot create file: %w", err)
	}

	var downloaded int64
	buf := make([]byte, 32*1024)
	for {
		n, readErr := resp.Body.Read(buf)
		if n > 0 {
			if _, writeErr := out.Write(buf[:n]); writeErr != nil {
				out.Close()
				os.Remove(partialPath)
				return fmt.Errorf("write error: %w", writeErr)
			}
			downloaded += int64(n)
			if progressFn != nil {
				progressFn(downloaded, total)
			}
		}
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			out.Close()
			os.Remove(partialPath)
			return fmt.Errorf("download interrupted: %w", readErr)
		}
	}
	out.Close()

	// Atomic rename
	if err := os.Rename(partialPath, destPath); err != nil {
		// Fallback: try copy if rename fails (cross-device)
		return fmt.Errorf("cannot finalize download: %w", err)
	}

	// Verify file exists at final path
	if _, err := os.Stat(destPath); err != nil {
		return fmt.Errorf("download completed but file not found at %s", destPath)
	}

	return nil
}

// IsDownloaded checks if a model's GGUF file exists locally.
func IsDownloaded(model *BonsaiModel) bool {
	_, err := os.Stat(model.LocalPath())
	return err == nil
}

// ModelFileName returns just the filename from a full model path.
func ModelFileName(path string) string {
	return filepath.Base(path)
}
