package main

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"time"
)

const (
	JPEG ImageType = iota
	PNG
	GIF
	UNKNOWN
)

var (
	ImageDetectionTimeout = 1000 * time.Millisecond
)

type ImageType int64
func (i ImageType) Ext() string {
	if i == JPEG { return ".jpeg" }
	if i == PNG  { return ".png"  }
	if i == GIF  { return ".gif"  }
	return ""
}

func Detect(in io.Reader) ImageType {
	_, format, err := image.Decode(in)
	if err == nil {
		if format == "jpeg" {
			return JPEG
		}
		if format == "png"  {
			return PNG
		}
		if format == "gif"  {
			return GIF
		}
	}
	return UNKNOWN
}
