package main

import (
	"log"
	"net/http"
	"time"

	repo "github.com/eng-yousef-khaled/expenses_api/internal/adapters/postgresql/sqlc"
	"github.com/eng-yousef-khaled/expenses_api/internal/application/auth"
	httpserver "github.com/eng-yousef-khaled/expenses_api/internal/inbound/http_server"
	"github.com/eng-yousef-khaled/expenses_api/internal/inbound/json"
	"github.com/eng-yousef-khaled/expenses_api/internal/outbound/gomailing"
	passwordhashing "github.com/eng-yousef-khaled/expenses_api/internal/outbound/password_hashing"
	"github.com/eng-yousef-khaled/expenses_api/internal/outbound/queue"
	userrepo "github.com/eng-yousef-khaled/expenses_api/internal/outbound/repo"
	"github.com/gocraft/work"
	"golang.org/x/crypto/bcrypt"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
)

const APP_NAME = "expanses_api"

func (app *application) mount() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID) // important for rate limiting
	r.Use(middleware.RealIP)    // import for rate limiting and analytics and tracing
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer) // recover from crashes
	r.Use(middleware.Timeout(time.Minute))
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		json.Write(w, 200, "Working ...")
	})
	redisAdapter := queue.NewRedisAdapter(queue.RedisConfig(app.config.caching), APP_NAME)

	mailAdapter := gomailing.NewGoMail(app.config.mail.server, app.config.mail.email, app.config.mail.password, int(app.config.mail.port))

	postgresUserRepo := userrepo.NewRepo(repo.New(app.db))
	bcryptPasswordHash := passwordhashing.CreateBcryptPasswordHash(bcrypt.DefaultCost)
	user_service := auth.NewService(postgresUserRepo, redisAdapter, bcryptPasswordHash, mailAdapter)
	jobHandler := &userrepo.JobHandler{Service: user_service}
	pool := work.NewWorkerPool(userrepo.JobHandler{}, 10, APP_NAME, redisAdapter.Pool())
	go pool.Start()
	pool.Job("send_verification_mail", jobHandler.ProccessSendMail)
	auth_handler := httpserver.CreateHandler(user_service)
	// Users
	// httpServer := server.NewHttpServer(user_service)

	r.Post("/auth/register", auth_handler.RegisterUser)
	return r
}

func (app *application) run(h http.Handler) error {
	server := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Minute,
	}
	log.Printf("Server start at address http://localhost%s", app.config.addr)
	return server.ListenAndServe()
}

type application struct {
	config config
	// Global DB connection
	db *pgx.Conn
}

type config struct {
	addr    string
	db      dbConfig
	mail    mailConfig
	caching cachingConfig
}

type dbConfig struct {
	dsn string
}

type mailConfig struct {
	server   string
	port     int16
	email    string
	password string
}

type cachingConfig struct {
	MaxActive int16
	MaxIdle   int16
	Wait      bool
	Dial      any
}
