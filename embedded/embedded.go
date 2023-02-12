package embedded

import (
	"embed"
	"io"
)

//go:embed assets/*
var assets embed.FS

// Asset bridges the gap between go-bindata and go:embed
func Asset(path string) ([]byte, error) {
	f, err := assets.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	return io.ReadAll(f)
}
