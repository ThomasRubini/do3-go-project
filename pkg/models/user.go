package models

import "time"

// User represents a user profile with personal information
type User struct {
	FirstName    string
	LastName     string
	Age          int
	Weight       float64 // in kg
	Height       float64 // in cm
	Gender       string
	Goal         string  // e.g., "weight loss", "muscle gain", "maintenance"
	ActivityLevel string // e.g., "sedentary", "moderate", "active", "very active"
}

// BodyComposition represents physical measurements at a point in time
type BodyComposition struct {
	Date      time.Time
	Weight    float64
	BodyFat   float64 // percentage
	BMI       float64
	Timestamp time.Time
}

// CalculateBMI calculates the Body Mass Index
func (u *User) CalculateBMI() float64 {
	if u.Height <= 0 || u.Weight <= 0 {
		return 0
	}
	heightInMeters := u.Height / 100
	return u.Weight / (heightInMeters * heightInMeters)
}

// EstimateBodyFat estimates body fat percentage using BMI method
// Note: This is a simple estimation, not as accurate as other methods
func (u *User) EstimateBodyFat() float64 {
	bmi := u.CalculateBMI()
	age := float64(u.Age)
	
	// Simple body fat estimation based on BMI
	// This is a basic formula and should be replaced with more accurate methods
	var bodyFat float64
	if u.Gender == "male" {
		bodyFat = (1.20 * bmi) + (0.23 * age) - 16.2
	} else if u.Gender == "female" {
		bodyFat = (1.20 * bmi) + (0.23 * age) - 5.4
	}
	
	return bodyFat
}
