package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mexirica/chi-template/internal/configs"
	"github.com/mexirica/chi-template/internal/db"
	"github.com/mexirica/chi-template/internal/o11y"
	rc "github.com/mexirica/chi-template/internal/redis"
	"github.com/mexirica/chi-template/internal/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server.
// @termsOfService http://localhost:6000/terms/

// @contact.name API Support
// @contact.url http://www.localhost:6000/support
// @contact.email support@localhost:6000

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
func main() {
	ctx := context.Background()
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()

	fn := o11y.InitTracer(ctx)
	defer fn(ctx)

	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Error().Msg(fmt.Sprintf("Error loading config: %v", err))
		return
	}

	dburl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.DB_HOST, cfg.DB_PORT, cfg.DB_USER, cfg.DB_PASSWORD, cfg.DB_NAME)

	dbConn, err := db.Connect(dburl)
	if err != nil {
		log.Error().Msg(fmt.Sprintf("Error connecting to database: %v", err))
		panic("Error connecting to database")
	}

	defer dbConn.Close()

	fmt.Printf("redis %s:%s\n", cfg.REDIS_HOST, cfg.REDIS_PORT)
	redisClient, err := rc.InitRedisClient(cfg.REDIS_HOST, cfg.REDIS_PORT)
	if err != nil {
		log.Error().Msg(fmt.Sprintf("Error connecting to Redis: %v", err))
	}

	defer redisClient.Close()

	app := server.New(&cfg, redisClient, dbConn)

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		log.Info().Msg("Desligando o servidor...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := app.Shutdown(ctx); err != nil {
			log.Printf("Erro ao desligar o servidor: %v", err)
		}
		close(idleConnsClosed)
	}()

	app.Serve()

	<-idleConnsClosed
}
