package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dablon/runflow/internal/config"
	"github.com/dablon/runflow/internal/handlers"
	"github.com/dablon/runflow/internal/parser"
	"github.com/gin-gonic/gin"
)

func main() {
	_ = config.Load()
	
	p := parser.New()
	h := handlers.New(p)
	
	r := gin.Default()
	
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Next()
	})
	
	api := r.Group("/api/v1")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})
		
		api.POST("/runbooks", h.CreateRunbook)
		api.GET("/runbooks", h.ListRunbooks)
		api.GET("/runbooks/:id", h.GetRunbook)
		api.PUT("/runbooks/:id", h.UpdateRunbook)
		api.DELETE("/runbooks/:id", h.DeleteRunbook)
		
		api.POST("/runbooks/:id/execute", h.ExecuteRunbook)
		api.GET("/executions/:id", h.GetExecution)
		api.GET("/executions/:id/logs", h.GetExecutionLogs)
		
		api.POST("/providers", h.CreateProvider)
		api.GET("/providers", h.ListProviders)
		api.DELETE("/providers/:id", h.DeleteProvider)
	}
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	log.Printf("Starting RunFlow API on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

var _ = fmt.Println
