package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/psilin/sluggin/api"
	"github.com/spf13/viper"
)

func main() {
	// Process config
	path := flag.String("c", "./config", "Path to config file")
	flag.Parse()

	viper.SetConfigName("config.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(*path)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("fatal error config file: %v", err)
	}

	// Initialize environment object
	var env api.Env
	err = env.Init(viper.GetString("dsn"))
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Create router
	router := api.InitRouter(&env)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	// Cleanup environment object
	err = env.Teardown()
	if err != nil {
		log.Fatalf("%v", err)
	}

	log.Println("Server exiting")
}
