package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/mholt/archives"
)

// getIndex sends a GET request to the specified URL "https://ziglang.org/download/index.json",
// decodes the response body, and stores it in a map[string]map[string]interface{}.
// It handles errors related to getting and decoding the index.
func getIndex() (map[string]map[string]any, error) {
	// Send a GET request to the URL
	url := "https://ziglang.org/download/index.json"
	r, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get index.json: %w", err)
	}

	// Decode the response body and store it in the res map
	res := make(map[string]map[string]any)
	err = json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		return nil, fmt.Errorf("failed to decode index.json: %w", err)
	}
	return res, nil
}

// Fetch the specified version tarball and return its name and version
func fetch(version string) (name, ver string) {
	index, err := getIndex()
	if err != nil {
		p.Errorf("Failed to get index: %s\n", err)
		os.Exit(1)
	}

	v, ok := index[version]
	if !ok {
		p.Warnf("Version: %s Not found\n", version)
		os.Exit(1)
	}

	// Check if the distribution info exists in the version data
	tmpInfo, ok := v[distInfo]
	if !ok {
		p.Warnf("Unsupported dist: %s\n", distInfo)
		os.Exit(1)
	}

	dist := tmpInfo.(map[string]any)
	url := dist["tarball"].(string)
	shasum := dist["shasum"].(string)
	name = download(url, shasum)

	if version == "master" {
		ver = v["version"].(string)
	} else {
		ver = version
	}
	return
}

func extract(name string, dir string) {
	// Create a reader for the archive data
	r, err := os.Open(name)
	if err != nil {
		panic(err)
	}

	// Detect the format of the archive
	ctx := context.Background()
	format, rd, err := archives.Identify(ctx, name, r)
	if err != nil {
		panic(err)
	}
	ex, ok := format.(archives.Extractor)
	if !ok {
		p.Errorf("Failed to create extractor: %s\n", err)
		os.Exit(1)
	}
	err = ex.Extract(ctx, rd, func(ctx context.Context, info archives.FileInfo) error {
		// Create the full path for the extracted file or directory
		tmp := strings.Split(info.NameInArchive, "/")
		if len(tmp) < 2 {
			return fmt.Errorf("invalid file name: %s", info.NameInArchive)
		}
		name := filepath.Join(dir, strings.Join(tmp[1:], "/"))
		if info.IsDir() {
			return os.MkdirAll(name, 0755)
		}
		rc, err := info.Open()
		if err != nil {
			return err
		}
		f2, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY, 0755)
		if err != nil {
			return err
		}
		_, err = io.Copy(f2, rc)
		rc.Close()
		return err
	})

	if err != nil {
		p.Errorf("Failed to extract archive: %s\n", err)
		os.RemoveAll(dir)
		os.Exit(1)
	}
}
