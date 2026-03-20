package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"curriculum-service/internal/config"
	"curriculum-service/internal/http/handlers"
	"curriculum-service/internal/http/middleware"
	"curriculum-service/internal/http/router"
	"curriculum-service/internal/repo/postgres"
	"curriculum-service/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env")

	cfg, err := config.ReadEnv()
	if err != nil {
		log.Fatal("configuration error:", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := postgres.NewPool(
		ctx,
		cfg.DatabaseURL(),
		cfg.DB.MaxConns,
		cfg.DB.MinConns,
		cfg.DB.MaxConnLifetime,
		cfg.DB.HealthCheckPeriod,
	)
	if err != nil {
		log.Fatal("error connecting to the database:", err)
	}
	defer pool.Close()

	catalogRepo := postgres.NewCatalogRepo(pool)
	catalogUC := usecase.NewCatalog(catalogRepo)
	catalogH := handlers.NewCatalogHandler(catalogUC)
	engine := router.New(
		catalogH,
		[]gin.HandlerFunc{
			middleware.CORS(cfg.CORS.AllowedOrigins, cfg.CORS.AllowedMethods, cfg.CORS.AllowedHeaders),
			middleware.RequestID(),
			middleware.Logger(),
			middleware.Recover(),
		},
	)

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           engine,
		ReadHeaderTimeout: cfg.HTTP.ReadHeaderTimeout,
	}

	go func() {
		log.Println("curriculum-service started on", cfg.HTTPAddr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Println("server error:", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.HTTP.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Println("shutdown error:", err)
	}

	log.Println("curriculum-service stopped")
}
