package config

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	// Laufzeit-Parameter (Defaults werden in NewDefaultConfig gesetzt)
	ListenAddr   string `yaml:"listenAddr,omitempty"`
	Env          string `yaml:"env,omitempty"`
	SessionKey   string `yaml:"sessionKey,omitempty"`
	LDAPURL      string `yaml:"ldapURL,omitempty"`
	LDAPBaseDN   string `yaml:"ldapBaseDN,omitempty"`
	PrivacyLevel string `yaml:"privacyLevel,omitempty"` // low|high
	LogFile      string `yaml:"logFile,omitempty"`
	ConfigFile   string `yaml:"-"` // Pfad, aus dem geladen wurde (keine YAML-Ausgabe)

	// Beispiel-AD/DHCP Settings
	Realm     string `yaml:"realm,omitempty"`
	DomainLAN string `yaml:"domainLAN,omitempty"`
	DomainDMZ string `yaml:"domainDMZ,omitempty"`
	Workgroup string `yaml:"workgroup,omitempty"`
}

// Defaults setzen – immer gültige Konfiguration erzeugen
func NewDefaultConfig() *Config {
	return new(Config).SetDefaultOnEmpty()
}

func defaultIfEmpty(s, def string) string {
	if s == "" {
		return def
	}
	return s
}
func (c *Config) ConfigFileOrDefault() string {
	if c.ConfigFile != "" {
		return c.ConfigFile
	}
	return "config.yaml"
}

func (c *Config) SetDefaultOnEmpty() *Config {
	c.ListenAddr = defaultIfEmpty(c.ListenAddr, getenv("GO_AD_LISTEN", ":8080"))
	c.LogFile = defaultIfEmpty(c.LogFile, "logs/go-ad-admin.log")
	c.Realm = defaultIfEmpty(c.Realm, "WERUMINGER.LAN")
	c.DomainLAN = defaultIfEmpty(c.DomainLAN, "weruminger.lan")
	c.DomainDMZ = defaultIfEmpty(c.DomainDMZ, "weruminger.dmz")
	c.Workgroup = defaultIfEmpty(c.Workgroup, "WERUMINGER")
	c.Env = defaultIfEmpty(c.Env, getenv("GO_AD_ENV", "dev"))
	c.SessionKey = defaultIfEmpty(c.SessionKey, getenv("GO_AD_SESSION_KEY", randKey(32)))
	c.LDAPURL = defaultIfEmpty(c.LDAPURL, getenv("GO_AD_LDAP_URL", "ldap://127.0.0.1:389"))
	c.LDAPBaseDN = defaultIfEmpty(c.LDAPBaseDN, getenv("GO_AD_LDAP_BASEDN", "dc=weruminger, dc=eu"))
	c.PrivacyLevel = defaultIfEmpty(c.PrivacyLevel, getenv("GO_AD_PRIVACY", "low"))
	return c
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func (c *Config) Validate() error {
	if c.ListenAddr == "" {
		return errors.New("listenAddr must not be empty")
	}
	if c.Realm == "" {
		return errors.New("realm must not be empty")
	}
	return nil
}

func (c *Config) LoadYAML(path string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	if err := yaml.Unmarshal(b, c); err != nil {
		return fmt.Errorf("parse yaml: %w", err)
	}
	c.ConfigFile = path
	return c.Validate()
}

func (c *Config) SaveYAML(path string) error {
	if err := c.Validate(); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	b, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, b, 0o644); err != nil {
		return err
	}
	return nil
}
func randKey(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
