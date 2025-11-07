package web

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/Weruminger/go-ad-admin/internal/config"
	"github.com/Weruminger/go-ad-admin/internal/errs"
)

type Server struct {
	cfg config.Config
	tpl *template.Template
}

func NewServer(cfg config.Config) *Server {
	const glob = "web/templates/*.html"
	t, err := template.ParseGlob(glob)
	if err != nil || t == nil {
		t = template.Must(template.New("layout").Parse(`{{define "layout"}}<!doctype html><html><head><meta charset="utf-8"><title>go-ad-admin</title></head><body>
<main><h1>go-ad-admin</h1><p>Env: {{.Env}}</p></main>
</body></html>{{end}}`))
	} else {
		matches, _ := filepath.Glob(glob)
		if len(matches) == 0 {
			t = template.Must(template.New("layout").Parse(`{{define "layout"}}<!doctype html><html><head><meta charset="utf-8"><title>go-ad-admin</title></head><body>
<main><h1>go-ad-admin</h1><p>Env: {{.Env}}</p></main>
</body></html>{{end}}`))
		}
	}
	return &Server{cfg: cfg, tpl: t}
}

func (s *Server) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleIndex)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(204)
	})
	return withReqID(mux)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if len(q) > 256 {
		writeError(w, r, errs.New("web.Index", errs.InvalidInput, fmt.Errorf("q>256"), map[string]any{"len": len(q)}))
		return
	}
	_ = s.tpl.ExecuteTemplate(w, "layout", map[string]any{"Env": s.cfg.Env})
}

func ListenAndServe(cfg config.Config) error {
	return http.ListenAndServe(cfg.ListenAddr, NewServer(cfg).routes())
}
