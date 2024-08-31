package plastk

import (
	"fmt"
	"image/color"

	"golang.org/x/image/font"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

var (
	menuSeparation []int
)

type MenuBarColumn struct {
	// three column types currently avaliable: button, dropdown, separator
	ColumnType string
	// title of column
	ColumnName string
	// if column type is button, other than 0 is ignored.
	// separator never has column.
	ColumnBase []MenuBarColumn
}

func TransformColumn([]MenuBarColumn) []string {
	fmt.Println("todo")
	return []string{"todo"}
}

func DrawMenuBar(targImage *ebiten.Image, targColor color.Color, targFont font.Face, menuHeight int, menu ...[]MenuBarColumn) {
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
		mousex, mousey := ebiten.CursorPosition()
		relx += 5

		fontWidth := font.MeasureString(targFont, menu[i][0].ColumnName)
		relx2 := (fontWidth.Round() + 10) + (relx - 5)

		menuColumnButton := ebiten.NewImage(int(fontWidth.Round()+10), menuHeight)
		if ((menuSeparation[i] < mousex) && (mousex <= relx2)) && (mousey <= menuHeight) {
			menuColumnButton.Fill(color.Black)
		} else {
			menuColumnButton.Fill(targColor)
		}

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(relx-5), 0)

		image2merge.DrawImage(menuColumnButton, op)
		text.Draw(image2merge, menu[i][0].ColumnName, targFont, relx, menuHeight-4, color.White)

		relx += (fontWidth.Round() + 5)

		menuSeparation = append(menuSeparation, relx)
	}
	targImage.DrawImage(image2merge, nil)
}
