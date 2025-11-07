package modelx

import (
	"context"
	"os"
	"path/filepath"
)

type FileStore struct{}

func (FileStore) Scheme() string { return "file" }

func (FileStore) Load(ctx context.Context, uri string) ([]byte, error) {
	path := uri
	if len(uri) >= 7 && uri[:7] == "file://" {
		path = uri[len("file://"):]
	}
	return os.ReadFile(filepath.Clean(path))
}

func (FileStore) Save(ctx context.Context, uri string, data []byte) error {
	path := uri
	if len(uri) >= 7 && uri[:7] == "file://" {
		path = uri[len("file://"):]
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Clean(path), data, 0o600)
}
