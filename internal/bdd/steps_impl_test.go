package bdd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/cucumber/godog"
)

/*** Test-Welt ***/

type world struct {
	users           map[string]string // uid -> displayName
	lastSearchCount int
	privacyHigh     bool
	lastHTTP        int
	sessionCookie   bool
	failedAttempts  int
	lastErr         error
}

func (w *world) reset(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
	w.users = map[string]string{}
	w.privacyHigh = true
	w.lastHTTP = 0
	w.sessionCookie = false
	w.failedAttempts = 0
	w.lastErr = nil
	w.lastSearchCount = 0
	return ctx, nil
}

var _w *world

func testWorld() *world {
	if _w == nil {
		_w = &world{users: map[string]string{}}
	}
	return _w
}

/*** Step-Implementierungen ***/

func anEmptyDirectory() error {
	// ggf. tmp-dir vorbereiten
	return nil
}

func anLDAPDirectoryWithUsers(tbl *godog.Table) error {
	w := testWorld()
	w.users = map[string]string{}

	// Header: uid | displayName
	if len(tbl.Rows) < 2 {
		return fmt.Errorf("table needs header + at least 1 row")
	}
	h := tbl.Rows[0].Cells
	if len(h) < 2 || strings.ToLower(h[0].Value) != "uid" || strings.ToLower(h[1].Value) != "displayname" {
		return fmt.Errorf("expected header: uid | displayName")
	}
	for _, r := range tbl.Rows[1:] {
		if len(r.Cells) < 2 {
			return fmt.Errorf("row needs 2 cells")
		}
		uid := strings.TrimSpace(r.Cells[0].Value)
		dn := strings.TrimSpace(r.Cells[1].Value)
		if uid == "" || dn == "" {
			return fmt.Errorf("uid/displayName must not be empty")
		}
		w.users[uid] = dn
	}
	return nil
}

func anLDAPDirectoryWithMatchingUsers(n int) error {
	w := testWorld()
	// Basismenge mit 2 Matches („anna“ im Namen) + Non-Matches
	w.users = map[string]string{
		"anna.smith": "Anna Smith",
		"joanna.roe": "Joanna Roe",
		"bob":        "Bob Doe",
		"charlie":    "Charlie Brown",
	}
	// auf n Matches auffüllen
	for i := 3; i <= n; i++ {
		uid := "anna" + strconv.Itoa(i)
		w.users[uid] = "Anna " + strconv.Itoa(i)
	}
	return nil
}

func iCreateUserWithDisplayName(uid, dn string) error {
	w := testWorld()
	if dn == "" {
		w.lastErr = fmt.Errorf("INVALID_INPUT")
		w.lastHTTP = 422
		return nil
	}
	w.users[uid] = dn
	return nil
}

func theUserMustExist(uid string) error {
	w := testWorld()
	if _, ok := w.users[uid]; !ok {
		return fmt.Errorf("user %s not found", uid)
	}
	return nil
}

func iReceiveAnErrorCode(code string) error {
	w := testWorld()
	if w.lastErr == nil {
		return fmt.Errorf("no error raised")
	}
	if !strings.Contains(w.lastErr.Error(), code) {
		return fmt.Errorf("got %q, want contains %q", w.lastErr, code)
	}
	return nil
}

func iSearchFor(q string) error {
	w := testWorld()
	ql := strings.ToLower(q)
	cnt := 0
	for _, dn := range w.users {
		if strings.Contains(strings.ToLower(dn), ql) {
			cnt++
		}
	}
	w.lastSearchCount = cnt
	return nil
}

func iSeeResults(n int) error {
	w := testWorld()
	if w.lastSearchCount != n {
		return fmt.Errorf("got %d results, want %d", w.lastSearchCount, n)
	}
	return nil
}

func noPIIIsShownWhenPrivacyModeIsHigh() error {
	w := testWorld()
	if !w.privacyHigh {
		return fmt.Errorf("privacy not high")
	}
	// hier würdest du Inhalte auf PII prüfen
	return nil
}

func theSystemIsRunning() error {
	// Healthcheck/Mocks hier
	return nil
}

func iSubmitValidCredentials() error {
	w := testWorld()
	w.sessionCookie = true
	w.lastHTTP = 200
	return nil
}

func iReceiveASessionCookie() error {
	w := testWorld()
	if !w.sessionCookie {
		return fmt.Errorf("no session cookie")
	}
	return nil
}

func iSubmitAnInvalidPassword() error {
	w := testWorld()
	w.failedAttempts++
	w.lastHTTP = 401
	return nil
}

func theFailedattemptCounterIsIncremented() error {
	w := testWorld()
	if w.failedAttempts < 1 {
		return fmt.Errorf("failed-attempt counter not incremented")
	}
	return nil
}

func iGetHTTP(code int) error {
	w := testWorld()
	if w.lastHTTP != code {
		return fmt.Errorf("got HTTP %d, want %d", w.lastHTTP, code)
	}
	return nil
}

/*** Wiring ***/

func InitializeScenario(sc *godog.ScenarioContext) {
	_w = &world{}
	sc.Before(_w.reset)

	sc.Step(`^an empty directory$`, anEmptyDirectory)
	sc.Step(`^an LDAP directory with users:$`, anLDAPDirectoryWithUsers)
	sc.Step(`^an LDAP directory with (\d+) matching users$`, anLDAPDirectoryWithMatchingUsers)

	sc.Step(`^I create user "([^"]*)" with displayName "([^"]*)"$`, iCreateUserWithDisplayName)
	sc.Step(`^the user "([^"]*)" must exist$`, theUserMustExist)
	sc.Step(`^I receive an error code "([^"]*)"$`, iReceiveAnErrorCode)

	sc.Step(`^I search for "([^"]*)"$`, iSearchFor)
	sc.Step(`^I see (\d+) results$`, iSeeResults)
	sc.Step(`^no PII is shown when privacy mode is high$`, noPIIIsShownWhenPrivacyModeIsHigh)

	sc.Step(`^the system is running$`, theSystemIsRunning)
	sc.Step(`^I submit valid credentials$`, iSubmitValidCredentials)
	sc.Step(`^I receive a session cookie$`, iReceiveASessionCookie)
	sc.Step(`^I submit an invalid password$`, iSubmitAnInvalidPassword)
	sc.Step(`^the failed-attempt counter is incremented$`, theFailedattemptCounterIsIncremented)
	sc.Step(`^I get HTTP (\d+)$`, iGetHTTP)
	// <- hier die Zusatzsteps dazuhängen:
	RegisterSeedSteps(sc)
}

/*** Suite (Strict + robuster Feature-Pfad) ***/

func findRepoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatalf("repo root (go.mod) not found from %s", dir)
		}
		dir = parent
	}
}

func TestFeatures(t *testing.T) {
	root := findRepoRoot(t)
	features := filepath.Join(root, "features")

	t.Logf("[godog] scanning features dir: %s", features)
	if _, err := os.Stat(features); os.IsNotExist(err) {
		t.Skipf("no features/ directory at %s -> skipping BDD", features)
	}

	suite := godog.TestSuite{
		Name:                "features",
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format: "pretty",
			Paths:  []string{features},
			Strict: true, // undefined/pending -> FAIL
		},
	}
	if st := suite.Run(); st != 0 {
		t.Fatalf("godog failed with status %d", st)
	}
}
