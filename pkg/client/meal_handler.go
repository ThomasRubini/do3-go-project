package client

import (
	"fmt"
	"nutritionapp/pkg/server"
)

func (c *Client) handleMeal(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: meal [add|list]")
		return
	}

	switch args[0] {
	case "add":
		fmt.Print("Meal name (breakfast/lunch/dinner/snack): ")
		name := c.readString()

		_, err := makeRequest(c, server.ReqAddMeal, server.AddMealData{Name: name})
		if err != nil {
			fmt.Printf("Error adding meal: %s\n", err)
			return
		}

		fmt.Printf("Added %s meal\n", name)

	case "list":
		resp, err := makeRequestTyped[server.MealListResponse](c, server.ReqListMeals, nil)
		if err != nil {
			fmt.Printf("Error fetching meal list: %s\n", err)
			return
		}

		c.displayMeals(*resp)
	}
}

func (c *Client) displayMeals(response server.MealListResponse) {
	if len(response.Meals) == 0 {
		fmt.Println("No meals recorded today.")
		return
	}

	fmt.Println("\n=== Today's Meals ===")
	for _, meal := range response.Meals {
		fmt.Printf("\n%s (at %s)\n", meal.Name, meal.Time)
		if len(meal.FoodItems) == 0 {
			fmt.Println("  No food items recorded")
			continue
		}
		for _, item := range meal.FoodItems {
			fmt.Printf("  - %s (%.0fg)\n", item.Name, item.Quantity)
		}
	}
}
