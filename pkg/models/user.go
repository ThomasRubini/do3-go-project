package models

// User represents a user profile
type User struct {
	FirstName string
	LastName  string
	Age       int
	Weight    float64
	Height    float64
	Gender    string
	Goal      string
	DailyLog  *DailyLog
}

// CalculateBMI calculates the user's BMI
func (u *User) CalculateBMI() float64 {
	heightInMeters := u.Height / 100
	return u.Weight / (heightInMeters * heightInMeters)
}

// CalculateBodyFat estimates body fat percentage using BMI
func (u *User) CalculateBodyFat() float64 {
	bmi := u.CalculateBMI()
	age := float64(u.Age)
	genderFactor := 0.0
	if u.Gender == "male" {
		genderFactor = 1.0
	}

	// Using the Deurenberg formula
	return (1.20 * bmi) + (0.23 * age) - (10.8 * genderFactor) - 5.4
}

func (u *User) EstimateBodyFat() float64 {
	bmi := u.CalculateBMI()
	if bmi <= 0 {
		return 0
	}

	// Very rough estimation based on BMI
	if u.Gender == "male" {
		return (1.20 * bmi) + (0.23 * float64(u.Age)) - 16.2
	}
	return (1.20 * bmi) + (0.23 * float64(u.Age)) - 5.4
}
