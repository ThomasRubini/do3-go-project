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
func NewServer(userDB db.UserDatabase, foodProcessor *fdc.FoodProcessor, requests chan Request) *Server {
	return &Server{
		userDB:        userDB,
		foodProcessor: foodProcessor,
		requests:      requests,
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
	case ReqCreateProfile:
		resp = s.handleCreateProfile(req)
	case ReqGetProfile:
		resp = s.handleGetProfile(req)
	case ReqUpdateProfile:
		resp = s.handleUpdateProfile(req)
	case ReqAddMeal:
		resp = s.handleAddMeal(req)
	case ReqListMeals:
		resp = s.handleListMeals(req)
	case ReqSearchFood:
		resp = s.handleSearchFood(req)
	case ReqAddFood:
		resp = s.handleAddFood(req)
	case ReqGetReport:
		resp = s.handleGetReport(req)
	default:
		resp = Response{Error: fmt.Errorf("unknown request type: %s", req.Type)}
	}

	req.Return <- resp
}
