package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/budsx/retail-management/model"
	"github.com/budsx/retail-management/utils"
)

func (svc *Service) RegisterUser(ctx context.Context, user model.User) error {
	svc.logger.Info(fmt.Sprintf("[REQUEST] Registering user: %+v", user.Username))

	if user.Username == "" || user.Password == "" {
		svc.logger.Error("Bad Request")
		return errors.New("bad request")
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		svc.logger.Error(fmt.Sprintf("[ERROR] Failed to hash password: %s", err.Error()))
		return err
	}
	user.Password = hashedPassword

	err = svc.repo.Postgres.RegisterUser(ctx, user)
	if err != nil {
		svc.logger.Error(fmt.Sprintf("[ERROR] Failed to register user: %s", err.Error()))
		return err
	}

	svc.logger.Info("[RESPONSE] User registered successfully")
	return nil
}
