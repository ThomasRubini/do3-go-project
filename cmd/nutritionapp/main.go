package main

import (
	"bufio"
	"fmt"
	"nutritionapp/pkg/cli"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

// Message represents a command from the user
type Message struct {
	Command string
	Args    []string
}

func main() {
	// Load .env if it exists
	_ = godotenv.Load()

	// Set up data directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %v\n", err)
		os.Exit(1)
	}
	dataDir := filepath.Join(homeDir, ".nutritionapp")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		fmt.Printf("Error creating data directory: %v\n", err)
		os.Exit(1)
	}

	// Channel for communication between UI routine and main routine
	cmdChan := make(chan Message)

	// Create CLI commander
	commander, err := cli.NewCommander(dataDir)
	if err != nil {
		fmt.Printf("Error initializing commander: %v\n", err)
		os.Exit(1)
	}

	// Start UI routine
	go handleUserInput(cmdChan)

	// Main application loop
	for msg := range cmdChan {
		if strings.ToLower(msg.Command) == "exit" {
			fmt.Println("Exiting application...")
			close(cmdChan)
			return
		}
		commander.HandleCommand(msg.Command, msg.Args)
	}
}

func handleUserInput(cmdChan chan Message) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Welcome to NutritionApp!")
	fmt.Println("Type 'help' for available commands or 'exit' to quit")

	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			continue
		}

		// Clean the input
		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		// Parse the command and arguments
		parts := strings.Fields(input)
		msg := Message{
			Command: parts[0],
			Args:    parts[1:],
		}

		// Send the message through the channel
		cmdChan <- msg

		if strings.ToLower(msg.Command) == "exit" {
			return
		}
	}
}
