package main

import (
	"log"
	"os"

	"github.com/OpenIndustrial/cloud/internal/api"
	"github.com/OpenIndustrial/cloud/internal/identity"
	"github.com/OpenIndustrial/cloud/internal/kernel/resource"
	"github.com/OpenIndustrial/cloud/internal/persistence/postgres"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	// --- Configuration ---
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	// --- Database Connection ---
	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// --- Repository Instantiation ---
	// UPDATED: Create all necessary repositories with their correct constructors.
	userRepo := postgres.NewUserRepository(db)
	roleRepo := postgres.NewRoleRepository(db)
	groupRepo := postgres.NewGroupRepository(db) // For Identity Kernel
	tenantRepo := postgres.NewTenantRepository(db)
	permRepo := postgres.NewPermissionRepository(db)

	// Repositories for the Resource Kernel
	resourceRepo := postgres.NewResourceRepository(db)
	attrDefRepo := postgres.NewAttributeDefinitionRepository(db)
	resAttrRepo := postgres.NewResourceAttributeRepository(db)

	// --- Service Instantiation ---
	// UPDATED: Inject correct repositories into each service.
	identityService := identity.NewService(tenantRepo, userRepo, roleRepo, groupRepo, jwtSecret)
	resourceService := resource.NewService(resourceRepo, attrDefRepo, resAttrRepo)

	// --- HTTP Server Setup ---
	router := gin.Default()
	authMiddleware := api.NewAuthMiddleware(jwtSecret)

	// --- Handler Instantiation & Route Registration ---
	apiV1 := router.Group("/api/v1")

	// UPDATED: Instantiate and register the new IdentityHandler.
	identityHandler := api.NewIdentityHandler(identityService, permRepo,authMiddleware)
	identityHandler.RegisterRoutes(apiV1)

	// UPDATED: Fix the ResourceHandler instantiation.
	resourceHandler := api.NewResourceHandler(resourceService, permRepo, authMiddleware)
	resourceHandler.RegisterRoutes(apiV1)

	// --- Start Server ---
	log.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}