package server

import (
	"fmt"
	"nutritionapp/pkg/models"
	"time"
)

func (s *Server) handleAddMeal(untypedData any) Response {
	data, ok := untypedData.(AddMealData)
	if !ok {
		return Response{Error: fmt.Errorf("invalid request data")}
	}

	dailyLog := s.userDB.GetDailyLog(time.Now())
	meal := models.Meal{
		Name:  data.Name,
		Time:  time.Now(),
		Foods: make([]models.FoodQuantity, 0),
	}

	dailyLog.Meals = append(dailyLog.Meals, &meal)
	if err := s.userDB.SaveDailyLog(dailyLog); err != nil {
		return Response{Error: fmt.Errorf("failed to save meal: %v", err)}
	}

	return Response{}
}

func (s *Server) handleListMeals(untypedData any) Response {
	dailyLog := s.userDB.GetDailyLog(time.Now())
	if dailyLog == nil {
		return Response{Error: fmt.Errorf("no meals found")}
	}

	var meals []MealInfo
	for i, meal := range dailyLog.Meals {
		var foodItems []FoodItemInfo
		for _, food := range meal.Foods {
			foodItems = append(foodItems, FoodItemInfo{
				Name:     food.Food.Name,
				Quantity: food.Quantity,
				Calories: food.Food.Calories,
				Proteins: food.Food.Proteins,
				Carbs:    food.Food.Carbs,
				Fats:     food.Food.Fats,
				Fiber:    food.Food.Fiber,
			})
		}

		meals = append(meals, MealInfo{
			Index:     i,
			Name:      meal.Name,
			Time:      meal.Time.Format("15:04"),
			FoodItems: foodItems,
		})
	}

	return Response{
		Data: MealListResponse{
			Meals: meals,
		},
	}
}
