/**
@ Image to tensor using Golag
@ slimdestro
*/

package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"
	"github.com/golang/tensorflow/tensorflow/go"
)

func DestroFnImgToTensor(img image.Image) (*tensorflow.Tensor, error) {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	var imgData []float32
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			imgData = append(imgData, float32(r), float32(g), float32(b))
		}
	}
	tensor, err := tensorflow.NewTensor(imgData)
	if err != nil {
		return nil, err
	}
	return tensor, nil
}

func main() {
	// Create a new image
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{R: 255, G: 0, B: 0, A: 255}}, image.ZP, draw.Src)
	// Save the image
	f, _ := os.Create("image.png")
	defer f.Close()
	png.Encode(f, img)

	// Convert the image to a tensor
	tensor, err := DestroFnImgToTensor(img)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(tensor)
}