package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/berkkaradalan/stackflow/config"
	"github.com/berkkaradalan/stackflow/database"
	"github.com/berkkaradalan/stackflow/handler"
	repository "github.com/berkkaradalan/stackflow/repository/postgres"
	"github.com/berkkaradalan/stackflow/routes"
	"github.com/berkkaradalan/stackflow/service"
	"github.com/berkkaradalan/stackflow/utils"
)

func main() {

	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config: ", err)
	}

	pool, err := database.Connect(ctx, cfg.Env)
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	defer pool.Close()

	if err := database.Migrate(ctx, pool, cfg); err != nil {
		log.Fatal("Failed to run migrations: ", err)
	}

	jwtManager, err := utils.NewJWTManager(
		cfg.Env.JWTSecret,
		cfg.Env.JWTAccessExpiry,
		cfg.Env.JWTRefreshExpiry,
	)
	if err != nil {
		log.Fatal("Failed to initialize JWT manager: ", err)
	}

	userRepo := repository.NewUserRepository(pool)
	inviteTokenRepo := repository.NewInviteTokenRepository(pool)

	authService := service.NewAuthService(userRepo, jwtManager)
	userService := service.NewUserService(userRepo, inviteTokenRepo)

	authHandler := handler.NewAuthHandler(authService, userService)
	userHandler := handler.NewUserHandler(userService)

	router := routes.SetupRouter(jwtManager, authHandler, userHandler)


	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Env.HostName, cfg.Env.HostPort),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("Starting server on :%s", cfg.Env.HostPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Forced shutdown: %v", err)
	}

	log.Println("Server exited")
}