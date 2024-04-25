# Example usage

``` go
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	gomagick "github.com/antalmiklos/go-magick"
	"gopkg.in/gographics/imagick.v3/imagick"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("ImageMagick not properly installed, please consult the manual!", r)
		}
	}()
	defer func() {
		fmt.Println("terminating")
		imagick.Terminate()
		fmt.Println("terminated")
	}()

	outfile, err := os.Create("out.png")
	if err != nil {
		log.Fatal(err)
	}
	defer outfile.Close()
	c, err := gomagick.NewConverter(outfile, gomagick.ConverterOptions{})
	if err != nil {
		log.Fatal(err)
	}

	c.WithTargetFormat(gomagick.FORMAT_PNG)
	c.WithQuality(30)

	defer c.Destroy()
	b, _ := os.ReadFile("a.cr3")
	c.Read(b)
	fmt.Println("starting conversion")
	cstart := time.Now()
	if err := c.Convert(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("conversion done in: ", time.Since((cstart)))
	fmt.Println("starting scaling")
	cstart = time.Now()
	if err := c.ScalePercent(50); err != nil {
		log.Fatal(err)
	}
	fmt.Println("scaling done in: ", time.Since((cstart)))
	fmt.Println("starting rotation")
	cstart = time.Now()
	if err := c.Rotate(gomagick.CCW); err != nil {
		log.Fatal(err)
	}
	fmt.Println("rotation done in: ", time.Since((cstart)))
	c.Encode()
	stats, err := os.Stat(outfile.Name())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("file created %s\nsize: %dMb\n", stats.Name(), stats.Size()/1024/1024)
}
```
