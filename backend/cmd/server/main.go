package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"inventory-management/backend/internal/config"
	"inventory-management/backend/internal/database"
	"inventory-management/backend/internal/handlers"
	"inventory-management/backend/internal/models"
)

func main() {
	cfg := config.Load()

	dbConfig := database.NewConfig(
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
	)

	db, err := database.Connect(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&models.Product{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	productHandler := handlers.NewProductHandler(db)

	api := router.Group("/api")
	{
		api.GET("/products", productHandler.List)
		api.GET("/products/:id", productHandler.Get)
		api.POST("/products", productHandler.Create)
		api.PUT("/products/:id", productHandler.Update)
		api.DELETE("/products/:id", productHandler.Delete)
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	staticDir := cfg.StaticDir
	if _, err := os.Stat(staticDir); err == nil {
		assetsDir := filepath.Join(staticDir, "assets")
		if _, err := os.Stat(assetsDir); err == nil {
			router.Static("/assets", assetsDir)
		}

		router.StaticFile("/favicon.svg", filepath.Join(staticDir, "favicon.svg"))
		router.StaticFile("/icons.svg", filepath.Join(staticDir, "icons.svg"))

		router.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path
			if strings.HasPrefix(path, "/api/") {
				c.JSON(http.StatusNotFound, gin.H{"error": "api endpoint not found"})
				return
			}

			filePath := filepath.Join(staticDir, path)
			if _, err := os.Stat(filePath); err == nil {
				c.File(filePath)
				return
			}

			c.File(filepath.Join(staticDir, "index.html"))
		})

		log.Printf("Serving static files from: %s", staticDir)
	} else {
		log.Printf("Static directory not found: %s (skipping static file serving)", staticDir)
	}

	log.Printf("Server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
