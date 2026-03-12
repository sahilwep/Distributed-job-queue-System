package routers

import (
	"Distributed_job_queue_system/internal/api/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	}) // GET /health

	router.POST("/jobs", handlers.CreateJob) // POST /jobs

	router.GET("/jobs/:id", handlers.GetJob) // GET /jobs/{id}
}
