package main

import (
	"context"
	statushandler "curriculum-service/internal/http/handlers/status"
	"curriculum-service/internal/repo/postgres/catalog"
	statusrepo "curriculum-service/internal/repo/postgres/status"
	statususecase "curriculum-service/internal/usecase/status"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"curriculum-service/internal/http/handlers"
	"curriculum-service/internal/http/middleware"
	"curriculum-service/internal/http/router"
	"curriculum-service/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")

	cfg, err := ReadEnv()
	if err != nil {
		log.Fatal("configuration error:", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := NewPool(
		ctx,
		cfg.DB,
	)
	if err != nil {
		log.Fatal("error connecting to the database:", err)
	}
	defer pool.Close()

	db, err := NewDB(ctx, cfg.DB)
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}

	catalogRepo := catalog.NewRepo(pool)
	catalogUC := usecase.NewCatalog(catalogRepo)
	catalogHandler := handlers.NewCatalogHandler(catalogUC)

	// status courses
	statusRepo := statusrepo.NewRepo(db)
	statusUseCase := statususecase.New(statusRepo)
	statusHandler := statushandler.NewHandler(statusUseCase)

	handler := router.Handler{
		Catalog: catalogHandler,
		Status:  statusHandler,
	}

	engine := router.New(
		handler,
		[]gin.HandlerFunc{
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
