package main

import (
	"adire-apparel/config"
	"adire-apparel/database"
	"adire-apparel/lib"
	"adire-apparel/middlewares"
	"adire-apparel/routes"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
	config.InitializeConfig()

	err := database.InitializeDatabase()
	defer database.CloseDatabase()
	if err != nil {
		log.Fatal("Database error:", err)
	}

	lib.InitialiseJWT(string(config.AppConfig.JWTSecret))

	corsConfig := cors.Config{
		AllowOrigins:     []string{config.AppConfig.ClientUrl},
		AllowMethods:     config.AppConfig.AllowMethods,
		AllowHeaders:     config.AppConfig.AllowHeaders,
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}

	if config.AppConfig.IsDevMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	app := gin.Default()
	app.Use(gin.Logger())
	app.Use(cors.New(corsConfig))
	app.Use(middlewares.ErrorHandlerMiddleware())
	app.Use(middlewares.AuthMiddleware())
	app.Use(lib.ErrorHandler())

	app.MaxMultipartMemory = 10 << 20 // 10

	hub := lib.NewHub()
	go hub.Run()

	websocket := lib.NewWebSocketHandler(hub)

	prefix := config.AppConfig.Version
	router := app.Group(prefix)

	router.GET("/ws", websocket.HandleWebSocket)

	router.GET("/health", func(ctx *gin.Context) {
		lib.Success(ctx, "Adire Apparel API is healthy", map[string]interface{}{
			"version": config.AppConfig.Version,
		})
	})

	routes.AuthRoutes(router)
	routes.CategoryRoutes(router)
	routes.OrderRoutes(router)
	routes.ProductRoutes(router)
	routes.TransactionRoutes(router)
	routes.UserRoutes(router)

	app.NoRoute(lib.GlobalNotFound())

	server := &http.Server{
		Addr:           fmt.Sprintf(":%s", config.AppConfig.Port),
		Handler:        app,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	log.Printf("Server starting on port http://localhost:%s/%s", config.AppConfig.Port, config.AppConfig.Version)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Server failed to start:", err)
	}
}
