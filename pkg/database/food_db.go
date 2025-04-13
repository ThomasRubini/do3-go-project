package database

import (
	"encoding/json"
	"fmt"
	"nutritionapp/pkg/models"
	"os"
	"path/filepath"
	"strings"
)

type FoodDatabase struct {
	Foods []models.Food
	path  string
}

func NewFoodDatabase(dbPath string) (*FoodDatabase, error) {
	db := &FoodDatabase{
		path: dbPath,
	}

	// Create database directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %v", err)
	}

	// Load existing database or create new one
	if err := db.load(); err != nil {
		// If file doesn't exist, create initial database
		if os.IsNotExist(err) {
			db.Foods = getInitialFoods()
			if err := db.save(); err != nil {
				return nil, fmt.Errorf("failed to save initial database: %v", err)
			}
		} else {
			return nil, fmt.Errorf("failed to load database: %v", err)
		}
	}

	return db, nil
}

func (db *FoodDatabase) load() error {
	data, err := os.ReadFile(db.path)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &db.Foods)
}

func (db *FoodDatabase) save() error {
	data, err := json.MarshalIndent(db.Foods, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal foods: %v", err)
	}

	return os.WriteFile(db.path, data, 0644)
}

func (db *FoodDatabase) Search(query string) []models.Food {
	query = strings.ToLower(query)
	var results []models.Food

	for _, food := range db.Foods {
		if strings.Contains(strings.ToLower(food.Name), query) {
			results = append(results, food)
		}
	}

	return results
}

func (db *FoodDatabase) AddFood(food models.Food) error {
	food.ID = fmt.Sprintf("f%d", len(db.Foods)+1)
	db.Foods = append(db.Foods, food)
	return db.save()
}

// getInitialFoods returns a list of common foods to populate the database
func getInitialFoods() []models.Food {
	return []models.Food{
		{
			ID:       "f1",
			Name:     "Chicken Breast",
			Calories: 165,
			Proteins: 31,
			Carbs:    0,
			Fats:     3.6,
			Fiber:    0,
		},
		{
			ID:       "f2",
			Name:     "Oatmeal",
			Calories: 367,
			Proteins: 13.5,
			Carbs:    68,
			Fats:     7,
			Fiber:    10.5,
		},
		{
			ID:       "f3",
			Name:     "Banana",
			Calories: 89,
			Proteins: 1.1,
			Carbs:    22.8,
			Fats:     0.3,
			Fiber:    2.6,
		},
		{
			ID:       "f4",
			Name:     "Egg",
			Calories: 155,
			Proteins: 12.6,
			Carbs:    1.1,
			Fats:     11.3,
			Fiber:    0,
		},
		{
			ID:       "f5",
			Name:     "Salmon",
			Calories: 208,
			Proteins: 22,
			Carbs:    0,
			Fats:     13,
			Fiber:    0,
		},
	}
}
