package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"joylanguageschool.ru/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
)

type contextKey string

var contextKeyUser = contextKey("user")

type config struct {
	port int
	env  string
	db   struct {
		dsn             string
		SetMaxOpenConns int
		SetMaxIdleConns int
		SetMaxIdleTime  string
	}
}

type application struct {
	config        config
	logger        *log.Logger
	session       *sessions.Session
	posts         *mysql.PostModel
	templateCache map[string]*template.Template
	users         *mysql.UserModel
}

func main() {
	// Define configuration settings
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "HTTP network address")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	var secret = "secret"
	flag.StringVar(&secret, "secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret Key")

	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("JOY_DB_DSN"), "MYSQL DSN")

	flag.IntVar(&cfg.db.SetMaxOpenConns, "db-max-open-conns", 25, "MYSQL max open connections")
	flag.IntVar(&cfg.db.SetMaxIdleConns, "db-max-idle-conns", 25, "MYSQL max idle connections")
	flag.StringVar(&cfg.db.SetMaxIdleTime, "db-max-idle-time", "15m", "MYSQL max connection idle time")

	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()

	// Initialize a new template cache
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		logger.Fatal(err)
	}

	logger.Printf("database connection pool established")

	session := sessions.New([]byte(secret))
	session.Lifetime = 12 * time.Hour

	// Application dependencies
	app := &application{
		config:        cfg,
		logger:        logger,
		session:       session,
		posts:         &mysql.PostModel{DB: db},
		templateCache: templateCache,
		users:         &mysql.UserModel{DB: db},
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("starting %s server on %s", cfg.env, srv.Addr)
	err = srv.ListenAndServe()
	if err != nil {
		logger.Fatal(err)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
