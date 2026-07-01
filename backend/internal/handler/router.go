package handler

import (
	"net/http"
	"net/http/pprof"

	firebase "firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"
	"github.com/matuxaar/BioMech-api/internal/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	maxUploadSizeMB int64,
	uploadsDir string,
) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery(), middleware.Logger(), middleware.RequestID(), middleware.CORS(), middleware.Metrics(), middleware.RateLimit(100, 200))
	r.MaxMultipartMemory = maxUploadSizeMB << 20

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	r.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	pprofGroup := r.Group("/debug/pprof")
	pprofGroup.Use(middleware.InternalAuth())
	{
		pprofGroup.GET("/", gin.WrapH(pprof.Handler("")))
		pprofGroup.GET("/cmdline", gin.WrapH(http.HandlerFunc(pprof.Cmdline)))
		pprofGroup.GET("/profile", gin.WrapH(http.HandlerFunc(pprof.Profile)))
		pprofGroup.GET("/symbol", gin.WrapH(http.HandlerFunc(pprof.Symbol)))
		pprofGroup.GET("/trace", gin.WrapH(http.HandlerFunc(pprof.Trace)))
		pprofGroup.GET("/allocs", gin.WrapH(pprof.Handler("allocs")))
		pprofGroup.GET("/block", gin.WrapH(pprof.Handler("block")))
		pprofGroup.GET("/goroutine", gin.WrapH(pprof.Handler("goroutine")))
		pprofGroup.GET("/heap", gin.WrapH(pprof.Handler("heap")))
		pprofGroup.GET("/mutex", gin.WrapH(pprof.Handler("mutex")))
		pprofGroup.GET("/threadcreate", gin.WrapH(pprof.Handler("threadcreate")))
	}

	r.POST("/api/v1/training/jobs/:id/status", middleware.InternalAuth(), trainingHandler.UpdateJobStatus)

	uploads := r.Group("/uploads")
	uploads.Use(middleware.AuthRequired(firebaseApp))
	uploads.Static("", uploadsDir)

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
			emg.POST("/sessions/:id/samples", middleware.RateLimit(30, 50), emgHandler.AddSample)
			emg.POST("/sessions/:id/samples/batch", middleware.RateLimit(10, 20), emgHandler.AddSamplesBatch)
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
			training.GET("/files/:id/download", trainingFileHandler.Download)
			training.DELETE("/files/:id", trainingFileHandler.Delete)
		}

		protected.POST("/predict", middleware.RateLimit(20, 30), trainingHandler.Predict)
		protected.GET("/predict/ws", wsHandler.PredictStream)
		protected.GET("/stats/dashboard", statsHandler.Dashboard)
	}

	return r
}
