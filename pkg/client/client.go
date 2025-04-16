package client

import (
	"bufio"
	"fmt"
	"nutritionapp/pkg/server"
	"os"
	"strings"
)

// Client handles user interaction through the terminal
type Client struct {
	requests chan server.Request
	reader   *bufio.Reader
}

// NewClient creates a new client instance
func NewClient(requests chan server.Request) *Client {
	return &Client{
		requests: requests,
		reader:   bufio.NewReader(os.Stdin),
	}
}

// Start begins the client's main loop
func (c *Client) Start() {
	fmt.Println("Welcome to NutritionApp!")
	fmt.Println("Type 'help' for available commands or 'exit' to quit")

	for {
		fmt.Print("> ")
		input, err := c.reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			continue
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		parts := strings.Fields(input)
		command := strings.ToLower(parts[0])
		args := parts[1:]

		if command == "exit" {
			fmt.Println("Goodbye!")
			close(c.requests)
			return
		}

		c.handleCommand(command, args)
	}
}

// handleCommand routes commands to their appropriate handlers
func (c *Client) handleCommand(command string, args []string) {
	switch command {
	case "help":
		c.showHelp()
	case "profile":
		c.handleProfile(args)
	case "meal":
		c.handleMeal(args)
	case "food":
		c.handleFood(args)
	case "report":
		c.handleReport()
	default:
		fmt.Printf("Unknown command: %s\n", command)
	}
}

// showHelp displays available commands
func (c *Client) showHelp() {
	fmt.Println("\nAvailable commands:")
	fmt.Println("  profile        - Show current profile")
	fmt.Println("  profile create - Create a new profile")
	fmt.Println("  meal add       - Add a new meal")
	fmt.Println("  meal list      - List today's meals")
	fmt.Println("  food search    - Search for food items")
	fmt.Println("  report         - Show daily nutritional report")
	fmt.Println("  help           - Show this help message")
	fmt.Println("  exit           - Exit the application")
}
