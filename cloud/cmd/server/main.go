package main

import (
	"log"
	"os"

	"github.com/OpenIndustrial/cloud/internal/api"
	"github.com/OpenIndustrial/cloud/internal/identity"
	"github.com/OpenIndustrial/cloud/internal/resource"
	"github.com/OpenIndustrial/cloud/internal/persistence/postgres"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	// 0. Load configuration directly from environment variables
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	// 1. Initialize database connection
	db, err := sqlx.Connect("postgres", databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Dependency Injection: Wire everything together
	// 1. Create repository implementations
	tenantRepo := postgres.NewTenantRepository(db)
	userRepo := postgres.NewUserRepository(db)
	roleRepo := postgres.NewRoleRepository(db)
	permRepo := postgres.NewPermissionRepository(db)

	// Repositories for the new Resource Kernel
	resourceRepo := postgres.NewPgResourceRepository(db)
	groupRepo := postgres.NewPgGroupRepository(db)
	authzRepo := groupRepo // PgGroupRepository implements both interfaces

	// 2. Create the services
	identityService := identity.NewService(tenantRepo, userRepo, roleRepo)
	resourceService := resource.NewService(resourceRepo, groupRepo)

	// 3. Create the shared middleware and HTTP handlers
	authMiddleware := api.NewAuthMiddleware(jwtSecret)
	identityHandler := api.NewIdentityHandler(identityService, permRepo, authMiddleware)
	resourceHandler := api.NewResourceHandler(resourceService, authzRepo, permRepo, authMiddleware)

	// 4. Set up the router and register routes
	router := gin.Default()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	apiV1 := router.Group("/api/v1")

	api.RegisterHealthRoutes(router)

	identityHandler.RegisterRoutes(apiV1)
	resourceHandler.RegisterRoutes(apiV1)

	// 5. Start the server
	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}