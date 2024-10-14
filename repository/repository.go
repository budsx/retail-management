package repository

import (
	"github.com/budsx/retail-management/repository/postgres"
)

type Repository struct {
	Postgres postgres.PostgresRepository
}

type DBConfig struct {
	Host, User, Password, DBName string
}

type RepoConfig struct {
	DBConfig
}

func NewRetailManagementRepository(conf RepoConfig) (*Repository, error) {
	postgres, err := postgres.NewPostgres(conf.Host, conf.User, conf.Password, conf.DBName)
	if err != nil {
		return nil, err
	}
	return &Repository{
		Postgres: postgres,
	}, nil
}

func (r *Repository) Close() {
	r.Postgres.Close()
}
