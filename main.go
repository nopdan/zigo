package main

import (
	"fmt"
	"os"
)

const version = "2.0.1"

func main() {
	if len(os.Args) == 1 {
		help()
		return
	}

	switch os.Args[1] {
	case "help", "-h":
		help()
	case "use":
		if len(os.Args) < 3 {
			fmt.Printf("Usage: %s use <version>\n", os.Args[0])
			return
		}
		Use(os.Args[2])
	case "list", "ls":
		List()
	case "remove", "rm":
		switch len(os.Args) {
		case 2:
			fmt.Printf("Usage: %s remove|rm <version>\n", os.Args[0])
		case 3:
			success := Remove(os.Args[2], false)
			if success {
				cInfo.Println("Done.")
			}
		case 4:
			if os.Args[3] == "--force" || os.Args[3] == "-f" {
				success := Remove(os.Args[2], true)
				if success {
					cInfo.Println("Done.")
				}
			} else {
				p.Warnf("Undefined argument: %s\n", os.Args[3])
			}
		default:
			p.Warnf("Undefined argument: %s\n", os.Args[3:])
		}
	case "clean":
		Clean()
	case "fetch":
		if len(os.Args) < 3 {
			fmt.Printf("Usage: %s fetch <version>\n", os.Args[0])
			return
		}
		Install(os.Args[2], false)
	default:
		Install(os.Args[1], true)
	}
}

// Print help message
func help() {
	cError.Printf("zigo %s (Download and manage Zig compilers)\n", version)
	fmt.Printf("zigo path: %s\n", zigoPath)
	fmt.Println()
	cInfo.Printf("Root Command:\n")
	fmt.Printf(
		"  %-22s Download the specified version of zig compiler and set it as default\n",
		"zigo <version>",
	)
	fmt.Println()
	cInfo.Printf("Sub Commands:\n")
	fmt.Printf("  %-22s Download the specified version of zig compiler\n", "fetch <version>")
	fmt.Printf("  %-22s Set the specific installed version as default\n", "use <version>")
	fmt.Printf("  %-22s List installed compiler versions\n", "ls, list")
	fmt.Printf("  %-22s Remove the specified compiler, -f means force\n", "rm <version> [-f]")
	fmt.Printf("  %-22s Clean up unused dev version compilers\n", "clean")
	fmt.Printf("  %-22s Print help message\n", "help, -h")
}
