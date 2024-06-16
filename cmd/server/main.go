package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	esv7 "github.com/elastic/go-elasticsearch/v7"
	chi "github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	elasticsearch "github.com/colmmurphy91/muzz/internal/adapter/elasticsearch"
	matchStore "github.com/colmmurphy91/muzz/internal/adapter/mysql/match"
	swipeStore "github.com/colmmurphy91/muzz/internal/adapter/mysql/swipe"
	userStore "github.com/colmmurphy91/muzz/internal/adapter/mysql/user"
	"github.com/colmmurphy91/muzz/internal/api/discover"
	authhttp "github.com/colmmurphy91/muzz/internal/api/login"
	swipeHttp "github.com/colmmurphy91/muzz/internal/api/swipe"
	"github.com/colmmurphy91/muzz/internal/api/user"
	"github.com/colmmurphy91/muzz/internal/usecase/auth"
	discoverService "github.com/colmmurphy91/muzz/internal/usecase/discover"
	swipeService "github.com/colmmurphy91/muzz/internal/usecase/swipe"
	userM "github.com/colmmurphy91/muzz/internal/usecase/user"
	tools "github.com/colmmurphy91/muzz/tools"
	"github.com/colmmurphy91/muzz/tools/envvar"
)

type serverConfig struct {
	Address     string
	DB          *sqlx.DB
	ES          *esv7.Client
	MiddleWares []func(next http.Handler) http.Handler
	Logger      *zap.SugaredLogger
	PrivateKey  string
	PublicKey   string
}

func main() {
	var env, address string

	fmt.Println("i am here")

	flag.StringVar(&env, "env", "env.example", "Environment Variables filename")
	flag.StringVar(&address, "address", ":8080", "HTTP Server Address")
	flag.Parse()

	errC, err := run(env, address)
	if err != nil {
		log.Fatalf("Could not run %s", err)
	}

	if err := <-errC; err != nil {
		log.Fatalf("error while running %s", err)
	}
}

func LoggingMiddleware(log *zap.SugaredLogger) func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			log.Info(request.Method,
				zap.Time("time", time.Now()),
				zap.String("url", request.URL.String()))
			next.ServeHTTP(writer, request)
		})
	}
}

func run(env, _ string) (<-chan error, error) {
	logger, err := tools.New("event-thor-service")
	if err != nil {
		return nil, fmt.Errorf("zap.NewProduction %w", err)
	}

	logger.Infof("Loading environment variables from: %s", env)

	if err := envvar.Load(env); err != nil {
		return nil, fmt.Errorf("envar.Load %w", err)
	}

	fmt.Println(env)

	conf := envvar.New()

	var db *sqlx.DB

	db, err = tools.NewDBConnection(conf)
	if err != nil {
		return nil, fmt.Errorf("failed to create db connection: %w", err)
	}

	es, err := tools.NewElasticSearch(conf)
	if err != nil {
		return nil, fmt.Errorf("failed to create es connection: %w", err)
	}

	errC := make(chan error, 1)

	port := conf.Get("PORT")
	addr := fmt.Sprintf(":%s", port)

	srv := newServer(serverConfig{
		Address:     addr,
		DB:          db,
		ES:          es,
		MiddleWares: []func(next http.Handler) http.Handler{LoggingMiddleware(logger)},
		Logger:      logger,
	})

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-ctx.Done()

		logger.Info("Shutdown signal received")

		ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		defer func() {
			loggerErr := logger.Sync()
			if loggerErr != nil {
				logger.Info("Shutdown error")
			}

			db.Close()
			stop()
			cancel()
			close(errC)
		}()

		srv.SetKeepAlivesEnabled(false)

		if err := srv.Shutdown(ctxTimeout); err != nil {
			errC <- err
		}

		logger.Info("Shutdown completed")
	}()

	go func() {
		logger.Info("Listening and service", zap.String("address", addr))

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errC <- err
		}
	}()

	return errC, nil
}

func newServer(conf serverConfig) *http.Server {
	r := chi.NewRouter()

	for _, mw := range conf.MiddleWares {
		r.Use(mw)
	}

	var (
		store       = userStore.NewStore(conf.Logger, conf.DB)
		swipeStorer = swipeStore.NewStore(conf.Logger, conf.DB)
		matchStorer = matchStore.NewStore(conf.Logger, conf.DB)
		index       = elasticsearch.NewUser(conf.ES)
	)

	swipeS := swipeService.NewService(swipeStorer, matchStorer)

	discoverS := discoverService.NewService(index, swipeStorer)

	userManager := userM.NewManager(store, index)
	authService := auth.NewAuthService("my-secret-key", store)

	tools.InitAuth("my-secret-key")

	authhttp.NewHandler(conf.Logger, authService).Register(r)

	user.NewHandler(conf.Logger, userManager).Register(r)

	r.Group(func(r chi.Router) {
		r.Use(tools.AuthMiddleware)
		discover.NewHandler(conf.Logger, discoverS).Register(r)
	})

	r.Group(func(r chi.Router) {
		r.Use(tools.AuthMiddleware)
		swipeHttp.NewHandler(conf.Logger, swipeS).Register(r)
	})

	return &http.Server{
		Handler:           r,
		Addr:              conf.Address,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       10 * time.Second,
	}
}
