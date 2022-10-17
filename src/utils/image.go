package utils

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"math"
)

func ResizeImage(file io.Reader, width int, height int) (*image.RGBA, error) {
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	minX := img.Bounds().Min.X
	minY := img.Bounds().Min.Y
	maxX := img.Bounds().Max.X
	maxY := img.Bounds().Max.Y
	for (maxX-minX)%width != 0 {
		maxX--
	}
	for (maxY-minY)%height != 0 {
		maxY--
	}

	scaleX := (maxX - minX) / width
	scaleY := (maxY - minY) / height

	imgRect := image.Rect(0, 0, width, height)
	resImg := image.NewRGBA(imgRect)

	draw.Draw(resImg, resImg.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)

	for x := 0; x < width; x += 1 {
		for y := 0; y < height; y += 1 {
			averageColor := getAverageColor(img, minX+x*scaleX, minX+(x+1)*scaleX, minY+y*scaleY, minY+(y+1)*scaleY)
			resImg.Set(x, y, averageColor)
		}
	}

	return resImg, nil
}

func ImageToBytes(img image.Image) ([]byte, int, error) {
	buff := bytes.NewBuffer(nil)
	err := png.Encode(buff, img)
	if err != nil {
		return nil, 0, err
	}

	return buff.Bytes(), buff.Len(), nil
}

func getAverageColor(img image.Image, minX int, maxX int, minY int, maxY int) color.Color {
	var averageRed float64
	var averageGreen float64
	var averageBlue float64
	var averageAlpha float64
	scale := 1.0 / float64((maxX-minX)*(maxY-minY))

	for i := minX; i < maxX; i++ {
		for k := minY; k < maxY; k++ {
			r, g, b, a := img.At(i, k).RGBA()
			averageRed += float64(r) * scale
			averageGreen += float64(g) * scale
			averageBlue += float64(b) * scale
			averageAlpha += float64(a) * scale
		}
	}

	averageRed = math.Sqrt(averageRed)
	averageGreen = math.Sqrt(averageGreen)
	averageBlue = math.Sqrt(averageBlue)
	averageAlpha = math.Sqrt(averageAlpha)

	averageColor := color.RGBA{
		R: uint8(averageRed),
		G: uint8(averageGreen),
		B: uint8(averageBlue),
		A: uint8(averageAlpha)}

	return averageColor
}
