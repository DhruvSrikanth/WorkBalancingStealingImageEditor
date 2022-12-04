// Package png allows for loading png images and applying
// image flitering effects on them.
package png

import (
	"image/color"
)

// Apply the effect to a segment of the image
func (img *Image) ApplyEffect(effect string, startY, endY int) {
	// Apply the effect
	switch effect {
	case "G":
		img.Grayscale(startY, endY)
	case "S":
		img.Sharpen(startY, endY)
	case "B":
		img.Blur(startY, endY)
	case "E":
		img.EdgeDetect(startY, endY)
	default:
		panic("Invalid effect")
	}
}

// Grayscale applies a grayscale filtering effect to the image
func (img *Image) Grayscale(startY, endY int) {
	// Bounds returns defines the dimensions of the image. Always
	// use the bounds Min and Max fields to get out the width
	// and height for the image
	bounds := img.out.Bounds()
	for y := startY; y < endY; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			//Returns the pixel (i.e., RGBA) value at a (x,y) position
			// Note: These get returned as int32 so based on the math you'll
			// be performing you'll need to do a conversion to float64(..)
			r, g, b, a := img.in.At(x, y).RGBA()

			//Note: The values for r,g,b,a for this assignment will range between [0, 65535].
			//For certain computations (i.e., convolution) the values might fall outside this
			// range so you need to clamp them between those values.
			greyC := clamp(float64(r+g+b) / 3)

			//Note: The values need to be stored back as uint16 (I know weird..but there's valid reasons
			// for this that I won't get into right now).
			img.out.Set(x, y, color.RGBA64{greyC, greyC, greyC, uint16(a)})
		}
	}
}

// Sharpen filter
func (img *Image) Sharpen(startY, endY int) {
	// Sharpen kernel
	kernel := [][]float64{
		{0, -1, 0},
		{-1, 5, -1},
		{0, -1, 0},
	}

	// Apply the kernel
	conv2D(img, kernel, startY, endY)
}

// Blur filter
func (img *Image) Blur(startY, endY int) {
	// Blur kernel
	kernel := [][]float64{
		{1 / 9.0, 1 / 9.0, 1 / 9.0},
		{1 / 9.0, 1 / 9.0, 1 / 9.0},
		{1 / 9.0, 1 / 9.0, 1 / 9.0},
	}

	// Apply the kernel
	conv2D(img, kernel, startY, endY)
}

// Edge filter
func (img *Image) EdgeDetect(startY, endY int) {
	// Edge kernel
	kernel := [][]float64{
		{-1, -1, -1},
		{-1, 8, -1},
		{-1, -1, -1},
	}

	// Apply the kernel
	conv2D(img, kernel, startY, endY)
}

// 2D convolution filter
func conv2D(img *Image, kernel [][]float64, startY, endY int) {
	bounds := img.out.Bounds()
	for y := startY; y < endY; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// For each pixel, compute the inner product of the image and kernel
			rOut, gOut, bOut, aOut := frobeniusNorm(img, kernel, x, y)

			// Clamp the values
			rOutC := clamp(rOut)
			gOutC := clamp(gOut)
			bOutC := clamp(bOut)

			// Set the pixel value
			img.out.Set(x, y, color.RGBA64{rOutC, gOutC, bOutC, uint16(aOut)})

		}
	}
}

// Frobenius inner product of image and kernel
func frobeniusNorm(img *Image, kernel [][]float64, x, y int) (float64, float64, float64, uint32) {
	// Kernel dimensions
	m := len(kernel)
	n := len(kernel[0])

	// Image dimensions
	bounds := img.out.Bounds()

	// Image shift
	shiftY := len(kernel) / 2
	shiftX := len(kernel[0]) / 2

	rOut := 0.0
	gOut := 0.0
	bOut := 0.0
	aOut := uint32(0)
	var imgX, imgY int
	for j := 0; j < m; j++ {
		imgY = y + j - shiftY
		for i := 0; i < n; i++ {
			imgX = x + i - shiftX

			// If the pixel is outside the image, use 0s i.e skip
			if imgY < bounds.Min.Y || imgY > bounds.Max.Y-1 || imgX < bounds.Min.X || imgX > bounds.Max.X-1 {
				continue
			}

			// Get the pixel value at the current position
			rIn, gIn, bIn, aIn := img.in.At(imgX, imgY).RGBA()

			// Multiply the pixel value by the kernel value
			rOut += float64(rIn) * kernel[j][i]
			gOut += float64(gIn) * kernel[j][i]
			bOut += float64(bIn) * kernel[j][i]

			// Alpha remains the same for each pixel (0,0) offset index
			if j == shiftY && i == shiftX {
				aOut = aIn
			}
		}
	}
	return rOut, gOut, bOut, aOut
}
