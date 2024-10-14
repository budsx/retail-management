package main

import (
	"net/http"

	"github.com/budsx/retail-management/controller"
	"github.com/budsx/retail-management/repository"
	"github.com/budsx/retail-management/services"
	"github.com/gorilla/mux"
)

func main() {
	// conf, err := config.NewConfig()
	// if err != nil {
	// 	log.Fatalf("Config error: %s", err)
	// }

	// logger := utils.NewLogger(conf.Log.Level)

	repoConf := repository.RepoConfig{
		DBConfig: repository.DBConfig{
			// Host: ,
			// User: ,
			// Password: ,
			// DBName: ,
		},
	}

	repo, err := repository.NewRetailManagementRepository(repoConf)
	if err != nil {
		return
	}
	defer repo.Close()

	service := services.NewRetailManagementService(*repo)
	controller := controller.NewRetailManagementController(service)

	r := mux.NewRouter()
	// Health Check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		
	})
	r.HandleFunc("/products/{id}", controller.GetProductByID).Methods("GET")
	r.HandleFunc("/stock", controller.GetStock).Methods("GET")

	// Run Server
	// Graceful Shutdown
}
