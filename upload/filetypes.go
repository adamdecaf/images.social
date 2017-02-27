package upload

import (
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"strings"
)

// ImageType represents the underlying content type from an image. It's intended to
// be used for generating the extension via Ext().
type ImageType string

// Ext generates a string representation of the full file extension based on the
// content type that was created with ImageType.
func (i ImageType) Ext() string {
	return fmt.Sprintf(".%s", strings.ToLower(string(i)))
}

// Detect attempts to determine the content type of the provided input Reader.
// It is a blocking call (be warned if the Reader is slow).
func Detect(in io.Reader) (ImageType, error) {
	_, format, err := image.Decode(in)
	if err == nil {
		if format == "jpeg" {
			return ImageType("jpeg"), nil
		}
		if format == "png" {
			return ImageType("png"), nil
		}
		if format == "gif" {
			return ImageType("gif"), nil
		}
	}
	return ImageType("unknown"), errors.New("Unknown image type, detected")
}
