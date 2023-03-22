/**
	@ OCR golang
	@ slimdestro
*/
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"os"
	"strings"

	"github.com/otiai10/gosseract"
)

func main() {
	// Read the image file
	file, err := os.Open("image.jpg")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	// Decode the image
	img, err := jpeg.Decode(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Create a new image with a white background
	bounds := img.Bounds()
	white := color.RGBA{255, 255, 255, 255}
	newImg := image.NewRGBA(bounds)
	draw.Draw(newImg, bounds, &image.Uniform{white}, image.ZP, draw.Src)

	// Draw the original image on top of the white background
	draw.Draw(newImg, bounds, img, image.ZP, draw.Over)

	// Encode the new image
	out, err := os.Create("new_image.jpg")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer out.Close()
	jpeg.Encode(out, newImg, nil)

	// Use Gosseract to extract the text from the image
	client := gosseract.NewClient()
	defer client.Close()
	client.SetImage("new_image.jpg")
	text, _ := client.Text()

	// Clean up the text
	text = strings.Replace(text, "\n", " ", -1)
	text = strings.Replace(text, "\t", " ", -1)
	text = strings.TrimSpace(text)

	// Write the text to a file
	err = ioutil.WriteFile("text.txt", []byte(text), 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}