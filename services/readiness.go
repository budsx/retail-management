package services

import "context"

func (svc *Service) Readiness(ctx context.Context) error {
	return svc.repo.Postgres.Close()
}