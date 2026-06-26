package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"strconv"

	"github.com/eng-yousef-khaled/expenses_api/internal/env"
	"github.com/gomodule/redigo/redis"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		os.Exit(0)
	}
	mail_port, err := strconv.Atoi(env.GetString("MAIL_PORT", "432"))

	if err != nil {
		log.Printf("Invalid Mail Port")
		os.Exit(0)
	}
	cfg := config{
		addr: ":8080",
		db: dbConfig{
			dsn: env.GetString("GOOSE_DBSTRING", "user=expenses_app_user dbname=expenses_app_db password=EX560$#% sslmode=disable"),
		},
		mail: mailConfig{
			server:   env.GetString("MAIL_SERVER", "smtp.gmail.com"),
			port:     int16(mail_port),
			email:    env.GetString("MAIL_EMAIL", "example@gmail.com"),
			password: env.GetString("MAIL_PASSWORD", ""),
		},
		caching: cachingConfig{
			MaxActive: 5,
			MaxIdle:   5,
			Wait:      true,
			Dial:      func() (redis.Conn, error) { return redis.Dial("tcp", "127.0.0.1:3605") },
		},
		secretEncryptionKey: env.GetString("SECRET_ENCRYPTION_KEY", "secret_encryption_key"),
	}

	// logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	log.Printf("server mail config is %s", cfg.mail.server)

	conn, dbErr := pgx.Connect(ctx, cfg.db.dsn)
	if dbErr != nil {
		log.Printf("Connection to db fail with error %s", dbErr)
		os.Exit(0)
	}
	defer conn.Close(ctx)
	logger.Info("connected to database")
	app := application{
		config: cfg,
		db:     conn,
	}
	app.run(app.mount())
}
