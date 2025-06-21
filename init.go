package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	zigoPath string
	current  string // Current version
	master   string // Master version
)

func init() {
	// Read ZIGO_PATH environment variable if it exists
	zigoPath = os.Getenv("ZIGO_PATH")
	if zigoPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			cError.Println("Failed to get user's home directory")
			panic(err)
		}
		zigoPath = filepath.Join(homeDir, ".zig")
	}

	err := os.MkdirAll(zigoPath, 0755)
	if err != nil {
		cError.Println("Failed to create zigo directory")
		panic(err)
	}

	// Create config file if it doesn't exist
	configPath := filepath.Join(zigoPath, "config")
	_, err = os.Stat(configPath)
	if os.IsNotExist(err) {
		err = os.WriteFile(configPath, []byte("\n"), 0666)
		if err != nil {
			cError.Println("Failed to create config file")
			panic(err)
		}
	} else if err != nil {
		cError.Println("Failed to check if config file exists")
		panic(err)
	} else {
		config, err := os.ReadFile(configPath)
		if err != nil {
			cError.Println("Failed to read config file")
			panic(err)
		}
		lines := strings.Split(string(config), "\n")
		if len(lines) >= 2 {
			current = lines[0]
			master = lines[1]
		}
	}
}

func modifyConfig() {
	configPath := filepath.Join(zigoPath, "config")
	err := os.WriteFile(configPath, []byte(current+"\n"+master), 0666)
	if err != nil {
		cError.Println("Failed to modify config file")
		panic(err)
	}
}

var (
	// Version is the current version of the program.
	version  = "2.0.0"
	distInfo = getDistInfo()
	cacheDir = getCacheDir()
)

// getDistInfo returns a string representing the distribution information of the system.
func getDistInfo() string {
	arch := runtime.GOARCH
	switch arch {
	case "amd64":
		arch = "x86_64"
	case "amd64p32":
		arch = "x86"
	case "arm64":
		arch = "aarch64"
	}

	os := runtime.GOOS
	if os == "darwin" {
		os = "macos"
	}
	return arch + "-" + os
}

func getCacheDir() string {
	// Get the user's cache directory
	dir, err := os.UserCacheDir()
	if err != nil {
		cError.Println("Failed to get user's cache directory")
		panic(err)
	}

	dir = filepath.Join(dir, "zigo")
	if err := os.MkdirAll(dir, 0755); err != nil {
		cError.Println("Failed to create cache directory")
		panic(err)
	}
	return dir
}
