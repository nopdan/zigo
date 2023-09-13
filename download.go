package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	"time"

	"github.com/cavaliergopher/grab/v3"
	"github.com/schollz/progressbar/v3"
)

// Download a file from a given URL and verifies its integrity using SHA256 checksum.
func (i *Info) download() {
	client := grab.NewClient()
	req, _ := grab.NewRequest("", i.URL)
	req.NoStore = true // don't store the downloaded file

	// Start download
	resp := client.Do(req)
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

	// Check errors
	if err := resp.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Download failed: %v\n", err)
		os.Exit(1)
	}

	// Verify sha256 checksum
	if i.Shasum == "" {
		return
	}
	i.data, _ = resp.Bytes()
	h := sha256.New()
	h.Write(i.data)
	if fmt.Sprintf("%x", h.Sum(nil)) != i.Shasum {
		fmt.Printf("sha256 mismatch. want: %s, got: %x\n", i.Shasum, h.Sum(nil))
		os.Exit(1)
	}
}
