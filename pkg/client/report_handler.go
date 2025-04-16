package client

import (
	"fmt"
	"nutritionapp/pkg/server"
)

func (c *Client) handleReport() {
	resp, err := makeRequestTyped[server.ReportResponse](c, server.ReqGetReport, nil)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fmt.Println("\n=== Daily Nutritional Report ===")
	fmt.Printf("Calories: %.0f kcal\n", resp.Calories)
	fmt.Printf("Proteins: %.1f g\n", resp.Proteins)
	fmt.Printf("Carbs: %.1f g\n", resp.Carbs)
	fmt.Printf("Fats: %.1f g\n", resp.Fats)
	fmt.Printf("Fiber: %.1f g\n", resp.Fiber)
}
