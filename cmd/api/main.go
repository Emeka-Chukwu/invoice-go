package main

import (
	"database/sql"
	"fmt"
	"go-invoice/app"
	activitylog_repository "go-invoice/internal/activitylog/repository"
	payment_repository "go-invoice/internal/payments/repository"
	"go-invoice/util"
	"go-invoice/worker"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/hibiken/asynq"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Msg("Cannot connect to db:")
	}
	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	fmt.Println(config.DBSource)
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal().Msg("Cannot connect to db:")
	}
	runDBMigration(config.MigrationURL, config.DBSource)

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}
	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)
	go runTaskProcessor(config, redisOpt, conn)
	runGinServer(config, conn, taskDistributor)

}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().Msg("cannot create new migrate instance:" + fmt.Sprintf("%v", err))
	}
	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Msg("failed to run migration up:")
	}
	log.Info().Msg("db migrated successfully")
}

func runGinServer(config util.Config, conn *sql.DB, taskDistributor worker.TaskDistributor) {
	server, err := app.NewServer(config, conn, taskDistributor)
	if err != nil {
		log.Fatal().Msg("cannot create server")
	}
	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Msg("cannot start server")
	}
}

func runTaskProcessor(config util.Config, redisOpt asynq.RedisClientOpt, store *sql.DB) {
	payStore := payment_repository.NewPaymentRepo(store)
	activityStore := activitylog_repository.NewAuthRepository(store)
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, payStore, activityStore, config)
	log.Info().Msg("start task processor")
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start task processor")
	}
}
