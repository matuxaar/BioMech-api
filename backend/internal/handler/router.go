package handler

import (
	firebase "firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"
	"github.com/motvii/desertacia/internal/middleware"
)

func SetupRouter(
	firebaseApp *firebase.App,
	authHandler *AuthHandler,
	userHandler *UserHandler,
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
		api.POST("/auth/firebase", authHandler.SyncUser)

		protected := api.Group("")
		protected.Use(middleware.AuthRequired(firebaseApp))
		{
			protected.GET("/me", userHandler.Me)

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
