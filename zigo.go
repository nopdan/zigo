package main

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
)

// Get a list of installed verList of the zig compiler.
// The returned list is sorted in ascending order.
func verList() []string {
	files, err := os.ReadDir(zigoPath)
	if err != nil {
		cError.Println("Failed to read zig directory")
		panic(err)
	}

	ret := make([]string, 0, len(files))
	for _, f := range files {
		// Skip the "current" dir
		if f.Name() == "current" || f.Name() == "config" {
			continue
		}
		ret = append(ret, f.Name())
	}

	// Sort the versions in ascending order using the cmpVersion function
	slices.SortFunc(ret, cmpVersion)
	return ret
}

// compare versions
// 0.11.0
// 0.12.0-dev.1127+32bc07767
func cmpVersion(a, b string) int {
	simp := func(s string) string {
		// del "-dev"
		s = strings.Replace(s, "-dev", "", 1)
		// del hash
		idx := strings.LastIndexByte(s, '+')
		if idx != -1 {
			s = s[:idx]
		}
		return s
	}

	parse := func(s string) [4]int {
		s = simp(s)
		var ver [4]int
		sli := strings.Split(s, ".")
		for i, v := range sli {
			integer, _ := strconv.Atoi(v)
			ver[i] = integer
		}
		return ver
	}

	verA := parse(a)
	verB := parse(b)

	for i := 0; i < 4; i++ {
		if verA[i] < verB[i] {
			return -1
		}
		if verA[i] > verB[i] {
			return 1
		}
	}
	return 0
}

// Create a new symlink, removing the old one if it exists.
// dir: the parent directory of the symlink
func symlink(dir, version string) {
	dst := filepath.Join(dir, "current")
	src := filepath.Join(dir, version)

	// Remove the existing symlink
	err := os.RemoveAll(dst)
	if err != nil {
		cError.Printf("Failed to remove symlink %s\n", dst)
		panic(err)
	}

	// Create a new symlink
	err = os.Symlink(src, dst)
	if err != nil {
		cError.Printf("Failed to create symlink from %s to %s\n", src, dst)
		panic(err)
	}
}

// List prints the list of installed versions in the Config object.
func List() {
	// Check if there are no installed versions
	versions := verList()
	if len(versions) == 0 {
		cInfo.Println("No installed versions")
		return
	}

	// Print the master version if it exists
	if master != "" {
		// Print the master version with an asterisk if it is the current version
		if current == "master" {
			fmt.Printf("* %s => %s\n", "master", master)
		} else {
			fmt.Printf("  %s => %s\n", "master", master)
		}
	}

	// Print the list of installed versions
	for _, v := range versions {
		// Print the version with an asterisk if it is the current version
		if v == current {
			fmt.Printf("* %s\n", v)
			continue
		}
		fmt.Printf("  %s\n", v)
	}
}

// Set the specific installed version as default
func Use(version string) {
	isMaster := version == "master"
	if isMaster {
		version = master
	} else if version == current {
		// Check if the specified version is the current version
		return
	}

	// Check if the specified version is in the list of available versions
	if !slices.Contains(verList(), version) {
		cWarn.Printf("Version: %s Not found\n", version)
		os.Exit(1)
	}
	symlink(zigoPath, version)

	if isMaster {
		cInfo.Printf("Using master => %s\n", master)
		current = "master"
	} else {
		cInfo.Printf("Using %s\n", version)
		current = version
	}
	modifyConfig()
}

func Install(version string, asDefault bool) {
	versions := verList()
	if version == "master" {
		name, ver := fetch(version)
		if slices.Contains(versions, ver) {
			if asDefault {
				master = ver
				Use("master")
			}
		} else {
			cInfo.Printf("Installing master => %s... \n", ver)
			dir := filepath.Join(zigoPath, ver)
			extract(name, dir)
			if asDefault {
				master = ver
				Use("master")
			}
			cInfo.Println("Done.")
		}
	} else {
		if slices.Contains(versions, version) {
			if asDefault {
				Use(version)
			}
		} else {
			name, ver := fetch(version)
			cInfo.Printf("Installing %s... \n", ver)
			dir := filepath.Join(zigoPath, ver)
			extract(name, dir)
			if asDefault {
				Use(version)
			}
			cInfo.Println("Done.")
		}
	}
}

func Remove(version string, force bool) bool {
	if (version == current || version == master) && !force {
		cWarn.Printf("Cannot remove the version you are using.\n")
		return false
	}

	if version == "master" {
		if master == "" {
			cWarn.Printf("Version: master Not found\n")
			return false
		}
		cInfo.Printf("Removing master => %s... \n", master)
		err := os.RemoveAll(filepath.Join(zigoPath, master))
		if err == nil {
			master = ""
			if current == "master" {
				current = ""
				_ = os.RemoveAll(filepath.Join(zigoPath, "current"))
			}
		} else {
			cWarn.Printf("%s\n", err)
			return false
		}
	} else {
		if !slices.Contains(verList(), version) {
			cWarn.Printf("Version: %s Not found\n", version)
			return false
		}
		cInfo.Printf("Removing %s... \n", version)
		err := os.RemoveAll(filepath.Join(zigoPath, version))
		if err == nil {
			if current == version {
				current = ""
				_ = os.RemoveAll(filepath.Join(zigoPath, "current"))
			}
			if master == version {
				master = ""
			}
		} else {
			cWarn.Printf("%s\n", err)
			return false
		}
	}
	modifyConfig()
	return true
}

func Clean() {
	for _, v := range verList() {
		if strings.Contains(v, "dev") && v != master && v != current {
			Remove(v, false)
		}
	}
	cInfo.Println("Done.")
}
