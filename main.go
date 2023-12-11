package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"strconv"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
)

const (
	sampleText = "Photon"
)

var (
	screenWidth     = 640
	screenHeight    = 480
	mplusNormalFont font.Face
	mplusBigFont    font.Face
	photonline      []string
	cursornowx      = int(1)
	cursornowy      = int(1)
	justadot        = ebiten.NewImage(10, 10)
	returncode      = "\n"
)

//func

func init() {
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}
	mplusBigFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    48,
		DPI:     dpi,
		Hinting: font.HintingFull, // Use quantization to save glyph cache images.
	})
	if err != nil {
		log.Fatal(err)
	}

	// Adjust the line height.
	mplusBigFont = text.FaceWithLineHeight(mplusBigFont, 54)
}

type Game struct {
	counter        int
	kanjiText      string
	kanjiTextColor color.RGBA
}

func (g *Game) Update() error {
	/*\
	 * detect mouse wheel actions.
	\*/
	_, dy := ebiten.Wheel()
	if dy > 0 {
		cursornowy--
	}
	if dy < 0 {
		cursornowy++
	}
	return nil

	// detect keyboard actions.

}

func (g *Game) Draw(screen *ebiten.Image) {
	const x = 20
	screen.Fill(color.RGBA{0x80, 0x80, 0x80, 0xff})

	// Draw info
	msg := fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS())
	text.Draw(screen, msg, mplusNormalFont, x, 40, color.White)

	// Draw the sample text
	text.Draw(screen, sampleText, mplusNormalFont, x, 80, color.White)

	// Draw Kanji text lines
	text.Draw(screen, strconv.Itoa(cursornowy), mplusBigFont, x, 160, color.White)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	screenWidth, screenHeight := ebiten.WindowSize()
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("photon")

	err := os.WriteFile("./output.txt", []byte("超unko"+returncode+"unko"), 0644)
	if err != nil {
		log.Fatal(err)
	}

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

/*for index, runeValue := range "超unko" {
println("位置:", index, "文字:", string([]rune{runeValue}))
}*/
