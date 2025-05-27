package userservices

import (
	"pos-master/config"
	"pos-master/models"
	pb "pos-master/proto/auth"
	"pos-master/utils"
)

func EditUser(req *pb.EditUserRequest) error {
	updates := map[string]interface{}{}

	if req.Fullname != "" {
		updates["firstname"] = req.Fullname
	}

	if req.Email != "" {
		updates["email"] = req.Email
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

	return nil
}
