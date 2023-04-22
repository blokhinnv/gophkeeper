// Package server provides functionality for running a web server
// with endpoints for storing and retrieving different types of data.
// It utilizes the Gin web framework and the MongoDB database.
package server

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"gophkeeper/internal/server/config"
	"gophkeeper/internal/server/controller"
	"gophkeeper/internal/server/middleware"
	"gophkeeper/internal/server/service"
)

// RunServer starts the server and listens for incoming requests.
func RunServer(cfg *config.ServerConfig) {
	// Connect to MongoDB.
	log.Printf("Starting server with config %+v\n", cfg)
	ctx := context.Background()
	client, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI(cfg.MongoURI),
	)
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// Create service and controller instances.
	var (
		textService service.StorageService = service.NewStorageService(
			client.Database("gophkeeper"),
		)
		utilsService service.UtilsService = service.NewUtilsService(
			client,
		)
		authService service.AuthService = service.NewAuthService(
			client.Database("gophkeeper").Collection("users"),
			cfg.SigningKey,
			cfg.ExpireDuration,
		)

		storageController controller.StorageController = controller.NewStorageController(
			textService,
		)
		utilsController controller.UtilsController = controller.NewUtilsController(utilsService)
		authController  controller.AuthController  = controller.NewAuthController(authService)
	)

	// Set up routes and middleware.
	r := gin.Default()

	public := r.Group("/api")
	public.GET("/ping", utilsController.Ping)
	public.PUT("/user/register", authController.Register)
	public.PUT("/user/login", authController.Login)

	protected := r.Group("/api/store")
	protected.Use(middleware.JWTAuthMiddleware([]byte(cfg.SigningKey)))
	protected.PUT("/text", storageController.Store)
	protected.PUT("/credentials", storageController.Store)
	protected.PUT("/binary", storageController.Store)
	protected.PUT("/cards", storageController.Store)
	protected.POST("/text", storageController.Update)
	protected.POST("/credentials", storageController.Update)
	protected.POST("/binary", storageController.Update)
	protected.POST("/cards", storageController.Update)
	protected.GET("/text", storageController.GetAll)
	protected.GET("/credentials", storageController.GetAll)
	protected.GET("/binary", storageController.GetAll)
	protected.GET("/cards", storageController.GetAll)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
