# NutritionApp

A command-line nutrition tracking application written in Go that helps users monitor their daily nutritional intake and physical progress.

## Features

- Personal profile management (weight, height, age, goals)
- Body composition tracking (BMI, body fat estimation)
- Food and meal tracking
- Daily nutritional totals
- Command-line interface for easy interaction

## Project Structure

```
nutritionapp/
├── cmd/
│   └── nutritionapp/      # Main application
├── pkg/
│   ├── models/            # Data models
│   ├── database/          # Data storage
│   ├── utils/            # Utility functions
│   └── cli/              # CLI interface
└── README.md
```

## Getting Started

1. Clone the repository
2. Run `go mod tidy` to install dependencies
3. Build and run the application:
   ```bash
   go run cmd/nutritionapp/main.go
   ```

## Available Commands

- `help` - Show available commands
- `exit` - Exit the application

More commands will be added for:
- User profile management
- Food search and tracking
- Meal management
- Nutritional reports
