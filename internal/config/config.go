package config

import (
	"crypto/rand"
	"encoding/hex"
	"os"
)

type Config struct {
	ListenAddr   string
	Env          string
	SessionKey   string
	LDAPURL      string
	LDAPBaseDN   string
	PrivacyLevel string // low|high
}

func FromEnv() Config {
	return Config{
		ListenAddr:   getenv("GOAD_LISTEN", ":8080"),
		Env:          getenv("GOAD_ENV", "dev"),
		SessionKey:   getenv("GOAD_SESSION_KEY", randKey(32)),
		LDAPURL:      getenv("GOAD_LDAP_URL", "ldap://127.0.0.1:389"),
		LDAPBaseDN:   getenv("GOAD_LDAP_BASEDN", "dc=example,dc=com"),
		PrivacyLevel: getenv("GOAD_PRIVACY", "low"),
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func randKey(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
