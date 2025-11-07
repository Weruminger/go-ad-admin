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

	c := FromEnv()
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
