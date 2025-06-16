package userservices

import (
	"pos-master/config"
	"pos-master/models"
	pb "pos-master/proto/auth"
	eventservices "pos-master/services/event_services"
	"pos-master/utils"
)

func EditUser(req *pb.EditUserRequest) error {
	updates := map[string]interface{}{}

	if req.Fullname != "" {
		updates["full_name"] = req.Fullname
	}

	if req.Email != "" {
		updates["email"] = req.Email
	}

	var currentUser models.User

	result := config.DB.Where("id = ?", req.Id).First(&currentUser)
	if result.Error != nil {
		return utils.CapitalizeError("unable to find user with that ID")
	}

	if req.Status != currentUser.Status {
		updates["status"] = req.Status
	}

	var role models.Role
	config.DB.Where("name = ?", req.Role).First(&role)

	if req.Role != "" {
		updates["role_id"] = role.ID
	}

	if len(updates) == 0 {
		return utils.CapitalizeError("no changes detected")
	}

	err := config.DB.Model(&models.User{}).
		Where("id = ?", req.Id).
		Updates(updates).Error

	if err != nil {
		return utils.CapitalizeError("failed to update user")
	}
	eventservices.RegisterEvent("User edited successfully", map[string]interface{}{
		"user id":   req.Id,
		"Full name": req.Fullname,
		"Email":     req.Email,
		"Role":      req.Role,
		"Status":    req.Status,
	})

	return nil
}
