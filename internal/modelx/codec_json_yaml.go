package modelx

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

type JSON struct{}

func (JSON) Format() string                  { return "json" }
func (JSON) Marshal(v any) ([]byte, error)   { return json.MarshalIndent(v, "", "  ") }
func (JSON) Unmarshal(b []byte, v any) error { return json.Unmarshal(b, v) }

type YAML struct{}

func (YAML) Format() string                  { return "yaml" }
func (YAML) Marshal(v any) ([]byte, error)   { return yaml.Marshal(v) }
func (YAML) Unmarshal(b []byte, v any) error { return yaml.Unmarshal(b, v) }
