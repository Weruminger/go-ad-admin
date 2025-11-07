package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/Weruminger/go-ad-admin/internal/config"
	"github.com/Weruminger/go-ad-admin/internal/errs"
	"github.com/Weruminger/go-ad-admin/internal/testx"
)

func TestWriteError_Mapping(t *testing.T) {
	type tc struct {
		err      error
		wantCode int
		wantJSON string
	}
	cases := []tc{
		{errs.New("web.Search", errs.InvalidInput, fmt.Errorf("bad q"), nil), http.StatusUnprocessableEntity, `"code":"INVALID_INPUT"`},
		{errs.New("ldap.Search", errs.Timeout, fmt.Errorf("ctx"), nil), http.StatusServiceUnavailable, `"code":"TIMEOUT"`},
		{errs.New("ldap.Get", errs.NotFound, fmt.Errorf("dn"), nil), http.StatusNotFound, `"code":"NOT_FOUND"`},
		{fmt.Errorf("raw"), http.StatusInternalServerError, `"code":"INTERNAL"`},
	}
	_ = NewServer(config.FromEnv()) // just ensure it builds

	for _, c := range cases {
		rec := testx.NewRecorder()
		req := testx.NewRequest("GET", "/trigger", nil)
		// inject req id
		req.Header.Set("X-Request-ID", "t-req-id")
		writeError(rec, req, c.err)

		if rec.Code != c.wantCode {
			t.Fatalf("status: got %d want %d", rec.Code, c.wantCode)
		}
		if !json.Valid(rec.Body.Bytes()) || !strings.Contains(rec.BodyString(), c.wantJSON) {
			t.Fatalf("payload: %s missing %s", rec.BodyString(), c.wantJSON)
		}
		if h := rec.Header().Get("Content-Type"); !strings.HasPrefix(h, "application/json") {
			t.Fatalf("content-type not json: %q", h)
		}
		if rid := rec.Header().Get("X-Request-ID"); rid != "t-req-id" {
			t.Fatalf("missing/propagated request id header, got %q", rid)
		}
	}
}
