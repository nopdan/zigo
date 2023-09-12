package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) == 1 {
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
	case "help", "-h":
		fmt.Printf("Download and manage Zig binaries\n\n")
		fmt.Printf("Root Command:\n")
		fmt.Printf("  %-22s download and set the compiler as default\n", "zigo <version>")
		fmt.Println()
		fmt.Printf("Sub Commands:\n")
		fmt.Printf("  %-22s list installed compiler versions\n", "list, ls")
		fmt.Printf("  %-22s remove compiler\n", "remove, rm <version>")
		fmt.Printf("  %-22s \n", "help, -h")
		fmt.Println()
	default:
		c.Use(os.Args[1])
	}
}
