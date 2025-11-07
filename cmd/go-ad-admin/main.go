package main

import (
	"log"

	"github.com/Weruminger/go-ad-admin/internal/config"
	"github.com/Weruminger/go-ad-admin/internal/web"
)

func main() {
	cfg := config.FromEnv()
	log.Println("go-ad-admin starting on", cfg.ListenAddr)
	if err := web.ListenAndServe(cfg); err != nil {
		log.Fatal(err)
	}
}
