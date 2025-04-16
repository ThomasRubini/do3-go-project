package server

import (
	"fmt"
	"nutritionapp/pkg/db"
	"nutritionapp/pkg/fdc"
)

type Server struct {
	userDB        db.UserDatabase
	foodProcessor *fdc.FoodProcessor
	requests      chan Request
}

// NewServer creates a new server instance
func NewServer(userDB db.UserDatabase, foodProcessor *fdc.FoodProcessor) *Server {
	return &Server{
		userDB:        userDB,
		foodProcessor: foodProcessor,
		requests:      make(chan Request),
	}
}

// Start starts the server
func (s *Server) Start() {
	for req := range s.requests {
		go s.handleRequest(req)
	}
}

// SendRequest sends a request to the server
func (s *Server) SendRequest(reqType string, data interface{}) Response {
	resp := make(chan Response)
	s.requests <- Request{
		Type:   reqType,
		Data:   data,
		Return: resp,
	}
	return <-resp
}

func (s *Server) handleRequest(req Request) {
	var resp Response

	switch req.Type {
	case "create_profile":
		resp = s.handleCreateProfile(req)
	case "get_profile":
		resp = s.handleGetProfile(req)
	case "update_profile":
		resp = s.handleUpdateProfile(req)
	case "add_meal":
		resp = s.handleAddMeal(req)
	case "list_meals":
		resp = s.handleListMeals(req)
	case "search_food":
		resp = s.handleSearchFood(req)
	case "add_food":
		resp = s.handleAddFood(req)
	case "get_report":
		resp = s.handleGetReport(req)
	default:
		resp = Response{Error: fmt.Errorf("unknown request type: %s", req.Type)}
	}

	req.Return <- resp
}
