package cli

import (
	"fmt"
	"nutritionapp/pkg/database"
	"nutritionapp/pkg/fdc"
	"nutritionapp/pkg/models"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Commander struct {
	CurrentUser   *models.User
	DailyLog      *models.DailyLog
	userDB        *database.UserDatabase
	foodProcessor *fdc.FoodProcessor
	foodChan      chan foodSearchResult
}

type foodSearchResult struct {
	foods []models.Food
	err   error
}

func NewCommander(dataDir string) (*Commander, error) {
	// Get API key from environment
	apiKey := os.Getenv("FDC_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("FDC_API_KEY environment variable not set")
	}

	// Initialize food processor
	foodProcessor := fdc.NewFoodProcessor(apiKey)

	userDB, err := database.NewUserDatabase(filepath.Join(dataDir, "user.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize user database: %v", err)
	}

	cmd := &Commander{
		userDB:        userDB,
		CurrentUser:   userDB.GetUser(),
		DailyLog:      userDB.GetDailyLog(time.Now()),
		foodProcessor: foodProcessor,
		foodChan:      make(chan foodSearchResult),
	}

	return cmd, nil
}

func (c *Commander) HandleCommand(command string, args []string) {
	switch strings.ToLower(command) {
	case "profile":
		c.handleProfile(args)
	case "meal":
		c.handleMeal(args)
	case "food":
		c.handleFood(args)
	case "report":
		c.showDailyReport()
	case "help":
		c.showHelp()
	default:
		fmt.Printf("Unknown command: %s\n", command)
	}
}

func (c *Commander) handleProfile(args []string) {
	if len(args) == 0 {
		if c.CurrentUser == nil {
			fmt.Println("No profile set. Use 'profile create' to create one.")
			return
		}
		c.showProfile()
		return
	}

	switch args[0] {
	case "create":
		c.createProfile()
	case "update":
		c.updateProfile()
	default:
		fmt.Println("Unknown profile command. Use 'help' for usage.")
	}
}

func (c *Commander) createProfile() {
	var user models.User

	fmt.Print("First Name: ")
	fmt.Scanln(&user.FirstName)

	fmt.Print("Last Name: ")
	fmt.Scanln(&user.LastName)

	fmt.Print("Age: ")
	var age string
	fmt.Scanln(&age)
	user.Age, _ = strconv.Atoi(age)

	fmt.Print("Weight (kg): ")
	var weight string
	fmt.Scanln(&weight)
	user.Weight, _ = strconv.ParseFloat(weight, 64)

	fmt.Print("Height (cm): ")
	var height string
	fmt.Scanln(&height)
	user.Height, _ = strconv.ParseFloat(height, 64)

	fmt.Print("Gender (male/female): ")
	fmt.Scanln(&user.Gender)

	fmt.Print("Goal (weight loss/muscle gain/maintenance): ")
	fmt.Scanln(&user.Goal)

	c.CurrentUser = &user
	// Save to database
	if err := c.userDB.SaveUser(&user); err != nil {
		fmt.Printf("Warning: Failed to save profile: %v\n", err)
	}
	fmt.Println("Profile created successfully!")
	c.showProfile()
}

func (c *Commander) showProfile() {
	if c.CurrentUser == nil {
		fmt.Println("No profile set.")
		return
	}

	user := c.CurrentUser
	fmt.Println("\n=== Profile ===")
	fmt.Printf("Name: %s %s\n", user.FirstName, user.LastName)
	fmt.Printf("Age: %d\n", user.Age)
	fmt.Printf("Weight: %.1f kg\n", user.Weight)
	fmt.Printf("Height: %.1f cm\n", user.Height)
	fmt.Printf("Gender: %s\n", user.Gender)
	fmt.Printf("Goal: %s\n", user.Goal)
	fmt.Printf("BMI: %.1f\n", user.CalculateBMI())
	fmt.Printf("Estimated Body Fat: %.1f%%\n", user.EstimateBodyFat())
}

func (c *Commander) handleMeal(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: meal [add|list]")
		return
	}

	switch args[0] {
	case "add":
		c.addMeal()
	case "list":
		c.listMeals()
	default:
		fmt.Println("Unknown meal command. Use 'help' for usage.")
	}
}

func (c *Commander) addMeal() {
	var meal models.Meal

	fmt.Print("Meal name (breakfast/lunch/dinner/snack): ")
	fmt.Scanln(&meal.Name)

	meal.Time = time.Now()
	meal.ID = fmt.Sprintf("%s-%d", meal.Name, meal.Time.Unix())

	c.DailyLog.Meals = append(c.DailyLog.Meals, meal)
	fmt.Printf("Added %s meal\n", meal.Name)
}

func (c *Commander) listMeals() {
	if len(c.DailyLog.Meals) == 0 {
		fmt.Println("No meals recorded today.")
		return
	}

	fmt.Println("\n=== Today's Meals ===")
	for _, meal := range c.DailyLog.Meals {
		fmt.Printf("\n%s (at %s)\n", meal.Name, meal.Time.Format("15:04"))
		if len(meal.FoodItems) == 0 {
			fmt.Println("  No food items recorded")
			continue
		}
		for _, item := range meal.FoodItems {
			fmt.Printf("  - %s (%.0fg)\n", item.Food.Name, item.Quantity)
		}
	}
}

func (c *Commander) showDailyReport() {
	if len(c.DailyLog.Meals) == 0 {
		fmt.Println("No meals recorded today.")
		return
	}

	var totals models.NutritionalTotals
	for _, meal := range c.DailyLog.Meals {
		mealTotals := meal.CalculateTotals()
		totals.Calories += mealTotals.Calories
		totals.Proteins += mealTotals.Proteins
		totals.Carbs += mealTotals.Carbs
		totals.Fats += mealTotals.Fats
		totals.Fiber += mealTotals.Fiber
	}

	fmt.Println("\n=== Daily Nutritional Report ===")
	fmt.Printf("Calories: %.0f kcal\n", totals.Calories)
	fmt.Printf("Proteins: %.1f g\n", totals.Proteins)
	fmt.Printf("Carbs: %.1f g\n", totals.Carbs)
	fmt.Printf("Fats: %.1f g\n", totals.Fats)
	fmt.Printf("Fiber: %.1f g\n", totals.Fiber)
}

func (c *Commander) updateProfile() {
	if c.CurrentUser == nil {
		fmt.Println("No profile exists. Use 'profile create' first.")
		return
	}

	fmt.Println("\nUpdating profile (press Enter to keep current value):")

	fmt.Printf("First Name [%s]: ", c.CurrentUser.FirstName)
	if name := readString(); name != "" {
		c.CurrentUser.FirstName = name
	}

	fmt.Printf("Last Name [%s]: ", c.CurrentUser.LastName)
	if name := readString(); name != "" {
		c.CurrentUser.LastName = name
	}

	fmt.Printf("Age [%d]: ", c.CurrentUser.Age)
	if age := readInt(); age > 0 {
		c.CurrentUser.Age = age
	}

	fmt.Printf("Weight (kg) [%.1f]: ", c.CurrentUser.Weight)
	if weight := readFloat(); weight > 0 {
		c.CurrentUser.Weight = weight
	}

	fmt.Printf("Height (cm) [%.1f]: ", c.CurrentUser.Height)
	if height := readFloat(); height > 0 {
		c.CurrentUser.Height = height
	}

	fmt.Printf("Gender (male/female) [%s]: ", c.CurrentUser.Gender)
	if gender := readString(); gender != "" {
		c.CurrentUser.Gender = gender
	}

	fmt.Printf("Goal (weight loss/muscle gain/maintenance) [%s]: ", c.CurrentUser.Goal)
	if goal := readString(); goal != "" {
		c.CurrentUser.Goal = goal
	}

	fmt.Println("Profile updated successfully!")
	c.showProfile()
}

func (c *Commander) handleFood(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: food [add|search]")
		return
	}

	switch args[0] {
	case "add":
		c.addFood()
	case "search":
		c.searchFood()
	default:
		fmt.Println("Unknown food command. Use 'help' for usage.")
	}
}

func (c *Commander) addFood() {
	if len(c.DailyLog.Meals) == 0 {
		fmt.Println("No meals available. Add a meal first using 'meal add'")
		return
	}

	// List available meals
	fmt.Println("\nAvailable meals:")
	for i, meal := range c.DailyLog.Meals {
		fmt.Printf("%d. %s (%s)\n", i+1, meal.Name, meal.Time.Format("15:04"))
	}

	// Select meal
	fmt.Print("Select meal number: ")
	mealIndex := readInt() - 1
	if mealIndex < 0 || mealIndex >= len(c.DailyLog.Meals) {
		fmt.Println("Invalid meal number")
		return
	}

	// Create food item
	var food models.Food
	fmt.Print("Food name: ")
	food.Name = readString()

	fmt.Print("Calories (per 100g): ")
	food.Calories = readFloat()

	fmt.Print("Proteins (g per 100g): ")
	food.Proteins = readFloat()

	fmt.Print("Carbs (g per 100g): ")
	food.Carbs = readFloat()

	fmt.Print("Fats (g per 100g): ")
	food.Fats = readFloat()

	fmt.Print("Fiber (g per 100g): ")
	food.Fiber = readFloat()

	// Get quantity
	fmt.Print("Quantity (g): ")
	quantity := readFloat()

	// Add to meal
	c.DailyLog.Meals[mealIndex].AddFoodItem(food, quantity)
	// Save to database
	if err := c.userDB.SaveDailyLog(c.DailyLog); err != nil {
		fmt.Printf("Warning: Failed to save meal: %v\n", err)
	}
	fmt.Printf("Added %.0fg of %s to %s\n", quantity, food.Name, c.DailyLog.Meals[mealIndex].Name)
}

func (c *Commander) searchFood() {
	fmt.Print("Enter food name to search: ")
	query := readString()

	// Start async search
	go func() {
		foods, err := c.foodProcessor.SearchFoods(query)
		c.foodChan <- foodSearchResult{foods: foods, err: err}
	}()

	fmt.Println("Searching...")

	// Wait for results
	result := <-c.foodChan
	if result.err != nil {
		fmt.Printf("Error searching foods: %v\n", result.err)
		return
	}

	if len(result.foods) == 0 {
		fmt.Println("No foods found matching your search.")
		return
	}

	fmt.Println("\nSearch results:")
	for i, food := range result.foods {
		fmt.Printf("%d. %s\n", i+1, food.Name)
		fmt.Printf("   Per 100g: %.1f kcal, %.1fg protein, %.1fg carbs, %.1fg fat, %.1fg fiber\n",
			food.Calories, food.Proteins, food.Carbs, food.Fats, food.Fiber)
	}

	// Ask if user wants to add any of these foods
	fmt.Print("\nEnter number to add food (or 0 to cancel): ")
	choice := readInt()
	if choice <= 0 || choice > len(result.foods) {
		return
	}

	selectedFood := result.foods[choice-1]

	// Get quantity
	fmt.Print("Enter quantity in grams: ")
	quantity := readFloat()
	if quantity <= 0 {
		fmt.Println("Invalid quantity")
		return
	}

	// Add to meal
	if len(c.DailyLog.Meals) == 0 {
		fmt.Println("No meals available. Add a meal first using 'meal add'")
		return
	}

	// List available meals
	fmt.Println("\nAvailable meals:")
	for i, meal := range c.DailyLog.Meals {
		fmt.Printf("%d. %s (%s)\n", i+1, meal.Name, meal.Time.Format("15:04"))
	}

	// Select meal
	fmt.Print("Select meal number: ")
	mealIndex := readInt() - 1
	if mealIndex < 0 || mealIndex >= len(c.DailyLog.Meals) {
		fmt.Println("Invalid meal number")
		return
	}

	// Add to meal
	c.DailyLog.Meals[mealIndex].AddFoodItem(selectedFood, quantity)

	// Save to database
	if err := c.userDB.SaveDailyLog(c.DailyLog); err != nil {
		fmt.Printf("Warning: Failed to save meal: %v\n", err)
	}

	fmt.Printf("Added %.0fg of %s to %s\n", quantity, selectedFood.Name, c.DailyLog.Meals[mealIndex].Name)
}

// Helper functions for reading input
func readString() string {
	var input string
	fmt.Scanln(&input)
	return strings.TrimSpace(input)
}

func readInt() int {
	var input string
	fmt.Scanln(&input)
	val, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil {
		return 0
	}
	return val
}

func readFloat() float64 {
	var input string
	fmt.Scanln(&input)
	val, err := strconv.ParseFloat(strings.TrimSpace(input), 64)
	if err != nil {
		return 0
	}
	return val
}

func (c *Commander) showHelp() {
	fmt.Println("\nAvailable commands:")
	fmt.Println("  profile        - Show current profile")
	fmt.Println("  profile create - Create a new profile")
	fmt.Println("  profile update - Update existing profile")
	fmt.Println("  meal add       - Add a new meal")
	fmt.Println("  meal list      - List today's meals")
	fmt.Println("  food add       - Add food to a meal")
	fmt.Println("  food search    - Search for food items")
	fmt.Println("  report         - Show daily nutritional report")
	fmt.Println("  help          - Show this help message")
	fmt.Println("  exit          - Exit the application")
}
