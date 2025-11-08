package config

import (
	"os"
	"testing"
	// "github.com/Weruminger/go-ad-admin/internal/config"
)

func TestFromEnvOverrides(t *testing.T) {
	var err error
	if err = os.Setenv("GO_AD_LISTEN", ":8080"); err != nil {
		t.Fatal(err)
	}
	if err = os.Setenv("GO_AD_ENV", "dev"); err != nil {
		t.Fatal(err)
	}
	if err = os.Setenv("GO_AD_SESSION_KEY", "abc123"); err != nil {
		t.Fatal(err)
	}
	defer os.Clearenv()

	c := new(Config).SetDefaultOnEmpty()
	if c.ListenAddr != ":8080" {
		t.Fatalf("got %s", c.ListenAddr)
	}
	if c.Env != "dev" {
		t.Fatalf("got %s", c.Env)
	}
	if c.SessionKey != "abc123" {
		t.Fatalf("got %s", c.SessionKey)
	}
}
