package main

import (
	"log"

	"github.com/Weruminger/go-ad-admin/internal/config"
	"github.com/Weruminger/go-ad-admin/internal/web"
)

var (
	version   = "dev"
	commit    = "none"
	buildDate = "unknown"
)

func main() {
	cfg := *(new(config.Config))
	log.Printf("go-ad-admin %s (commit=%s, build=%s) on %s", version, commit, buildDate, cfg.ListenAddr)
	if err := web.ListenAndServe(cfg); err != nil {
		log.Fatalf("fatal: version=%s commit=%s build=%s err=%v", version, commit, buildDate, err)
	}
}
