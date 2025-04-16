package server

import (
	"nutritionapp/pkg/models"
	"time"
)

func (s *Server) handleGetReport(untypedData any) Response {
	dailyLog := s.userDB.GetDailyLog(time.Now())
	var totals models.NutritionalTotals

	for _, meal := range dailyLog.Meals {
		mealTotals := meal.CalculateTotals()
		totals.Calories += mealTotals.Calories
		totals.Proteins += mealTotals.Proteins
		totals.Carbs += mealTotals.Carbs
		totals.Fats += mealTotals.Fats
		totals.Fiber += mealTotals.Fiber
	}

	return Response{
		Data: ReportResponse{
			Calories: totals.Calories,
			Proteins: totals.Proteins,
			Carbs:    totals.Carbs,
			Fats:     totals.Fats,
			Fiber:    totals.Fiber,
		},
	}
}
