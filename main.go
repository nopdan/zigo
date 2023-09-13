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
	case "list", "ls":
		c.List()
	case "remove", "rm":
		if len(os.Args) < 3 {
			fmt.Println("usage: ", os.Args[0], "remove <version>")
			fmt.Println("you can use 'ls' or 'list' to list all installed versions")
			return
		}
		c.Remove(os.Args[2])
	case "move", "mv":
		if len(os.Args) < 3 {
			fmt.Println("current install dir is: ", c.ZigDIR)
			return
		}
		c.Move(os.Args[2])
	case "help", "-h":
		help()
	default:
		c.Use(os.Args[1])
	}
}

// Print help message
func help() {
	fmt.Printf("zigo v1.0 (Download and manage Zig compilers)\n\n")
	fmt.Printf("Root Command:\n")
	fmt.Printf("  %-22s Download the specified version of zig compiler and set it as default\n", "zigo <version>")
	fmt.Println()
	fmt.Printf("Sub Commands:\n")
	fmt.Printf("  %-22s List installed compiler versions\n", "list, ls")
	fmt.Printf("  %-22s Remove the specified compiler\n", "remove, rm <version>")
	fmt.Printf("  %-22s Move the zig installation directory\n", "move, mv <install dir>")
	fmt.Printf("  %-22s Print help message\n", "help, -h")
}
