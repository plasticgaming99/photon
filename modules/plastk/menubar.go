package plastk

import (
	"image/color"

	"golang.org/x/image/font"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

func DrawMenuBar(targImage *ebiten.Image, targColor color.Color, targFont font.Face, menuHeight int, menu ...[]string) {
	var (
		targImgWidth  = targImage.Bounds().Dx()
		targImgHeight = targImage.Bounds().Dy()
	)
	image2merge := ebiten.NewImage(targImgWidth, targImgHeight)
	menuBarImage := ebiten.NewImage(targImgWidth, menuHeight)
	menuBarImage.Fill(targColor)

	image2merge.DrawImage(menuBarImage, nil)

	relx := 5
	for i := 1; len(menu) < i; i++ {
		fontWidth := font.MeasureString(targFont, menu[i][0])
		menuColumnButton := ebiten.NewImage(int(fontWidth), menuHeight)
		menuColumnButton.Fill(color.Gray16{14})
		text.Draw(menuBarImage, menu[i][0], targFont, relx, 10, color.Gray16{0})
	}

	targImage.DrawImage(image2merge, nil)
}
