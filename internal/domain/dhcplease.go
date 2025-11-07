package domain

import (
	"context"
	"fmt"
	"net"
	"regexp"
	"time"

	"github.com/Weruminger/go-ad-admin/internal/errs"
	"github.com/Weruminger/go-ad-admin/internal/modelx"
)

var reHost = regexp.MustCompile(`^[A-Za-z0-9]([A-Za-z0-9-]{0,61}[A-Za-z0-9])?$`)

type DHCPLease struct {
	*modelx.Base `json:"-" yaml:"-"`
	Kind         string    `json:"kind" yaml:"kind"`       // "DHCPLease"
	Version      string    `json:"version" yaml:"version"` // "v1"
	MAC          string    `json:"mac" yaml:"mac"`
	IP           string    `json:"ip" yaml:"ip"`
	Host         string    `json:"host" yaml:"host"`
	Start        time.Time `json:"start" yaml:"start"`
	End          time.Time `json:"end" yaml:"end"`
}

func NewDHCPLease(b *modelx.Base) *DHCPLease {
	return &DHCPLease{Base: b, Kind: "DHCPLease", Version: "v1"}
}

func (d *DHCPLease) Init() *DHCPLease { return d }

func (d *DHCPLease) Validate() *DHCPLease {
	if d.Err() != nil {
		return d
	}
	op := errs.Op("dhcplease.Validate")
	if _, err := net.ParseMAC(d.MAC); err != nil {
		d.setErr(op, errs.InvalidInput, fmt.Errorf("mac invalid: %w", err), nil)
	}
	ip := net.ParseIP(d.IP)
	if ip == nil || ip.To4() == nil {
		d.setErr(op, errs.InvalidInput, fmt.Errorf("ip must be IPv4"), nil)
	}
	if !reHost.MatchString(d.Host) {
		d.setErr(op, errs.InvalidInput, fmt.Errorf("host RFC-952/1123 invalid"), nil)
	}
	if !d.Start.Before(d.End) {
		d.setErr(op, errs.InvalidInput, fmt.Errorf("start must be before end"), nil)
	}
	return d
}

func (d *DHCPLease) Load(ctx context.Context, uri string) *DHCPLease {
	if d.Err() != nil {
		return d
	}
	op := errs.Op("dhcplease.Load")
	store, _, err := d.pickStore(uri)
	if err != nil {
		d.setErr(op, errs.InvalidInput, err, nil)
		return d
	}
	raw, err := store.Load(ctx, uri)
	if err != nil {
		d.setErr(op, errs.NotFound, err, nil)
		return d
	}
	cdc, err := d.pickCodec(modelxFormatFromURI(uri, d.Base))
	if err != nil {
		d.setErr(op, errs.InvalidInput, err, nil)
		return d
	}
	if err := cdc.Unmarshal(raw, d); err != nil {
		d.setErr(op, errs.InvalidInput, err, nil)
		return d
	}
	return d.Validate()
}

func (d *DHCPLease) Save(ctx context.Context, uri, format string) *DHCPLease {
	if d.Err() != nil {
		return d
	}
	op := errs.Op("dhcplease.Save")
	d = d.Validate()
	if d.Err() != nil {
		return d
	}
	cdc, err := d.pickCodec(format)
	if err != nil {
		d.setErr(op, errs.InvalidInput, err, nil)
		return d
	}
	raw, err := cdc.Marshal(d)
	if err != nil {
		d.setErr(op, errs.Internal, err, nil)
		return d
	}
	store, _, err := d.pickStore(uri)
	if err != nil {
		d.setErr(op, errs.InvalidInput, err, nil)
		return d
	}
	if err := store.Save(ctx, uri, raw); err != nil {
		d.setErr(op, errs.Unavailable, err, nil)
		return d
	}
	return d
}

func (d *DHCPLease) Serialize(format string) (string, error) {
	cdc, err := d.pickCodec(format)
	if err != nil {
		return "", errs.Wrap("dhcplease.Serialize", err, errs.InvalidInput)
	}
	b, err := cdc.Marshal(d)
	if err != nil {
		return "", errs.Wrap("dhcplease.Serialize", err, errs.Internal)
	}
	return string(b), nil
}

func (d *DHCPLease) Deserialize(format, data string) *DHCPLease {
	if d.Err() != nil {
		return d
	}
	cdc, err := d.pickCodec(format)
	if err != nil {
		d.setErr("dhcplease.Deserialize", errs.InvalidInput, err, nil)
		return d
	}
	if err := cdc.Unmarshal([]byte(data), d); err != nil {
		d.setErr("dhcplease.Deserialize", errs.InvalidInput, err, nil)
		return d
	}
	return d.Validate()
}
