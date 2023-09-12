package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	"time"

	"github.com/cavaliergopher/grab/v3"
	"github.com/schollz/progressbar/v3"
)

func (i *Info) download() {
	// create client
	client := grab.NewClient()
	req, _ := grab.NewRequest("", i.URL)
	req.NoStore = true

	// start download
	fmt.Printf("Downloading %v...\n", req.URL())
	resp := client.Do(req)

	// start UI loop
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

	// check for errors
	if err := resp.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Download failed: %v\n", err)
		os.Exit(1)
	}
	// check sha256
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
