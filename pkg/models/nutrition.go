package models

import "time"

// Food represents a food item with its nutritional information per 100g
type Food struct {
	ID          string
	Name        string
	Calories    float64
	Proteins    float64
	Carbs       float64
	Fats        float64
	Fiber       float64
	Vitamins    map[string]float64
	Minerals    map[string]float64
}

// Meal represents a collection of consumed foods with their quantities
type Meal struct {
	ID        string
	Name      string    // e.g., "Breakfast", "Lunch", "Dinner", "Snack"
	Time      time.Time
	FoodItems []ConsumedFood
}

// ConsumedFood represents a food item and its consumed quantity
type ConsumedFood struct {
	Food     Food
	Quantity float64 // in grams
}

// DailyLog represents all meals and nutritional totals for a day
type DailyLog struct {
	Date   time.Time
	Meals  []Meal
	Totals NutritionalTotals
}

// NutritionalTotals represents the total nutrients consumed
type NutritionalTotals struct {
	Calories float64
	Proteins float64
	Carbs    float64
	Fats     float64
	Fiber    float64
}

// CalculateNutrients calculates the actual nutrients based on quantity consumed
func (cf *ConsumedFood) CalculateNutrients() NutritionalTotals {
	factor := cf.Quantity / 100.0 // Convert to proportion of 100g
	return NutritionalTotals{
		Calories: cf.Food.Calories * factor,
		Proteins: cf.Food.Proteins * factor,
		Carbs:    cf.Food.Carbs * factor,
		Fats:     cf.Food.Fats * factor,
		Fiber:    cf.Food.Fiber * factor,
	}
}

// AddFoodItem adds a food item to the meal
func (m *Meal) AddFoodItem(food Food, quantity float64) {
	m.FoodItems = append(m.FoodItems, ConsumedFood{
		Food:     food,
		Quantity: quantity,
	})
}

// CalculateTotals calculates the total nutrients for the meal
func (m *Meal) CalculateTotals() NutritionalTotals {
	var totals NutritionalTotals
	for _, item := range m.FoodItems {
		nutrients := item.CalculateNutrients()
		totals.Calories += nutrients.Calories
		totals.Proteins += nutrients.Proteins
		totals.Carbs += nutrients.Carbs
		totals.Fats += nutrients.Fats
		totals.Fiber += nutrients.Fiber
	}
	return totals
}
