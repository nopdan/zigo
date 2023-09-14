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
		if len(os.Args) < 3 {
			fmt.Printf("usage: %s remove <version>\n", os.Args[0])
			fmt.Println("you can use 'ls' or 'list' to list all installed versions")
			return
		}
		c.Remove(os.Args[2])
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
		c.Update(os.Args[1])
	}
}

// Print help message
func help() {
	fmt.Printf("zigo v1.1 (Download and manage Zig compilers)\n\n")
	fmt.Printf("Root Command:\n")
	fmt.Printf("  %-22s Download the specified version of zig compiler and set it as default\n", "zigo <version>")
	fmt.Println()
	fmt.Printf("Sub Commands:\n")
	fmt.Printf("  %-22s Set the specific installed version as default\n", "use <version>")
	fmt.Printf("  %-22s List installed compiler versions\n", "ls, list")
	fmt.Printf("  %-22s Remove the specified compiler\n", "rm, remove <version>")
	fmt.Printf("  %-22s Clean up unused dev version compilers\n", "clean")
	fmt.Printf("  %-22s Move the zig installation directory\n", "mv, move <directory>")
	fmt.Printf("  %-22s Print help message\n", "help, -h")
}
