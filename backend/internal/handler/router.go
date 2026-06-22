package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/motvii/desertacia/internal/middleware"
	"github.com/motvii/desertacia/pkg/jwt"
)

func SetupRouter(
	jwtManager *jwt.Manager,
	authHandler *AuthHandler,
	deviceHandler *DeviceHandler,
	emgHandler *EMGHandler,
	trainingHandler *TrainingHandler,
) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CORS())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
		}

		protected := api.Group("")
		protected.Use(middleware.AuthRequired(jwtManager))
		{
			protected.GET("/me", authHandler.Refresh)

			devices := protected.Group("/devices")
			{
				devices.POST("", deviceHandler.Create)
				devices.GET("", deviceHandler.List)
				devices.GET("/:id", deviceHandler.GetByID)
			}

			emg := protected.Group("/emg")
			{
				emg.POST("/sessions", emgHandler.StartSession)
				emg.POST("/sessions/:id/end", emgHandler.EndSession)
				emg.GET("/sessions", emgHandler.ListSessions)
				emg.GET("/sessions/:id", emgHandler.GetSession)
				emg.POST("/sessions/:id/samples", emgHandler.AddSample)
				emg.POST("/sessions/:id/samples/batch", emgHandler.AddSamplesBatch)
				emg.GET("/sessions/:id/samples", emgHandler.GetSamples)
			}

			training := protected.Group("/training")
			{
				training.POST("/jobs", trainingHandler.CreateJob)
				training.GET("/jobs", trainingHandler.ListJobs)
				training.GET("/jobs/:id", trainingHandler.GetJob)
			}
		}
	}

	return r
}
