package config

import (
	"os"
	"testing"
)

func TestFromEnvDefaults(t *testing.T) {
	os.Clearenv()
	c := FromEnv()
	if c.ListenAddr != ":8080" {
		t.Fatalf("default listen expected :8080, got %s", c.ListenAddr)
	}
	if c.Env != "dev" {
		t.Fatalf("default env expected dev, got %s", c.Env)
	}
	if len(c.SessionKey) < 64 {
		t.Fatalf("session key too short")
	}
}

func TestFromEnvOverrides(t *testing.T) {
	os.Setenv("GOAD_LISTEN", ":9090")
	os.Setenv("GOAD_ENV", "prod")
	os.Setenv("GOAD_SESSION_KEY", "abc123")
	defer os.Clearenv()

	c := FromEnv()
	if c.ListenAddr != ":9090" {
		t.Fatalf("got %s", c.ListenAddr)
	}
	if c.Env != "prod" {
		t.Fatalf("got %s", c.Env)
	}
	if c.SessionKey != "abc123" {
		t.Fatalf("got %s", c.SessionKey)
	}
}
