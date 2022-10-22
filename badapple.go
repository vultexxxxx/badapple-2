package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"strconv"
)

var finalXRes int
var finalYRes int

type BadApple struct {
	image.Image
	Custom map[image.Point]color.Color
}

func (apple *BadApple) Set(x, y int, c color.Color) {
	apple.Custom[image.Point{x, y}] = c
}

func (apple BadApple) At(x, y int) color.Color {
	if color := apple.Custom[image.Point{x, y}]; color != nil {
		return color
	}
	return color.RGBA{0, 0, 0, 255}
}

func (apple BadApple) Bounds() image.Rectangle {
	return image.Rectangle{image.Point{0, 0}, image.Point{finalXRes, finalYRes}}
}

func getBitmap(img image.Image, xRes, yRes int) [][]uint8 {
	var bitmap [][]uint8
	for i := 0; i < yRes; i++ {
		var row []uint8
		for j := 0; j < xRes; j++ {
			r, g, b, _ := img.At(j, i).RGBA()
			avg := (r + g + b) / 3
			if avg >= 128 {
				row = append(row, 1)
			} else {
				row = append(row, 0)
			}
		}
		bitmap = append(bitmap, row)
	}
	return bitmap
}

func (apple *BadApple) createFrame(bitmap [][]uint8, img image.Image, frameNumer, smallXRes, smallYRes int, outputDir string) {
	for columnNumber, row := range bitmap {
		for rowNumber, pixel := range row {
			if pixel == 1 {
				x := rowNumber * smallXRes
				y := columnNumber * smallYRes
				for i := 0; i < smallYRes; i++ {
					for j := 0; j < smallXRes; j++ {
						apple.Set(x, y, img.At(i, j))
						x++
					}
					y++
				}
			}
		}
	}
	output, err := os.Create(fmt.Sprintf("%s/%04d.png", outputDir, frameNumer))
	if err != nil {
		panic(err)
	}
	err = png.Encode(output, apple)
	if err != nil {
		panic(err)
	}

}

func badAppleSquared(numOfFrames int, largeFramesPath, smallFramesPath, outputDir string) {
	var smallXRes, smallYRes, xRes, yRes int
	for a := 1; a <= numOfFrames; a++ {
		file, err := os.Open(fmt.Sprintf("%s/%04d.png", largeFramesPath, a))
		if err != nil {
			panic(err)
		}
		smallFile, err := os.Open(fmt.Sprintf("%s/%04d.png", smallFramesPath, a))
		if err != nil {
			panic(err)
		}
		img, err := png.Decode(file)
		if err != nil {
			panic(err)
		}
		smallImg, err := png.Decode(smallFile)
		if err != nil {
			panic(err)
		}

		if a == 1 {
			imageBounds := img.Bounds()
			smallImageBounds := smallImg.Bounds()
			xRes = imageBounds.Max.X
			yRes = imageBounds.Max.Y
			smallXRes = smallImageBounds.Max.X
			smallYRes = smallImageBounds.Max.Y
			if xRes <= smallXRes || yRes <= smallYRes {
				log.Println("Error with resolution size. Did you point to the correct smallframe and regular frame directories?")
			}
		}

		badapple := BadApple{img, make(map[image.Point]color.Color, 0)}

		bitmap := getBitmap(img, xRes, yRes)

		badapple.createFrame(bitmap, smallImg, a, smallXRes, smallYRes, outputDir)
	}
}

func main() {
	log.Println("HELLO")
	var err error
	numOfFrames := flag.Int("f", 6572, "Number of frames the video was")
	framesPath := flag.String("fp", "frames", "Path to larger frames to form bitmap")
	smallFramesPath := flag.String("sfp", "smallframes", "Path to small frames")
	outputDir := flag.String("o", "finalframes", "Path to output frames")
	res := flag.String("r", "1200x900", "Resolution of output frames")
	flag.Parse()
	resArr := bytes.Split([]byte(*res), []byte("x"))
	if len(resArr) != 2 {
		log.Print("Entered in resolution incorrectly. Should be of the format intxint")
		return
	}
	finalXRes, err = strconv.Atoi(string(resArr[0]))
	if err != nil {
		panic(err)
	}
	finalYRes, err = strconv.Atoi(string(resArr[1]))
	if err != nil {
		panic(err)
	}
	badAppleSquared(*numOfFrames, *framesPath, *smallFramesPath, *outputDir)
	log.Println("GOODBYE")
}
