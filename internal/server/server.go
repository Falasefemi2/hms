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
	deptRepo := repository.NewDepartmentRepository(s.db.Pool())
	doctorRepo := repository.NewDoctorRepository(s.db.Pool())
	nurseRepo := repository.NewNurseRepository(s.db.Pool())
	patientRepo := repository.NewPatientRepository(s.db.Pool())
	availabilityRepo := repository.NewAvailabilityRepository(s.db.Pool())
	hospitalConfigRepo := repository.NewHospitalConfigRepository(s.db.Pool())
	appointmentRepo := repository.NewAppointmentRepository(s.db.Pool())
	consultationRepo := repository.NewConsultationRepository(s.db.Pool())

	userService := service.NewUserService(userRepo)
	deptService := service.NewDepartmentService(deptRepo)
	doctorService := service.NewDoctorService(doctorRepo, userRepo)
	nurseService := service.NewNurseService(nurseRepo, userRepo)
	patientService := service.NewPatientService(patientRepo, userRepo)
	availabilityService := service.NewAvailabilityService(availabilityRepo, doctorRepo)
	hospitalConfigService := service.NewHospitalConfigService(hospitalConfigRepo)
	appointmentService := service.NewAppointmentService(appointmentRepo, patientRepo, doctorRepo)
	consultationService := service.NewConsultationService(consultationRepo, appointmentRepo, patientRepo, doctorRepo)

	userHandler := handlers.NewUserHandler(userService)
	deptHandler := handlers.NewDeptHandler(deptService)
	doctorHandler := handlers.NewDoctorHandler(doctorService)
	nurseHandler := handlers.NewNurseHandler(nurseService)
	patientHandler := handlers.NewPatientHandlers(patientService)
	availabilityHandler := handlers.NewAvailabilityHandlers(availabilityService)
	hospitalConfigHandler := handlers.NewHospitalConfigHandler(hospitalConfigService)
	appointmentHandler := handlers.NewAppointmentHandler(appointmentService)
	consultationHandler := handlers.NewConsultationHandler(consultationService)

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
		r.Route("/departments", func(r chi.Router) {
			r.Post("/", deptHandler.CreateDepartment)
			r.Get("/", deptHandler.GetAllDepartments)
			r.Get("/{id}", deptHandler.GetDepartment)
			r.Put("/{id}", deptHandler.UpdateDepartment)
			r.Delete("/{id}", deptHandler.DeleteDepartment)
		})
		r.Route("/doctors", func(r chi.Router) {
			r.Post("/", doctorHandler.CreateDoctor)
			r.Route("/availability", func(r chi.Router) {
				r.Post("/", availabilityHandler.CreateAvailability)
			})
		})
		r.Route("/nurses", func(r chi.Router) {
			r.Post("/", nurseHandler.CreateNurse)
		})
		r.Route("/hospital-configs", func(r chi.Router) {
			r.Post("/", hospitalConfigHandler.CreateHospitalConfig)
			r.Get("/", hospitalConfigHandler.GetAllHospitalConfigs)
			r.Get("/{id}", hospitalConfigHandler.GetHospitalConfig)
			r.Put("/{id}", hospitalConfigHandler.UpdateHospitalConfig)
			r.Delete("/{id}", hospitalConfigHandler.DeleteHospitalConfig)
		})
	})

	r.Route("/patients", func(r chi.Router) {
		r.Use(middleware.JWTAuth)
		r.Use(middleware.PatientOnly)
		r.Route("/patientprofile", func(r chi.Router) {
			r.Post("/", patientHandler.PatientProfile)
		})
	})

	r.Route("/appointments", func(r chi.Router) {
		r.Use(middleware.JWTAuth)
		r.Post("/", appointmentHandler.CreateAppointment)
		r.Get("/{id}", appointmentHandler.GetAppointment)
		r.Put("/{id}", appointmentHandler.UpdateAppointment)
	})

	r.Route("/consultations", func(r chi.Router) {
		r.Use(middleware.JWTAuth)
		r.Post("/", consultationHandler.CreateConsultation)
		r.Get("/{id}", consultationHandler.GetConsultation)
		r.Put("/{id}", consultationHandler.UpdateConsultation)
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
