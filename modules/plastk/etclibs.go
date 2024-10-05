package plastk

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// renew image. if same size, return same image
func ReNewImg(image *ebiten.Image, targetwidth int, targetheight int, targetcolor color.Color) *ebiten.Image {
	widthnow := image.Bounds().Dx()
	heightnow := image.Bounds().Dy()
	if widthnow != targetwidth || heightnow != targetheight {
		newimage := ebiten.NewImage(targetwidth, targetheight)
		newimage.Fill(targetcolor)
		return newimage
	}
	return image
}
