package main

import (
	"fmt"
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
	r.HandleFunc("/health", controller.Health).Methods("GET")

	// User
	r.HandleFunc("/user/register", controller.RegisterUser).Methods("POST")
	r.HandleFunc("/user/login", controller.Login).Methods("POST")

	// Private Route
	private := r.PathPrefix("/v1").Subrouter()
	private.Use(middleware.TokenValidationMiddleware)
	private.HandleFunc("/user/validate", controller.ValidateToken).Methods("GET")

	// Product
	private.HandleFunc("/product/{id}", controller.GetProductByID).Methods("GET")
	private.HandleFunc("/product", controller.AddProduct).Methods("POST")
	private.HandleFunc("/product/{id}", controller.EditProduct).Methods("PUT")
	private.HandleFunc("/products", controller.GetProducts).Methods("GET")

	// Warehouse
	private.HandleFunc("/warehouse", controller.AddWarehouseByUserID).Methods("POST")
	private.HandleFunc("/warehouse/{id}", controller.EditWarehouseByUserID).Methods("PUT")
	private.HandleFunc("/warehouses", controller.GetWarehousesByUserID).Methods("GET")

	// Location
	private.HandleFunc("/location", controller.AddLocation).Methods("POST")
	private.HandleFunc("/location/{id}", controller.EditLocationByUserID).Methods("PUT")
	private.HandleFunc("/location/{id}", controller.DeleteLocationByUserID).Methods("DELETE")

	// Stock
	private.HandleFunc("/stock-transactions", controller.CreateStockTransaction).Methods("POST")
	private.HandleFunc("/stock-transactions", controller.GetStockTransactions).Methods("GET")
	private.HandleFunc("/stock-transactions/{id}", controller.GetStockTransactionByID).Methods("GET")
	private.HandleFunc("/total-stocks", controller.GetTotalStocks).Methods("GET")
	private.HandleFunc("/total-stock/{location_id}", controller.GetTotalStockByLocation).Methods("GET")

	// Run Server
	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf(":%s", conf.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Info("HTTP server error: %s", err)
		}
	}()

	logger.Info(fmt.Sprintf("Server started on port %s", conf.Port))

	// Graceful Shutdown
	utils.OnShutdown(srv)
}
