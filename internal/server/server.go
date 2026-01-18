package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/falasefemi2/hms/docs"
	"github.com/falasefemi2/hms/internal/database"
	"github.com/falasefemi2/hms/internal/handlers"
	"github.com/falasefemi2/hms/internal/middleware"
	"github.com/falasefemi2/hms/internal/repository"
	"github.com/falasefemi2/hms/internal/service"
)

type Server struct {
	db     *database.DB
	server *http.Server
}

func NewServer(db *database.DB) *Server {
	return &Server{db: db}
}

func (s *Server) Start(port string) error {
	r := chi.NewRouter()

	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	userRepo := repository.NewUserRepository(s.db.Pool())
	userService := service.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome to the HMS API")
	})
	r.Get("/health", handlers.HealthCheck)
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/signup", userHandler.SignUpPatient)
		r.Post("/login", userHandler.Login)
	})

	r.Route("/admin", func(r chi.Router) {
		r.Use(middleware.JWTAuth)
		r.Use(middleware.AdminOnly)

		r.Route("/users", func(r chi.Router) {
			r.Post("/", userHandler.CreateUser)
			r.Get("/", userHandler.ListUsers)
			r.Get("/{id}", userHandler.GetUser)
		})
	})

	s.server = &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	log.Printf("Server starting on port %s", port)
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) Shutdown(timeout time.Duration) error {
	if s.server == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}

	s.db.Close()
	log.Println("Server stopped gracefully")
	return nil
}
