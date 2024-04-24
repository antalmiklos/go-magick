package examples

import (
	"fmt"
	"log"
	"os"
	"time"

	gomagick "github.com/antalmiklos/go-magick"
	"github.com/pkg/browser"
	"gopkg.in/gographics/imagick.v3/imagick"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("ImageMagick not properly installed, please consult the manual!", r)
		}
	}()
	defer imagick.Terminate()

	outfile, err := os.Create("out.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer outfile.Close()
	c, err := gomagick.NewConverter(outfile, gomagick.ConverterOptions{
		Compression:        imagick.COMPRESSION_NO,
		CompressionQuality: 100,
		TargetFormat:       gomagick.FORMAT_JPG,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer c.Destroy()
	b, _ := os.ReadFile("a.cr3")
	c.Read(b)
	fmt.Println("starting conversion")
	cstart := time.Now()
	if err := c.Convert(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("conversion done in: ", time.Since((cstart)))
	fmt.Println("starting rotation")
	cstart = time.Now()
	if err := c.Rotate(gomagick.CCW); err != nil {
		log.Fatal(err)
	}
	fmt.Println("rotation done in: ", time.Since((cstart)))
	fmt.Println("starting scaling")
	cstart = time.Now()
	c.ScalePercent(80)
	fmt.Println("scaling done in: ", time.Since((cstart)))
	browser.OpenFile(outfile.Name())
}
