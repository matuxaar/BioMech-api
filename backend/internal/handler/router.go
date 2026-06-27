package handler

import (
	firebase "firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"
	"github.com/matuxaar/BioMech-api/internal/middleware"
)

func SetupRouter(
	firebaseApp *firebase.App,
	authHandler *AuthHandler,
	userHandler *UserHandler,
	deviceHandler *DeviceHandler,
	emgHandler *EMGHandler,
	trainingHandler *TrainingHandler,
	statsHandler *StatsHandler,
	wsHandler *WSHandler,
	trainingFileHandler *TrainingFileHandler,
) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CORS())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.Static("/uploads", "./uploads")

	protected := r.Group("/api/v1")
	protected.Use(middleware.AuthRequired(firebaseApp))
	{
		protected.POST("/auth/firebase", authHandler.SyncUser)
		protected.GET("/me", userHandler.Me)
		protected.PUT("/me", userHandler.Update)
		protected.POST("/me/avatar", userHandler.UploadAvatar)

		protected.GET("/users/:id", userHandler.GetByID)

		devices := protected.Group("/devices")
		{
			devices.POST("", deviceHandler.Create)
			devices.GET("", deviceHandler.List)
			devices.GET("/:id", deviceHandler.GetByID)
			devices.PUT("/:id", deviceHandler.Update)
			devices.DELETE("/:id", deviceHandler.Delete)
			devices.GET("/:id/actions", deviceHandler.GetActions)
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
			training.POST("/upload", trainingHandler.Upload)
			training.POST("/files", trainingFileHandler.Upload)
			training.GET("/files", trainingFileHandler.List)
			training.GET("/files/:id", trainingFileHandler.Get)
			training.DELETE("/files/:id", trainingFileHandler.Delete)
		}

		protected.POST("/predict", trainingHandler.Predict)
		protected.GET("/predict/ws", wsHandler.PredictStream)
		protected.GET("/stats/dashboard", statsHandler.Dashboard)
	}

	return r
}
