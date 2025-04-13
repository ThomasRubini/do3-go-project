package database

import (
	"encoding/json"
	"fmt"
	"nutritionapp/pkg/models"
	"os"
	"path/filepath"
	"time"
)

type UserDatabase struct {
	User      *models.User
	DailyLogs map[string]*models.DailyLog // key: date in YYYY-MM-DD format
	path      string
}

func NewUserDatabase(dbPath string) (*UserDatabase, error) {
	db := &UserDatabase{
		path:      dbPath,
		DailyLogs: make(map[string]*models.DailyLog),
	}

	// Create database directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %v", err)
	}

	// Load existing database or create new one
	if err := db.load(); err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to load database: %v", err)
		}
		// If file doesn't exist, start with empty database
		if err := db.save(); err != nil {
			return nil, fmt.Errorf("failed to save initial database: %v", err)
		}
	}

	return db, nil
}

func (db *UserDatabase) load() error {
	data, err := os.ReadFile(db.path)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, db)
}

func (db *UserDatabase) save() error {
	data, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal database: %v", err)
	}

	return os.WriteFile(db.path, data, 0644)
}

func (db *UserDatabase) SaveUser(user *models.User) error {
	db.User = user
	return db.save()
}

func (db *UserDatabase) GetUser() *models.User {
	return db.User
}

func (db *UserDatabase) SaveDailyLog(log *models.DailyLog) error {
	dateKey := log.Date.Format("2006-01-02")
	db.DailyLogs[dateKey] = log
	return db.save()
}

func (db *UserDatabase) GetDailyLog(date time.Time) *models.DailyLog {
	dateKey := date.Format("2006-01-02")
	if log, exists := db.DailyLogs[dateKey]; exists {
		return log
	}
	return &models.DailyLog{
		Date:  date,
		Meals: make([]models.Meal, 0),
	}
}
