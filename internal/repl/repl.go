package repl

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func readLine() (string, error) {
	fmt.Print("> ") // Prompt for input
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return line[:len(line)-1], nil // Remove the trailing newline
}

// Returns "false" when user exits
func processLine(line string) bool {
	fmt.Println("You entered:", line)
	s := strings.Split(line, " ")

	command := s[0]
	args := s[1:]

	switch command {
	case "exit":
		return false
	case "echo":
		fmt.Println("Echo:", strings.Join(args, " "))
	case "help":
		fmt.Println("Available commands: help, exit")
	default:
		fmt.Printf("Unknown command: %s. Type 'help' for list of commands\n", command)
	}
	return true
}

// Returns the exit code (0 for non error, non-zero for errors) when user exits
func StartRepl() int {
	for {
		// Read line
		line, err := readLine()
		if err != nil {
			fmt.Println("Error reading line:", err)
			return 1
		}

		// Skip empty lines
		if line == "" {
			continue
		}

		// Process line
		shouldContinue := processLine(line)
		if !shouldContinue {
			return 0
		}
	}
}
