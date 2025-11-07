package domain

import (
	"context"
	"strings"

	"github.com/Weruminger/go-ad-admin/internal/errs"
	"github.com/Weruminger/go-ad-admin/internal/modelx"
)

type FeatureSpec struct {
	*modelx.Base `json:"-" yaml:"-"`
	Kind         string            `json:"kind" yaml:"kind"`
	Version      string            `json:"version" yaml:"version"`
	Meta         map[string]string `json:"meta,omitempty" yaml:"meta,omitempty"`
	Data         map[string]any    `json:"data,omitempty" yaml:"data,omitempty"`
}

func NewFeatureSpec(b *modelx.Base) *FeatureSpec {
	return &FeatureSpec{Base: b, Kind: "FeatureSpec", Version: "v1", Meta: map[string]string{}, Data: map[string]any{}}
}

func (f *FeatureSpec) Init() *FeatureSpec { return f }

func (f *FeatureSpec) Load(ctx context.Context, uri string) *FeatureSpec {
	if f.Err() != nil {
		return f
	}
	store, _, err := f.Base.PickStore(uri)
	if err != nil {
		f.Base.SetErr("feature.Load", errs.InvalidInput, err, map[string]any{"uri": uri})
		return f
	}
	raw, err := store.Load(ctx, uri)
	if err != nil {
		f.Base.SetErr("feature.Load", errs.NotFound, err, map[string]any{"uri": uri})
		return f
	}
	format := modelxFormatFromURI(uri, f.Base)
	cdc, err := f.Base.PickCodec(format)
	if err != nil {
		f.Base.SetErr("feature.Load", errs.InvalidInput, err, map[string]any{"fmt": format})
		return f
	}
	if err := cdc.Unmarshal(raw, f); err != nil {
		f.Base.SetErr("feature.Load", errs.InvalidInput, err, map[string]any{"fmt": cdc.Format()})
		return f
	}
	return f
}

func (f *FeatureSpec) Save(ctx context.Context, uri string, format string) *FeatureSpec {
	if f.Err() != nil {
		return f
	}
	cdc, err := f.Base.PickCodec(format)
	if err != nil {
		f.Base.SetErr("feature.Save", errs.InvalidInput, err, map[string]any{"fmt": format})
		return f
	}
	raw, err := cdc.Marshal(f)
	if err != nil {
		f.Base.SetErr("feature.Save", errs.Internal, err, map[string]any{"fmt": cdc.Format()})
		return f
	}
	store, _, err := f.Base.PickStore(uri)
	if err != nil {
		f.Base.SetErr("feature.Save", errs.InvalidInput, err, map[string]any{"uri": uri})
		return f
	}
	if err := store.Save(ctx, uri, raw); err != nil {
		f.Base.SetErr("feature.Save", errs.Unavailable, err, map[string]any{"uri": uri})
		return f
	}
	return f
}

func (f *FeatureSpec) Serialize(format string) (string, error) {
	cdc, err := f.Base.PickCodec(format)
	if err != nil {
		return "", errs.Wrap("feature.Serialize", err, errs.InvalidInput)
	}
	b, err := cdc.Marshal(f)
	if err != nil {
		return "", errs.Wrap("feature.Serialize", err, errs.Internal)
	}
	return string(b), nil
}

func (f *FeatureSpec) Deserialize(format, data string) *FeatureSpec {
	if f.Err() != nil {
		return f
	}
	cdc, err := f.Base.PickCodec(format)
	if err != nil {
		f.Base.SetErr("feature.Deserialize", errs.InvalidInput, err, map[string]any{"fmt": format})
		return f
	}
	if err := cdc.Unmarshal([]byte(data), f); err != nil {
		f.Base.SetErr("feature.Deserialize", errs.InvalidInput, err, map[string]any{"fmt": cdc.Format()})
	}
	return f
}

func modelxFormatFromURI(uri string, _ *modelx.Base) string {
	low := strings.ToLower(uri)
	switch {
	case strings.HasSuffix(low, ".yaml"), strings.HasSuffix(low, ".yml"):
		return "yaml"
	case strings.HasSuffix(low, ".json"):
		return "json"
	default:
		return "json"
	}
}
