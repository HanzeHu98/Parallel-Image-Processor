// Package png allows for loading png images and applying
// image flitering effects on them.
package png

import (
	"image/color"
	"image"
)

// Grayscale applies a grayscale filtering effect to the image
func (img *Image) Grayscale(subBounds image.Rectangle) {

	// Only calculate over the given subBounds of this image, to make this code more
	// parallelizable.
	for y := subBounds.Min.Y; y < subBounds.Max.Y; y++ {
		for x := subBounds.Min.X; x < subBounds.Max.X; x++ {
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

// Sharpen applies a sharpening filtering effect to the image
func (img *Image) Sharpen(subBounds image.Rectangle){
	filter := [9]float64{0, -1, 0, -1, 5, -1, 0, -1, 0}

	// Only calculate over the given subBounds of this image, to make this code more
	// parallelizable.
	for y := subBounds.Min.Y; y < subBounds.Max.Y; y++{
		for x := subBounds.Min.X; x < subBounds.Max.X; x++ {
			// Get square around pixel, apply convolution filter over each of the three
			// color channels, and then set pixel to new summed value.
			red, green, blue, a:= img.getSquare(x, y)
			r := applyFilter(filter, red)
			g := applyFilter(filter, green)
			b := applyFilter(filter, blue)

			img.out.Set(x, y, color.RGBA64{r, g, b, uint16(a)})
		}
	}
}

// EdgeDetection applies a EdgeDetection filtering effect to the image
func (img *Image) EdgeDetection(subBounds image.Rectangle){
	filter := [9]float64{-1, -1, -1, -1, 8, -1, -1, -1, -1}

	// Only calculate over the given subBounds of this image, to make this code more
	// parallelizable.
	for y := subBounds.Min.Y; y < subBounds.Max.Y; y++{
		for x := subBounds.Min.X; x < subBounds.Max.X; x++ {
			// Get square around pixel, apply convolution filter over each of the three
			// color channels, and then set pixel to new summed value.
			red, green, blue, a:= img.getSquare(x, y)
			r := applyFilter(filter, red)
			g := applyFilter(filter, green)
			b := applyFilter(filter, blue)

			img.out.Set(x, y, color.RGBA64{r, g, b, uint16(a)})
		}
	}
}

// Blur applies a Blur filtering effect to the image
func (img *Image) Blur(subBounds image.Rectangle){
	filter := [9]float64{1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0}
	
	// Only calculate over the given subBounds of this image, to make this code more
	// parallelizable.
	for y := subBounds.Min.Y; y < subBounds.Max.Y; y++{
		for x := subBounds.Min.X; x < subBounds.Max.X; x++ {
			// Get square around pixel, apply convolution filter over each of the three
			// color channels, and then set pixel to new summed value.
			red, green, blue, a:= img.getSquare(x, y)
			r := applyFilter(filter, red)
			g := applyFilter(filter, green)
			b := applyFilter(filter, blue)

			img.out.Set(x, y, color.RGBA64{r, g, b, uint16(a)})
		}
	}
}

// Apply the convolution filter over the square and sum to a single value,
// clamp it down so it can be added to the output image
func applyFilter(filter [9]float64, origin []uint32) uint16 {
	var sum float64 = 0
	for i := 0; i < len(filter); i++ {
		sum = sum + filter[i]*float64(origin[i])
	}
	return clamp(sum)
}


// Get the 3 * 3 square around a pixel so we could apply convolution filter over it
// If the square is at the edge of the image, zero pad it. Return 3 uint32 array of 
// length 9, one for each of rbg, as well as the a value.
func (img *Image) getSquare(x int, y int) ([]uint32, []uint32, []uint32, uint32){
	var red, green, blue []uint32

	bounds := img.out.Bounds()
	for row := x - 1; row <= x + 1; row ++{
		for col := y - 1; col <= y + 1; col ++{
			// Zero pad if the image is at the edge, otherwise append it to the array
			if row < bounds.Min.X || row >= bounds.Max.X || col < bounds.Min.Y || col >= bounds.Max.Y{
				red = append(red, 0)
				green = append(green, 0)
				blue = append(blue, 0)
			} else {
				r, g, b, _ := img.in.At(row, col).RGBA()
				red = append(red, r)
				green = append(green, g)
				blue = append(blue, b)
			}
		}
	}
	_, _, _, a := img.in.At(x, y).RGBA()
	return red, green, blue, a
}

// Swap the in and out image squares of the Image struct, so we can write over/processes
// the next effect immediately without saving the image.
func (img *Image) SwapInOut() {
	temp := img.in
	img.in = img.out
	img.out = temp
}