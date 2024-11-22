package services

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/budsx/retail-management/middleware"
	"github.com/budsx/retail-management/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestService_GetTotalStockByLocation(t *testing.T) {
	srv := NewTestServer(t)
	defer srv.MockCtrl.Finish()

	// Mock user context
	testUser := model.User{
		UserID:   1,
		Username: "testuser",
	}
	ctx := middleware.SetUserInfoToContext(context.Background(), testUser)

	testStocks := []model.ProductStock{
		{
			ProductID:  1,
			TotalStock: 100,
		},
		{
			ProductID:  2,
			TotalStock: 150,
		},
	}

	tests := []struct {
		name       string
		locationID int64
		mock       func()
		want       []model.ProductStock
		wantErr    bool
	}{
		{
			name:       "success",
			locationID: 1,
			mock: func() {
				srv.MockRepo.EXPECT().
					GetTotalStockByLocation(
						gomock.Any(),
						gomock.Any(),
					).
					Return(testStocks, nil)
			},
			want:    testStocks,
			wantErr: false,
		},
		{
			name:       "empty stock",
			locationID: 2,
			mock: func() {
				srv.MockRepo.EXPECT().
					GetTotalStockByLocation(
						gomock.Any(),
						gomock.Eq(int64(2)),
					).
					Return([]model.ProductStock{}, nil)
			},
			want:    []model.ProductStock{},
			wantErr: false,
		},
		{
			name:       "location not found",
			locationID: 999,
			mock: func() {
				srv.MockRepo.EXPECT().
					GetTotalStockByLocation(
						gomock.Any(),
						gomock.Eq(int64(999)),
					).
					Return(nil, errors.New("location not found"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:       "database error",
			locationID: 1,
			mock: func() {
				srv.MockRepo.EXPECT().
					GetTotalStockByLocation(
						gomock.Any(),
						gomock.Eq(int64(1)),
					).
					Return(nil, errors.New("database error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			tt.mock()

			// Execute test
			got, err := srv.Service.GetTotalStockByLocation(ctx, tt.locationID)

			// Assert results
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestService_GetTotalStocks(t *testing.T) {
	srv := NewTestServer(t)
	defer srv.MockCtrl.Finish()

	// Mock user context
	testUser := model.User{
		UserID:   1,
		Username: "testuser",
	}
	ctx := middleware.SetUserInfoToContext(context.Background(), testUser)

	testStocks := []model.ProductStock{
		{
			ProductID:  1,
			TotalStock: 200,
		},
		{
			ProductID:  2,
			TotalStock: 300,
		},
		{
			ProductID:  3,
			TotalStock: 150,
		},
	}

	tests := []struct {
		name    string
		mock    func()
		want    []model.ProductStock
		wantErr bool
	}{
		{
			name: "success",
			mock: func() {
				srv.MockRepo.EXPECT().
					GetTotalStocks(gomock.Any()).
					Return(testStocks, nil)
			},
			want:    testStocks,
			wantErr: false,
		},
		{
			name: "empty stock",
			mock: func() {
				srv.MockRepo.EXPECT().
					GetTotalStocks(gomock.Any()).
					Return([]model.ProductStock{}, nil)
			},
			want:    []model.ProductStock{},
			wantErr: false,
		},
		{
			name: "database error",
			mock: func() {
				srv.MockRepo.EXPECT().
					GetTotalStocks(gomock.Any()).
					Return(nil, errors.New("database error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := srv.Service.GetTotalStocks(ctx)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// Additional test for context without user info
func TestService_GetTotalStocks_NoUserContext(t *testing.T) {
	srv := NewTestServer(t)
	defer srv.MockCtrl.Finish()

	// Using context without user info
	ctx := context.Background()

	_, err := srv.Service.GetTotalStocks(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user info not found in context")
}

// Test for context timeout
func TestService_GetTotalStocks_ContextTimeout(t *testing.T) {
	srv := NewTestServer(t)
	defer srv.MockCtrl.Finish()

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Microsecond)
	defer cancel()

	// Add user info to context
	ctx = middleware.SetUserInfoToContext(ctx, model.User{UserID: 1, Username: "testuser"})

	// Wait for context to timeout
	time.Sleep(1 * time.Millisecond)

	srv.MockRepo.EXPECT().
		GetTotalStocks(gomock.Any()).
		Return(nil, context.DeadlineExceeded)

	_, err := srv.Service.GetTotalStocks(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context deadline exceeded")
}

// Test concurrent access
func TestService_GetTotalStocks_ConcurrentAccess(t *testing.T) {
	srv := NewTestServer(t)
	defer srv.MockCtrl.Finish()

	// Mock user context
	testUser := model.User{
		UserID:   1,
		Username: "testuser",
	}
	ctx := middleware.SetUserInfoToContext(context.Background(), testUser)

	// Setup mock for multiple calls
	srv.MockRepo.EXPECT().
		GetTotalStocks(gomock.Any()).
		Return([]model.ProductStock{}, nil).
		AnyTimes()

	// Run concurrent requests
	concurrentRequests := 10
	var wg sync.WaitGroup
	for i := 0; i < concurrentRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := srv.Service.GetTotalStocks(ctx)
			assert.NoError(t, err)
		}()
	}
	wg.Wait()
}

// Helper function to create test stocks
func createTestStocks(count int) []model.ProductStock {
	stocks := make([]model.ProductStock, count)
	for i := 0; i < count; i++ {
		stocks[i] = model.ProductStock{
			ProductID:  int64(i + 1),
			TotalStock: int64((i + 1) * 100),
		}
	}
	return stocks
}

func TestService_GetTotalStocks_Performance(t *testing.T) {
	srv := NewTestServer(t)
	defer srv.MockCtrl.Finish()

	ctx := middleware.SetUserInfoToContext(context.Background(), model.User{UserID: 1})

	// Test with large dataset
	largeStocks := createTestStocks(1000)

	srv.MockRepo.EXPECT().
		GetTotalStocks(gomock.Any()).
		Return(largeStocks, nil)

	start := time.Now()
	got, err := srv.Service.GetTotalStocks(ctx)
	duration := time.Since(start)

	assert.NoError(t, err)
	assert.Equal(t, largeStocks, got)
	assert.Less(t, duration, 100*time.Millisecond)
}
