package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gdbu/scribe"
	generator "github.com/itsmontoya/shield-generator/lib"
)

var out *scribe.Scribe

func main() {
	var (
		output   string
		filename string

		key   string
		value string
		color string

		err error
	)

	out = scribe.New("Shield Generator")
	flag.StringVar(&output, "output", "image", "output type (options are image and url)")
	flag.StringVar(&filename, "filename", "./badge.png", "filename to save PNG")
	flag.StringVar(&key, "key", "coverage", "key of badge")
	flag.StringVar(&value, "value", "0%", "value of badge")
	flag.StringVar(&color, "color", "red", "color of value section")
	flag.Parse()

	var url string
	if url, err = generator.GenerateURL(key, value, color); err != nil {
		handleError("Error generating URL: %v", err)
	}

	switch output {
	case "url":
		fmt.Println(url)
	case "image":
		if err = generateImage(url, filename); err != nil {
			handleError("Error generating image: %v", err)
		}

	default:
		handleError("Invalid output, <%s> is not supported", output)
	}
}

func generateImage(url, filename string) (err error) {
	var f *os.File
	if f, err = os.Create(filename); err != nil {
		return
	}
	defer f.Close()
	return generator.GenerateImage(url, f)
}

func handleError(layout string, args ...interface{}) {
	out.Errorf(layout, args...)
	os.Exit(1)
}
