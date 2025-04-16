package server

import (
	"fmt"
	"nutritionapp/pkg/models"
)

func (s *Server) handleCreateProfile(req Request) Response {
	data, ok := req.Data.(CreateProfileData)
	if !ok {
		return Response{Error: fmt.Errorf("invalid request data")}
	}

	user := &models.User{
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Age:       data.Age,
		Weight:    data.Weight,
		Height:    data.Height,
		Gender:    data.Gender,
		Goal:      data.Goal,
	}

	if err := s.userDB.SaveUser(user); err != nil {
		return Response{Error: fmt.Errorf("failed to save user: %v", err)}
	}

	return Response{}
}

func (s *Server) handleGetProfile(req Request) Response {
	user := s.userDB.GetUser()
	if user == nil {
		return Response{Error: fmt.Errorf("no profile exists")}
	}

	return Response{
		Data: ProfileResponseData{
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Age:         user.Age,
			Weight:      user.Weight,
			Height:      user.Height,
			Gender:      user.Gender,
			Goal:        user.Goal,
			BMI:         user.CalculateBMI(),
			BodyFatPerc: user.EstimateBodyFat(),
		},
	}
}
func (s *Server) handleUpdateProfile(req Request) Response {
	data, ok := req.Data.(UpdateProfileData)
	if !ok {
		return Response{Error: fmt.Errorf("invalid request data")}
	}

	user := s.userDB.GetUser()
	if user == nil {
		return Response{Error: fmt.Errorf("no profile exists")}
	}

	user.FirstName = data.FirstName
	user.LastName = data.LastName
	user.Age = data.Age
	user.Weight = data.Weight
	user.Height = data.Height
	user.Gender = data.Gender
	user.Goal = data.Goal

	if err := s.userDB.SaveUser(user); err != nil {
		return Response{Error: fmt.Errorf("failed to update user: %v", err)}
	}

	return Response{
		Data: ProfileResponseData{
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Age:         user.Age,
			Weight:      user.Weight,
			Height:      user.Height,
			Gender:      user.Gender,
			Goal:        user.Goal,
			BMI:         user.CalculateBMI(),
			BodyFatPerc: user.EstimateBodyFat(),
		},
	}
}
