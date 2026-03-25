package main

import (
	"context"
	coursehandler "curriculum-service/internal/http/handlers/course"
	durationcategoryhandler "curriculum-service/internal/http/handlers/durationcategory"
	levelhandler "curriculum-service/internal/http/handlers/level"
	localehandler "curriculum-service/internal/http/handlers/locale"
	modulehandler "curriculum-service/internal/http/handlers/module"
	statushandler "curriculum-service/internal/http/handlers/status"
	taghandler "curriculum-service/internal/http/handlers/tag"
	topichandler "curriculum-service/internal/http/handlers/topic"
	courserepo "curriculum-service/internal/repo/postgres/course"
	durationcategoryrepo "curriculum-service/internal/repo/postgres/durationcategory"
	levelrepo "curriculum-service/internal/repo/postgres/level"
	localerepo "curriculum-service/internal/repo/postgres/locale"
	modulerepo "curriculum-service/internal/repo/postgres/module"
	statusrepo "curriculum-service/internal/repo/postgres/status"
	tagrepo "curriculum-service/internal/repo/postgres/tag"
	topicrepo "curriculum-service/internal/repo/postgres/topic"
	courseusecase "curriculum-service/internal/usecase/course"
	durationcategoryusecase "curriculum-service/internal/usecase/durationcategory"
	levelusecase "curriculum-service/internal/usecase/level"
	localeusecase "curriculum-service/internal/usecase/locale"
	moduleusecase "curriculum-service/internal/usecase/module"
	statususecase "curriculum-service/internal/usecase/status"
	tagusecase "curriculum-service/internal/usecase/tag"
	topicusecase "curriculum-service/internal/usecase/topic"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"curriculum-service/internal/http/middleware"
	"curriculum-service/internal/http/router"
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
		log.Fatal("error connecting to the database pool:", err)
	}
	defer pool.Close()

	db, err := NewDB(ctx, cfg.DB)
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}

	// status course
	statusRepo := statusrepo.NewRepo(db)
	statusUseCase := statususecase.New(statusRepo)
	statusHandler := statushandler.NewHandler(statusUseCase)

	// level course
	levelRepo := levelrepo.NewRepo(db)
	levelUseCase := levelusecase.New(levelRepo)
	levelHandler := levelhandler.NewHandler(levelUseCase)

	// duration category
	durationCategoryRepo := durationcategoryrepo.NewRepo(db)
	durationCategoryUseCase := durationcategoryusecase.New(durationCategoryRepo)
	durationCategoryHandler := durationcategoryhandler.NewHandler(durationCategoryUseCase)

	// topic course
	topicRepo := topicrepo.NewRepo(db)
	topicUseCase := topicusecase.New(topicRepo)
	topicHandler := topichandler.NewHandler(topicUseCase)

	// tag course
	tagRepo := tagrepo.NewRepo(db)
	tagUseCase := tagusecase.New(tagRepo)
	tagHandler := taghandler.NewHandler(tagUseCase)

	courseRepo := courserepo.NewRepo(db)
	courseUseCase := courseusecase.New(courseRepo)
	courseHandler := coursehandler.New(courseUseCase)

	LocaleRepo := localerepo.NewRepo(db)
	localeUseCase := localeusecase.New(LocaleRepo)
	localeHandler := localehandler.NewHandler(localeUseCase)

	moduleRepo := modulerepo.New(db)
	moduleUseCase := moduleusecase.New(moduleRepo)
	moduleHandler := modulehandler.NewHandler(moduleUseCase)

	handler := router.Handler{
		Status:           statusHandler,
		Level:            levelHandler,
		DurationCategory: durationCategoryHandler,
		Topic:            topicHandler,
		Tag:              tagHandler,
		Course:           courseHandler,
		Locale:           localeHandler,
		Module:           moduleHandler,
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
		Addr:              cfg.ListenAddr(),
		Handler:           engine,
		ReadHeaderTimeout: cfg.HTTP.ReadHeaderTimeout,
	}

	go func() {
		log.Println("curriculum-service started on", cfg.ListenAddr())
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
