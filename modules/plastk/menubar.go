package plastk

import (
	"image/color"

	"golang.org/x/image/font"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

var (
	menuSeparation []int
)

func DrawMenuBar(targImage *ebiten.Image, targColor color.Color, targFont font.Face, menuHeight int, menu ...[]string) {
	menuSeparation = nil
	menuSeparation = append(menuSeparation, int(0))
	var (
		targImgWidth  = targImage.Bounds().Dx()
		targImgHeight = targImage.Bounds().Dy()
	)
	image2merge := ebiten.NewImage(targImgWidth, targImgHeight)
	menuBarImage := ebiten.NewImage(targImgWidth, menuHeight)
	menuBarImage.Fill(targColor)

	image2merge.DrawImage(menuBarImage, nil)

	relx := 0
	for i := 0; len(menu) > i; i++ {
		relx += 5

		fontWidth := font.MeasureString(targFont, menu[i][0])

		menuColumnButton := ebiten.NewImage(int(fontWidth), menuHeight)
		menuColumnButton.Fill(color.Gray16{14})
		text.Draw(image2merge, menu[i][0], targFont, relx, menuHeight-4, color.White)

		relx += (fontWidth.Round() + 5)

		menuSeparation = append(menuSeparation, relx)
	}
	targImage.DrawImage(image2merge, nil)
}
