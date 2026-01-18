
package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/falasefemi2/hms/internal/database"
)

type Server struct {
	db *database.DB
}

func NewServer(db *database.DB) *Server {
	return &Server{db: db}
}

func (s *Server) Start(port string) error {
	router := http.NewServeMux()

	// Add your routes here
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome to the HMS API")
	})

	log.Printf("Server starting on port %s", port)
	return http.ListenAndServe(":"+port, router)
}
