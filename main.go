package main

import (
	"log"
	"net/http"
	"time"

	"github.com/budsx/retail-management/config"
	"github.com/budsx/retail-management/controller"
	"github.com/budsx/retail-management/middleware"
	"github.com/budsx/retail-management/repository"
	"github.com/budsx/retail-management/services"
	"github.com/budsx/retail-management/utils"
	"github.com/gorilla/mux"
)

func main() {
	conf, err := config.NewConfig()
	if err != nil {
		log.Printf("Config error: %s", err)
	}

	logger := utils.NewLogger(conf.Log.Level)
	repoConf := repository.RepoConfig{
		DBConfig: repository.DBConfig{
			Host:     conf.DBHost,
			User:     conf.DBUser,
			Password: conf.DBPass,
			DBName:   conf.DBName,
		},
	}
	repo, err := repository.NewRetailManagementRepository(repoConf)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer repo.Close()

	service := services.NewRetailManagementService(*repo, logger)
	controller := controller.NewRetailManagementController(service)

	r := mux.NewRouter()

	// Health Check & Readiness
	r.HandleFunc("/health", controller.Health)
	r.HandleFunc("/readiness", func(w http.ResponseWriter, r *http.Request) {
	})

	// User
	r.HandleFunc("/user/register", controller.RegisterUser).Methods("POST")
	r.HandleFunc("/user/login", controller.Login).Methods("POST")
	r.Handle("/user/validate", middleware.TokenValidationMiddleware(http.HandlerFunc(controller.ValidateToken))).Methods("GET")

	// Product
	r.HandleFunc("/product/{id}", controller.GetProductByID).Methods("GET")
	r.HandleFunc("/product", controller.AddProduct).Methods("POST")
	r.HandleFunc("/product/{id}", controller.EditProduct).Methods("PUT")
	r.HandleFunc("/products", controller.GetProducts).Methods("GET")

	// TODO: Warehouse
	r.HandleFunc("/warehouse", controller.AddWarehouseByUserID).Methods("POST")
	r.HandleFunc("/warehouse/{id}", controller.EditWarehouseByUserID).Methods("PUT")
	r.HandleFunc("/warehouses", controller.GetWarehousesByUserID).Methods("GET")

	// TODO: Location
	r.HandleFunc("/location", controller.AddLocation).Methods("POST")
	r.HandleFunc("/location/{id}", controller.EditLocationByUserID).Methods("PUT")
	r.HandleFunc("/location/{id}", controller.DeleteLocationByUserID).Methods("DELETE")

	// TODO: Stock
	// Stock Inventory Management
	// User can create stock adjustment, view inventory transaction, and read total stock from
	// all location and single location. You must prevent that total stock can't be zero.

	// AdStockTransaction
	// UpdateStockTransaction
	// GetStockTransaction
	// GetTotalStockTransactionByLocation
	// GetTotalStockTransactions

	// Run Server
	srv := &http.Server{
		Handler:      r,
		Addr:         conf.Port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Info("HTTP server error: %s", err)
		}
	}()
	logger.Info("Server started ...")

	// Graceful Shutdown
	utils.OnShutdown(srv)
}
