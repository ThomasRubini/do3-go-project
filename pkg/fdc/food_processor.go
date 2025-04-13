package fdc

import (
	"fmt"
	"nutritionapp/pkg/models"
	"sync"
)

// FoodProcessor handles food data retrieval and caching
type FoodProcessor struct {
	client *Client
	cache  map[int]*models.Food
	mu     sync.RWMutex
}

func NewFoodProcessor(apiKey string) *FoodProcessor {
	return &FoodProcessor{
		client: NewClient(apiKey),
		cache:  make(map[int]*models.Food),
	}
}

// SearchFoods searches for foods and returns them in our app's format
func (fp *FoodProcessor) SearchFoods(query string) ([]models.Food, error) {
	resp, err := fp.client.SearchFoods(query, "Foundation,SR Legacy") // Filter to standard reference foods
	if err != nil {
		return nil, fmt.Errorf("search failed: %v", err)
	}

	var foods []models.Food
	for _, item := range resp.Foods {
		food := fp.convertFoodItem(item)
		foods = append(foods, food)
		
		// Cache the food item
		fp.mu.Lock()
		fp.cache[item.FdcId] = &food
		fp.mu.Unlock()
	}

	return foods, nil
}

// GetFoodDetails gets detailed food information by FDC ID
func (fp *FoodProcessor) GetFoodDetails(fdcId int) (*models.Food, error) {
	// Check cache first
	fp.mu.RLock()
	if food, exists := fp.cache[fdcId]; exists {
		fp.mu.RUnlock()
		return food, nil
	}
	fp.mu.RUnlock()

	// If not in cache, fetch from API
	item, err := fp.client.GetFoodDetails(fdcId)
	if err != nil {
		return nil, fmt.Errorf("failed to get food details: %v", err)
	}

	food := fp.convertFoodItem(*item)
	
	// Cache the result
	fp.mu.Lock()
	fp.cache[fdcId] = &food
	fp.mu.Unlock()

	return &food, nil
}

// convertFoodItem converts FDC food item to our app's format
func (fp *FoodProcessor) convertFoodItem(item FoodItem) models.Food {
	food := models.Food{
		ID:       fmt.Sprintf("fdc_%d", item.FdcId),
		Name:     item.Description,
		Calories: 0,
		Proteins: 0,
		Carbs:    0,
		Fats:     0,
		Fiber:    0,
	}

	// Map nutrient values
	for _, nutrient := range item.Nutrients {
		switch nutrient.Name {
		case "Energy":
			food.Calories = nutrient.Amount
		case "Protein":
			food.Proteins = nutrient.Amount
		case "Carbohydrate, by difference":
			food.Carbs = nutrient.Amount
		case "Total lipid (fat)":
			food.Fats = nutrient.Amount
		case "Fiber, total dietary":
			food.Fiber = nutrient.Amount
		}
	}

	return food
}
