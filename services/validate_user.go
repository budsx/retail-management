package services

import (
	"context"
	"fmt"

	"github.com/budsx/retail-management/model"
	"github.com/budsx/retail-management/utils"
)

func (svc *Service) ValidateUser(ctx context.Context, req model.Credentials) (model.User, error) {
	svc.logger.Info(fmt.Sprintf("[REQUEST] User validated: %s", req.Username))
	user, err := svc.repo.Postgres.GetUserByUsername(ctx, req.Username)
	if err != nil {
		svc.logger.Error(fmt.Sprintf("[ERROR] Failed to find user: %s", err.Error()))
		return model.User{}, fmt.Errorf("invalid username or password")
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		svc.logger.Error("[ERROR] Invalid password")
		return model.User{}, fmt.Errorf("invalid username or password")
	}

	svc.logger.Info(fmt.Sprintf("[RESPONSE] User validated: %s", user.Username))
	return user, nil
}
