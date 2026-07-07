package app

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"

	"marketplace/backend/internal/auth"
	"marketplace/backend/internal/config"
	"marketplace/backend/internal/platform/id"
	"marketplace/backend/internal/products"
)

type Server struct {
	HTTP *http.Server
	DB   *pgxpool.Pool
}

func NewServer(ctx context.Context, cfg config.Config, logger *slog.Logger) (*Server, error) {
	db, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(ctx); err != nil {
		db.Close()
		return nil, err
	}

	now := func() time.Time { return time.Now().UTC() }
	jwtManager := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTExpiration, now)
	authService := auth.NewService(auth.NewPostgresRepository(db), jwtManager, id.NewUUID, now)
	productService := products.NewService(products.NewPostgresRepository(db), id.NewUUID, now)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(requestLogger(logger))
	router.Use(cors(cfg.FrontendOrigin))
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	authHandler := auth.NewHandler(authService)
	productHandler := products.NewHandler(productService)

	router.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/register", authHandler.Register)
		r.Post("/auth/login", authHandler.Login)
		r.Get("/products", productHandler.PublicList)
		r.Get("/products/{productID}", productHandler.PublicGet)
		r.Group(func(protected chi.Router) {
			protected.Use(auth.RequireSeller(jwtManager))
			protected.Post("/seller/products", productHandler.Create)
			protected.Get("/seller/products", productHandler.List)
			protected.Get("/seller/products/{productID}", productHandler.Get)
			protected.Post("/seller/products/{productID}/publish", productHandler.Publish)
		})
	})

	return &Server{
		HTTP: &http.Server{
			Addr:              ":" + cfg.AppPort,
			Handler:           router,
			ReadHeaderTimeout: 5 * time.Second,
		},
		DB: db,
	}, nil
}

func (s *Server) Close() {
	if s.DB != nil {
		s.DB.Close()
	}
}

func cors(frontendOrigin string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", frontendOrigin)
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func requestLogger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			logger.Info("http request",
				"method", r.Method,
				"path", r.URL.Path,
				"duration_ms", time.Since(start).Milliseconds(),
			)
		})
	}
}
