package userservices

import (
	"fmt"
	"log"
	database "pos-master/config"
	"pos-master/models"
	pb "pos-master/proto/auth"
	emailservices "pos-master/services/emailservices"
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
	result := database.DB.Where("name = ?", req.Role).First(&role)
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

	// Start transaction for user registration
	tx := database.DB.Begin()
	if tx.Error != nil {
		errorString := fmt.Sprintf("Unable to start transaction: %v", tx.Error)
		return &pb.AuthResponse{
			Success: false,
			Message: errorString,
			Status:  "failed",
		}
	}

	result = tx.Create(&user)
	if result.Error != nil {
		tx.Rollback()
		errorString := fmt.Sprintf("Unable to register user %v", result.Error)
		return &pb.AuthResponse{
			Success: false,
			Message: errorString,
			Status:  "failed",
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return &pb.AuthResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to commit transaction: %v", err),
			Status:  "failed",
		}
	}

	// Send welcome email
	go func() {
		if err := emailservices.SendWelcomeEmail(req.Email, req.Fullname); err != nil {
			log.Printf("Failed to send welcome email to %s: %v", req.Email, err)
		}
	}()

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
