package domain

import (
	"context"
	"testing"

	"github.com/Weruminger/go-ad-admin/internal/errs"
	"github.com/Weruminger/go-ad-admin/internal/modelx"
)

func TestFeature_Load_NotFound(t *testing.T) {
	base := modelx.NewBase("yaml", []modelx.Codec{modelx.YAML{}}, []modelx.Store{modelx.FileStore{}})
	f := NewFeatureSpec(base).Load(context.Background(), "file://./nope/does-not-exist.yaml")
	if f.Err() == nil {
		t.Fatal("expected error")
	}
	if !errs.IsCode(f.Err(), errs.NotFound) {
		t.Fatalf("want NOT_FOUND, got %v", f.Err())
	}
}

func TestFeature_Save_And_Load_RoundTrip(t *testing.T) {
	base := modelx.NewBase("json", []modelx.Codec{modelx.JSON{}}, []modelx.Store{modelx.FileStore{}})
	tmp := t.TempDir() + "/feature.json"
	fw := NewFeatureSpec(base)
	fw.Meta["module"] = "auth"
	fw = fw.Save(context.Background(), "file://"+tmp, "json")
	if fw.Err() != nil {
		t.Fatalf("save err: %v", fw.Err())
	}

	fr := NewFeatureSpec(base).Load(context.Background(), "file://"+tmp)
	if fr.Err() != nil {
		t.Fatalf("load err: %v", fr.Err())
	}
	if fr.Meta["module"] != "auth" {
		t.Fatalf("roundtrip mismatch")
	}
}

func TestFeature_Deserialize_Invalid(t *testing.T) {
	base := modelx.NewBase("yaml", []modelx.Codec{modelx.YAML{}}, []modelx.Store{modelx.FileStore{}})
	f := NewFeatureSpec(base).Deserialize("yaml", ":\n- broken")
	if f.Err() == nil {
		t.Fatal("expected invalid input error")
	}
	if !errs.IsCode(f.Err(), errs.InvalidInput) {
		t.Fatalf("want INVALID_INPUT, got %v", f.Err())
	}
}
