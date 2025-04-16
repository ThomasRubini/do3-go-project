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
	data := req.Data

	switch req.Type {
	case ReqCreateProfile:
		resp = s.handleCreateProfile(data)
	case ReqGetProfile:
		resp = s.handleGetProfile(data)
	case ReqUpdateProfile:
		resp = s.handleUpdateProfile(data)
	case ReqAddMeal:
		resp = s.handleAddMeal(data)
	case ReqListMeals:
		resp = s.handleListMeals(data)
	case ReqSearchFood:
		resp = s.handleSearchFood(data)
	case ReqAddFood:
		resp = s.handleAddFood(data)
	case ReqGetReport:
		resp = s.handleGetReport(data)
	default:
		resp = Response{Error: fmt.Errorf("unknown request type: %s", req.Type)}
	}

	req.Return <- resp
}
