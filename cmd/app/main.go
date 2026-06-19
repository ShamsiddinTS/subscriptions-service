// @title API сервиса управления подписками
// @version 1.0
// @description REST API для управления онлайн-подписками пользователей.
// @description Поддерживает CRUDL-операции и расчет суммарной стоимости подписок за выбранный период.
// @host localhost:8080
// @BasePath /api

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

	"github.com/ShamsiddinTS/subscriptions-service/internal/config"
	"github.com/ShamsiddinTS/subscriptions-service/internal/database"
	"github.com/ShamsiddinTS/subscriptions-service/internal/handler"
	loggerPkg "github.com/ShamsiddinTS/subscriptions-service/internal/logger"
	"github.com/ShamsiddinTS/subscriptions-service/internal/migrations"
	"github.com/ShamsiddinTS/subscriptions-service/internal/repository"
	routerPkg "github.com/ShamsiddinTS/subscriptions-service/internal/router"
	"github.com/ShamsiddinTS/subscriptions-service/internal/service"
	"github.com/joho/godotenv"
	"go.uber.org/zap"

	_ "github.com/ShamsiddinTS/subscriptions-service/docs"
)

func main() {
	// ***** ENV *****
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file")
	}

	// ***** Config *****
	cfg := config.ReadConfig()

	// ***** Logger *****
	appLogger, err := loggerPkg.New()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := appLogger.Sync(); err != nil {
			log.Println("failed to sync logger:", err)
		}
	}()

	appLogger.Info("logger initialized")

	// ***** DB *****
	ctx := context.Background()

	db, err := database.New(cfg.Db)
	if err != nil {
		appLogger.Fatal("failed to connect to database", zap.Error(err))
	}

	defer func() {
		if err := db.Close(); err != nil {
			appLogger.Error("failed to close database connection", zap.Error(err))
		}
	}()

	if err := db.PingContext(ctx); err != nil {
		appLogger.Fatal("database ping failed", zap.Error(err))
	}

	appLogger.Info("database connected")

	// ***** Migrations *****
	migrationDir, err := migrations.ResolveMigrationsDir(cfg.Db.MigrationsDir)
	if err != nil {
		appLogger.Fatal("failed to resolve migrations directory", zap.Error(err))
	}

	appLogger.Info(
		"migrations directory resolved",
		zap.String("dir", migrationDir),
	)

	migCtx, migCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer migCancel()

	if err := migrations.Run(migCtx, db, migrationDir); err != nil {
		appLogger.Fatal("failed to run migrations", zap.Error(err))
	}

	appLogger.Info("migrations completed")

	// ***** Dependency Injection *****
	subscriptionRepo := repository.NewSubscriptionRepository(db)

	subscriptionService := service.NewSubscriptionService(
		subscriptionRepo,
	)

	subscriptionHandler := handler.NewSubscriptionHandler(
		subscriptionService,
		appLogger,
	)

	engine := routerPkg.SetupRouter(subscriptionHandler)

	// ***** HTTP Server *****
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Endpoint.Host, cfg.Endpoint.Port),
		Handler:      engine,
		ReadTimeout:  time.Duration(cfg.Endpoint.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Endpoint.WriteTimeout) * time.Second,
	}

	// ***** Run Server *****
	go func() {
		appLogger.Info(
			"SERVER STARTED",
			zap.String("addr", server.Addr),
		)

		if err := server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			appLogger.Fatal("failed to start server", zap.Error(err))
		}
	}()

	// ***** Graceful Shutdown *****
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	<-ch

	appLogger.Info("shutdown signal received")

	shutdownCtx, cancelShutdown := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancelShutdown()

	if err := server.Shutdown(shutdownCtx); err != nil {
		appLogger.Fatal("failed to shutdown server", zap.Error(err))
	}

	appLogger.Info(
		"server stopped gracefully",
		zap.String("addr", server.Addr),
	)
}
