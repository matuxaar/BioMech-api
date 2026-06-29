package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	firebase "firebase.google.com/go/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/matuxaar/BioMech-api/internal/config"
	"github.com/matuxaar/BioMech-api/internal/handler"
	"github.com/matuxaar/BioMech-api/internal/middleware"
	"github.com/matuxaar/BioMech-api/internal/migrations"
	"github.com/matuxaar/BioMech-api/internal/repository"
	"github.com/matuxaar/BioMech-api/internal/service"
	"google.golang.org/api/option"
)

func main() {
	migrations.SetupLogger()
	cfg := config.Load()

	middleware.InitCORS(cfg.CORSOrigins)
	middleware.InitAuth(cfg.DevMode)

	db, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(context.Background()); err != nil {
		slog.Error("failed to ping database", "error", err)
		os.Exit(1)
	}
	slog.Info("connected to database")

	if err := migrations.Run(db, cfg.MigrationsDir); err != nil {
		slog.Error("migration failed", "error", err)
		os.Exit(1)
	}
	slog.Info("all migrations applied")

	var firebaseApp *firebase.App
	if !cfg.DevMode {
		if _, statErr := os.Stat(cfg.FirebaseCredsFile); statErr == nil {
			opt := option.WithCredentialsFile(cfg.FirebaseCredsFile)
			firebaseApp, err = firebase.NewApp(context.Background(), nil, opt)
			if err != nil {
				slog.Error("failed to initialize Firebase app", "error", err)
				os.Exit(1)
			}
		} else {
			slog.Warn("Firebase credentials file not found, auth disabled", "path", cfg.FirebaseCredsFile)
		}
	}

	userRepo := repository.NewUserRepository(db)
	deviceRepo := repository.NewDeviceRepository(db)
	emgRepo := repository.NewEMGRepository(db)
	trainingRepo := repository.NewTrainingRepository(db)
	trainingFileRepo := repository.NewTrainingFileRepository(db)

	mlClient := service.NewMLClient(cfg.MLServiceURL, cfg.MLRequestTimeout)

	authService := service.NewAuthService(userRepo)
	deviceService := service.NewDeviceService(deviceRepo)
	emgService := service.NewEMGService(emgRepo, deviceRepo)
	trainingService := service.NewTrainingService(trainingRepo, emgRepo, deviceRepo, mlClient)
	trainingFileService := service.NewTrainingFileService(trainingFileRepo, cfg.TrainingDir)
	statsRepo := repository.NewStatsRepository(db)
	statsService := service.NewStatsService(statsRepo)

	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(authService, cfg.AvatarsDir)
	deviceHandler := handler.NewDeviceHandler(deviceService)
	emgHandler := handler.NewEMGHandler(emgService)
	trainingHandler := handler.NewTrainingHandler(trainingService)
	trainingFileHandler := handler.NewTrainingFileHandler(trainingFileService)
	statsHandler := handler.NewStatsHandler(statsService)
	wsHandler := handler.NewWSHandler(mlClient)

	router := handler.SetupRouter(firebaseApp, authHandler, userHandler, deviceHandler, emgHandler, trainingHandler, statsHandler, wsHandler, trainingFileHandler, cfg.MaxUploadSizeMB, cfg.UploadsDir)

	srv := &http.Server{
		Addr:              ":" + cfg.ServerPort,
		Handler:           router,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		IdleTimeout:       cfg.IdleTimeout,
	}

	go func() {
		slog.Info("server starting", "port", cfg.ServerPort, "dev_mode", cfg.DevMode)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server...")
	handler.CloseAllWS()
	middleware.StopRateLimiters()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("forced shutdown", "error", err)
	}
}
