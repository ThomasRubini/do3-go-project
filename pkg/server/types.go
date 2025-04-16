package server

// Request represents a request from client to server
type Request struct {
	Type   string
	Data   any
	Return chan Response
}

type Response struct {
	Data  any
	Error error
}

// Request Types
const (
	ReqCreateProfile = "create_profile"
	ReqGetProfile    = "get_profile"
	ReqUpdateProfile = "update_profile"
	ReqAddMeal       = "add_meal"
	ReqListMeals     = "list_meals"
	ReqSearchFood    = "search_food"
	ReqAddFood       = "add_food"
	ReqGetReport     = "get_report"
)

// Request Data Types
type CreateProfileData struct {
	FirstName string
	LastName  string
	Age       int
	Weight    float64
	Height    float64
	Gender    string
	Goal      string
}

type UpdateProfileData struct {
	CreateProfileData
}

type AddMealData struct {
	Name string
}

type SearchFoodData struct {
	Query string
}

type AddFoodData struct {
	MealIndex int
	FoodID    string
	Quantity  float64
}

// Response Types
type ProfileResponseData struct {
	FirstName   string
	LastName    string
	Age         int
	Weight      float64
	Height      float64
	Gender      string
	Goal        string
	BMI         float64
	BodyFatPerc float64
}

type SearchFoodResponseData struct {
	Foods []FoodItem
}

type FoodItem struct {
	ID       string
	Name     string
	Calories float64
	Proteins float64
	Carbs    float64
	Fats     float64
	Fiber    float64
}

type MealListResponse struct {
	Meals []MealInfo
	Error string
}

type MealInfo struct {
	Index     int
	Name      string
	Time      string
	FoodItems []FoodItemInfo
}

type FoodItemInfo struct {
	Name     string
	Quantity float64
	Calories float64
	Proteins float64
	Carbs    float64
	Fats     float64
	Fiber    float64
}

type ReportResponse struct {
	Calories float64
	Proteins float64
	Carbs    float64
	Fats     float64
	Fiber    float64
	Error    string
}
