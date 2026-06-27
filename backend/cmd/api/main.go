package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/matuxaar/BioMech-api/internal/config"
	"github.com/matuxaar/BioMech-api/internal/handler"
	"github.com/matuxaar/BioMech-api/internal/repository"
	"github.com/matuxaar/BioMech-api/internal/service"
	"google.golang.org/api/option"
)

var firebaseAuth *auth.Client

func runMigrations(db *pgxpool.Pool) {
	migrations := []string{
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS nickname VARCHAR(50) UNIQUE`,
		`CREATE INDEX IF NOT EXISTS idx_users_nickname ON users(nickname)`,
		`ALTER TABLE devices ADD COLUMN IF NOT EXISTS ble_service_uuid VARCHAR(64) DEFAULT ''`,
		`ALTER TABLE devices ADD COLUMN IF NOT EXISTS ble_command_char_uuid VARCHAR(64) DEFAULT ''`,
		`ALTER TABLE devices ADD COLUMN IF NOT EXISTS ble_status_char_uuid VARCHAR(64) DEFAULT ''`,
		`ALTER TABLE devices ADD COLUMN IF NOT EXISTS ble_emg_char_uuid VARCHAR(64) DEFAULT ''`,
		`ALTER TABLE device_actions ADD COLUMN IF NOT EXISTS action_code INT NOT NULL DEFAULT 0`,
	}
	
	for _, m := range migrations {
		if _, err := db.Exec(context.Background(), m); err != nil {
			log.Printf("migration warning: %v", err)
		}
	}
	log.Println("migrations applied")
}

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

	runMigrations(db)

	opt := option.WithCredentialsFile(cfg.FirebaseCredsFile)
	firebaseApp, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("failed to initialize Firebase app: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	deviceRepo := repository.NewDeviceRepository(db)
	emgRepo := repository.NewEMGRepository(db)
	trainingRepo := repository.NewTrainingRepository(db)
	trainingFileRepo := repository.NewTrainingFileRepository(db)

	mlClient := service.NewMLClient(cfg.MLServiceURL)

	authService := service.NewAuthService(userRepo)
	deviceService := service.NewDeviceService(deviceRepo)
	emgService := service.NewEMGService(emgRepo, deviceRepo)
	trainingService := service.NewTrainingService(trainingRepo, emgRepo, mlClient)
	trainingFileService := service.NewTrainingFileService(trainingFileRepo)
	statsRepo := repository.NewStatsRepository(db)
	statsService := service.NewStatsService(statsRepo)

	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(authService)
	deviceHandler := handler.NewDeviceHandler(deviceService)
	emgHandler := handler.NewEMGHandler(emgService)
	trainingHandler := handler.NewTrainingHandler(trainingService)
	trainingFileHandler := handler.NewTrainingFileHandler(trainingFileService)
	statsHandler := handler.NewStatsHandler(statsService)

	wsHandler := handler.NewWSHandler(mlClient)

	router := handler.SetupRouter(firebaseApp, authHandler, userHandler, deviceHandler, emgHandler, trainingHandler, statsHandler, wsHandler, trainingFileHandler)

	srv := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 0,
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
