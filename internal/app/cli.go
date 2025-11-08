package app

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	. "github.com/Weruminger/go-ad-admin/internal/config"
	"github.com/Weruminger/go-ad-admin/internal/web"
	"github.com/spf13/pflag"
)

type InitResult struct {
	UsedConfig string
	LogPath    string
}

type App struct {
	Cfg     *Config
	Version string
	result  InitResult
}

func NewApp() *App {

	return &App{Cfg: NewDefaultConfig(), Version: ComputeVersion()}
}

// parsed Flags
type cliFlags struct {
	showHelp   bool
	configPath string
	listenAddr string
	logFile    string
	realm      string
	domainLAN  string
	domainDMZ  string
	workgroup  string

	env          string
	sessionKey   string
	LdapURL      string
	LdapBaseDN   string
	privacyLevel string
}

func parseFlags(args []string) (*cliFlags, error) {
	// pflag mit stdlib Flagset verbinden (damit Testbarkeit & default help)
	fs := pflag.NewFlagSet("go-ad-admin", pflag.ContinueOnError)
	fs.SetInterspersed(true)
	fs.SortFlags = true

	var f cliFlags
	fs.BoolVarP(&f.showHelp, "help", "h", false, "show help and exit")
	fs.StringVar(&f.configPath, "config", "", "YAML config file path")
	fs.StringVar(&f.listenAddr, "listen", "", "listen address (e.g. :8080)")
	fs.StringVar(&f.logFile, "log", "", "log file path")
	fs.StringVar(&f.realm, "realm", "", "AD realm (e.g. WERUMINGER.LAN)")
	fs.StringVar(&f.domainLAN, "domain-lan", "", "DNS zone LAN")
	fs.StringVar(&f.domainDMZ, "domain-dmz", "", "DNS zone DMZ")
	fs.StringVar(&f.workgroup, "workgroup", "", "Samba workgroup")

	fs.StringVar(&f.env, "env", "", "environment: dev, test, prod")
	fs.StringVar(&f.sessionKey, "session", "", "session key")
	fs.StringVar(&f.LdapURL, "ldap-url", "", "url to ldap server")
	fs.StringVar(&f.LdapBaseDN, "ldap-base-dn", "", "base dn for ldap server")
	fs.StringVar(&f.privacyLevel, "privacy", "", "privacy low or high")

	// pflag schluckt stdlib flags:
	fs.AddGoFlagSet(flag.CommandLine)
	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	return &f, nil
}

func usage() {
	_, _ = fmt.Fprintf(os.Stdout, `%s
MIT License – https://opensource.org/licenses/MIT
Repo: https://github.com/Weruminger/go-ad-admin

Usage:
  go-ad-admin [--config file.yaml] [--listen :8080] [--log logs/app.log] [--realm WERUMINGER.LAN] ...

Flags:
  -h, --help            show help and exit
      --config string   YAML config file path
      --listen string   listen address (default ":8080")
      --log string      log file path (default "logs/go-ad-admin.log")
      --realm string    AD realm
      --domain-lan str  DNS zone (LAN)
      --domain-dmz str  DNS zone (DMZ)
      --workgroup str   Samba workgroup
      --env string      environment: dev, test, prod
      --session string  session key
      --ldap-url string url to ldap server
      --ldap-base-dn    string base dn for ldap server
      --privacy string  privacy low or high
`, VersionBanner())
}

func (a *App) Initial() (bool, error) {
	var errA, errB error
	f, err := parseFlags(os.Args[1:])
	if err != nil {
		usage()
		return false, err
	}
	if f.showHelp {
		usage()
		return false, nil
	}

	// 1) Wenn keine Args ODER --config gesetzt → YAML einlesen (wenn vorhanden)

	noArgs := len(os.Args[1:]) == 0
	if f.configPath != "" || noArgs {
		path := f.configPath
		if path == "" {
			path = "config.yaml"
		}
		if _, errA = os.Stat(path); errA == nil {
			if errB = a.Cfg.LoadYAML(path); errB != nil {
				return false, fmt.Errorf("config invalid: %w", err)
			}
			a.result.UsedConfig = path
		}
	}

	// 2) Flags überschreiben ggf. YAML/Defaults
	if f.listenAddr != "" {
		a.Cfg.ListenAddr = f.listenAddr
	}
	if f.logFile != "" {
		a.Cfg.LogFile = f.logFile
	}
	if f.realm != "" {
		a.Cfg.Realm = f.realm
	}
	if f.domainLAN != "" {
		a.Cfg.DomainLAN = f.domainLAN
	}
	if f.domainDMZ != "" {
		a.Cfg.DomainDMZ = f.domainDMZ
	}
	if f.workgroup != "" {
		a.Cfg.Workgroup = f.workgroup
	}
	if f.env != "" {
		a.Cfg.Env = f.env
	}
	if f.sessionKey != "" {
		a.Cfg.SessionKey = f.sessionKey
	}
	if f.LdapURL != "" {
		a.Cfg.LDAPURL = f.LdapURL
	}
	if f.LdapBaseDN != "" {
		a.Cfg.LDAPBaseDN = f.LdapBaseDN
	}
	if f.privacyLevel != "" {
		a.Cfg.PrivacyLevel = f.privacyLevel
	}
	if f.configPath != "" {
		a.Cfg.ConfigFile = f.configPath
	}

	// 4) Config am Ende zurückschreiben (immer gültig)
	if err := a.Cfg.SaveYAML(a.Cfg.ConfigFileOrDefault()); err != nil {
		return false, fmt.Errorf("config save: %w", err)
	}

	a.result.LogPath = a.Cfg.LogFile
	return true, nil
}

func (a *App) Run() (bool, error) {
	// Hier startest du Services, Logger, HTTP, etc.
	// Für Demo nur OK zurück:
	log.Printf("go-ad-admin %s (commit=%s, build=%d) on %s", a.Version, a.Cfg.Env, epoch2010Seconds(), a.Cfg.ListenAddr)
	if err := web.ListenAndServe(*a.Cfg); err != nil {
		log.Fatalf("fatal: version=%s commit=%s build=%d err=%v", a.Version, a.Cfg.Env, epoch2010Seconds(), err)
	}

	return true, nil
}

func (a *App) Status() (string, error) {
	// Report für 1st Level: Version, Config, Log
	used := a.Cfg.ConfigFileOrDefault()
	return fmt.Sprintf("OK version=%s config=%s log=%s", a.Version, used, a.Cfg.LogFile), nil
}

func epoch2010Seconds() int64 {
	ref := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	now := time.Now().UTC()

	return int64(now.Sub(ref).Seconds())
}
