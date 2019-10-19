package main

import (
	"fmt"
	"image"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
)

func main() {

	args := os.Args[1:]

	if len(args) < 5 || args[0] == "help" {
		getHelp()
		return
	}

	outName := args[0]
	backgroundImage := args[1]
	markImage := args[2]

	// Resize the mark to fit these dimenstions.
	markWidth, markHeight := parseCoordinates(args[3], "x")

	// Coordinate to super-impose on. e.g. 200x500
	locationX, locationY := parseCoordinates(args[4], "x")

	src, err := imaging.Open(backgroundImage)
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}

	mark, err := imaging.Open(markImage)
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}

	markFit := imaging.Fit(mark, markWidth, markHeight, imaging.Lanczos)

	dst := imaging.Paste(src, markFit, image.Pt(locationX, locationY))

	err = imaging.Save(dst, "data/"+outName)
	if err != nil {
		log.Fatalf("failed to save image: %v", err)
	}

	fmt.Println("Done!")
}

func parseCoordinates(input, delimiter string) (int, int) {

	array := strings.Split(input, delimiter)

	x, err := strconv.Atoi(array[0])

	if err != nil {
		log.Fatalf("failed to parse coordinates: %v", err)
	}

	y, err := strconv.Atoi(array[1])

	if err != nil {
		log.Fatalf("failed to parse coordinates: %v", err)
	}

	return x, y
}

func getHelp() {
	fmt.Print(
		`This is a simple program that imposes an image over another. 
Usage:
	go run main.go <output-name> <background-name> <watermark-name> <watermark-dimensions> <location>
	e.g. go run main.go zebra.png sample1.png mark.png 200x200 100x100
Arguments:
	- output-name: What should the output image be called?
	- background-name: Name of background image.
	- watermark-name: Name of smaller image.
	- watermark-dimensions: Watermark would be resized to this before imposing.
	- location: x & y coordinates to place the watermark.
`)
}
