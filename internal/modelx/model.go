package modelx

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"sync"

	"github.com/Weruminger/go-ad-admin/internal/errs"
)

type Codec interface {
	Marshal(v any) ([]byte, error)
	Unmarshal(b []byte, v any) error
	Format() string
}

type Store interface {
	Load(ctx context.Context, uri string) ([]byte, error)
	Save(ctx context.Context, uri string, data []byte) error
	Scheme() string
}

type Base struct {
	mu      sync.RWMutex
	lastErr error
	codecs  map[string]Codec
	stores  map[string]Store
	format  string
}

func NewBase(defaultFormat string, codecs []Codec, stores []Store) *Base {
	cm := make(map[string]Codec)
	for _, c := range codecs {
		cm[strings.ToLower(c.Format())] = c
	}
	sm := make(map[string]Store)
	for _, s := range stores {
		sm[strings.ToLower(s.Scheme())] = s
	}
	return &Base{codecs: cm, stores: sm, format: strings.ToLower(defaultFormat)}
}

func (b *Base) Err() error {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.lastErr
}

func (b *Base) setErr(op errs.Op, code errs.Code, err error, fields map[string]any) {
	if err == nil {
		return
	}
	w := errs.New(op, code, err, fields)
	b.mu.Lock()
	b.lastErr = w
	b.mu.Unlock()
}

func (b *Base) pickCodec(format string) (Codec, error) {
	f := strings.ToLower(strings.TrimSpace(format))
	if f == "" {
		f = b.format
	}
	if f == "" {
		return nil, fmt.Errorf("no format specified")
	}
	c, ok := b.codecs[f]
	if !ok {
		return nil, fmt.Errorf("unknown format %q", f)
	}
	return c, nil
}

func (b *Base) pickStore(u string) (Store, *url.URL, error) {
	pu, err := url.Parse(u)
	if err != nil {
		return nil, nil, err
	}
	s, ok := b.stores[strings.ToLower(pu.Scheme)]
	if !ok {
		return nil, nil, fmt.Errorf("no store for scheme %q", pu.Scheme)
	}
	return s, pu, nil
}
