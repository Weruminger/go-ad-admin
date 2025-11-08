package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	. "github.com/Weruminger/go-ad-admin/internal/config"
)

func TestParseFlags_Help(t *testing.T) {
	f, err := parseFlags([]string{"-h"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.showHelp {
		t.Fatal("expected showHelp to be true")
	}
}

func TestParseFlags_AllFlags(t *testing.T) {
	args := []string{
		"--config", "test.yaml",
		"--listen", ":9090",
		"--log", "test.log",
		"--realm", "TEST.LAN",
		"--domain-lan", "lan.test",
		"--domain-dmz", "dmz.test",
		"--workgroup", "TESTGROUP",
		"--env", "test",
		"--session", "secretkey",
		"--ldap-url", "ldap://localhost",
		"--ldap-base-dn", "dc=test,dc=lan",
		"--privacy", "high",
	}

	f, err := parseFlags(args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if f.configPath != "test.yaml" {
		t.Errorf("configPath: got %q, want %q", f.configPath, "test.yaml")
	}
	if f.listenAddr != ":9090" {
		t.Errorf("listenAddr: got %q, want %q", f.listenAddr, ":9090")
	}
	if f.logFile != "test.log" {
		t.Errorf("logFile: got %q, want %q", f.logFile, "test.log")
	}
	if f.realm != "TEST.LAN" {
		t.Errorf("realm: got %q, want %q", f.realm, "TEST.LAN")
	}
	if f.domainLAN != "lan.test" {
		t.Errorf("domainLAN: got %q, want %q", f.domainLAN, "lan.test")
	}
	if f.domainDMZ != "dmz.test" {
		t.Errorf("domainDMZ: got %q, want %q", f.domainDMZ, "dmz.test")
	}
	if f.workgroup != "TESTGROUP" {
		t.Errorf("workgroup: got %q, want %q", f.workgroup, "TESTGROUP")
	}
	if f.env != "test" {
		t.Errorf("env: got %q, want %q", f.env, "test")
	}
	if f.sessionKey != "secretkey" {
		t.Errorf("sessionKey: got %q, want %q", f.sessionKey, "secretkey")
	}
	if f.LdapURL != "ldap://localhost" {
		t.Errorf("LdapURL: got %q, want %q", f.LdapURL, "ldap://localhost")
	}
	if f.LdapBaseDN != "dc=test,dc=lan" {
		t.Errorf("LdapBaseDN: got %q, want %q", f.LdapBaseDN, "dc=test,dc=lan")
	}
	if f.privacyLevel != "high" {
		t.Errorf("privacyLevel: got %q, want %q", f.privacyLevel, "high")
	}
}

func TestParseFlags_InvalidFlag(t *testing.T) {
	_, err := parseFlags([]string{"--invalid-flag"})
	if err == nil {
		t.Fatal("expected error for invalid flag")
	}
}

func TestNewApp(t *testing.T) {
	app := NewApp()
	if app == nil {
		t.Fatal("NewApp returned nil")
	}
	if app.Cfg == nil {
		t.Fatal("expected Config to be initialized")
	}
	if app.Version == "" {
		t.Error("expected Version to be set")
	}
}

func TestApp_Initial_WithFlags(t *testing.T) {
	// Temporäres Verzeichnis für Config-Datei
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test-config.yaml")

	// os.Args überschreiben für Test
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"test", "--config", configPath, "--env", "test", "--listen", ":9999"}
	testFileContent := `listenAddr: :8080
env: dev
sessionKey: cf63708453807afb4f689bbc6f2fca208e5824eee5c7bb91f9ff867e5463ba5b
ldapURL: ldap://127.0.0.1:389
ldapBaseDN: dc=weruminger, dc=eu
privacyLevel: low
logFile: logs/go-ad-admin.log
realm: WERUMINGER.LAN
domainLAN: weruminger.lan
domainDMZ: weruminger.dmz
workgroup: WERUMINGER
`
	if err := os.MkdirAll("config", 0o755); err != nil {
		fmt.Printf("mkdir error: %v\n", err)
		return
	}

	// YAML-String in Datei schreiben
	if err := os.WriteFile(configPath, []byte(testFileContent), 0o644); err != nil {
		fmt.Printf("write error: %v\n", err)
		return
	}

	app := NewApp()
	ok, err := app.Initial()
	if err != nil {
		t.Fatalf("Initial() error: %v", err)
	}
	if !ok {
		t.Fatal("Initial() returned false")
	}

	if app.Cfg.Env != "test" {
		t.Errorf("Env: got %q, want %q", app.Cfg.Env, "test")
	}
	if app.Cfg.ListenAddr != ":9999" {
		t.Errorf("ListenAddr: got %q, want %q", app.Cfg.ListenAddr, ":9999")
	}
	if app.result.UsedConfig != configPath {
		t.Errorf("UsedConfig: got %q, want %q", app.result.UsedConfig, configPath)
	}
}

func TestApp_Initial_ShowHelp(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"test", "-h"}

	app := NewApp()
	ok, err := app.Initial()
	if err != nil {
		t.Fatalf("Initial() error: %v", err)
	}
	if ok {
		t.Fatal("Initial() should return false when help is shown")
	}
}

func TestApp_Initial_ParseError(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"test", "--invalid-flag"}

	app := NewApp()
	ok, err := app.Initial()
	if err == nil {
		t.Fatal("expected error for invalid flag")
	}
	if ok {
		t.Fatal("Initial() should return false on parse error")
	}
}

func TestApp_Status(t *testing.T) {
	app := NewApp()
	app.Version = "test-version"
	app.Cfg.LogFile = "test.log"

	status, err := app.Status()
	if err != nil {
		t.Fatalf("Status() error: %v", err)
	}

	if !strings.Contains(status, "test-version") {
		t.Errorf("status should contain version, got: %s", status)
	}
	if !strings.Contains(status, "test.log") {
		t.Errorf("status should contain log file, got: %s", status)
	}
	if !strings.Contains(status, "OK") {
		t.Errorf("status should contain OK, got: %s", status)
	}
}

func TestEpoch2010Seconds(t *testing.T) {
	seconds := epoch2010Seconds()
	if seconds <= 0 {
		t.Errorf("epoch2010Seconds should return positive value, got: %d", seconds)
	}

	// Sollte eine große Zahl sein (Jahre seit 2000)
	const minExpected = 788918400 // ca. 25 Jahre in Sekunden
	if seconds < minExpected {
		t.Errorf("epoch2010Seconds seems too small: %d", seconds)
	}
}

func TestUsage(t *testing.T) {
	// Teste, dass usage() nicht panicked
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("usage() panicked: %v", r)
		}
	}()

	// Redirect stdout um Output zu testen
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() { os.Stdout = oldStdout }()

	usage()

	w.Close()
	os.Stdout = oldStdout

	// Lese gepipten Output
	buf := make([]byte, 1024)
	n, _ := r.Read(buf)
	output := string(buf[:n])

	if !strings.Contains(output, "go-ad-admin") {
		t.Error("usage output should contain 'go-ad-admin'")
	}
	if !strings.Contains(output, "--help") {
		t.Error("usage output should contain '--help'")
	}
}

func TestApp_Initial_FlagOverridesYAML(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "override-test.yaml")

	// Erstelle eine YAML-Datei mit bestimmten Werten
	cfg := NewDefaultConfig()
	cfg.Env = "yaml-env"
	cfg.ListenAddr = ":8888"
	if err := cfg.SaveYAML(configPath); err != nil {
		t.Fatalf("failed to save test config: %v", err)
	}

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	// Flag sollte YAML überschreiben
	os.Args = []string{"test", "--config", configPath, "--env", "flag-env"}

	app := NewApp()
	ok, err := app.Initial()
	if err != nil {
		t.Fatalf("Initial() error: %v", err)
	}
	if !ok {
		t.Fatal("Initial() returned false")
	}

	// Flag sollte YAML-Wert überschrieben haben
	if app.Cfg.Env != "flag-env" {
		t.Errorf("Flag should override YAML: got %q, want %q", app.Cfg.Env, "flag-env")
	}
	// Nicht überschriebener Wert sollte aus YAML kommen
	if app.Cfg.ListenAddr != ":8888" {
		t.Errorf("Non-overridden value should come from YAML: got %q, want %q", app.Cfg.ListenAddr, ":8888")
	}
}
