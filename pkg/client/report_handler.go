package client

import (
	"fmt"
	"nutritionapp/pkg/server"
)

func (c *Client) handleReport() {
	respChannel := make(chan server.Response)
	c.requests <- server.Request{
		Type:   server.ReqGetReport,
		Return: respChannel,
	}

	resp := <-respChannel
	if resp.Error != nil {
		fmt.Printf("Error: %s\n", resp.Error)
		return
	}

	result, ok := resp.Data.(server.ReportResponse)
	if !ok {
		fmt.Println("Error: Invalid response data")
		return
	}

	fmt.Println("\n=== Daily Nutritional Report ===")
	fmt.Printf("Calories: %.0f kcal\n", result.Calories)
	fmt.Printf("Proteins: %.1f g\n", result.Proteins)
	fmt.Printf("Carbs: %.1f g\n", result.Carbs)
	fmt.Printf("Fats: %.1f g\n", result.Fats)
	fmt.Printf("Fiber: %.1f g\n", result.Fiber)
}
