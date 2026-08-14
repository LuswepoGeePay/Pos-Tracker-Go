package userservices

import (
	"fmt"
	"log"
	"strings"

	database "pos-master/config"
	"pos-master/models"
	pb "pos-master/proto/auth"
	emailservices "pos-master/services/emailservices"
	eventservices "pos-master/services/event_services"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(req *pb.RegisterRequest) *pb.AuthResponse {
	fullname := strings.TrimSpace(req.Fullname)
	email := strings.TrimSpace(req.Email)
	roleName := strings.TrimSpace(req.Role)

	if fullname == "" {
		return &pb.AuthResponse{
			Success: false,
			Message: "fullname is required",
			Status:  "failure",
		}
	}
	if email == "" {
		return &pb.AuthResponse{
			Success: false,
			Message: "email is required",
			Status:  "failure",
		}
	}
	if strings.TrimSpace(req.Password) == "" {
		return &pb.AuthResponse{
			Success: false,
			Message: "password is required",
			Status:  "failure",
		}
	}
	if roleName == "" {
		return &pb.AuthResponse{
			Success: false,
			Message: "role is required",
			Status:  "failure",
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return &pb.AuthResponse{
			Success: false,
			Message: "Failed to hash password",
			Status:  "failure",
		}
	}

	var role models.Role
	result := database.DB.Where("name = ?", roleName).First(&role)
	if result.Error != nil {
		return &pb.AuthResponse{
			Success: false,
			Message: "Role not found",
			Status:  "failure",
		}
	}

	userID := uuid.New()
	user := models.User{
		ID:       userID,
		FullName: fullname,
		Email:    email,
		Password: string(hashedPassword),
		RoleID:   role.ID,
	}

	tx := database.DB.Begin()
	if tx.Error != nil {
		return &pb.AuthResponse{
			Success: false,
			Message: fmt.Sprintf("Unable to start transaction: %v", tx.Error),
			Status:  "failure",
		}
	}

	result = tx.Create(&user)
	if result.Error != nil {
		tx.Rollback()
		return &pb.AuthResponse{
			Success: false,
			Message: fmt.Sprintf("Unable to register user: %v", result.Error),
			Status:  "failure",
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return &pb.AuthResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to commit transaction: %v", err),
			Status:  "failure",
		}
	}

	go func() {
		if err := emailservices.SendWelcomeEmail(email, fullname); err != nil {
			log.Printf("Failed to send welcome email to %s: %v", email, err)
		}
	}()

	eventservices.RegisterEvent("User registered successfully", map[string]interface{}{
		"User ID":  userID,
		"Fullname": fullname,
		"Email":    email,
		"Role":     roleName,
	})

	return &pb.AuthResponse{
		Success: true,
		Message: "Registration successful",
		Status:  "success",
	}
}
