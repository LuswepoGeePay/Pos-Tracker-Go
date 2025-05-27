package authservices

import (
	"fmt"
	"pos-master/config"
	"pos-master/models"
	pb "pos-master/proto/auth"
	"pos-master/utils"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func LoginUser(req *pb.LoginRequest) (*pb.AuthResponse, error) {
	var user models.User

	// First check if user exists
	result := config.DB.Preload("Role.Permissions").Where("email = ?", req.Email).First(&user)
	if result.Error != nil {
		return nil, utils.CapitalizeError("invalid credentials")
	}

	// Separate password check to avoid timing attacks
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, utils.CapitalizeError("invalid credentials")
	}

	token, tokenExpiry, err := GenerateJWT(user.ID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to create token: %v", err)
	}
	var permissions []string
	for _, permission := range user.Role.Permissions {
		permissions = append(permissions, permission.Name)
	}

	return &pb.AuthResponse{
		Success:     true,
		Status:      "success",
		Message:     "Login successful",
		Token:       token,
		Id:          user.ID.String(),
		Role:        user.Role.Name,
		Permissions: permissions,
		TokenExpiry: tokenExpiry.Format(time.RFC3339),
	}, nil
}
