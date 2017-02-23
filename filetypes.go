package main

import (
	"fmt"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"time"
)

// filetypes
// - match on header info
// - verify images are valid

type ImageType int64
func (i ImageType) Ext() string {
	if i == JPEG { return ".jpeg" }
	if i == PNG  { return ".png"  }
	if i == GIF  { return ".gif"  }
	return ""
}

const (
	JPEG ImageType = iota
	PNG
	GIF
	UNKNOWN
)

var (
	ImageDetectionTimeout = 1000 * time.Millisecond
)

func Detect(in io.Reader) ImageType {
	out := make(chan ImageType)

	// JPEG
	go func() {
		_, err := jpeg.Decode(in)
		if err == nil {
			out <- JPEG
		}
		fmt.Println(err)
	}()

	// PNG
	go func() {
		_, err := png.Decode(in)
		if err == nil {
			out <- PNG
		}
		fmt.Println(err)
	}()

	// GIF
	go func() {
		_, err := gif.Decode(in)
		if err == nil {
			out <- GIF
		}
		fmt.Println(err)
	}()

	select {
	case tpe := <-out:
		return tpe
	case <-time.After(ImageDetectionTimeout):
		return UNKNOWN
	}
}
