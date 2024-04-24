package gomagick

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"os"

	"golang.org/x/image/draw"
	"gopkg.in/gographics/imagick.v3/imagick"
)

const FORMAT_JPG = "jpg"
const FORMAT_PNG = "png"
const FORMAT_TIF = "tiff"

type Converter interface {
	// eg os.file
	io.Writer
	io.Closer
	io.Reader
	Scaler
	Destroy()
	//	Wand() *imagick.MagickWand
	Convert() error
}

type imageConverter struct {
	wand   *imagick.MagickWand
	opts   ConverterOptions
	writer io.WriteCloser
}

type ConverterOptions struct {
	Compression        imagick.CompressionType
	CompressionQuality uint
	TargetFormat       string
	//	CompressionFilter  imagick.FilterType
}

/*
Defaults to JPEG format and compression with 100 quality
*/
func DefaultOptions() ConverterOptions {
	o := ConverterOptions{
		TargetFormat:       FORMAT_JPG,
		Compression:        imagick.COMPRESSION_JPEG,
		CompressionQuality: 100,
		//	CompressionFilter:  imagick.FILTER_LANCZOS,
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

func (i *imageConverter) Write(p []byte) (n int, err error) {
	return i.writer.Write(p)
}

func (i *imageConverter) Close() error {
	return i.writer.Close()
}

// Reads p into wand as imageblob
func (i *imageConverter) Read(p []byte) (n int, err error) {
	return len(p), i.wand.ReadImageBlob(p)
}

func (i *imageConverter) ImageBlob() ([]byte, error) {
	blob := i.wand.GetImageBlob()
	if len(blob) == 0 {
		return nil, fmt.Errorf("empty image in wand")
	}
	return blob, nil
}

// Deallocates all memory associated with the converter and destroys the MagicWand
func (i *imageConverter) Destroy() {
	imgNumber := i.wand.GetNumberImages()
	for j := 0; j < int(imgNumber); j++ {
		img := i.wand.GetImageFromMagickWand()
		if img == nil {
			break
		}
		i.wand.DestroyImage(img)
	}
	i.wand.Destroy()
}

func (j *imageConverter) WriteToFile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	j.wand.ResetIterator()
	blob := j.wand.GetImageBlob()
	if _, err := f.Write(blob); err != nil {
		return err
	}
	return nil
}

func (i *imageConverter) Convert() error {
	if err := i.wand.SetFormat(i.opts.TargetFormat); err != nil {
		return err
	}
	if err := i.wand.SetImageCompression(i.opts.Compression); err != nil {
		return err
	}
	if err := i.wand.SetImageCompressionQuality(i.opts.CompressionQuality); err != nil {
		return err
	}
	return nil
}

func (i imageConverter) ScaleXY(x, y int) image.Image {
	blob, _ := i.ImageBlob()
	src, _ := jpeg.Decode(bytes.NewReader(blob))
	dst := image.NewRGBA(image.Rect(0, 0, src.Bounds().Max.X/2, src.Bounds().Max.Y/2))
	draw.NearestNeighbor.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)
	jpeg.Encode(i.writer, dst, &jpeg.Options{Quality: int(100)})
	return dst
}

func (i imageConverter) ScalePercent(p int) image.Image {
	blob, _ := i.ImageBlob()
	src, _ := jpeg.Decode(bytes.NewReader(blob))
	nx, ny := newXYPercent(src.Bounds().Max.X, src.Bounds().Max.Y, float64(p))
	dst := image.NewRGBA(image.Rect(0, 0, nx, ny))
	draw.NearestNeighbor.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)
	jpeg.Encode(i.writer, dst, &jpeg.Options{Quality: int(100)})
	return dst
}

func newXYPercent(x, y int, percent float64) (int, int) {
	nx := int((percent / 100) * float64(x))
	ny := int((percent / 100) * float64(y))
	fmt.Println(nx, ny)
	return nx, ny
}

type Scaler interface {
	ScaleXY(x, y int) image.Image
	ScalePercent(p int) image.Image
}
