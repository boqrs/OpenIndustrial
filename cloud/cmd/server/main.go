package main

import (
	"log"
	"os"

	"github.com/OpenIndustrial/cloud/internal/api"
	"github.com/OpenIndustrial/cloud/internal/identity"
	"github.com/OpenIndustrial/cloud/internal/kernel/resource"
	"github.com/OpenIndustrial/cloud/internal/persistence/postgres"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" // Still needed by the postgres driver
	//"gorm.io/gorm"
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

	// --- Database Connection (Using GORM) ---
	gormDB, err := postgres.NewDB(dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database using GORM: %v", err)
	}

	// --- Repository Instantiation (Using GORM DB) ---
	// Identity Kernel Repositories
	userRepo := postgres.NewUserRepository(gormDB)
	roleRepo := postgres.NewRoleRepository(gormDB)
	groupRepo := postgres.NewGroupRepository(gormDB)
	tenantRepo := postgres.NewTenantRepository(gormDB)
	permRepo := postgres.NewPermissionRepository(gormDB)

	// Resource Kernel Repositories
	resourceRepo := postgres.NewResourceRepository(gormDB)
	attrDefRepo := postgres.NewAttributeDefinitionRepository(gormDB)
	resAttrRepo := postgres.NewResourceAttributeRepository(gormDB)

	// --- Service Instantiation ---
	identityService := identity.NewService(tenantRepo, userRepo, roleRepo, groupRepo, jwtSecret)
	// Corrected: resource.NewService only takes 3 arguments.
	resourceService := resource.NewService(resourceRepo, attrDefRepo, resAttrRepo)

	// --- HTTP Server Setup ---
	router := gin.Default()
	authMiddleware := api.NewAuthMiddleware(jwtSecret)

	// --- Handler Instantiation & Route Registration ---
	apiV1 := router.Group("/api/v1")

	// Corrected: NewIdentityHandler takes the identity.Service interface.
	identityHandler := api.NewIdentityHandler(identityService, permRepo, authMiddleware)
	identityHandler.RegisterRoutes(apiV1)

	// Corrected: NewResourceHandler takes 3 arguments.
	resourceHandler := api.NewResourceHandler(resourceService, permRepo, authMiddleware)
	resourceHandler.RegisterRoutes(apiV1)

	// --- Start Server ---
	log.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}