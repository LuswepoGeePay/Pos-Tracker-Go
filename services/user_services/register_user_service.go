package userservices

import (
	"fmt"
	"pos-master/config"
	"pos-master/models"
	pb "pos-master/proto/auth"
	eventservices "pos-master/services/event_services"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(req *pb.RegisterRequest) *pb.AuthResponse {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		return &pb.AuthResponse{
			Success: false,
			Message: "Failed to hash password",
		}
	}

	var role models.Role
	result := config.DB.Where("name = ?", req.Role).First(&role)
	if result.Error != nil {
		return &pb.AuthResponse{
			Success: false,
			Message: "Role not found",
		}
	}

	userID := uuid.New()
	user := models.User{
		ID:       userID,
		FullName: req.Fullname,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     role,
	}

	result = config.DB.Create(&user)
	if result.Error != nil {

		errorString := fmt.Sprintf("Unable to register user %v", result.Error)
		return &pb.AuthResponse{
			Success: false,
			Message: errorString,
			Status:  "failed",
		}
	}
	eventservices.RegisterEvent("User registered successfully", map[string]interface{}{
		"User ID":  userID,
		"Fullname": req.Fullname,
		"Email":    req.Email,
		"Role":     req.Role,
	})

	return &pb.AuthResponse{
		Success: true,
		Message: "Registration successful",
		Status:  "success",
	}
}
