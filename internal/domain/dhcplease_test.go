package domain

import (
	"context"
	"testing"
	"time"

	"github.com/Weruminger/go-ad-admin/internal/errs"
	"github.com/Weruminger/go-ad-admin/internal/modelx"
)

func baseYAML() *modelx.Base {
	return modelx.NewBase("yaml", []modelx.Codec{modelx.YAML{}}, []modelx.Store{modelx.FileStore{}})
}

func TestDHCPLease_Validate_Error(t *testing.T) {
	d := NewDHCPLease(baseYAML())
	d.MAC = "zz:zz:zz:zz:zz:zz"
	d.IP = "not-ip"
	d.Host = "bad host"
	d.Start = time.Now()
	d.End = d.Start.Add(-1 * time.Hour)
	d = d.Validate()
	if d.Err() == nil || !errs.IsCode(d.Err(), errs.InvalidInput) {
		t.Fatalf("expected INVALID_INPUT, got %v", d.Err())
	}
}

func TestDHCPLease_Save_Load(t *testing.T) {
	b := baseYAML()
	path := t.TempDir() + "/lease.yaml"
	start := time.Now().UTC().Truncate(time.Second)
	end := start.Add(2 * time.Hour)

	exp := NewDHCPLease(b)
	exp.MAC = "bc:24:11:9d:ca:fa"
	exp.IP = "10.0.10.6"
	exp.Host = "dc1"
	exp.Start = start
	exp.End = end

	exp.Save(context.Background(), "file://"+path, "yaml")
	if exp.Err() != nil {
		t.Fatalf("save err: %v", exp.Err())
	}

	got := NewDHCPLease(b).Load(context.Background(), "file://"+path)
	if got.Err() != nil {
		t.Fatalf("load err: %v", got.Err())
	}
	if got.MAC != exp.MAC || got.IP != exp.IP || got.Host != exp.Host {
		t.Fatalf("mismatch")
	}
}
