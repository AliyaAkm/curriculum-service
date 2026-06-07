package main

import (
	"context"
	achievementhandler "curriculum-service/internal/http/handlers/achievement"
	certificatehandler "curriculum-service/internal/http/handlers/certificate"
	codeattempthandler "curriculum-service/internal/http/handlers/codeattempt"
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
	studentstatshandler "curriculum-service/internal/http/handlers/studentstats"
	taghandler "curriculum-service/internal/http/handlers/tag"
	topichandler "curriculum-service/internal/http/handlers/topic"
	dictionarycache "curriculum-service/internal/repo/cache/dictionary"
	achievementrepo "curriculum-service/internal/repo/postgres/achievement"
	certificaterepo "curriculum-service/internal/repo/postgres/certificate"
	codeattemptrepo "curriculum-service/internal/repo/postgres/codeattempt"
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
	studentstatsrepo "curriculum-service/internal/repo/postgres/studentstats"
	tagrepo "curriculum-service/internal/repo/postgres/tag"
	topicrepo "curriculum-service/internal/repo/postgres/topic"
	"curriculum-service/internal/service/aianalytics"
	cacheclient "curriculum-service/internal/service/cache"
	"curriculum-service/internal/service/coderunner"
	"curriculum-service/internal/service/notification"
	"curriculum-service/internal/service/storage"
	achievementusecase "curriculum-service/internal/usecase/achievement"
	certificateusecase "curriculum-service/internal/usecase/certificate"
	codeattemptusecase "curriculum-service/internal/usecase/codeattempt"
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
	studentstatsusecase "curriculum-service/internal/usecase/studentstats"
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
	dictionaryCache := cacheclient.NewJSONCache(cacheclient.RedisConfig{
		Addr:      cfg.Redis.Addr,
		Password:  cfg.Redis.Password,
		DB:        cfg.Redis.DB,
		KeyPrefix: cfg.Redis.KeyPrefix,
	})
	defer func() {
		if err := dictionaryCache.Close(); err != nil {
			log.Println("redis close error:", err)
		}
	}()

	notificationClient, err := notification.NewClient(notification.ClientConfig{
		BaseURL:        cfg.Notification.URL,
		Timeout:        cfg.Notification.Timeout,
		InternalAPIKey: cfg.Notification.InternalAPIKey,
	})
	if err != nil {
		log.Fatal("error configuring notification service client:", err)
	}

	achievementRepo := achievementrepo.NewRepo(db)
	achievementUseCase := achievementusecase.New(achievementRepo)
	achievementHandler := achievementhandler.NewHandler(achievementUseCase, jwtMgr)

	// daily streak
	streakRepo := streakrepo.NewRepo(db)
	streakUseCase := streakusecase.New(streakRepo, notificationClient, achievementUseCase)
	streakHandler := streakhandler.NewHandler(streakUseCase, jwtMgr)

	// status course
	statusRepo := statusrepo.NewRepo(db)
	cachedStatusRepo := dictionarycache.NewStatus(statusRepo, dictionaryCache, cfg.Redis.DictionaryTTL)
	statusUseCase := statususecase.New(cachedStatusRepo)
	statusHandler := statushandler.NewHandler(statusUseCase)

	// level course
	levelRepo := levelrepo.NewRepo(db)
	cachedLevelRepo := dictionarycache.NewLevel(levelRepo, dictionaryCache, cfg.Redis.DictionaryTTL)
	levelUseCase := levelusecase.New(cachedLevelRepo)
	levelHandler := levelhandler.NewHandler(levelUseCase)

	// duration category
	durationCategoryRepo := durationcategoryrepo.NewRepo(db)
	cachedDurationCategoryRepo := dictionarycache.NewDurationCategory(durationCategoryRepo, dictionaryCache, cfg.Redis.DictionaryTTL)
	durationCategoryUseCase := durationcategoryusecase.New(cachedDurationCategoryRepo)
	durationCategoryHandler := durationcategoryhandler.NewHandler(durationCategoryUseCase)

	// topic course
	topicRepo := topicrepo.NewRepo(db)
	cachedTopicRepo := dictionarycache.NewTopic(topicRepo, dictionaryCache, cfg.Redis.DictionaryTTL)
	topicUseCase := topicusecase.New(cachedTopicRepo)
	topicHandler := topichandler.NewHandler(topicUseCase)

	// tag course
	tagRepo := tagrepo.NewRepo(db)
	cachedTagRepo := dictionarycache.NewTag(tagRepo, dictionaryCache, cfg.Redis.DictionaryTTL)
	tagUseCase := tagusecase.New(cachedTagRepo)
	tagHandler := taghandler.NewHandler(tagUseCase)

	// отзывы
	reviewRepo := reviewrepo.NewRepo(db)
	reviewUseCase := reviewusecase.New(reviewRepo)
	reviewHandler := reviewhandler.NewHandler(reviewUseCase, validate)

	LocaleRepo := localerepo.NewRepo(db)
	cachedLocaleRepo := dictionarycache.NewLocale(LocaleRepo, dictionaryCache, cfg.Redis.DictionaryTTL)
	localeUseCase := localeusecase.New(cachedLocaleRepo)
	localeHandler := localehandler.NewHandler(localeUseCase)

	moduleRepo := modulerepo.New(db)
	moduleUseCase := moduleusecase.New(moduleRepo)
	moduleHandler := modulehandler.NewHandler(moduleUseCase, jwtMgr)

	storageClient, err := storage.NewClient(storage.ClientConfig{
		BaseURL: cfg.Storage.URL,
		Timeout: cfg.Storage.Timeout,
	})
	if err != nil {
		log.Fatal("error configuring storage service client:", err)
	}

	courseRepo := courserepo.NewRepo(db)
	courseUseCase := courseusecase.New(courseRepo, reviewRepo, moduleRepo, notificationClient)
	courseHandler := coursehandler.New(courseUseCase, jwtMgr)

	lessonRepo := lessonrepo.NewRepo(db)
	lessonUseCase := lessonusecase.New(lessonRepo)
	lessonHandler := lessonhandler.NewHandler(lessonUseCase, localeUseCase, storageClient, jwtMgr)

	progressRepo := progressrepo.NewRepo(db)
	progressUseCase := progressusecase.New(progressRepo, notificationClient, achievementUseCase)
	progressHandler := progresshandler.NewHandler(progressUseCase, jwtMgr)

	certificateRepo := certificaterepo.NewRepo(db)
	certificateUseCase := certificateusecase.NewUseCase(certificateRepo, storageClient)
	certificateHandler := certificatehandler.NewHandler(certificateUseCase, jwtMgr)

	quizRepo := quizrepo.NewRepo(db)
	quizUseCase := quizusecase.New(quizRepo, notificationClient, achievementUseCase)
	quizHandler := quizhandler.NewHandler(quizUseCase, jwtMgr)

	coursePointRepo := coursepointrepo.NewRepo(db)
	coursePointUseCase := coursepointusecase.New(coursePointRepo)
	coursePointHandler := coursepointhandler.New(coursePointUseCase)

	codeRunnerClient, err := coderunner.NewClient(coderunner.ClientConfig{
		BaseURL: cfg.CodeRunner.URL,
		Timeout: cfg.CodeRunner.Timeout,
	})
	if err != nil {
		log.Fatal("error configuring code runner service client:", err)
	}
	aiAnalyticsClient, err := aianalytics.NewClient(aianalytics.ClientConfig{
		BaseURL: cfg.AI.URL,
		Timeout: cfg.AI.Timeout,
	})
	if err != nil {
		log.Fatal("error configuring ai analytics service client:", err)
	}
	studentStatsRepo := studentstatsrepo.NewRepo(db)
	studentStatsUseCase := studentstatsusecase.New(studentStatsRepo, aiAnalyticsClient)
	studentStatsHandler := studentstatshandler.NewHandler(studentStatsUseCase, jwtMgr)

	practiceRepo := practicerepo.NewRepo(db)
	practiceUseCase := practiceusecase.New(practiceRepo)
	practiceHandler := practicehandler.NewHandler(practiceUseCase, jwtMgr)

	codeAttemptRepo := codeattemptrepo.NewRepo(db)
	codeAttemptUseCase := codeattemptusecase.New(codeAttemptRepo, codeRunnerClient, practiceUseCase)
	codeAttemptHandler := codeattempthandler.NewHandler(codeAttemptUseCase, jwtMgr)

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
		Progress:         progressHandler,
		Certificate:      certificateHandler,
		Quiz:             quizHandler,
		Review:           reviewHandler,
		CoursePoint:      coursePointHandler,
		Streak:           streakHandler,
		CodeAttempt:      codeAttemptHandler,
		StudentStats:     studentStatsHandler,
		Practice:         practiceHandler,
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
