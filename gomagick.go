package gomagick

import (
	"io"

	"gopkg.in/gographics/imagick.v3/imagick"
)

type fileformat string

const FORMAT_JPG fileformat = "jpg"
const FORMAT_PNG fileformat = "png"
const FORMAT_TIF fileformat = "tiff"

const CW float64 = 90
const CCW float64 = 270

/*
Defaults to JPEG format and compression with 100 quality
*/
func DefaultOptions() ConverterOptions {
	o := ConverterOptions{
		TargetFormat:       FORMAT_JPG,
		Compression:        imagick.COMPRESSION_JPEG,
		CompressionQuality: 100,
	}
	return o
}

func NewConverter(output io.WriteCloser, opts ConverterOptions) (Converter, error) {
	imagick.Initialize()
	/* defer func() {
		if r := recover(); r != nil {
			log.Fatal("could not initialise MagickWand")
		}
	}() */
	wand := imagick.NewMagickWand()
	c := &imageConverter{
		wand:   wand,
		writer: output,
		opts:   opts,
	}

	return c, nil
}
