package services

import (
	"context"
	"fmt"

	"github.com/budsx/retail-management/model"
)

func (svc *Service) RegisterUser(ctx context.Context, user model.User) error {
	svc.logger.Info(fmt.Sprintf("[REQUEST] Registering user: %+v", user.Username))

	err := svc.repo.Postgres.RegisterUser(ctx, user)
	if err != nil {
		svc.logger.Error(fmt.Sprintf("[ERROR] Failed to register user: %s", err.Error()))
		return err
	}

	svc.logger.Info("[RESPONSE] User registered successfully")
	return nil
}
