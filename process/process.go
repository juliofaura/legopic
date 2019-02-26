package process

import (
	"image"
)

func TurnToBW(in image.Image) (out *image.Gray) {
	out = image.NewGray(in.Bounds())
	for x := 0; x < in.Bounds().Dx(); x++ {
		for y := 0; y < in.Bounds().Dy(); y++ {
			out.Set(x, y, in.At(x, y))
		}
	}
	return
}
