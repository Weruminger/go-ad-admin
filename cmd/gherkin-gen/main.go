package bdd

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/cucumber/godog"
)

// einfache Welt für Szenario-Status
type world struct {
	users   map[string]string
	lastErr error
}

func (w *world) reset(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
	w.users = map[string]string{}
	w.lastErr = nil
	return ctx, nil
}

// Step-Registrierung
func InitializeScenario(sc *godog.ScenarioContext) {
	w := &world{}
	sc.Before(w.reset)

	sc.Step(`^an empty directory$`, func(ctx context.Context) error {
		// hier ggf. tmp dir erzeugen; für das Beispiel reicht "ok"
		return nil
	})

	sc.Step(`^I create user "([^"]+)" with displayName "([^"]*)"$`,
		func(ctx context.Context, uid, dn string) error {
			if dn == "" {
				w.lastErr = fmt.Errorf("INVALID_INPUT")
				return nil
			}
			w.users[uid] = dn
			return nil
		})

	sc.Step(`^the user "([^"]+)" must exist$`,
		func(ctx context.Context, uid string) error {
			if _, ok := w.users[uid]; !ok {
				return fmt.Errorf("user %s not found", uid)
			}
			return nil
		})

	sc.Step(`^I receive an error code "([^"]+)"$`,
		func(ctx context.Context, code string) error {
			if w.lastErr == nil {
				return fmt.Errorf("no error raised")
			}
			if !strings.Contains(w.lastErr.Error(), code) {
				return fmt.Errorf("got %q, want contains %q", w.lastErr, code)
			}
			return nil
		})
}

// Godog in "go test" integrieren
func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		Name:                "features",
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format: "pretty",
			// Pfad relativ zu internal/bdd/
			Paths: []string{"../../features"},
		},
	}
	if st := suite.Run(); st != 0 {
		t.Fatalf("godog failed with status %d", st)
	}
}
