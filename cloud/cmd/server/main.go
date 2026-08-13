package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/OpenIndustrial/cloud/internal/api"
	"github.com/OpenIndustrial/cloud/internal/config"
	"github.com/OpenIndustrial/cloud/internal/device"
	"github.com/OpenIndustrial/cloud/internal/event"
	"github.com/OpenIndustrial/cloud/internal/identity"
	"github.com/OpenIndustrial/cloud/internal/infrastructure/redis"
	"github.com/OpenIndustrial/cloud/internal/kernel/resource"
	"github.com/OpenIndustrial/cloud/internal/notification"
	"github.com/OpenIndustrial/cloud/internal/persistence/postgres"
	"github.com/gin-gonic/gin"
	redisV8 "github.com/go-redis/redis/v8"
	_ "github.com/lib/pq" // Still needed by the postgres driver
)

func main() {
	// --- Configuration ---
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// --- Database Connection ---
	gormDB, err := postgres.NewDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// --- Redis Client ---
	redisClient := redisV8.NewClient(&redisV8.Options{
		Addr: cfg.RedisAddr,
	})
	if _, err := redisClient.Ping(redisClient.Context()).Result(); err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}

	// --- Event Bus & Publisher ---
	eventBus := event.NewInMemoryBus()
	var eventPublisher event.Publisher = redis.NewStreamPublisher(redisClient)

	// --- Repository Instantiation ---
	userRepo := postgres.NewUserRepository(gormDB)
	roleRepo := postgres.NewRoleRepository(gormDB)
	groupRepo := postgres.NewGroupRepository(gormDB)
	tenantRepo := postgres.NewTenantRepository(gormDB)
	permRepo := postgres.NewPermissionRepository(gormDB)
	resourceRepo := postgres.NewResourceRepository(gormDB)
	attrDefRepo := postgres.NewAttributeDefinitionRepository(gormDB)
	resAttrRepo := postgres.NewResourceAttributeRepository(gormDB)
    resConRepo := postgres.NewResourceConnectionsRepository(gormDB)
	// --- Service Instantiation ---
	identityService := identity.NewService(tenantRepo, userRepo, roleRepo, groupRepo, cfg.JWTSecret, eventPublisher)
	resourceService := resource.NewService(resourceRepo, attrDefRepo, resAttrRepo, resConRepo)
	deviceService := device.NewService(resourceService)

	// --- Register Event Handlers ---
	// We are building a single-process service, so we register all handlers here.
	registerNotificationHandlers(eventBus)
	// As we add more domains (e.g., analytics), we'll add more registration calls here.
	// registerAnalyticsHandlers(eventBus)

	// --- Start Background Event Subscriber ---
	// We create a cancellable context for graceful shutdown.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// The subscriber will listen for events and dispatch them to our eventBus.
	// We run it in a separate goroutine so it doesn't block the main thread.
	subscriber := redis.NewStreamSubscriber(redisClient, eventBus, "openindustrial:events", "notification-group", "notification-worker-1")
	go subscriber.Start(ctx)

	// --- HTTP Server Setup ---
	router := gin.Default()
	authMiddleware := api.NewAuthMiddleware(cfg.JWTSecret)
	apiV1 := router.Group("/api/v1")
	identityHandler := api.NewIdentityHandler(identityService, permRepo, authMiddleware)
	identityHandler.RegisterRoutes(apiV1)
	resourceHandler := api.NewResourceHandler(resourceService, permRepo, authMiddleware)
	resourceHandler.RegisterRoutes(apiV1)
	deviceHandler := api.NewDeviceHandler(deviceService, authMiddleware)
	deviceHandler.RegisterRoutes(apiV1)

	// --- Graceful Shutdown Handling ---
	go func() {
		serverAddr := fmt.Sprintf(":%s", cfg.Port)
		log.Printf("Starting server on %s", serverAddr)
		if err := router.Run(serverAddr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the application.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down application...")

	// Cancel the context to stop the background subscriber.
	cancel()

	// You could add a timeout here to wait for the subscriber to finish.
	log.Println("Application gracefully shut down.")
}

// registerNotificationHandlers initializes and registers all handlers for the notification domain.
func registerNotificationHandlers(bus event.Bus) {
	userCreatedHandler := notification.NewUserCreatedHandler()
	bus.Subscribe(event.IdentityUserCreated, userCreatedHandler.Handle)
}