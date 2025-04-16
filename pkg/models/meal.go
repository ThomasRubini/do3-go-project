package models

import "time"

// Meal represents a meal with a list of foods
type Meal struct {
	Name  string
	Time  time.Time
	Foods []FoodQuantity
}

// FoodQuantity represents a food item with its quantity
type FoodQuantity struct {
	Food     *Food
	Quantity float64
}

// Food represents a food item with its nutritional values
type Food struct {
	ID       string
	Name     string
	Calories float64
	Proteins float64
	Carbs    float64
	Fats     float64
	Fiber    float64
}

// AddFood adds a food item to the meal
func (m *Meal) AddFood(food *Food, quantity float64) {
	m.Foods = append(m.Foods, FoodQuantity{
		Food:     food,
		Quantity: quantity,
	})
}

func (m *Meal) CalculateTotals() NutritionalTotals {
	var totals NutritionalTotals
	for _, item := range m.Foods {
		multiplier := item.Quantity / 100 // Convert from per 100g to actual quantity
		totals.Calories += item.Food.Calories * multiplier
		totals.Proteins += item.Food.Proteins * multiplier
		totals.Carbs += item.Food.Carbs * multiplier
		totals.Fats += item.Food.Fats * multiplier
		totals.Fiber += item.Food.Fiber * multiplier
	}
	return totals
}

type NutritionalTotals struct {
	Calories float64
	Proteins float64
	Carbs    float64
	Fats     float64
	Fiber    float64
}
