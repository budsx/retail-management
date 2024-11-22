package services

import (
	"testing"

	"github.com/budsx/retail-management/utils"
	"github.com/budsx/retail-management/repository"
	mocks "github.com/budsx/retail-management/repository/postgres"
	"github.com/golang/mock/gomock"
)

type TestServer struct {
	MockCtrl     *gomock.Controller
	MockRepo     *mocks.MockPostgresRepository
	MockLogger   *utils.Logger
	Service      RetailManagementService
}

func NewTestServer(t *testing.T) *TestServer {
	mockCtrl := gomock.NewController(t)
	mockRepo := mocks.NewMockPostgresRepository(mockCtrl)
	mockLogger := utils.NewLogger("info")

	svc := NewRetailManagementService(repository.Repository{
		Postgres: mockRepo,
	}, mockLogger)

	return &TestServer{
		MockCtrl:   mockCtrl,
		MockRepo:   mockRepo,
		MockLogger: mockLogger,
		Service:    svc,
	}
}