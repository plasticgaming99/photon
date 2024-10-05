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
	prevID         string

	pressedprevX int
	pressedprevY int

	buttonD *ebiten.Image
)

func init() {
	buttonD = ebiten.NewImage(1, 1)
	buttonD.Fill(color.RGBA{100, 100, 100, 0})
}

type MenuBarColumn struct {
	// three column types currently avaliable: button, dropdown, separator
	ColumnType string
	// title of column
	ColumnName string
	// Information of column
	ColumnDesc string
	// if column type is button, other than 0 is ignored.
	// separator never has column.
	ColumnBase []MenuBarColumn
	// Detect if Button
	ColumnID string
}

func TransformColumn([]MenuBarColumn) []string {
	fmt.Println("todo")
	return []string{"todo"}
}

func DrawMenuBar(targImage *ebiten.Image, targColor color.Color, targFont *text.GoTextFace, menuHeight int, menu ...[]MenuBarColumn) {
	menuSeparation = nil
	menuSeparation = append(menuSeparation, int(0))
	var (
		targImgWidth = targImage.Bounds().Dx()
		//targImgHeight = targImage.Bounds().Dy()
	)
	menuBarImage := ebiten.NewImage(targImgWidth, menuHeight)
	menuBarImage.Fill(targColor)

	targImage.DrawImage(menuBarImage, nil)

	relx := 0
	buttonSize := 200
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
				img := ebiten.NewImage(buttonSize, menuHeight*(len(menu[i])-1)+1)
				img.Fill(color.Black)
				targImage.DrawImage(img, op)
				img.Deallocate()
			}
			if clickrepeated && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
				clickrepeated = false
			}

		} else {
			menuColumnButton.Fill(targColor)
		}

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(relx-5), 0)

		targImage.DrawImage(menuColumnButton, op)
		//text.Draw(image2merge, menu[i][0].ColumnName, targFont, relx, menuHeight-4, color.White)
		op2 := &text.DrawOptions{}
		op2.GeoM.Translate(float64(relx), (float64(menuHeight)-(fontHeight))/2)
		text.Draw(targImage, menu[i][0].ColumnName, targFont, op2)

		relx += (int(fontWidth) + 5)

		menuSeparation = append(menuSeparation, relx)

		if clicktoggle && prevmenu == i {
			for ii := 1; ii < len(menu[i]); ii++ {
				if ((mousex >= currentbuttonbase) && (mousex <= currentbuttonbase+buttonSize)) && ((mousey >= menuHeight*ii) && (mousey <= menuHeight*(ii+1))) {
					buttonD = ReNewImg(buttonD, buttonSize, menuHeight, color.RGBA{100, 100, 100, 0})
					op := &ebiten.DrawImageOptions{}
					op.GeoM.Translate(float64(currentbuttonbase), float64(menuHeight*ii))
					targImage.DrawImage(buttonD, op)
					pressedprevX, pressedprevY = i, ii
					prevID = menu[i][ii].ColumnID
				}
				op := &text.DrawOptions{}
				op.GeoM.Translate(float64(currentbuttonbase), float64(menuHeight*ii))
				text.Draw(targImage, menu[i][ii].ColumnName, targFont, op)
			}
		}
	}
	//targImage.DrawImage(image2merge, nil)
}

func MenuBarDetectClicked(column int, line int) bool {
	if column == pressedprevX && line == pressedprevY {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			pressedprevX, pressedprevY = 0, 0
			return true
		}
	}
	return false
}

// make detection easier
func MenuBarDetectClickedByID(columnID string) bool {
	if columnID == prevID {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			prevID = ""
			return true
		}
	}
	return false
}
