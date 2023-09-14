package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"slices"
)

type Config struct {
	ZigDIR  string `json:"install-dir"` // directory to save zig compiler
	Current string `json:"current"`     // current version
	Master  string `json:"master"`      // the version that master links to

	versions []string // installed versions
	zigoPath string   // zigo.json path
}

// NewConfig returns a new Config object.
func NewConfig() *Config {
	// Get the user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("failed to get home dir\n")
		panic(err)
	}

	// Create a new Config object
	c := new(Config)
	c.zigoPath = filepath.Join(homeDir, ".config/zigo.json")

	// Create zigo.json if it doesn't exist
	_, err = os.Stat(c.zigoPath)
	if err != nil {
		if os.IsNotExist(err) {
			c.ZigDIR = filepath.Join(homeDir, ".zig/")
			c.write()
		} else {
			panic(err)
		}
	}

	// Read zigo.json
	data, err := os.ReadFile(c.zigoPath)
	if err != nil {
		fmt.Printf("failed to read zigo.json\n")
		panic(err)
	}
	err = json.Unmarshal(data, &c)
	if err != nil {
		fmt.Printf("failed to parse zigo.json\n")
		panic(err)
	}

	// Create zig dir if it doesn't exist
	_, err = os.Stat(c.ZigDIR)
	if err != nil {
		if os.IsNotExist(err) {
			err2 := os.MkdirAll(c.ZigDIR, 0755)
			if err2 != nil {
				fmt.Printf("failed to create zig dir\n")
				panic(err2)
			}
		} else {
			panic(err)
		}
	}

	// Get installed versions
	files, err := os.ReadDir(c.ZigDIR)
	if err != nil {
		fmt.Printf("failed to read zig dir\n")
		panic(err)
	}
	c.versions = make([]string, 0, len(files))
	for _, f := range files {
		if f.Name() == "current" {
			continue
		}
		c.versions = append(c.versions, f.Name())
	}

	return c
}

// List prints the list of installed versions in the Config object.
func (c *Config) List() {
	// Check if there are no installed versions
	if len(c.versions) == 0 {
		fmt.Println("no installed versions")
		return
	}

	// Print the master version if it exists
	if c.Master != "" {
		// Print the master version with an asterisk if it is the current version
		if c.Current == "master" {
			fmt.Printf("* %s => %s\n", "master", c.Master)
		} else {
			fmt.Printf("  %s => %s\n", "master", c.Master)
		}
	}

	// Print the list of installed versions
	for _, v := range c.versions {
		// Print the version with an asterisk if it is the current version
		if v == c.Current {
			fmt.Printf("* %s\n", v)
			continue
		}
		fmt.Printf("  %s\n", v)
	}
}

// Install the specified version of the application.
// If the version is already installed, it is set as the current version,
// if not, it is downloaded, installed, and set as the current version.
func (c *Config) Install(version string) {
	// Check if the specified version is in the list of installed versions
	if slices.Contains(c.versions, version) {
		c.Use(version)
		return
	}

	// Get the version info from url
	info := NewInfo(version)
	// No update on master
	if info.IsMaster && slices.Contains(c.versions, info.Version) {
		c.Master = info.Version
		c.Use("master")
		return
	}

	// Download and install to ZigDIR
	info.Install(c.ZigDIR)
	c.link(info.Version)
	c.Current = version
	if info.IsMaster {
		c.Master = info.Version
	}
	c.write()
}

// Use set the specific installed version as default.
func (c *Config) Use(version string) {
	specific := version
	if specific == "master" {
		specific = c.Master
	}
	// Check if the specified version is in the list of available versions
	if !slices.Contains(c.versions, specific) {
		fmt.Printf("version: %s not found\n", specific)
		os.Exit(1)
	}
	c.link(specific)
	c.Current = version
	c.write()
	if version == "master" {
		fmt.Printf("using master => %s\n", c.Master)
	} else {
		fmt.Printf("using %s\n", version)
	}
}

// Remove the specified version of the compilers.
func (c *Config) Remove(version string) {
	// Check if the specified version is the current version
	if c.Current == version {
		fmt.Printf("cannot remove current version\n")
		return
	}

	// Check if the specified version is the master version
	if c.Current == "master" && c.Master == version {
		fmt.Printf("cannot remove this version (pointed to by the master version)\n")
		return
	}

	// Determine the directory of the version to be removed
	dir := filepath.Join(c.ZigDIR, version)
	if version == "master" {
		dir = filepath.Join(c.ZigDIR, c.Master)
	}

	// Remove the version directory
	if version == "master" {
		fmt.Printf("removing master => %s... ", c.Master)
	} else {
		fmt.Printf("removing %s... ", version)
	}
	err := os.RemoveAll(dir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Update the master version if the master version was removed
	if version == "master" {
		c.Master = ""
	}
	fmt.Println("done.")
	c.write()
}

// Remove all installed compilers.
func (c *Config) RemoveAll() {
	yes := []string{"y", "Y", "yes", "Yes"}
	var input string
	fmt.Printf("Are you sure you want to delete all compilers? (Yes/No): ")
	fmt.Scanf("%s", &input)
	if !slices.Contains(yes, input) {
		return
	}
	c.Current = ""
	c.Master = ""
	for _, version := range c.versions {
		c.Remove(version)
	}
	c.write()
}

// Clean up unused dev version compilers.
func (c *Config) Clean() {
	for _, version := range c.versions {
		if strings.Contains(version, "dev") && version != c.Master {
			c.Remove(version)
		}
	}
}

// Move zig install directory to the specified location.
func (c *Config) Move(dir string) {
	// Rename the directory to the specified location
	err := os.Rename(c.ZigDIR, dir)
	if err != nil {
		fmt.Printf("failed to move zig dir\n")
		os.Exit(1)
	}

	fmt.Printf("moved %s to %s\n", c.ZigDIR, dir)
	c.ZigDIR = dir
	c.write()

	// Link the appropriate directory based on the current branch
	if c.Current == "master" {
		c.link(c.Master)
	} else if c.Current != "" {
		c.link(c.Current)
	}
}

// Save the configuration to the zigo.json file.
func (c *Config) write() {
	// Marshal the configuration into JSON format.
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		fmt.Printf("failed to marshal zigo.json\n")
		panic(err)
	}

	// Write the JSON data to the zigo.json file.
	err = os.WriteFile(c.zigoPath, data, 0644)
	if err != nil {
		fmt.Printf("failed to write zigo.json\n")
		panic(err)
	}
}

// Link "dir" to "current"
func (c *Config) link(dir string) {
	// Remove the existing symlink
	os.RemoveAll(filepath.Join(c.ZigDIR, "current"))

	// Create a new symlink from "dir" to "current"
	err := os.Symlink(filepath.Join(c.ZigDIR, dir), filepath.Join(c.ZigDIR, "current"))
	if err != nil {
		fmt.Printf("failed to create symlink %v\n", err)
		os.Exit(1)
	}
}
