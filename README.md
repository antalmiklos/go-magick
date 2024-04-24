# Example usage

``` go
func convertStuff() {
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
	c, err := gomagick.NewConverter(outfile, gomagick.DefaultOptions())
	if err != nil {
		log.Fatal(err)
	}
	b, _ := os.ReadFile("a.jpg")
	c.Read(b)
	c.ScalePercent(100)
	//	c.Destroy()
}
```