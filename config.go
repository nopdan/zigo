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
	zigoPath string   // zigo.json
}

func NewConfig() *Config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("failed to get home dir\n")
		panic(err)
	}
	c := new(Config)
	c.zigoPath = filepath.Join(homeDir, ".config/zigo.json")

	// create zigo.json if not exists
	_, err = os.Stat(c.zigoPath)
	if err != nil {
		if os.IsNotExist(err) {
			c.ZigDIR = filepath.Join(homeDir, ".zig/")
			c.write()
		} else {
			panic(err)
		}
	}

	// read zigo.json
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

	// create zig dir
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

	// get installed versions
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

func (c *Config) List() {
	if len(c.versions) == 0 {
		fmt.Println("no installed versions")
		return
	}
	if c.Master != "" {
		if c.Current == "master" {
			fmt.Printf("* %s => %s\n", "master", c.Master)
		} else {
			fmt.Printf("  %s => %s\n", "master", c.Master)
		}
	}

	for _, v := range c.versions {
		if v == c.Current {
			fmt.Printf("* %s\n", v)
			continue
		}
		fmt.Printf("  %s\n", v)
	}
}

func (c *Config) Use(version string) {
	if slices.Contains(c.versions, version) {
		c.link(version)
		c.write()
		fmt.Printf("using %s\n", version)
		c.Current = version
		return
	}
	info := newInfo(version)
	if !slices.Contains(c.versions, info.Version) {
		info.download()
		info.install(c.ZigDIR)
	}
	if len(info.data) == 0 {
		fmt.Printf("using master => %s\n", info.Version)
	}
	c.link(info.Version)
	if info.IsMaster {
		c.Master = info.Version
	}
	c.write()
	c.Current = version
}

func (c *Config) Remove(version string) {
	if c.Current == version {
		fmt.Printf("cannot remove current version\n")
		return
	}
	if c.Current == "master" && c.Master == version {
		fmt.Printf("cannot remove this version (pointed to by the master version)\n")
		return
	}
	dir := filepath.Join(c.ZigDIR, version)
	if version == "master" {
		dir = filepath.Join(c.ZigDIR, c.Master)
	}
	err := os.RemoveAll(dir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if version == "master" {
		fmt.Printf("removed master => %s\n", c.Master)
		c.Master = ""
	} else {
		fmt.Printf("removed %s\n", version)
	}
	c.write()
}

func (c *Config) write() {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		fmt.Printf("failed to marshal zigo.json\n")
		panic(err)
	}
	err = os.WriteFile(c.zigoPath, data, 0644)
	if err != nil {
		fmt.Printf("failed to write zigo.json\n")
		panic(err)
	}
}

func (c *Config) link(dir string) {
	os.RemoveAll(filepath.Join(c.ZigDIR, "current"))
	err := os.Symlink(filepath.Join(c.ZigDIR, dir), filepath.Join(c.ZigDIR, "current"))
	if err != nil {
		fmt.Printf("failed to create symlink %v\n", err)
		os.Exit(1)
	}
}
