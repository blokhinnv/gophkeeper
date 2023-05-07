// Package server provides functionality for running a web server
// with endpoints for storing and retrieving different types of data.
// It utilizes the Gin web framework and the MongoDB database.
package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/blokhinnv/gophkeeper/internal/server/config"
	"github.com/blokhinnv/gophkeeper/internal/server/controller"
	_ "github.com/blokhinnv/gophkeeper/internal/server/docs"
	"github.com/blokhinnv/gophkeeper/internal/server/middleware"
	"github.com/blokhinnv/gophkeeper/internal/server/service"
)

// RunServer starts the server and listens for incoming requests.
//
//	@title Gophkeeper server
//	@version 1.0
//	@description Gophkeeper server which allows user to store the sensitive data.
//	@BasePath /
//	@schemes http
//	@securityDefinitions.apikey bearerAuth
//	@in header
//	@name Authorization
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
		storageService service.StorageService = service.NewStorageService(
			client.Database(cfg.DBName), cfg.EncryptionKey,
		)
		utilsService service.UtilsService = service.NewUtilsService(
			client,
		)
		authService service.AuthService = service.NewAuthService(
			client.Database(cfg.DBName).Collection("users"),
			cfg.SigningKey,
			cfg.ExpireDuration,
		)
		syncService service.SyncService = service.NewSyncService()

		storageController controller.StorageController = controller.NewStorageController(
			storageService, syncService,
		)
		utilsController controller.UtilsController = controller.NewUtilsController(utilsService)
		authController  controller.AuthController  = controller.NewAuthController(authService)
		syncController  controller.SyncController  = controller.NewSyncController(syncService)
	)

	// Set up routes and middleware.
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	public := r.Group("/api")
	public.GET("/ping", utilsController.Ping)
	public.PUT("/user/register", authController.Register)
	public.PUT("/user/login", authController.Login)

	protected := r.Group("/api/store")
	protected.Use(middleware.JWTAuthMiddleware([]byte(cfg.SigningKey)))
	protected.PUT("/:collectionName", storageController.Store)
	protected.POST("/:collectionName", storageController.Update)
	protected.GET("/:collectionName", storageController.GetAll)
	protected.DELETE("/:collectionName", storageController.Delete)

	sync := r.Group("/api/sync")
	sync.Use(middleware.JWTAuthMiddleware([]byte(cfg.SigningKey)))
	sync.POST("/register", syncController.Register)
	sync.POST("/unregister", syncController.Unregister)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// r.Run(fmt.Sprintf("127.0.0.1%v", cfg.Port))
	srv := &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%v", cfg.Port),
		Handler: r,
	}
	go func() {
		var err error
		if cfg.UseHTTPS {
			err = srv.ListenAndServeTLS(cfg.CertFile, cfg.KeyFile)
		} else {
			err = srv.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %s\n", err)
		}
	}()
	quit := make(chan os.Signal, 2)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Bye!")

}
