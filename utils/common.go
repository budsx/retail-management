package utils

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// Graceful Shutdown
func OnShutdown(srv *http.Server) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	log.Println("Shutting down server...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Server forced to shutdown: %v", err)
	} else {
		log.Println("Server gracefully stopped")
	}
}
