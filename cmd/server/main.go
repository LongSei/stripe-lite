package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	v1 "github.com/thep1ckaxe91/stripe-lite/internal/api/v1"
	"github.com/thep1ckaxe91/stripe-lite/internal/config"
	"github.com/thep1ckaxe91/stripe-lite/internal/docs"
)

const version = "1.0.0"

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Configuration bootstrap failure: %v", err)
	}

	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// Mount documentation & OpenAPI schemas
	if err := docs.RegisterRoutes(router); err != nil {
		log.Fatalf("Failed to initialize documentation routes: %v", err)
	}

	// Root-level health probe alias for container orchestrators (k8s/nomad)
	router.GET("/healthz", func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, cfg.BasePath+"/healthz")
	})

	// Instantiate strictly-typed domain server
	domainServer := v1.NewServer(version)
	strictHandler := v1.NewStrictHandler(domainServer, nil)

	// Register generated Gin routes with strict server contract
	v1.RegisterHandlersWithOptions(router, strictHandler, v1.GinServerOptions{
		BaseURL: cfg.BasePath,
	})

	serverAddr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	srv := &http.Server{
		Addr:              serverAddr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		log.Printf("[Server] Starting Stripe-lite API gateway on %s (base: %s, env: %s)",
			serverAddr, cfg.BasePath, cfg.Environment)
		log.Printf("[Server] Swagger UI available at http://%s/docs", serverAddr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server listen error: %v", err)
		}
	}()

	// Graceful termination handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("[Server] Graceful shutdown sequence initiated...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("[Server] Forced server termination: %v", err)
	}

	log.Println("[Server] Gateway shut down cleanly.")
}
