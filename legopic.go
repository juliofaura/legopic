package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/juliofaura/legopic/server"
)

///////////////////////////////////////////////////
// Main function
///////////////////////////////////////////////////
func main() {

	args := os.Args
	if len(args) != 2 {
		fmt.Println(
			`Error calling legopic
Usage: legopic webport
Example: legopic 8100`,
		)
		return
	}
	server.WEBPORT = args[1]

	server.StartWeb()

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Legopic console: ")
		s, _ := reader.ReadString('\n')
		command := strings.Fields(s)
		if len(command) >= 1 {
			switch command[0] {
			case "exit":
				fmt.Println("Have a nice day!")
				os.Exit(0)
			default:
				fmt.Printf("Unknown command %v\n", command)
			}
		}
	}
}
