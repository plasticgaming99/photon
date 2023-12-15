package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"strconv"
	"time"

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
	screenWidth      = 640
	screenHeight     = 480
	mplusNormalFont  font.Face
	mplusBigFont     font.Face
	HackGenFont      font.Face
	smallHackGenFont font.Face
	photonline       []string
	cursornowx       = int(1)
	cursornowy       = int(1)
	justadot         = ebiten.NewImage(10, 10)
	returncode       = "\n"
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

	//unko
	ttbytes, _ := os.ReadFile("/usr/share/fonts/TTF/Hack-Regular.ttf")

	tt, err = opentype.Parse(ttbytes)
	if err != nil {
		log.Fatal(err)
	}

	HackGenFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	smallHackGenFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    16,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
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
	if (dy > 0) && (cursornowy > 1) {
		cursornowy--
	}
	if dy < 0 {
		cursornowy++
	}
	/*\
	 * Detect touch on buttons.
	\*/
	mousex, mousey := ebiten.CursorPosition()
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if mousex < 80 && mousey < 20 {
			for !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
				time.Sleep(50 * time.Millisecond)
			}
			fmt.Println("unko!!!")
			return nil
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screenWidth, screenHeight := ebiten.WindowSize()

	screenHeight -= 20

	/* init sidebar image. */
	sidebar := ebiten.NewImage(60, screenHeight)
	sidebar.Fill(color.RGBA{57, 57, 57, 255})
	/* init information bar image */
	infoBar := ebiten.NewImage(screenWidth, 20)
	infoBar.Fill(color.RGBA{87, 97, 87, 255})
	/* init top-op-bar image */
	topopbar := ebiten.NewImage(screenWidth, 20)
	topopbar.Fill(color.RGBA{100, 100, 100, 255})

	const x = 20
	screen.Fill(color.RGBA{61, 61, 61, 255})

	sidebarop := &ebiten.DrawImageOptions{}
	sidebarop.GeoM.Translate(float64(0), float64(20))
	screen.DrawImage(sidebar, sidebarop)

	/* Processing Info-Bar Image */
	infobarop := &ebiten.DrawImageOptions{}
	infobarop.GeoM.Translate(float64(0), float64(screenHeight))
	screen.DrawImage(infoBar, infobarop)

	// Draw info
	msg := fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS())
	text.Draw(screen, msg, mplusNormalFont, x, 40, color.White)

	// Draw the sample text
	text.Draw(screen, sampleText, mplusNormalFont, x, 80, color.White)

	// Draw Kanji text lines
	text.Draw(screen, "Col:"+strconv.Itoa(cursornowy), smallHackGenFont, screenWidth-(((len(strconv.Itoa(cursornowy))+4)*10)+8), screenHeight+16, color.White)

	//Final render --- Top operation-bar
	screen.DrawImage(topopbar, nil)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	screenWidth, screenHeight := ebiten.WindowSize()
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("PhotonText(kari)")

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
