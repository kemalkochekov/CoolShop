package app

import (
	"Coolshop/internal/auth"
	"Coolshop/internal/config"
	"Coolshop/internal/connection"
	"Coolshop/logger"
	"Coolshop/pkg/constant"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pressly/goose/v3"
)

type App struct {
	serviceProvider *serviceProvider
	server          *gin.Engine
	psqlDB          connection.DB
	redisDB         connection.Cache
	jwtGen          auth.JwtGen
	zapLogger       logger.LoggerInterface
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{server: gin.New()}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) initConfig(_ context.Context) error {
	err := config.LoadEnv(".env")
	if err != nil {
		return err
	}

	log.Println("Configuration loaded successfully.")

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(ctx context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initDB,
		a.initJwt,
		a.initLogger,
		a.initServer,
		a.runMigrations,
	}
	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) Run(ctx context.Context) error {
	return a.runServer(ctx)
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()

	log.Println("Service provider initialized.")

	return nil
}

func (a *App) runMigrations(ctx context.Context) error {
	dbConfig := a.serviceProvider.config.DbConfig()

	db, err := sql.Open("postgres", dbConfig.GenerateDSN())
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Println("Failed to close database connection")
		}
	}(db)

	if err := goose.Up(db, "./internal/migrations"); err != nil {
		log.Fatalf("Error setting up the database migrations: %v", err)
	}

	log.Println("Database migrations applied successfully.")

	return nil
}

func (a *App) initDB(ctx context.Context) error {
	cfg, err := a.serviceProvider.Config()
	if err != nil {
		return err
	}

	ctxTime, timeCancel := context.WithTimeout(ctx, constant.CtxTime*time.Second)

	defer timeCancel()

	a.psqlDB, err = connection.NewDB(ctxTime, cfg.DbConfig())
	if err != nil {
		return fmt.Errorf("connection.NewDB: %w", err)
	}

	a.redisDB, err = connection.NewCache(ctxTime, cfg.RDConfig())
	if err != nil {
		return fmt.Errorf("connection.NewCache: %w", err)
	}

	return nil
}

func (a *App) initJwt(ctx context.Context) error {
	cfg, err := a.serviceProvider.Config()
	if err != nil {
		return err
	}
	a.jwtGen = auth.NewTokenGenerator(cfg.SecretKeyByte())

	return nil
}
func (a *App) initLogger(ctx context.Context) error {
	zapLogger, err := logger.NewZapLogger()
	if err != nil {
		return err
	}

	a.zapLogger = zapLogger

	return nil
}

func (a *App) initServer(ctx context.Context) error {
	handler := a.serviceProvider.UserHandler(a.psqlDB, a.redisDB, a.jwtGen, a.zapLogger)
	a.initializeRoutes(handler)

	return nil
}

func (a *App) runServer(ctx context.Context) error {
	httpSrv := &http.Server{
		Addr:    a.serviceProvider.config.Address(),
		Handler: a.server,
	}
	go func() {
		if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to start server: %s", err.Error())
		}
	}()

	log.Println("Server is running on", a.serviceProvider.config.Address())

	defer func() {
		err := a.psqlDB.Close()
		if err != nil {
			log.Printf("[main] psqlDB.Close: %v\n", err)
		}
	}()

	<-ctx.Done()

	if err := httpSrv.Shutdown(ctx); err != nil {
		return err
	}

	log.Println("Server stopped gracefully.")

	return nil
}
