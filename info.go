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
	IsMaster bool
	Version  string
	URL      string
	Shasum   string
	Size     string
	FileName string
	data     []byte
}

// get download index.json
func getIndex() map[string]map[string]any {
	url := "https://ziglang.org/download/index.json"
	r, err := http.Get(url)
	if err != nil {
		fmt.Printf("failed to get index\n")
		panic(err)
	}
	res := make(map[string]map[string]any)
	err = json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		fmt.Printf("failed to parse index\n")
		panic(err)
	}
	return res
}

func newInfo(version string) *Info {
	index := getIndex()
	v, ok := index[version]
	if !ok {
		fmt.Printf("version: %s not found\n", version)
		os.Exit(1)
	}
	info := new(Info)
	info.Version = version
	// master version
	if version == "master" {
		info.IsMaster = true
		info.Version = v["version"].(string)
	}

	distInfo := getDistInfo()
	tmp, ok := v[distInfo]
	if !ok {
		fmt.Printf("unsupported dist: %s\n", distInfo)
		os.Exit(1)
	}

	dist := tmp.(map[string]any)
	info.URL = dist["tarball"].(string)
	info.Shasum = dist["shasum"].(string)
	info.Size = dist["size"].(string)

	re, _ := regexp.Compile("zig-.+")
	info.FileName = re.FindString(info.URL)
	return info
}

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

func (info *Info) install(ZigDIR string) {
	if info.IsMaster {
		fmt.Printf("installing master => %s\n", info.Version)
	} else {
		fmt.Printf("installing %s\n", info.Version)
	}
	// detect format and delete filename's extension
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

	r := bytes.NewReader(info.data)
	err := format.Extract(context.Background(), r, nil, func(ctx context.Context, f archiver.File) error {
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
		// delete version dir if install fails
		os.RemoveAll(filepath.Join(ZigDIR, info.Version))
		panic(err)
	}
	fmt.Printf("successfully installed\n")
}
