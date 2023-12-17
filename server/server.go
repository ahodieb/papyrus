package server

import (
	"log"
	"net/http"

	"github.com/ahodieb/papyrus/notes"
	"github.com/ahodieb/papyrus/server/templates"
)

type Server struct {
	m *notes.Manager
}

func New(m *notes.Manager) *Server {
	return &Server{m: m}
}

func (s *Server) getSections(w http.ResponseWriter, r *http.Request) {
	sections := s.m.Sections()

	if err := Json(w, http.StatusOK, sections); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) UI(w http.ResponseWriter, r *http.Request) {
	if err := Html(w, http.StatusOK, templates.Index); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) ListenAndServe() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.UI)
	mux.HandleFunc("/api/sections", s.getSections)

	log.Println("Server running: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
