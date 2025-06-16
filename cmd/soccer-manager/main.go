package main

import (
	"context"
	"errors"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/golang-migrate/migrate/v4"
	"github.com/hexley21/soccer-manager/cmd/util/shutdown"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/server"
	"github.com/hexley21/soccer-manager/pkg/config"
	"github.com/hexley21/soccer-manager/pkg/hasher/argon2"
	"github.com/hexley21/soccer-manager/pkg/infra/postgres"
	"github.com/hexley21/soccer-manager/pkg/logger/zap_logger"
	playground_validator "github.com/hexley21/soccer-manager/pkg/validator/playground_vlidator"
)

// @title Soccer Manager
// @version 1.0.0-alpha0
// @description Soccer manager api
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:80
// @BasePath /api
// @schemes http
// @securityDefinitions.apiKey AccessToken
// @in header
// @name Authorization
func main() {
	execFilepath, err := os.Executable()
	if err != nil {
		log.Fatalf("Could not read exec path: %v\n", err)
	}

	generalConfigPath := filepath.Join(filepath.Dir(execFilepath), "config/general.yml")
	serviceConfigPath := filepath.Join(filepath.Dir(execFilepath), "config/service.yml")
	envPath := filepath.Join(filepath.Dir(execFilepath), ".env")
	logPath := filepath.Join(filepath.Dir(execFilepath), "log/logs.log")

	cfg, err := config.LoadConfig(generalConfigPath, serviceConfigPath, envPath)
	if err != nil {
		log.Fatalf("Could not load config: %v\n", err)
	}

	zapLogger := zap_logger.New(logPath, cfg.Logging, cfg.IsProd)
	zapLogger.Debug(cfg)

	validator := playground_validator.New(zapLogger)

	snowflakeNode, err := snowflake.NewNode(cfg.Server.InstanceId)
	if err != nil {
		zapLogger.Fatal(err)
	}

	hasher := argon2.NewHasher(cfg.Argon2)

	pgCtx, pgCancel := context.WithTimeout(context.Background(), time.Second*10)
	defer pgCancel()

	pgPool, err := postgres.New(pgCtx, cfg.Postgres)
	if err != nil {
		zapLogger.Fatal(err)
	}

	migration, err := postgres.MigrateCfg(cfg.Postgres, "file://migrations")
	if err != nil {
		zapLogger.Fatal(err)
	}

	if err := migration.Migrate(cfg.Postgres.Version); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			zapLogger.Fatal(err)
		}
	}

	server := server.NewServer(cfg, zapLogger, validator, snowflakeNode, hasher, pgPool)

	shutCtx, shutCancel := context.WithCancelCause(context.Background())
	go shutdown.NotifyShutdown(shutCancel, zapLogger, server)

	zapLogger.Info("Tracker started...")
	if err := server.Run(); err != nil {
		zapLogger.Error(err)
	}

	// this is safe, program won't hang forever
	// the closer func passed to a notify shutdown has a timeout
	select {
	case <-shutCtx.Done():
		cause := context.Cause(shutCtx)
		if cause != nil {
			zapLogger.Errorf("Shutdown error: %v", cause)
		}
	}

	zapLogger.Info("Soccer manager stopped...")
}
