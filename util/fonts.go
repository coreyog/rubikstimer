package util

import (
	"github.com/coreyog/rubikstimer/embedded"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

// LoadTTF loads a TrueType font stored as an asset
func LoadTTF(path string, size float64) (font.Face, error) {
	rawFont, err := embedded.Asset(path)
	if err != nil {
		return nil, err
	}
	ttfont, err := truetype.Parse(rawFont)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(ttfont, &truetype.Options{
		Size:    size,
		Hinting: font.HintingFull,
	}), nil
}
