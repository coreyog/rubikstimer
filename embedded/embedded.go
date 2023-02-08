package embedded

import (
	"embed"
	"io"
)

//go:embed assets/*
var assets embed.FS

func Asset(path string) ([]byte, error) {
	f, err := assets.Open(path)
	if err != nil {
		return nil, err
	}

	return io.ReadAll(f)
}
