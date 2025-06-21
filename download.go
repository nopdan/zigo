package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cavaliergopher/grab/v3"
)

// Download a file from a given URL and verifies its integrity using SHA256 checksum.
func download(url, checksum string) string {
	// Create a new HTTP client for file downloading
	client := grab.NewClient()

	tmp := strings.Split(url, "/")
	filename := tmp[len(tmp)-1]
	name := filepath.Join(cacheDir, filename)

	// Create a new download request
	req, err := grab.NewRequest(name, url)
	if err != nil {
		p.Errorf("Failed to create download request: %s\n", err)
		os.Exit(1)
	}

	// Verify sha256 checksum
	if checksum != "" {
		h := sha256.New()
		sum, err := hex.DecodeString(checksum)
		if err == nil {
			req.SetChecksum(h, sum, true)
		}
	}

	// Start download
	resp := client.Do(req)

	if resp.IsComplete() {
		cInfo.Printf("Load cache from %s\n", name)
		return name
	}

	cInfo.Printf("Downloading %s...\n", filename)
	fmt.Printf("url: %s\n", url)
	fmt.Printf("save to: %s\n", name)
	for !resp.IsComplete() {
		time.Sleep(100 * time.Millisecond)
		progress(resp)
	}
	fmt.Println()

	// Check errors
	if err := resp.Err(); err != nil {
		p.Errorf("Failed to download file: %s\n", err)
		os.Exit(1)
	}

	cInfo.Println("Done.")
	return name
}

func progress(resp *grab.Response) {
	bytesComplete := decor(float64(resp.BytesComplete()))
	size := decor(float64(resp.Size()))
	bps := resp.BytesPerSecond()
	fmt.Printf("\rprogress: %s / %s | %.1f %% | %s/s %6s",
		bytesComplete, size,
		resp.Progress()*100,
		decor(bps),
		" ", // extra space to prevent truncation
	)
}

func decor(size float64) string {
	var unit string
	switch {
	case size < 1024:
		unit = "B"
	case size < 1024*1024:
		size /= 1024
		unit = "KiB"
	case size < 1024*1024*1024:
		size /= 1024 * 1024
		unit = "MiB"
	default:
		size /= 1024 * 1024 * 1024
		unit = "GiB"
	}
	return fmt.Sprintf("%.2f %s", size, unit)
}
