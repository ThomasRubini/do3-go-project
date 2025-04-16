package db

import (
	"database/sql"
	"encoding/json"
	"log"
	"nutritionapp/pkg/models"
	"time"
)

// UserDatabase defines the interface for database operations
type UserDatabase interface {
	GetUser() *models.User
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
	GetDailyLog(date time.Time) *models.DailyLog
	SaveDailyLog(log *models.DailyLog) error
	SaveUser(user *models.User) error
}

// SQLiteDB implements UserDatabase using SQLite3
type SQLiteDB struct {
	db *sql.DB
}

// NewSQLiteDB creates a new SQLite database instance
func NewSQLiteDB(path string) (*SQLiteDB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	if err := createTables(db); err != nil {
		return nil, err
	}
	return &SQLiteDB{db: db}, nil
}

func createTables(db *sql.DB) error {
	// Create users table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY,
			first_name TEXT NOT NULL,
			last_name TEXT NOT NULL,
			age INTEGER NOT NULL,
			weight REAL NOT NULL,
			height REAL NOT NULL,
			gender TEXT NOT NULL,
			goal TEXT NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	// Create daily_logs table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS daily_logs (
			id INTEGER PRIMARY KEY,
			date TEXT NOT NULL,
			meals TEXT NOT NULL,
			UNIQUE(date)
		)
	`)
	return err
}

// GetUser retrieves the user from the database
func (s *SQLiteDB) GetUser() *models.User {
	var user models.User
	err := s.db.QueryRow(`
		SELECT first_name, last_name, age, weight, height, gender, goal 
		FROM users 
		LIMIT 1
	`).Scan(&user.FirstName, &user.LastName, &user.Age, &user.Weight, &user.Height, &user.Gender, &user.Goal)

	if err != nil {
		return nil
	}
	return &user
}

// CreateUser creates a new user in the database
func (s *SQLiteDB) CreateUser(user *models.User) error {
	_, err := s.db.Exec(`
		INSERT INTO users (first_name, last_name, age, weight, height, gender, goal)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, user.FirstName, user.LastName, user.Age, user.Weight, user.Height, user.Gender, user.Goal)
	return err
}

// UpdateUser updates an existing user in the database
func (s *SQLiteDB) UpdateUser(user *models.User) error {
	_, err := s.db.Exec(`
		UPDATE users 
		SET first_name = ?, last_name = ?, age = ?, weight = ?, height = ?, gender = ?, goal = ?
		WHERE id = (SELECT id FROM users LIMIT 1)
	`, user.FirstName, user.LastName, user.Age, user.Weight, user.Height, user.Gender, user.Goal)
	return err
}

// GetDailyLog retrieves the daily log for a specific date
func (s *SQLiteDB) GetDailyLog(date time.Time) *models.DailyLog {
	dateStr := date.Format("2006-01-02")
	var mealsJSON string

	err := s.db.QueryRow(`
		SELECT meals FROM daily_logs WHERE date = ?
	`, dateStr).Scan(&mealsJSON)

	if err != nil {
		return &models.DailyLog{
			Date:  date,
			Meals: make([]*models.Meal, 0),
		}
	}

	// Parse JSON string into meals array
	var meals []*models.Meal
	if err := json.Unmarshal([]byte(mealsJSON), &meals); err != nil {
		return &models.DailyLog{
			Date:  date,
			Meals: make([]*models.Meal, 0),
		}
	}

	return &models.DailyLog{
		Date:  date,
		Meals: meals,
	}
}

// SaveDailyLog saves a daily log to the database
func (s *SQLiteDB) SaveDailyLog(log *models.DailyLog) error {
	dateStr := log.Date.Format("2006-01-02")

	// Convert meals to JSON
	mealsJSON, err := json.Marshal(log.Meals)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`
		INSERT OR REPLACE INTO daily_logs (date, meals)
		VALUES (?, ?)
	`, dateStr, string(mealsJSON))
	return err
}

// SaveUser saves a user to the database (alias for CreateUser)
func (s *SQLiteDB) SaveUser(user *models.User) error {
	return s.CreateUser(user)
}
