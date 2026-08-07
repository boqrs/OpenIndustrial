package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/OpenGongChang/OpenIndustrial/cloud/internal/bootstrap"
	"github.com/OpenGongChang/OpenIndustrial/cloud/internal/device"
	"github.com/OpenGongChang/OpenIndustrial/cloud/internal/event"
	"github.com/OpenGongChang/OpenIndustrial/cloud/internal/product"
	"github.com/OpenGongChang/OpenIndustrial/cloud/internal/workorder"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Setup context for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// 2. Initialize event bus and services container
	bus := event.NewMemoryBus()
	container := bootstrap.NewContainer(ctx, bus)
	log.Println("Application container initialized.")

	// 3. Setup HTTP server
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// All primary routes are grouped by organization
	orgRoutes := router.Group("/orgs/:org_id")

	// 4. Register API routes for each domain
	productAPI := product.NewAPI(container.ProductSvc)
	productAPI.RegisterRoutes(orgRoutes)

	deviceAPI := device.NewAPI(container.DeviceSvc)
	deviceAPI.RegisterRoutes(orgRoutes)

	workorderAPI := workorder.NewAPI(container.WorkorderSvc)
	workorderAPI.RegisterRoutes(orgRoutes)

	log.Println("API routes registered.")

	// 5. Start the server with graceful shutdown
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		log.Printf("Listening on %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Println("Shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the requests it is currently handling
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	fmt.Println("Server exiting")
}