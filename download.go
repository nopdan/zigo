package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/cavaliergopher/grab/v3"
)

// Download a file from a given URL and verifies its integrity using SHA256 checksum.
func (info *Info) download() {
	// Create a new HTTP client for file downloading
	client := grab.NewClient()

	// Get the user's cache directory
	dir, err := os.UserCacheDir()
	if err != nil {
		fmt.Printf("failed to get cache dir\n")
		panic(err)
	}
	dir = filepath.Join(dir, "zigo")

	// Check if the "zigo" directory exists
	fi, err := os.Stat(dir)
	if err != nil {
		// Create the "zigo" directory if it doesn't exist
		if os.IsNotExist(err) {
			os.MkdirAll(dir, 0755)
		} else {
			panic(err)
		}
	} else if !fi.IsDir() {
		// Remove the file with the same name if it exists
		os.Remove(dir)
		os.MkdirAll(dir, 0755)
	}

	// Create a new download request
	req, err := grab.NewRequest(dir, info.URL)
	if err != nil {
		fmt.Printf("failed to create request\n")
		panic(err)
	}

	// Verify sha256 checksum
	if info.Shasum != "" {
		h := sha256.New()
		sum, err := hex.DecodeString(info.Shasum)
		if err != nil {
			fmt.Printf("failed to decode sha256 checksum\n")
		} else {
			req.SetChecksum(h, sum, true)
		}
	}

	// Start download
	resp := client.Do(req)

	info.FileName = resp.Filename
	// Check if the download did resume
	if resp.DidResume {
		resp.Wait()
		fmt.Printf("load cache from %s\n", info.FileName)
	} else {
		fmt.Printf("downloading... %v\n", req.URL())
		for !resp.IsComplete() {
			time.Sleep(100 * time.Millisecond)
			progress(resp)
		}
		fmt.Println()
	}
	// Check errors
	if err := resp.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "download failed: %v\n", err)
		os.Exit(1)
	}

	if !resp.DidResume {
		fmt.Printf("done. save cache to %s\n", info.FileName)
	}
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
