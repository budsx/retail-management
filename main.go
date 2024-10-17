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
	r.HandleFunc("/readiness", controller.Readiness).Methods("GET")

	// User
	r.HandleFunc("/user/register", controller.RegisterUser).Methods("POST")
	r.HandleFunc("/user/login", controller.Login).Methods("POST")
	r.Handle("/user/validate", middleware.TokenValidationMiddleware(http.HandlerFunc(controller.ValidateToken))).Methods("GET")

	// Product
	r.Handle("/product/{id}", middleware.TokenValidationMiddleware(http.HandlerFunc(controller.GetProductByID))).Methods("GET")
	r.Handle("/product", middleware.TokenValidationMiddleware(http.HandlerFunc(controller.AddProduct))).Methods("POST")
	r.Handle("/product/{id}", middleware.TokenValidationMiddleware(http.HandlerFunc(controller.EditProduct))).Methods("PUT")
	r.Handle("/products", middleware.TokenValidationMiddleware(http.HandlerFunc(controller.GetProducts))).Methods("GET")

	// TODO: Warehouse
	r.Handle("/warehouse", middleware.TokenValidationMiddleware(http.HandlerFunc(controller.AddWarehouseByUserID))).Methods("POST")
	r.Handle("/warehouse/{id}", middleware.TokenValidationMiddleware(http.HandlerFunc(controller.EditWarehouseByUserID))).Methods("PUT")
	r.Handle("/warehouses", middleware.TokenValidationMiddleware(http.HandlerFunc(controller.GetWarehousesByUserID))).Methods("GET")

	// TODO: Location
	r.Handle("/location", middleware.TokenValidationMiddleware(http.HandlerFunc(controller.AddLocation))).Methods("POST")
	r.Handle("/location/{id}", middleware.TokenValidationMiddleware(http.HandlerFunc(controller.EditLocationByUserID))).Methods("PUT")
	r.Handle("/location/{id}", middleware.TokenValidationMiddleware(http.HandlerFunc(controller.DeleteLocationByUserID))).Methods("DELETE")

	// Adjust Stock
	r.Handle("/stock-transactions", middleware.TokenValidationMiddleware(http.HandlerFunc(controller.CreateStockTransaction))).Methods("POST")
	// View Transaction
	r.Handle("/stock-transactions", middleware.TokenValidationMiddleware(http.HandlerFunc(controller.GetStockTransactions))).Methods("GET")
	// View Transaction by ID
	r.Handle("/stock-transactions/{id}", middleware.TokenValidationMiddleware(http.HandlerFunc(controller.GetStockTransactionByID))).Methods("POST")
	// Total Stock All Location
	r.Handle("/total-stocks", middleware.TokenValidationMiddleware(http.HandlerFunc(controller.GetTotalStocks))).Methods("GET")
	// Total Stock By Location
	r.Handle("/total-stock/{location_id}", middleware.TokenValidationMiddleware(http.HandlerFunc(controller.GetTotalStockByLocation))).Methods("GET")

	log.Println(conf.Port)
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
