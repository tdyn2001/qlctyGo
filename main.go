package main

import (
	"context"
	"log"
	"net/http"
	"v2/controllers"
	"v2/initializers"
	"v2/kafkas"
	"v2/middleware"
	"v2/services"
	"v2/wss"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	server         *gin.Engine
	authController controllers.AuthController
	authService    services.AuthService
	wsServer       *wss.WsServer
)

func init() {
	initializers.ConnectDB(initializers.GetConfig())

	authService = services.NewAuthService(initializers.DB)
	authController = controllers.NewAuthController(authService)

	server = gin.Default()
	wsServer = wss.NewWebsocketServer()
}

func main() {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:8000", initializers.GetConfig().ClientOrigin}
	corsConfig.AllowCredentials = true

	initKafka()

	server.Use(cors.New(corsConfig))

	server.GET("/ws", middleware.AuthorizeJWT(), wsServer.SetupWSS)

	router := server.Group("/api")
	router.GET("/healthchecker", func(ctx *gin.Context) {
		message := "Welcome to Golang with Gorm and Postgres"
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
	})

	authController.AuthController(router)
	log.Fatal(server.Run(":" + initializers.GetConfig().ServerPort))
}

func initKafka() {
	go kafkas.Produce(context.Background(), "test-topic-1")
	go kafkas.Consume(context.Background(), "test-topic-1", "test-group")
}
