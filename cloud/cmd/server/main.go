package main

import (
	"log"

	"github.com/OpenIndustrial/cloud/internal/api"
	"github.com/OpenIndustrial/cloud/internal/identity"
	"github.com/OpenIndustrial/cloud/internal/persistence/postgres"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func main() {
	// Configuration would be loaded here via Viper in a real app
	dbConnectionString := "postgres://user:password@localhost:5432/dbname?sslmode=disable"

	// Initialize database connection
	db, err := sqlx.Connect("pgx", dbConnectionString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Dependency Injection: Wire everything together
	// 1. Create repository implementations
	tenantRepo := postgres.NewTenantRepository(db)
	userRepo := postgres.NewUserRepository(db)
	roleRepo := postgres.NewRoleRepository(db)
	permRepo := postgres.NewPermissionRepository(db) // 新增


	// 2. Create the service
	identityService := identity.NewService(tenantRepo, userRepo, roleRepo)

	// 3. Create the HTTP handler
	identityHandler := api.NewIdentityHandler(identityService, permRepo)

	// Initialize Gin router
	router := gin.Default()

	// Setup routes
	apiGroup := router.Group("/api/v1")
	identityHandler.RegisterRoutes(apiGroup)

	// Start the server
	log.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}