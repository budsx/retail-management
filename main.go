package main

import (
	"log"
	"net/http"
	"time"

	"github.com/budsx/retail-management/config"
	"github.com/budsx/retail-management/controller"
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

	// Health Check
	r.HandleFunc("/health", controller.Health)

	// Product
	r.HandleFunc("/product/{id}", controller.GetProductByID).Methods("GET")
	r.HandleFunc("/products", controller.GetProducts).Methods("GET")
	r.HandleFunc("/product", controller.AddProduct).Methods("POST")           // Menambahkan produk baru
	r.HandleFunc("/product/{id}", controller.EditProduct).Methods("PUT")  

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
