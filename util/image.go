package util

import (
	"bytes"
	"image"
	_ "image/png" // png decoder

	"github.com/coreyog/rubikstimer/embedded"

	"github.com/faiface/pixel"
)

// LoadPicture parses an Asset to a pixel Picture
func LoadPicture(path string) (pixel.Picture, error) {
	file, err := embedded.Asset(path)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(bytes.NewReader(file))
	if err != nil {
		return nil, err
	}

	return pixel.PictureDataFromImage(img), nil
}
