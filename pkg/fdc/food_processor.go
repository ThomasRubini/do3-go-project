package fdc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"nutritionapp/pkg/models"
)

type FoodProcessor struct {
	apiKey string
}

func NewFoodProcessor(apiKey string) *FoodProcessor {
	return &FoodProcessor{apiKey: apiKey}
}

func (fp *FoodProcessor) SearchFoods(query string) ([]models.Food, error) {
	baseURL := "https://api.nal.usda.gov/fdc/v1/foods/search"
	params := url.Values{}
	params.Add("api_key", fp.apiKey)
	params.Add("query", query)
	params.Add("pageSize", "10")
	params.Add("dataType", "SR Legacy")

	resp, err := http.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return nil, fmt.Errorf("failed to search foods: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		Foods []struct {
			FdcID         int    `json:"fdcId"`
			Description   string `json:"description"`
			FoodNutrients []struct {
				NutrientNumber string  `json:"nutrientNumber"`
				Value          float64 `json:"value"`
			} `json:"foodNutrients"`
		} `json:"foods"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	var foods []models.Food
	for _, f := range result.Foods {
		food := models.Food{
			ID:   fmt.Sprintf("fdc_%d", f.FdcID),
			Name: f.Description,
		}

		for _, n := range f.FoodNutrients {
			switch n.NutrientNumber {
			case "208":
				food.Calories = n.Value
			case "203":
				food.Proteins = n.Value
			case "205":
				food.Carbs = n.Value
			case "204":
				food.Fats = n.Value
			case "291":
				food.Fiber = n.Value
			}
		}

		foods = append(foods, food)
	}

	return foods, nil
}

func (fp *FoodProcessor) GetFoodDetails(fdcID string) (*models.Food, error) {
	baseURL := "https://api.nal.usda.gov/fdc/v1/food"
	params := url.Values{}
	params.Add("api_key", fp.apiKey)

	// Extract numeric ID from string (e.g., "fdc_123" -> "123")
	fdcNumericID := fdcID[4:]
	resp, err := http.Get(fmt.Sprintf("%s/%s?%s", baseURL, fdcNumericID, params.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to get food details: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		FdcID         int    `json:"fdcId"`
		Description   string `json:"description"`
		FoodNutrients []struct {
			NutrientNumber string  `json:"nutrientNumber"`
			Value          float64 `json:"value"`
		} `json:"foodNutrients"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	food := &models.Food{
		ID:   fmt.Sprintf("fdc_%d", result.FdcID),
		Name: result.Description,
	}

	for _, n := range result.FoodNutrients {
		switch n.NutrientNumber {
		case "208":
			food.Calories = n.Value
		case "203":
			food.Proteins = n.Value
		case "205":
			food.Carbs = n.Value
		case "204":
			food.Fats = n.Value
		case "291":
			food.Fiber = n.Value
		}
	}

	return food, nil
}
