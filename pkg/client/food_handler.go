package client

import (
	"fmt"
	"nutritionapp/pkg/server"
)

func (c *Client) handleFood(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: food [search]")
		return
	}

	if args[0] == "search" {
		fmt.Print("Enter food name to search: ")
		query := c.readString()

		resp, err := makeRequestTyped[server.SearchFoodResponseData](c, server.ReqSearchFood, server.SearchFoodData{Query: query})
		if err != nil {
			fmt.Printf("Error searching for food: %s\n", err)
			return
		}

		if len(resp.Foods) == 0 {
			fmt.Println("No foods found matching your search.")
			return
		}

		c.displayFoodResults(resp.Foods)

		// Handle food selection and addition to meal
		fmt.Print("\nEnter number to add food (or 0 to cancel): ")
		choice := c.readInt()
		if choice <= 0 || choice > len(resp.Foods) {
			return
		}

		selectedFood := resp.Foods[choice-1]

		fmt.Print("Enter quantity in grams: ")
		quantity := c.readFloat()
		if quantity <= 0 {
			fmt.Println("Invalid quantity")
			return
		}

		// Get meal list to add food
		mealListResp, err := makeRequestTyped[server.MealListResponse](c, server.ReqListMeals, nil)
		if err != nil {
			fmt.Printf("Error fetching meal list: %s\n", err)
			return
		}

		if len(mealListResp.Meals) == 0 {
			fmt.Println("No meals available. Add a meal first using 'meal add'")
			return
		}

		fmt.Println("\nAvailable meals:")
		for i, meal := range mealListResp.Meals {
			fmt.Printf("%d. %s (%s)\n", i+1, meal.Name, meal.Time)
		}

		fmt.Print("Select meal number: ")
		mealIndex := c.readInt() - 1
		if mealIndex < 0 || mealIndex >= len(mealListResp.Meals) {
			fmt.Println("Invalid meal number")
			return
		}

		// Add food to meal
		_, err = makeRequest(c, server.ReqAddFood, server.AddFoodData{
			MealIndex: mealIndex,
			FoodID:    selectedFood.ID,
			Quantity:  quantity,
		})
		if err != nil {
			fmt.Printf("Error adding food to meal: %s\n", err)
			return
		}

		fmt.Printf("Added %.0fg of %s to meal\n", quantity, selectedFood.Name)
	}
}

func (c *Client) displayFoodResults(foods []server.FoodItem) {
	fmt.Println("\nSearch results:")
	for i, food := range foods {
		fmt.Printf("%d. %s\n", i+1, food.Name)
		fmt.Printf("   Per 100g: %.1f kcal, %.1fg protein, %.1fg carbs, %.1fg fat, %.1fg fiber\n",
			food.Calories, food.Proteins, food.Carbs, food.Fats, food.Fiber)
	}
}
