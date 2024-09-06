package plastk

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	menuSeparation []int
	clickrepeated  bool
	clicktoggle    bool
	clickendurance bool
	prevmenu       int
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

func DrawMenuBar(targImage *ebiten.Image, targColor color.Color, targFont *text.GoTextFace, menuHeight int, menu ...[]MenuBarColumn) {
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
		currentbuttonbase := relx
		relx += 5

		fontWidth, fontHeight := text.Measure(menu[i][0].ColumnName, targFont, 0) //font.MeasureString(targFont, menu[i][0].ColumnName)
		relx2 := (int(fontWidth) + 10) + (relx - 5)

		menuColumnButton := ebiten.NewImage((int(fontWidth) + 10), menuHeight)
		menuColumnButton.Fill(color.Black)
		// i don't know how it works so i can't explain
		if (((menuSeparation[i] < mousex) && (mousex <= relx2)) && (mousey <= menuHeight)) || ((prevmenu == i) && clicktoggle) {
			if !clickendurance && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
				clicktoggle = !clicktoggle
				fmt.Println(clicktoggle)
				clickendurance = true
			} else if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
				clickendurance = false
				prevmenu = i
			}
			if prevmenu == i && clicktoggle {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(relx-5), float64(menuHeight))
				img := ebiten.NewImage(200, menuHeight*(len(menu[i])-1)+1)
				img.Fill(color.Black)
				targImage.DrawImage(img, op)
			}
			if clickrepeated && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
				clickrepeated = false
			}

		} else {
			menuColumnButton.Fill(targColor)
		}

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(relx-5), 0)

		image2merge.DrawImage(menuColumnButton, op)
		//text.Draw(image2merge, menu[i][0].ColumnName, targFont, relx, menuHeight-4, color.White)
		op2 := &text.DrawOptions{}
		op2.GeoM.Translate(float64(relx), (float64(menuHeight)-(fontHeight))/2)
		text.Draw(image2merge, menu[i][0].ColumnName, targFont, op2)

		relx += (int(fontWidth) + 5)

		menuSeparation = append(menuSeparation, relx)

		if clicktoggle && prevmenu == i {
			for ii := 1; ii < len(menu[i]); ii++ {
				op := &text.DrawOptions{}
				op.GeoM.Translate(float64(currentbuttonbase), float64(menuHeight*ii))
				text.Draw(targImage, menu[i][ii].ColumnName, targFont, op)
			}
		}
	}
	targImage.DrawImage(image2merge, nil)
}
