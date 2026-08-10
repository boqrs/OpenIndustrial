package main

import (
	"log"

	"github.com/OpenIndustrial/cloud/internal/api"
	"github.com/OpenIndustrial/cloud/internal/config" // UNCOMMENTED
	"github.com/OpenIndustrial/cloud/internal/identity"
	"github.com/OpenIndustrial/cloud/internal/persistence/postgres"
	"github.com/OpenIndustrial/cloud/internal/resource"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := postgres.NewDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	router := gin.Default()
	apiV1 := router.Group("/api/v1")

	// --- Repositories ---
	userRepo := postgres.NewUserRepository(db)
	tenantRepo := postgres.NewTenantRepository(db)
	roleRepo := postgres.NewRoleRepository(db)
	permRepo := postgres.NewPermissionRepository(db)
	resourceRepo := postgres.NewResourceRepository(db)
	groupRepo := postgres.NewGroupRepository(db)
	authzRepo := postgres.NewAuthorizationRepository(db)

	// --- Services ---
	// CORRECTED: Now passing jwtSecret as the last argument
	identityService := identity.NewService(tenantRepo, userRepo, roleRepo, groupRepo, cfg.JWTSecret)
	resourceService := resource.NewService(resourceRepo, groupRepo, authzRepo)

	// --- Handlers ---
	authMiddleware := api.NewAuthMiddleware(cfg.JWTSecret)
	identityHandler := api.NewIdentityHandler(identityService, permRepo, authMiddleware)
	resourceHandler := api.NewResourceHandler(resourceService, authzRepo, permRepo, authMiddleware)

	// --- Register Routes ---
	identityHandler.RegisterRoutes(apiV1)
	resourceHandler.RegisterRoutes(apiV1)

	log.Printf("server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}