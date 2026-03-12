package main

import (
	"Distributed_job_queue_system/internal/api/routers"
	"Distributed_job_queue_system/pkg/database"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	database.ConnectPostgres()

	router := gin.Default()

	routers.RegisterRoutes(router)

	log.Println("API Sever running on port 8080")

	err := router.Run(":8080")
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
