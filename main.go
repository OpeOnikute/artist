package main

import (
	"fmt"
	"image"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
)

const invalidCommand = "Please enter a valid command."
const helpSuggestion = "Use the 'help' command for reference."
const watermarkSize = "200x200"

func main() {

	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println(invalidCommand, helpSuggestion)
		return
	}

	switch args[0] {
	case "help":
		getHelp()
		break
	case "watermark":
		// TODO: Add flag for multiple images
		if len(args) < 3 {
			fmt.Println("The watermark command takes two arguments.", helpSuggestion)
			return
		}
		addWaterMark(args[1], args[2])
		break
	case "place":
		if len(args) < 6 {
			fmt.Println("The place command takes five arguments.", helpSuggestion)
			return
		}
		placeImage(args[0], args[1], args[2], args[3], args[4])
		break
	default:
		fmt.Println(invalidCommand, helpSuggestion)
		break
	}
}

func addWaterMark(bgImg, watermark string) {

	outName := fmt.Sprintf("watermark-new-%s", watermark)

	src := openImage(bgImg)
	mark := openImage(watermark)

	markFit := imaging.Fit(mark, 200, 200, imaging.Lanczos)

	bgDimensions := src.Bounds().Max
	markDimensions := markFit.Bounds().Max

	bgAspectRatio := math.Round(float64(bgDimensions.X) / float64(bgDimensions.Y))

	xPos, yPos := calcWaterMarkPosition(bgDimensions, markDimensions, bgAspectRatio)

	placeImage(outName, bgImg, watermark, watermarkSize, fmt.Sprintf("%dx%d", xPos, yPos))

	fmt.Printf("Added watermark '%s' to image '%s' with dimensions %s.\n", watermark, bgImg, watermarkSize)
}

func placeImage(outName, bgImg, markImg, markDimensions, locationDimensions string) {

	// Resize the mark to fit these dimenstions.
	markWidth, markHeight := parseCoordinates(markDimensions, "x")

	// Coordinate to super-impose on. e.g. 200x500
	locationX, locationY := parseCoordinates(locationDimensions, "x")

	src := openImage(bgImg)
	mark := openImage(markImg)

	markFit := imaging.Fit(mark, markWidth, markHeight, imaging.Lanczos)

	dst := imaging.Paste(src, markFit, image.Pt(locationX, locationY))

	err := imaging.Save(dst, fmt.Sprintf("data/%s", outName))

	if err != nil {
		log.Fatalf("failed to save image: %v", err)
	}

	fmt.Printf("Placed image '%s' on '%s'.\n", markImg, bgImg)
}

// Subtracts the dimensions of the watermark and padding based on the background's aspect ratio
func calcWaterMarkPosition(bgDimensions, markDimensions image.Point, aspectRatio float64) (int, int) {

	bgX := bgDimensions.X
	bgY := bgDimensions.Y
	markX := markDimensions.X
	markY := markDimensions.Y

	padding := 20 * int(aspectRatio)

	return bgX - markX - padding, bgY - markY - padding
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

func openImage(name string) image.Image {
	src, err := imaging.Open(name)
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}
	return src
}

func getHelp() {
	fmt.Print(
		`This is a simple program that imposes an image over another.
Usage:
	go run main.go <command> <..args>
	e.g. go run place main.go zebra.png sample1.png mark.png 200x200 100x100
Commmands:
	- watermark
	DESCRIPTION
	Add a watermark over an image. Location for watermark is the bottom right.
	ARGUMENTS
	- background-name: Name of background image.
	- watermark-name: Name of smaller image.
	- place
	DESCRIPTION
	Places an image over the other.
	ARGUMENTS
	- output-name: What should the output image be called?
	- background-name: Name of background image.
	- watermark-name: Name of smaller image.
	- watermark-dimensions: Watermark would be resized to this before imposing.
	- location: x & y coordinates to place the watermark.
	- help
	DESCRIPTION
	View this content. 
`)
}
