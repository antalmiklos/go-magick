package gomagick

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/gographics/imagick.v3/imagick"
)

func Test() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()
	defer imagick.Terminate()

	outfile, err := os.Create("out.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer outfile.Close()
	c, err := NewConverter(outfile, DefaultOptions())
	if err != nil {
		log.Fatal(err)
	}
	b, _ := os.ReadFile("a.jpg")
	c.Read(b)
	c.ScalePercent(100)
	//	c.Destroy()
}
