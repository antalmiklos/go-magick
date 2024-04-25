package gomagick

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/gographics/imagick.v3/imagick"
)

type (
	Scaler interface {
		// sets the image to supplied width and height
		ScaleXY(x, y int) error
		// sets the image width and height to scaled down values
		ScalePercent(p int) error
	}

	Rotator interface {
		Rotate(deg float64) error
	}

	Converter interface {
		// eg os.file
		io.Writer
		io.Closer
		io.Reader
		Scaler
		Rotator
		Destroy()
		Convert() error
		WithQuality(q int)
		WithCompression(c imagick.CompressionType)
		WithTargetFormat(format fileformat)
		Writer() io.WriteCloser
		// encodes converted image (jpg, png, tif) to writer
		Encode() error
	}

	ConverterOptions struct {
		Compression        imagick.CompressionType
		CompressionQuality uint
		TargetFormat       fileformat
	}

	imageConverter struct {
		wand   *imagick.MagickWand
		opts   ConverterOptions
		writer io.WriteCloser
	}
)

func (i *imageConverter) WithQuality(q int) {
	i.opts.CompressionQuality = uint(q)
}

func (i *imageConverter) WithCompression(c imagick.CompressionType) {
	i.opts.Compression = c
}

// sets output format and compression
func (i *imageConverter) WithTargetFormat(format fileformat) {
	i.opts.TargetFormat = format
	switch format {
	case FORMAT_JPG:
		i.opts.Compression = imagick.COMPRESSION_JPEG
	case FORMAT_PNG:
		i.opts.Compression = imagick.COMPRESSION_LZW
	case FORMAT_TIF:
		i.opts.Compression = imagick.COMPRESSION_LZW
	default:
		i.opts.Compression = imagick.COMPRESSION_NO
	}
}

func (i *imageConverter) Rotate(deg float64) error {
	pw := imagick.NewPixelWand()
	defer pw.Destroy()
	return i.wand.RotateImage(pw, float64(deg))
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
	i.writer.Close()
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

	if i.opts.TargetFormat == "" {
		i.opts.TargetFormat = FORMAT_JPG
	}
	if i.opts.CompressionQuality == 0 {
		i.opts.CompressionQuality = 100
	}
	if i.opts.Compression == 0 {
		i.opts.Compression = imagick.COMPRESSION_JPEG
	}

	if err := i.wand.SetFormat(string(i.opts.TargetFormat)); err != nil {
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

func (i *imageConverter) ScaleXY(x, y int) error {
	if x > int(i.wand.GetImageWidth()) {
		x = int(i.wand.GetImageWidth())
	}
	if y > int(i.wand.GetImageHeight()) {
		x = int(i.wand.GetImageHeight())
	}
	return i.wand.ScaleImage(uint(x), uint(y))
}

func (i *imageConverter) ScalePercent(p int) error {
	x := i.wand.GetImageWidth()
	y := i.wand.GetImageHeight()
	nx, ny := newXYPercent(int(x), int(y), float64(p))
	return i.wand.ScaleImage(uint(nx), uint(ny))
	/*
		 	dst := image.NewRGBA(image.Rect(0, 0, nx, ny))
			draw.NearestNeighbor.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)
			//	jpeg.Encode(i.writer, dst, &jpeg.Options{Quality: int(100)})
			return dst
		}
	*/
}

func (i *imageConverter) Writer() io.WriteCloser {
	return i.writer
}

func (i *imageConverter) Encode() error {
	blob := i.wand.GetImageBlob()
	_, err := i.Write(blob)
	return err
}

func newXYPercent(x, y int, percent float64) (int, int) {
	return int((percent / 100) * float64(x)), int((percent / 100) * float64(y))
}
