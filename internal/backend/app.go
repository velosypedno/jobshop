package backend

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger() *zap.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(encoderCfg)
	consoleWriter := zapcore.Lock(os.Stdout)

	level := zapcore.InfoLevel
	consoleCore := zapcore.NewCore(consoleEncoder, consoleWriter, level)
	core := zapcore.NewTee(consoleCore)
	logger := zap.New(core, zap.AddCaller())
	return logger
}

type DBConfig struct {
	Driver string `envconfig:"DB_DRIVER" required:"true"`
	Host   string `envconfig:"DB_HOST" required:"true"`
	Port   string `envconfig:"DB_PORT" required:"true"`
	User   string `envconfig:"DB_USER" required:"true"`
	Pass   string `envconfig:"DB_PASSWORD" required:"true"`
	Name   string `envconfig:"DB_NAME" required:"true"`
}

func (c DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Pass, c.Name,
	)
}

type SrvConfig struct {
	Port string `envconfig:"API_PORT" required:"true"`
}

type CORSConfig struct {
	AllowOrigins     []string `envconfig:"CORS_ALLOW_ORIGINS" required:"true"`
	AllowMethods     []string `envconfig:"CORS_ALLOW_METHODS" required:"true"`
	AllowHeaders     []string `envconfig:"CORS_ALLOW_HEADERS" required:"true"`
	ExposeHeaders    []string `envconfig:"CORS_EXPOSE_HEADERS" required:"true"`
	AllowCredentials bool     `envconfig:"CORS_ALLOW_CREDENTIALS" required:"true"`
	MaxAge           int      `envconfig:"CORS_MAX_AGE" required:"true"`
}

type Config struct {
	DB   DBConfig
	Srv  SrvConfig
	CORS CORSConfig
}

func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

const (
	shutdownTimeout = 20 * time.Second
	readTimeout     = 15 * time.Second
)

type App struct {
	cfg *Config
	DB  *sql.DB
	API *http.Server
	l   *zap.Logger
}

func New(cfg *Config, l *zap.Logger) *App {
	return &App{
		cfg: cfg,
		l:   l,
	}
}

func (a *App) Run(ctx context.Context) error {
	var err error

	// db
	a.DB, err = sql.Open(a.cfg.DB.Driver, a.cfg.DB.DSN())
	if err != nil {
		return err
	}
	a.l.Info("DB connected")

	// api
	router := a.setupRouter()
	a.API = &http.Server{
		Addr:        ":" + a.cfg.Srv.Port,
		Handler:     router,
		ReadTimeout: readTimeout,
	}
	go func() {
		if err := a.API.ListenAndServe(); err != nil {
			a.l.Error("API server error", zap.Error(err))
		}
	}()
	a.l.Info("API server started", zap.String("port", a.cfg.Srv.Port))

	// wait for shutdown signal
	<-ctx.Done()
	a.l.Info("Shutdown signal received")

	// graceful shutdown
	timeoutCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	a.l.Info("Shutting down...")
	err = a.shutdown(timeoutCtx)
	if err != nil {
		a.l.Error("Shutdown error", zap.Error(err))
	}
	return err
}

func (a *App) shutdown(timeoutCtx context.Context) error {
	var shutdownErr error

	// api
	if a.API != nil {
		if err := a.API.Shutdown(timeoutCtx); err != nil {
			wrapped := fmt.Errorf("shutdown api server: %w", err)
			a.l.Error("Shutdown API server error", zap.Error(wrapped))
			shutdownErr = wrapped
		} else {
			a.l.Info("API server shutdown successfully")
		}
	}

	// db
	if a.DB != nil {
		if err := a.DB.Close(); err != nil {
			wrapped := fmt.Errorf("shutdown db: %w", err)
			a.l.Error("Shutdown DB error", zap.Error(wrapped))
			if shutdownErr == nil {
				shutdownErr = wrapped
			}
		} else {
			a.l.Info("DB shutdown successfully")
		}
	}
	return shutdownErr
}

func (a *App) setupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     a.cfg.CORS.AllowOrigins,
		AllowMethods:     a.cfg.CORS.AllowMethods,
		AllowHeaders:     a.cfg.CORS.AllowHeaders,
		ExposeHeaders:    a.cfg.CORS.ExposeHeaders,
		AllowCredentials: a.cfg.CORS.AllowCredentials,
		MaxAge:           time.Duration(a.cfg.CORS.MaxAge) * time.Second,
	}))
	api := router.Group("/api")
	{
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "pong"})
		})
	}
	return router
}
