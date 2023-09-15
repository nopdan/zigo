package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		help()
		return
	}

	c := NewConfig()
	switch os.Args[1] {
	case "use":
		if len(os.Args) < 3 {
			fmt.Printf("usage: %s use <version>\n", os.Args[0])
			fmt.Println("you can use 'ls' or 'list' to list all installed versions")
			return
		}
		c.Use(os.Args[2])
	case "list", "ls":
		c.List()
	case "remove", "rm":
		switch len(os.Args) {
		case 2:
			fmt.Printf("usage: %s remove <version>\n", os.Args[0])
			fmt.Println("you can use 'ls' or 'list' to list all installed versions")
		case 3:
			if os.Args[2] == "--all" || os.Args[2] == "-a" {
				c.RemoveAll()
			} else {
				c.Remove(os.Args[2], false)
			}
		case 4:
			if os.Args[3] == "--force" || os.Args[3] == "-f" {
				c.Remove(os.Args[2], true)
			} else {
				fmt.Printf("undefined argument: %s\n", os.Args[3])
			}
		}
	case "clean":
		c.Clean()
	case "move", "mv":
		if len(os.Args) < 3 {
			fmt.Printf("current install dir is: %s", c.ZigDIR)
			return
		}
		c.Move(os.Args[2])
	case "help", "-h":
		help()
	default:
		c.Install(os.Args[1])
	}
}

// Print help message
func help() {
	fmt.Printf("zigo v1.4 (Download and manage Zig compilers)\n\n")
	fmt.Printf("Root Command:\n")
	fmt.Printf("  %-22s Download the specified version of zig compiler and set it as default\n", "zigo <version>")
	fmt.Println()
	fmt.Printf("Sub Commands:\n")
	fmt.Printf("  %-22s Set the specific installed version as default\n", "use <version>")
	fmt.Printf("  %-22s List installed compiler versions\n", "ls, list")
	fmt.Printf("  %-22s Remove the specified compiler, -f means force\n", "rm <version> [-f]")
	fmt.Printf("  %-22s Remove all installed compilers\n", "   [--all | -a]")
	fmt.Printf("  %-22s Clean up unused dev version compilers\n", "clean")
	fmt.Printf("  %-22s Move the zig installation directory\n", "mv, move <install-dir>")
	fmt.Printf("  %-22s Print help message\n", "help, -h")
}
