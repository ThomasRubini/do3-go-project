package server

import (
	"fmt"
	"time"
)

func (s *Server) handleSearchFood(untypedData any) Response {
	data, ok := untypedData.(SearchFoodData)
	if !ok {
		return Response{Error: fmt.Errorf("invalid request data")}
	}

	foods, err := s.foodProcessor.SearchFoods(data.Query)
	if err != nil {
		return Response{Error: fmt.Errorf("search failed: %v", err)}
	}

	var foodItems []FoodItem
	for _, f := range foods {
		foodItems = append(foodItems, FoodItem{
			ID:       f.ID,
			Name:     f.Name,
			Calories: f.Calories,
			Proteins: f.Proteins,
			Carbs:    f.Carbs,
			Fats:     f.Fats,
			Fiber:    f.Fiber,
		})
	}

	return Response{
		Data: SearchFoodResponseData{Foods: foodItems},
	}
}

func (s *Server) handleAddFood(untypedData any) Response {
	data, ok := untypedData.(AddFoodData)
	if !ok {
		return Response{Error: fmt.Errorf("invalid request data")}
	}

	dailyLog := s.userDB.GetDailyLog(time.Now())
	if data.MealIndex < 0 || data.MealIndex >= len(dailyLog.Meals) {
		return Response{Error: fmt.Errorf("invalid meal index")}
	}

	// Get food details from FDC
	fdcId := data.FoodID // Assuming format "fdc_123"
	food, err := s.foodProcessor.GetFoodDetails(fdcId)
	if err != nil {
		return Response{Error: fmt.Errorf("failed to get food details: %v", err)}
	}

	dailyLog.Meals[data.MealIndex].AddFood(food, data.Quantity)
	if err := s.userDB.SaveDailyLog(dailyLog); err != nil {
		return Response{Error: fmt.Errorf("failed to save food: %v", err)}
	}

	return Response{}
}
