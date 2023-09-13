package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/mholt/archiver/v4"
)

type Info struct {
	IsMaster bool   // true if this is a master version
	Version  string //
	URL      string // URL to download
	Shasum   string // SHA256 checksum
	Size     string // file size
	FileName string //

	data []byte // downloaded data
}

// getIndex retrieves the download index from the given URL and returns it as a map.
func getIndex() map[string]map[string]interface{} {
	// Send a GET request to the URL
	url := "https://ziglang.org/download/index.json"
	r, err := http.Get(url)
	if err != nil {
		fmt.Printf("failed to get index\n")
		panic(err)
	}

	// Decode the response body and store it in the res map
	res := make(map[string]map[string]interface{})
	err = json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		fmt.Printf("failed to parse index\n")
		panic(err)
	}
	return res
}

// newInfo returns a new Info instance based on the provided version.
// It retrieves the information from an index, validates the version,
// and populates the Info struct with the relevant data.
func newInfo(version string) *Info {
	// Check if the version exists in the index
	index := getIndex()
	v, ok := index[version]
	if !ok {
		fmt.Printf("version: %s not found\n", version)
		os.Exit(1)
	}

	// Create a new Info instance
	info := new(Info)
	info.Version = version

	// If the version is "master", set the IsMaster flag and update the version
	if version == "master" {
		info.IsMaster = true
		info.Version = v["version"].(string)
	}

	// Check if the distribution info exists in the version data
	distInfo := getDistInfo()
	tmp, ok := v[distInfo]
	if !ok {
		fmt.Printf("unsupported dist: %s\n", distInfo)
		os.Exit(1)
	}

	// Get the distribution data
	dist := tmp.(map[string]interface{})
	info.URL = dist["tarball"].(string)
	info.Shasum = dist["shasum"].(string)
	info.Size = dist["size"].(string)

	// Extract the file name from the URL using a regular expression
	re, _ := regexp.Compile("zig-.+")
	info.FileName = re.FindString(info.URL)
	return info
}

// getDistInfo returns a string representing the distribution information of the system.
func getDistInfo() string {
	arch := runtime.GOARCH
	switch arch {
	case "amd64":
		arch = "x86_64"
	case "arm64":
		arch = "aarch64"
	}

	os := runtime.GOOS
	if os == "darwin" {
		os = "macos"
	}
	return arch + "-" + os
}

// install the specified version of Zig to the given Zig directory.
func (info *Info) install(ZigDIR string) {
	if info.IsMaster {
		fmt.Printf("installing master => %s\n", info.Version)
	} else {
		fmt.Printf("installing %s\n", info.Version)
	}

	// Detect the format of the archive and remove the file extension
	var format archiver.Extractor
	if strings.HasSuffix(info.FileName, ".zip") {
		format = archiver.Zip{}
		info.FileName = strings.TrimSuffix(info.FileName, ".zip")
	} else {
		format = archiver.CompressedArchive{
			Compression: archiver.Xz{},
			Archival:    archiver.Tar{},
		}
		info.FileName = strings.TrimSuffix(info.FileName, ".tag.xz")
	}

	// Create a reader for the archive data
	r := bytes.NewReader(info.data)

	// Extract the archive using the detected format
	err := format.Extract(context.Background(), r, nil, func(ctx context.Context, f archiver.File) error {
		// Create the full path for the extracted file or directory
		subName := strings.TrimPrefix(f.NameInArchive, info.FileName)
		name := filepath.Join(ZigDIR, info.Version, subName)
		if f.IsDir() {
			return os.MkdirAll(name, 0755)
		}
		rc, err := f.Open()
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
		fmt.Printf("failed to install %s\n", info.Version)
		// Delete version directory if installation fails
		os.RemoveAll(filepath.Join(ZigDIR, info.Version))
		panic(err)
	}
	fmt.Printf("successfully installed\n")
}
