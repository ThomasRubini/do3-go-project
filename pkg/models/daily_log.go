package models

import "time"

// DailyLog represents a user's daily food log
type DailyLog struct {
	Date  time.Time
	Meals []*Meal
}

// NutritionTotals represents the total nutritional values
type NutritionTotals struct {
	Calories float64
	Proteins float64
	Carbs    float64
	Fats     float64
	Fiber    float64
}

// CalculateTotals calculates the total nutritional values for the day
func (dl *DailyLog) CalculateTotals() *NutritionTotals {
	totals := &NutritionTotals{}
	for _, meal := range dl.Meals {
		for _, food := range meal.Foods {
			totals.Calories += food.Food.Calories * food.Quantity
			totals.Proteins += food.Food.Proteins * food.Quantity
			totals.Carbs += food.Food.Carbs * food.Quantity
			totals.Fats += food.Food.Fats * food.Quantity
			totals.Fiber += food.Food.Fiber * food.Quantity
		}
	}
	return totals
}
