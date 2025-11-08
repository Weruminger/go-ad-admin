package app

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

// Diese Variablen werden via -ldflags injiziert.
// go build -ldflags "-X 'github.com/Weruminger/go-ad-admin/internal/app.VersionBase=1.2.3' -X 'github.com/Weruminger/go-ad-admin/internal/app.GitBranch=feature/x' -X 'github.com/Weruminger/go-ad-admin/internal/app.BuildEpoch=1731000000'"
var (
	VersionBase = "0.0.0" // aus version.txt
	GitBranch   = ""      // z.B. "main", "prod", "release/1.2.3", "feature/foo"
	BuildEpoch  = ""      // Epochtime als string
)

var reSemVer = regexp.MustCompile(`^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)$`)

func ComputeVersion() string {
	base := strings.TrimSpace(VersionBase)
	if !reSemVer.MatchString(base) {
		base = "0.0.0"
	}
	branch := strings.TrimSpace(GitBranch)
	epoch := strings.TrimSpace(BuildEpoch)

	// Branch-Policy:
	// - main / prod / release/* → exakt base
	// - sonst: base + .<epoch> + -rc_<branch_sanitized>
	if branch == "" || branch == "main" || branch == "prod" || strings.HasPrefix(branch, "release") {
		return base
	}

	// Conan2-kompatible Prerelease: nur [a-z0-9._], keine führenden Nullen problematisch, Kleinbuchstaben
	b := strings.ToLower(branch)
	b = strings.ReplaceAll(b, "/", "_")
	b = strings.ReplaceAll(b, "-", "_")
	b = strings.ReplaceAll(b, "+", "_")

	if epoch == "" {
		epoch = fmt.Sprintf("%d", time.Now().Unix())
	}

	return fmt.Sprintf("%s.%s-rc_%s", base, epoch, b)
}

func VersionBanner() string {
	return fmt.Sprintf("go-ad-admin %s (branch=%s build=%s)", ComputeVersion(), GitBranch, BuildEpoch)
}

// Hilfsfunktion: liest version.txt falls du nicht über ldflags injecten willst.
func ReadVersionTxt(path string) string {
	b, err := os.ReadFile(path)
	if err != nil {
		return VersionBase
	}
	v := strings.TrimSpace(string(b))
	if v == "" {
		return VersionBase
	}
	return v
}
