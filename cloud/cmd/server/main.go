package main

import (
	"log"

	"github.com/OpenIndustrial/cloud/internal/api"
	"github.com/OpenIndustrial/cloud/internal/config"
	//"github.com/OpenIndustrial/cloud/internal/persistence/postgres"
	"github.com/OpenIndustrial/cloud/internal/kernel/resource"
	"github.com/OpenIndustrial/cloud/internal/kernel/security"
	"github.com/OpenIndustrial/cloud/internal/persistence/postgres"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware is a placeholder for your actual authentication middleware.
// The compiler errors indicated that the ResourceHandler requires a `gin.HandlerFunc`,
// which is almost certainly for authentication and authorization.
// You should replace the contents of this function with your real auth logic,
// for example, validating a JWT and setting tenant/user info in the context.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("Executing placeholder authentication middleware...")
		// Example: c.Set("tenant_id", "some-tenant-id-from-token")
		c.Next()
	}
}

func main() {
	// 1. Load application configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	// 2. Connect to the database
	// FIX: Using `database.NewDb` as requested and passing the whole config,
	// as `cfg.Database` was undefined.
	db, err := postgres.NewDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// 3. Initialize ALL Repositories
	// FIX: The compiler showed that resource.New needs 4 repositories.
	resourceRepo := postgres.NewResourceRepository(db)
	attrDefRepo := postgres.NewAttributeDefinitionRepository(db)
	resAttrRepo := postgres.NewResourceAttributeRepository(db)
	resConnRepo := postgres.NewResourceConnectionsRepository(db)

	// FIX: The compiler showed that api.NewResourceHandler needs a PermissionRepository.
	permissionRepo := postgres.NewPermissionRepository(db)

	// Security Repositories
	credRepo := postgres.NewCredentialRepository(db)
	identityRepo := postgres.NewIdentityRepository(db)
	certRepo := postgres.NewCertificateRepository(db)

	// 4. Initialize Services
	// FIX: Calling `resource.New` with all 4 required repository dependencies.
	// I've removed the TransactionManager as it was undefined.
	resourceSvc := resource.NewService(resourceRepo, attrDefRepo, resAttrRepo, resConnRepo)

	// FIX: Renamed `security.NewService` to `security.New` and injected dependencies.
	securitySvc := security.NewService(resourceSvc, credRepo, identityRepo, certRepo)

	// 5. Setup HTTP Server & Handlers
	router := gin.Default()

	// FIX: Calling `api.NewResourceHandler` with all 3 required arguments,
	// including the placeholder authentication middleware.
	resourceHandler := api.NewResourceHandler(resourceSvc, permissionRepo, AuthMiddleware())
	securityHandler := api.NewSecurityHandler(securitySvc)

	// 6. Register Routes
	// FIX: Renamed `RegisterResourceRoutes` to `Register` as the original was undefined.
	resourceHandler.Register(router)
	securityHandler.RegisterSecurityRoutes(router)

	// 7. Start the server
	// FIX: Using `cfg.HTTP.Address` as a guess, since `cfg.Server` was undefined.
	// Please adjust if your config structure is different.
	log.Printf("server starting on %s", cfg.HTTP.Address)
	if err := router.Run(cfg.HTTP.Address); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}