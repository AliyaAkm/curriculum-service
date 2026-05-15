package main

import (
	"context"
	achievementhandler "curriculum-service/internal/http/handlers/achievement"
	coursehandler "curriculum-service/internal/http/handlers/course"
	coursepointhandler "curriculum-service/internal/http/handlers/coursepoint"
	durationcategoryhandler "curriculum-service/internal/http/handlers/durationcategory"
	lessonhandler "curriculum-service/internal/http/handlers/lesson"
	levelhandler "curriculum-service/internal/http/handlers/level"
	localehandler "curriculum-service/internal/http/handlers/locale"
	modulehandler "curriculum-service/internal/http/handlers/module"
	practicehandler "curriculum-service/internal/http/handlers/practice"
	progresshandler "curriculum-service/internal/http/handlers/progress"
	quizhandler "curriculum-service/internal/http/handlers/quiz"
	reviewhandler "curriculum-service/internal/http/handlers/review"
	statushandler "curriculum-service/internal/http/handlers/status"
	streakhandler "curriculum-service/internal/http/handlers/streak"
	taghandler "curriculum-service/internal/http/handlers/tag"
	topichandler "curriculum-service/internal/http/handlers/topic"
	achievementrepo "curriculum-service/internal/repo/postgres/achievement"
	courserepo "curriculum-service/internal/repo/postgres/course"
	coursepointrepo "curriculum-service/internal/repo/postgres/coursepoint"
	durationcategoryrepo "curriculum-service/internal/repo/postgres/durationcategory"
	lessonrepo "curriculum-service/internal/repo/postgres/lesson"
	levelrepo "curriculum-service/internal/repo/postgres/level"
	localerepo "curriculum-service/internal/repo/postgres/locale"
	modulerepo "curriculum-service/internal/repo/postgres/module"
	practicerepo "curriculum-service/internal/repo/postgres/practice"
	progressrepo "curriculum-service/internal/repo/postgres/progress"
	quizrepo "curriculum-service/internal/repo/postgres/quiz"
	reviewrepo "curriculum-service/internal/repo/postgres/review"
	statusrepo "curriculum-service/internal/repo/postgres/status"
	streakrepo "curriculum-service/internal/repo/postgres/streak"
	tagrepo "curriculum-service/internal/repo/postgres/tag"
	topicrepo "curriculum-service/internal/repo/postgres/topic"
	"curriculum-service/internal/service/storage"
	achievementusecase "curriculum-service/internal/usecase/achievement"
	courseusecase "curriculum-service/internal/usecase/course"
	coursepointusecase "curriculum-service/internal/usecase/coursepoint"
	durationcategoryusecase "curriculum-service/internal/usecase/durationcategory"
	lessonusecase "curriculum-service/internal/usecase/lesson"
	levelusecase "curriculum-service/internal/usecase/level"
	localeusecase "curriculum-service/internal/usecase/locale"
	moduleusecase "curriculum-service/internal/usecase/module"
	practiceusecase "curriculum-service/internal/usecase/practice"
	progressusecase "curriculum-service/internal/usecase/progress"
	quizusecase "curriculum-service/internal/usecase/quiz"
	reviewusecase "curriculum-service/internal/usecase/review"
	statususecase "curriculum-service/internal/usecase/status"
	streakusecase "curriculum-service/internal/usecase/streak"
	tagusecase "curriculum-service/internal/usecase/tag"
	topicusecase "curriculum-service/internal/usecase/topic"
	"errors"
	"github.com/go-playground/validator/v10"
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

	jwtMgr := middleware.New(
		[]byte(cfg.JWT.Secret),
		cfg.JWT.Issuer,
		cfg.JWT.Audience,
		cfg.JWT.AccessTTL,
	)

	validate := validator.New()

	// daily streak
	streakRepo := streakrepo.NewRepo(db)
	streakUseCase := streakusecase.New(streakRepo)
	streakHandler := streakhandler.NewHandler(streakUseCase, jwtMgr)

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

	// отзывы
	reviewRepo := reviewrepo.NewRepo(db)
	reviewUseCase := reviewusecase.New(reviewRepo)
	reviewHandler := reviewhandler.NewHandler(reviewUseCase, validate)

	LocaleRepo := localerepo.NewRepo(db)
	localeUseCase := localeusecase.New(LocaleRepo)
	localeHandler := localehandler.NewHandler(localeUseCase)

	moduleRepo := modulerepo.New(db)
	moduleUseCase := moduleusecase.New(moduleRepo)
	moduleHandler := modulehandler.NewHandler(moduleUseCase, jwtMgr)

	videoStore, err := storage.NewMinIO(storage.MinIOConfig{
		Endpoint:       cfg.MinIO.Endpoint,
		PublicEndpoint: cfg.MinIO.PublicEndpoint,
		AccessKey:      cfg.MinIO.AccessKey,
		SecretKey:      cfg.MinIO.SecretKey,
		Bucket:         cfg.MinIO.Bucket,
		Region:         cfg.MinIO.Region,
		UseSSL:         cfg.MinIO.UseSSL,
		PresignTTL:     cfg.MinIO.PresignTTL,
	})
	if err != nil {
		log.Fatal("error configuring minio storage:", err)
	}

	courseRepo := courserepo.NewRepo(db)
	courseUseCase := courseusecase.New(courseRepo, reviewRepo, moduleRepo)
	courseHandler := coursehandler.New(courseUseCase, jwtMgr)

	lessonRepo := lessonrepo.NewRepo(db)
	lessonUseCase := lessonusecase.New(lessonRepo)
	lessonHandler := lessonhandler.NewHandler(lessonUseCase, localeUseCase, videoStore, jwtMgr)

	practiceRepo := practicerepo.NewRepo(db)
	practiceUseCase := practiceusecase.New(practiceRepo)
	practiceHandler := practicehandler.NewHandler(practiceUseCase, jwtMgr)

	progressRepo := progressrepo.NewRepo(db)
	progressUseCase := progressusecase.New(progressRepo)
	progressHandler := progresshandler.NewHandler(progressUseCase, jwtMgr)

	achievementRepo := achievementrepo.NewRepo(db)
	achievementUseCase := achievementusecase.New(achievementRepo)
	achievementHandler := achievementhandler.NewHandler(achievementUseCase, jwtMgr)

	quizRepo := quizrepo.NewRepo(db)
	quizUseCase := quizusecase.New(quizRepo)
	quizHandler := quizhandler.NewHandler(quizUseCase, jwtMgr)

	coursePointRepo := coursepointrepo.NewRepo(db)
	coursePointUseCase := coursepointusecase.New(coursePointRepo)
	coursePointHandler := coursepointhandler.New(coursePointUseCase)

	handler := router.Handler{
		Achievement:      achievementHandler,
		Status:           statusHandler,
		Level:            levelHandler,
		DurationCategory: durationCategoryHandler,
		Topic:            topicHandler,
		Tag:              tagHandler,
		Course:           courseHandler,
		Locale:           localeHandler,
		Module:           moduleHandler,
		Lesson:           lessonHandler,
		Practice:         practiceHandler,
		Progress:         progressHandler,
		Quiz:             quizHandler,
		Review:           reviewHandler,
		CoursePoint:      coursePointHandler,
		Streak:           streakHandler,
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
