package userservices

import (
	"pos-master/config"
	"pos-master/models"
	pb "pos-master/proto/auth"
	"pos-master/utils"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func GetUsers(req *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {

	var users []models.User

	tx := config.DB.Begin()
	query := tx.Preload("Role").Model(&models.User{})

	var totalUsers int64
	err := query.Count(&totalUsers).Error

	if err != nil {
		return nil, utils.CapitalizeError("failed to count users")
	}

	totalPages := int32((totalUsers + int64(req.PageSize) - 1) / int64(req.PageSize))
	// Calculate offset for pagination
	offset := (req.Page - 1) * req.PageSize

	// Execute the final query with pagination and preloading
	err = query.Limit(int(req.PageSize)).
		Offset(int(offset)).
		Find(&users).Error

	if err != nil {
		return nil, utils.CapitalizeError("failed to retrieve users")
	}

	pbUsers := make([]*pb.User, len(users))

	for i, user := range users {

		pbUsers[i] = &pb.User{
			Id:       user.ID.String(),
			Fullname: user.FullName,
			Email:    user.Email,
			Role:     user.Role.Name,
			Status:   user.Status,
		}
	}

	return &pb.GetUsersResponse{
		User:        pbUsers,
		TotalPages:  totalPages,
		CurrentPage: req.Page,
		HasMore:     req.Page < totalPages,
	}, nil
}

func GetUser(userID string) (*pb.User, error) {

	var user models.User

	userid, err := uuid.Parse(userID)

	if err != nil {
		return nil, utils.CapitalizeError(err.Error())
	}

	err = config.DB.Preload("Role").Where("id = ?", userid).Find(&user).Error
	if err != nil {
		return nil, utils.CapitalizeError(err.Error())
	}

	return &pb.User{
		Email:    user.Email,
		Id:       user.ID.String(),
		Fullname: user.FullName,
		Role:     user.Role.Name,
		Status:   user.Status,
	}, nil
}

func ChangeEmailorPassword(req *pb.ChangeEmailOrPasswordRequest) error {

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return utils.CapitalizeError("invalid uuid. Check ID being sent")
	}

	var user models.User

	err = config.DB.Where("id = ?", userID).Find(&user).Error

	if err != nil {
		return utils.CapitalizeError("unable to find user")
	}

	tx := config.DB.Begin()

	if req.IsEmailRequest {
		if user.Email == req.OldEmail {
			result := tx.Model(&models.User{}).Update("email", req.NewEmail)
			if result.Error != nil {
				return utils.CapitalizeError("unable to update email")
			}

		} else {
			return utils.CapitalizeError("Old email does not match our records")
		}

	}

	if req.IsPasswordRequest {

		if req.ConfirmPassword != req.NewPassword {
			return utils.CapitalizeError("your passwords do not match")
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)

		if err != nil {
			return utils.CapitalizeError("unable to hash password")
		}

		result := tx.Model(&models.User{}).Update("password", hashedPassword)
		if result.Error != nil {
			return utils.CapitalizeError("unable to update password")
		}

	}

	return nil
}
