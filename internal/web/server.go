package web

import (
	"html/template"
	"log"
	"net/http"

	"github.com/Weruminger/go-ad-admin/internal/config"
)

type Server struct {
	cfg config.Config
	tpl *template.Template
}

func NewServer(cfg config.Config) *Server {
	tpl := template.Must(template.ParseFS(templates, "templates/*.html"))
	return &Server{cfg: cfg, tpl: tpl}
}

func (s *Server) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleIndex)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200); _, _ = w.Write([]byte("ok")) })
	return mux
}

func ListenAndServe(cfg config.Config) error {
	s := NewServer(cfg)
	log.Printf("listening on %s (env=%s)", cfg.ListenAddr, cfg.Env)
	return http.ListenAndServe(cfg.ListenAddr, s.routes())
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{"Env": s.cfg.Env}
	_ = s.tpl.ExecuteTemplate(w, "index.html", data)
}
