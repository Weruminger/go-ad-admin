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
		ListenAddr:   getenv("GO_AD_LISTEN", ":8080"),
		Env:          getenv("GO_AD_ENV", "dev"),
		SessionKey:   getenv("GO_AD_SESSION_KEY", randKey(32)),
		LDAPURL:      getenv("GO_AD_LDAP_URL", "ldap://127.0.0.1:389"),
		LDAPBaseDN:   getenv("GO_AD_LDAP_BASEDN", "dc=example,dc=com"),
		PrivacyLevel: getenv("GO_AD_PRIVACY", "low"),
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
