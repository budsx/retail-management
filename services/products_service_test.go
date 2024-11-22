package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/budsx/retail-management/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestService_AddProduct(t *testing.T) {
	srv := NewTestServer(t)
	defer srv.MockCtrl.Finish()

	fixedTime := time.Now()

	tests := []struct {
		name    string
		product model.Product
		mock    func()
		wantErr bool
	}{
		{
			name: "success",
			product: model.Product{
				ProductName: "Test Product",
				Description: "Test Description",
				Price:      100,
				SKU:        "TEST-SKU",
				CreatedAt:  fixedTime,
			},
			mock: func() {
				srv.MockRepo.EXPECT().
					WriteProduct(gomock.Any(), model.Product{
						ProductName: "Test Product",
						Description: "Test Description",
						Price:      100,
						SKU:        "TEST-SKU",
						CreatedAt:  fixedTime,
					}).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "empty product name",
			product: model.Product{
				ProductName: "",
				Description: "Test Description",
				Price:      100,
				SKU:        "TEST-SKU",
			},
			mock: func() {
				srv.MockRepo.EXPECT().
					WriteProduct(gomock.Any(), gomock.Any()).
					Return(errors.New("product name cannot be empty"))
			},
			wantErr: true,
		},
		{
			name: "duplicate SKU",
			product: model.Product{
				ProductName: "Test Product",
				Description: "Test Description",
				Price:      100,
				SKU:        "EXISTING-SKU",
			},
			mock: func() {
				srv.MockRepo.EXPECT().
					WriteProduct(gomock.Any(), gomock.Any()).
					Return(errors.New("duplicate SKU"))
			},
			wantErr: true,
		},
		{
			name: "database error",
			product: model.Product{
				ProductName: "Test Product",
				Description: "Test Description",
				Price:      100,
				SKU:        "TEST-SKU",
			},
			mock: func() {
				srv.MockRepo.EXPECT().
					WriteProduct(gomock.Any(), gomock.Any()).
					Return(errors.New("database connection error"))
			},
			wantErr: true,
		},
		{
			name: "invalid price",
			product: model.Product{
				ProductName: "Test Product",
				Description: "Test Description",
				Price:      -100,
				SKU:        "TEST-SKU",
			},
			mock: func() {
				srv.MockRepo.EXPECT().
					WriteProduct(gomock.Any(), gomock.Any()).
					Return(errors.New("price must be positive"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			tt.mock()

			// Execute test
			err := srv.Service.AddProduct(context.Background(), tt.product)

			// Assert results
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

// Optional: Add a matcher for more precise product comparison
type productMatcher struct {
	expected model.Product
}

func matchProduct(expected model.Product) gomock.Matcher {
	return &productMatcher{expected: expected}
}

func (m *productMatcher) Matches(x interface{}) bool {
	actual, ok := x.(model.Product)
	if !ok {
		return false
	}
	// Add your comparison logic here
	return actual.ProductName == m.expected.ProductName &&
		actual.Description == m.expected.Description &&
		actual.Price == m.expected.Price &&
		actual.SKU == m.expected.SKU
}

func (m *productMatcher) String() string {
	return "matches product"
}

// Example usage of the matcher:
func ExampleTestWithMatcher(t *testing.T) {
	srv := NewTestServer(t)
	defer srv.MockCtrl.Finish()

	expectedProduct := model.Product{
		ProductName: "Test Product",
		Description: "Test Description",
		Price:      100,
		SKU:        "TEST-SKU",
	}

	srv.MockRepo.EXPECT().
		WriteProduct(gomock.Any(), matchProduct(expectedProduct)).
		Return(nil)

	err := srv.Service.AddProduct(context.Background(), expectedProduct)
	assert.NoError(t, err)
}

func TestService_EditProduct(t *testing.T) {
	srv := NewTestServer(t)
	defer srv.MockCtrl.Finish()

	tests := []struct {
		name    string
		product model.Product
		mock    func()
		wantErr bool
	}{
		{
			name: "success",
			product: model.Product{
				ProductID:   1,
				ProductName: "Updated Product",
				Description: "Updated Description",
				Price:      200,
				SKU:        "UPD-SKU",
			},
			mock: func() {
				srv.MockRepo.EXPECT().
					UpdateProductByID(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "product not found",
			product: model.Product{
				ProductID:   999,
				ProductName: "Non-existent Product",
			},
			mock: func() {
				srv.MockRepo.EXPECT().
					UpdateProductByID(gomock.Any(), gomock.Any()).
					Return(errors.New("product not found"))
			},
			wantErr: true,
		},
		{
			name: "database error",
			product: model.Product{
				ProductID:   1,
				ProductName: "Error Product",
			},
			mock: func() {
				srv.MockRepo.EXPECT().
					UpdateProductByID(gomock.Any(), gomock.Any()).
					Return(errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := srv.Service.EditProduct(context.Background(), tt.product)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestService_GetProductByID(t *testing.T) {
	srv := NewTestServer(t)
	defer srv.MockCtrl.Finish()

	fixedTime := time.Now()
	testProduct := model.Product{
		ProductID:   1,
		ProductName: "Test Product",
		Description: "Test Description",
		Price:      100,
		SKU:        "TEST-SKU",
		CreatedAt:  fixedTime,
		UpdatedAt:  fixedTime,
	}

	tests := []struct {
		name    string
		id      int64
		mock    func()
		want    model.Product
		wantErr bool
	}{
		{
			name: "success",
			id:   1,
			mock: func() {
				srv.MockRepo.EXPECT().
					ReadProductByID(gomock.Any(), int64(1)).
					Return(testProduct, nil)
			},
			want:    testProduct,
			wantErr: false,
		},
		{
			name: "product not found",
			id:   999,
			mock: func() {
				srv.MockRepo.EXPECT().
					ReadProductByID(gomock.Any(), int64(999)).
					Return(model.Product{}, errors.New("product not found"))
			},
			want:    model.Product{},
			wantErr: true,
		},
		{
			name: "invalid id",
			id:   -1,
			mock: func() {
				srv.MockRepo.EXPECT().
					ReadProductByID(gomock.Any(), int64(-1)).
					Return(model.Product{}, errors.New("invalid product ID"))
			},
			want:    model.Product{},
			wantErr: true,
		},
		{
			name: "database error",
			id:   1,
			mock: func() {
				srv.MockRepo.EXPECT().
					ReadProductByID(gomock.Any(), int64(1)).
					Return(model.Product{}, errors.New("database error"))
			},
			want:    model.Product{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := srv.Service.GetProductByID(context.Background(), tt.id)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestService_GetProducts(t *testing.T) {
	srv := NewTestServer(t)
	defer srv.MockCtrl.Finish()

	fixedTime := time.Now()
	testProducts := []model.Product{
		{
			ProductID:   1,
			ProductName: "Product 1",
			Price:      100,
			CreatedAt:  fixedTime,
		},
		{
			ProductID:   2,
			ProductName: "Product 2",
			Price:      200,
			CreatedAt:  fixedTime,
		},
	}

	tests := []struct {
		name       string
		pagination model.Pagination
		mock       func()
		want       []model.Product
		wantErr    bool
	}{
		{
			name: "success - first page",
			pagination: model.Pagination{
				Page:  1,
				Limit: 10,
			},
			mock: func() {
				srv.MockRepo.EXPECT().
					ReadProductsWithPagination(gomock.Any(), int32(10), int32(0)).
					Return(testProducts, nil)
			},
			want:    testProducts,
			wantErr: false,
		},
		{
			name: "success - second page",
			pagination: model.Pagination{
				Page:  2,
				Limit: 10,
			},
			mock: func() {
				srv.MockRepo.EXPECT().
					ReadProductsWithPagination(gomock.Any(), int32(10), int32(10)).
					Return([]model.Product{}, nil)
			},
			want:    []model.Product{},
			wantErr: false,
		},
		{
			name: "invalid page",
			pagination: model.Pagination{
				Page:  0,
				Limit: 10,
			},
			mock: func() {
				srv.MockRepo.EXPECT().
					ReadProductsWithPagination(gomock.Any(), int32(10), int32(-10)).
					Return(nil, errors.New("invalid page number"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid limit",
			pagination: model.Pagination{
				Page:  1,
				Limit: 0,
			},
			mock: func() {
				srv.MockRepo.EXPECT().
					ReadProductsWithPagination(gomock.Any(), int32(0), int32(0)).
					Return(nil, errors.New("invalid limit"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "database error",
			pagination: model.Pagination{
				Page:  1,
				Limit: 10,
			},
			mock: func() {
				srv.MockRepo.EXPECT().
					ReadProductsWithPagination(gomock.Any(), int32(10), int32(0)).
					Return(nil, errors.New("database error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := srv.Service.GetProducts(context.Background(), tt.pagination)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// Helper function to check pagination calculation
func TestCalculateOffset(t *testing.T) {
	tests := []struct {
		name       string
		pagination model.Pagination
		want       int32
	}{
		{
			name: "first page",
			pagination: model.Pagination{
				Page:  1,
				Limit: 10,
			},
			want: 0,
		},
		{
			name: "second page",
			pagination: model.Pagination{
				Page:  2,
				Limit: 10,
			},
			want: 10,
		},
		{
			name: "third page with different limit",
			pagination: model.Pagination{
				Page:  3,
				Limit: 15,
			},
			want: 30,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := (tt.pagination.Page - 1) * tt.pagination.Limit
			assert.Equal(t, tt.want, got)
		})
	}
}

