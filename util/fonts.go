package util

import (
	"github.com/coreyog/rubikstimer/embedded"

	"github.com/faiface/pixel/text"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

// LoadTTF loads a TrueType font stored as an asset
func LoadTTF(path string, size float64) *text.Atlas {
	rawFont, err := embedded.Asset(path)
	if err != nil {
		panic(err)
	}

	ttfont, err := truetype.Parse(rawFont)
	if err != nil {
		panic(err)
	}

	fontface := truetype.NewFace(ttfont, &truetype.Options{
		Size:    size,
		Hinting: font.HintingFull,
	})

	return text.NewAtlas(fontface, text.ASCII)
}
