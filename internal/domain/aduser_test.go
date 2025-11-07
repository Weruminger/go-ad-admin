package domain

import (
	"context"
	"testing"
	"time"

	"github.com/Weruminger/go-ad-admin/internal/errs"
	"github.com/Weruminger/go-ad-admin/internal/modelx"
)

func baseJSON() *modelx.Base {
	return modelx.NewBase("json", []modelx.Codec{modelx.JSON{}}, []modelx.Store{modelx.FileStore{}})
}

func TestADUser_Validate_Errors(t *testing.T) {
	u := NewADUser(baseJSON())
	u.SAM = "!!!"  // invalid
	u.UPN = "user" // missing @
	u = u.Validate()
	if u.Err() == nil {
		t.Fatal("expected invalid input")
	}
	if !errs.IsCode(u.Err(), errs.InvalidInput) {
		t.Fatalf("want INVALID_INPUT got %v", u.Err())
	}
}

func TestADUser_Save_Load_Roundtrip(t *testing.T) {
	b := baseJSON()
	path := t.TempDir() + "/user.json"
	exp := NewADUser(b)
	exp.SAM = "rwerum"
	exp.UPN = "rwerum@WERUMINGER.LAN"
	ttl := time.Now().Add(24 * time.Hour).UTC().Truncate(time.Second)
	exp.ExpiresAt = &ttl

	exp.Save(context.Background(), "file://"+path, "json")
	if exp.Err() != nil {
		t.Fatalf("save: %v", exp.Err())
	}

	got := NewADUser(b).Load(context.Background(), "file://"+path)
	if got.Err() != nil {
		t.Fatalf("load: %v", got.Err())
	}
	if got.SAM != exp.SAM || got.UPN != exp.UPN {
		t.Fatalf("mismatch")
	}
}
