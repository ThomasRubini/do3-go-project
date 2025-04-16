package client

import (
	"fmt"
	"strconv"
	"strings"
)

// Helper functions for reading input
func (c *Client) readString() string {
	var input string
	fmt.Scanln(&input)
	return strings.TrimSpace(input)
}

func (c *Client) readInt() int {
	input := c.readString()
	val, err := strconv.Atoi(input)
	if err != nil {
		return 0
	}
	return val
}

func (c *Client) readFloat() float64 {
	input := c.readString()
	val, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return 0
	}
	return val
}
