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

	zigoPath string // zigo.json path
	cacheDir string // cache dir
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
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		fmt.Printf("failed to get cache dir\n")
		panic(err)
	}
	c.cacheDir = filepath.Join(cacheDir, "zigo")
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

	return c
}

// Get installed versions
func (c *Config) versions() []string {
	files, err := os.ReadDir(c.ZigDIR)
	if err != nil {
		fmt.Printf("failed to read zig dir\n")
		panic(err)
	}
	ret := make([]string, 0, len(files))
	for _, f := range files {
		if f.Name() == "current" {
			continue
		}
		ret = append(ret, f.Name())
	}
	return ret
}

// List prints the list of installed versions in the Config object.
func (c *Config) List() {
	// Check if there are no installed versions
	verions := c.versions()
	if len(verions) == 0 {
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
	for _, v := range verions {
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
	versions := c.versions()
	if slices.Contains(versions, version) {
		c.Use(version)
		return
	}

	// Get the version info from url
	info := NewInfo(version)
	// No update on master
	if info.IsMaster && slices.Contains(versions, info.Version) {
		c.Master = info.Version
		c.Use("master")
		return
	}

	// Download and install to ZigDIR
	info.Install(c.ZigDIR)
	if info.IsMaster {
		c.Master = info.Version
	}
	c.link(info.Version)
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
	if !slices.Contains(c.versions(), specific) {
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
// If force is true, the version will be removed regardless of any conditions.
// If force is false, the version will only be removed if it is not the current
// version or the version pointed to by the master version.
func (c *Config) Remove(version string, force bool) {
	// If force is true, remove the version regardless of any conditions.
	if force {
		c.remove(version)
	} else {
		// Check if the version is the current version.
		if version == c.Current {
			fmt.Printf("cannot remove the version you are using.\n")
			return
		}
		// Check if the version is the version pointed to by the master version.
		if version == c.Master && c.Current == "master" {
			fmt.Printf("cannot remove this version (pointed to by the master version)\n")
			return
		}

		c.remove(version)
	}

	// Get the list of available versions.
	versions := c.versions()

	// If the current version is not in the list of available versions,
	// update the current version to an empty string.
	if !slices.Contains(versions, c.Current) && c.Current != "master" {
		c.Current = ""
	}

	// If the master version is not in the list of available versions,
	// update the master version to an empty string.
	if !slices.Contains(versions, c.Master) {
		c.Master = ""
		if c.Current == "master" {
			c.Current = ""
		}
	}

	// Write the changes to the config file.
	c.write()
}

// remove the specified version of the compilers.
func (c *Config) remove(version string) {
	// Determine the directory of the version to be removed
	dir := filepath.Join(c.ZigDIR, version)
	// Handle the special case of removing the "master" version
	if version == "master" {
		if c.Master == "" {
			return
		}
		dir = filepath.Join(c.ZigDIR, c.Master)
		fmt.Printf("removing master => %s... ", c.Master)
	} else if slices.Contains(c.versions(), version) {
		fmt.Printf("removing %s... ", version)
	} else {
		fmt.Printf("version: %s not found\n", version)
		os.Exit(1)
	}
	// Remove the directory
	err := os.RemoveAll(dir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("done.")
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
	for _, version := range c.versions() {
		c.Remove(version, true)
	}
	c.write()
}

// Clean up unused dev version compilers.
func (c *Config) Clean() {
	for _, version := range c.versions() {
		if strings.Contains(version, "dev") &&
			version != c.Master && version != c.Current {

			c.Remove(version, false)
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
