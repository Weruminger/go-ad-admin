package domain

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/Weruminger/go-ad-admin/internal/errs"
	"github.com/Weruminger/go-ad-admin/internal/modelx"
)

var reSam = regexp.MustCompile(`^[A-Za-z0-9._-]{1,64}$`)

type ADUser struct {
	*modelx.Base `json:"-" yaml:"-"`
	Kind         string         `json:"kind" yaml:"kind"`       // "ADUser"
	Version      string         `json:"version" yaml:"version"` // "v1"
	SAM          string         `json:"sam" yaml:"sam"`         // sAMAccountName
	UPN          string         `json:"upn" yaml:"upn"`         // user@realm
	Display      string         `json:"display,omitempty" yaml:"display,omitempty"`
	Mail         string         `json:"mail,omitempty" yaml:"mail,omitempty"`
	Enabled      bool           `json:"enabled" yaml:"enabled"`
	ExpiresAt    *time.Time     `json:"expiresAt,omitempty" yaml:"expiresAt,omitempty"`
	Meta         map[string]any `json:"meta,omitempty" yaml:"meta,omitempty"`
}

func NewADUser(b *modelx.Base) *ADUser {
	return &ADUser{Base: b, Kind: "ADUser", Version: "v1", Enabled: true, Meta: map[string]any{}}
}

func (u *ADUser) Init() *ADUser { return u }

func (u *ADUser) Validate() *ADUser {
	if u.Err() != nil {
		return u
	}
	op := errs.Op("aduser.Validate")
	if strings.TrimSpace(u.SAM) == "" || !reSam.MatchString(u.SAM) {
		u.SetInvalid(op, "sam", "must be 1..64 chars [A-Za-z0-9._-]")
	}
	if !strings.Contains(u.UPN, "@") {
		u.SetInvalid(op, "upn", "must contain @realm")
	}
	return u
}

// helper to set INVALID_INPUT with field info
func (u *ADUser) SetInvalid(op errs.Op, field, msg string) {
	u.BaseSetErr(op, errs.InvalidInput, fmt.Errorf("%s: %s", field, msg), map[string]any{"field": field})
}

func (u *ADUser) BaseSetErr(op errs.Op, code errs.Code, err error, fields map[string]any) {
	u.setErr(op, code, err, fields)
}

func (u *ADUser) Load(ctx context.Context, uri string) *ADUser {
	if u.Err() != nil {
		return u
	}
	op := errs.Op("aduser.Load")
	store, _, err := u.pickStore(uri)
	if err != nil {
		u.setErr(op, errs.InvalidInput, err, map[string]any{"uri": uri})
		return u
	}
	raw, err := store.Load(ctx, uri)
	if err != nil {
		u.setErr(op, errs.NotFound, err, map[string]any{"uri": uri})
		return u
	}
	cdc, err := u.pickCodec(modelxFormatFromURI(uri, u.Base))
	if err != nil {
		u.setErr(op, errs.InvalidInput, err, nil)
		return u
	}
	if err := cdc.Unmarshal(raw, u); err != nil {
		u.setErr(op, errs.InvalidInput, err, nil)
		return u
	}
	return u.Validate()
}

func (u *ADUser) Save(ctx context.Context, uri, format string) *ADUser {
	if u.Err() != nil {
		return u
	}
	op := errs.Op("aduser.Save")
	u = u.Validate()
	if u.Err() != nil {
		return u
	}
	cdc, err := u.pickCodec(format)
	if err != nil {
		u.setErr(op, errs.InvalidInput, err, map[string]any{"fmt": format})
		return u
	}
	raw, err := cdc.Marshal(u)
	if err != nil {
		u.setErr(op, errs.Internal, err, nil)
		return u
	}
	store, _, err := u.pickStore(uri)
	if err != nil {
		u.setErr(op, errs.InvalidInput, err, map[string]any{"uri": uri})
		return u
	}
	if err := store.Save(ctx, uri, raw); err != nil {
		u.setErr(op, errs.Unavailable, err, nil)
		return u
	}
	return u
}

func (u *ADUser) Serialize(format string) (string, error) {
	cdc, err := u.pickCodec(format)
	if err != nil {
		return "", errs.Wrap("aduser.Serialize", err, errs.InvalidInput)
	}
	b, err := cdc.Marshal(u)
	if err != nil {
		return "", errs.Wrap("aduser.Serialize", err, errs.Internal)
	}
	return string(b), nil
}

func (u *ADUser) Deserialize(format, data string) *ADUser {
	if u.Err() != nil {
		return u
	}
	cdc, err := u.pickCodec(format)
	if err != nil {
		u.setErr("aduser.Deserialize", errs.InvalidInput, err, nil)
		return u
	}
	if err := cdc.Unmarshal([]byte(data), u); err != nil {
		u.setErr("aduser.Deserialize", errs.InvalidInput, err, nil)
		return u
	}
	return u.Validate()
}
