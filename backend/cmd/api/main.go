package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/motvii/desertacia/internal/config"
	"github.com/motvii/desertacia/internal/handler"
	"github.com/motvii/desertacia/internal/repository"
	"github.com/motvii/desertacia/internal/service"
	"github.com/motvii/desertacia/pkg/jwt"
)

func main() {
	cfg := config.Load()

	db, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(context.Background()); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	jwtManager := jwt.NewManager(cfg.JWTSecret, cfg.JWTExpiresIn)

	userRepo := repository.NewUserRepository(db)
	deviceRepo := repository.NewDeviceRepository(db)
	emgRepo := repository.NewEMGRepository(db)
	trainingRepo := repository.NewTrainingRepository(db)

	mlClient := service.NewMLClient(cfg.MLServiceURL)

	authService := service.NewAuthService(userRepo, jwtManager, cfg.RefreshExpiry)
	deviceService := service.NewDeviceService(deviceRepo)
	emgService := service.NewEMGService(emgRepo, deviceRepo)
	trainingService := service.NewTrainingService(trainingRepo, mlClient)

	authHandler := handler.NewAuthHandler(authService)
	deviceHandler := handler.NewDeviceHandler(deviceService)
	emgHandler := handler.NewEMGHandler(emgService)
	trainingHandler := handler.NewTrainingHandler(trainingService)

	router := handler.SetupRouter(jwtManager, authHandler, deviceHandler, emgHandler, trainingHandler)

	srv := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("server starting on :%s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("forced shutdown: %v", err)
	}
}
