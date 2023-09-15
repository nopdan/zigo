package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/cavaliergopher/grab/v3"
	"github.com/schollz/progressbar/v3"
)

// Download a file from a given URL and verifies its integrity using SHA256 checksum.
func (info *Info) download() {
	client := grab.NewClient()
	dir, err := os.UserCacheDir()
	if err != nil {
		fmt.Printf("failed to get cache dir\n")
		panic(err)
	}
	dir = filepath.Join(dir, "zigo")
	req, err := grab.NewRequest(dir, info.URL)
	if err != nil {
		fmt.Printf("failed to create request\n")
		panic(err)
	}

	// Start download
	resp := client.Do(req)

	// Check if the download did resume
	if resp.DidResume {
		resp.Wait()
		fmt.Printf("Load cache from %s\n", filepath.Join(dir, info.FileName))
	} else {
		fmt.Printf("Downloading %v...\n", req.URL())
		// Start UI loop
		t := time.NewTicker(200 * time.Millisecond)
		defer t.Stop()
		bar := progressbar.DefaultBytes(resp.Size())
	Loop:
		for {
			select {
			case <-t.C:
				bar.Set64(resp.BytesComplete())
			case <-resp.Done:
				break Loop
			}
		}
		bar.Set64(resp.BytesComplete())
	}

	// Check errors
	if err := resp.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Download failed: %v\n", err)
		os.Exit(1)
	}

	// Verify sha256 checksum
	if info.Shasum == "" {
		return
	}
	info.data, _ = resp.Bytes()
	h := sha256.New()
	h.Write(info.data)

	// Mismatch
	if fmt.Sprintf("%x", h.Sum(nil)) != info.Shasum {
		fmt.Printf("sha256 mismatch. want: %s, got: %x\n", info.Shasum, h.Sum(nil))
		// Remove downloaded file
		os.Remove(filepath.Join(dir, info.FileName))
		os.Exit(1)
	}

	fmt.Printf("Save cache to %s\n", filepath.Join(dir, info.FileName))
}
