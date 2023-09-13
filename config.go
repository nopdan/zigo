package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

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

// Update changes the current compiler version.
// If the specified version is already in the list of available versions,
// it links the version, sets it as the current version,
// and writes the updated configuration.
// If the version is not in the list, it creates a new info object,
// downloads and installs the version if necessary,
// and links the version.
// Finally, it sets the current version and writes the updated configuration.
func (c *Config) Update(version string) {
	// Check if the specified version is in the list of available versions
	if slices.Contains(c.versions, version) {
		c.link(version)
		fmt.Printf("using %s\n", version)
		c.Current = version
		c.write()
		return
	}

	info := newInfo(version)
	// Check if the specified version is not in the list of available versions
	if !slices.Contains(c.versions, info.Version) {
		// Download and install the version
		info.download()
		info.install(c.ZigDIR)
	}

	// Without download
	if info.data == nil {
		// Print the message indicating the master version being used
		fmt.Printf("using master => %s\n", info.Version)
	}
	c.link(info.Version)
	if info.IsMaster {
		c.Master = info.Version
	}
	c.Current = version
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
	err := os.RemoveAll(dir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Update the master version if the master version was removed
	if version == "master" {
		fmt.Printf("removed master => %s\n", c.Master)
		c.Master = ""
	} else {
		fmt.Printf("removed %s\n", version)
	}
	c.write()
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
