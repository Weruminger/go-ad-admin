package bdd

import (
	"fmt"
	"strings"

	"github.com/cucumber/godog"
)

/*** Helpers (werden von Steps aufgerufen) ***/

func seedUsers(m map[string]string) {
	w := testWorld()
	if w.users == nil {
		w.users = map[string]string{}
	}
	for uid, dn := range m {
		w.users[uid] = dn
	}
}

func seedUsersWithPrefix(n int, prefix, displayBase string) {
	w := testWorld()
	if w.users == nil {
		w.users = map[string]string{}
	}
	for i := 1; i <= n; i++ {
		uid := fmt.Sprintf("%s%d", prefix, i)
		dn := fmt.Sprintf("%s %d", displayBase, i)
		w.users[uid] = dn
	}
}

func setPrivacyHigh(high bool) {
	w := testWorld()
	w.privacyHigh = high
}

func clearSearchResults() {
	w := testWorld()
	w.lastSearchCount = 0
	// Wenn du später echte Ergebnislisten prüfst, könntest du hier z.B. w.lastResults = nil setzen.
}

/*** Zusatz-Steps ***/

// Given privacy mode is "high"|"low"
func privacyModeIs(mode string) error {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case "high":
		setPrivacyHigh(true)
	case "low":
		setPrivacyHigh(false)
	default:
		return fmt.Errorf("unknown privacy mode %q (use high|low)", mode)
	}
	return nil
}

// Given LDAP contains users:
//
//	| uid | displayName |
//
// (Alias zu anLDAPDirectoryWithUsers, aber „additiv“ statt überschreibend)
func ldapContainsUsers(tbl *godog.Table) error {
	if len(tbl.Rows) < 2 {
		return fmt.Errorf("table needs header + at least 1 row")
	}
	h := tbl.Rows[0].Cells
	if len(h) < 2 || strings.ToLower(h[0].Value) != "uid" || strings.ToLower(h[1].Value) != "displayname" {
		return fmt.Errorf("expected header: uid | displayName")
	}
	m := make(map[string]string, len(tbl.Rows)-1)
	for _, r := range tbl.Rows[1:] {
		if len(r.Cells) < 2 {
			return fmt.Errorf("row needs 2 cells")
		}
		uid := strings.TrimSpace(r.Cells[0].Value)
		dn := strings.TrimSpace(r.Cells[1].Value)
		if uid == "" || dn == "" {
			return fmt.Errorf("uid/displayName must not be empty")
		}
		m[uid] = dn
	}
	seedUsers(m)
	return nil
}

// Given LDAP contains 5 users with prefix "anna" and display base "Anna"
func ldapContainsNWithPrefixAndBase(n int, prefix, displayBase string) error {
	if n < 0 {
		return fmt.Errorf("n must be >= 0")
	}
	seedUsersWithPrefix(n, prefix, displayBase)
	return nil
}

// Then privacy mode should be "high"|"low"
func privacyModeShouldBe(mode string) error {
	wantHigh := strings.ToLower(strings.TrimSpace(mode)) == "high"
	w := testWorld()
	if w.privacyHigh != wantHigh {
		return fmt.Errorf("privacy=%v, want %v", w.privacyHigh, wantHigh)
	}
	return nil
}

/*** Wiring dieser Zusatz-Steps ***/

func init() {
	// nichts
}

func RegisterSeedSteps(sc *godog.ScenarioContext) {
	sc.Step(`^privacy mode is "([^"]*)"$`, privacyModeIs)
	sc.Step(`^privacy mode should be "([^"]*)"$`, privacyModeShouldBe)

	sc.Step(`^LDAP contains users:$`, ldapContainsUsers)
	sc.Step(`^LDAP contains (\d+) users with prefix "([^"]*)" and display base "([^"]*)"$`,
		ldapContainsNWithPrefixAndBase)
}
