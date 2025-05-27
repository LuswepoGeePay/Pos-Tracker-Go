package userservices

import (
	"pos-master/config"
	"pos-master/models"
	pb "pos-master/proto/auth"

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

	user := models.User{
		ID:       uuid.New(),
		FullName: req.Fullname,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     role,
	}

	result = config.DB.Create(&user)
	if result.Error != nil {
		return &pb.AuthResponse{
			Success: false,
			Message: "Registration failed",
			Status:  "failed",
		}
	}

	return &pb.AuthResponse{
		Success: true,
		Message: "Registration successful",
		Status:  "success",
	}
}
