package main

/*
 * PhotonText very alpha by plasticgaming99
 * (c)opyright plasticgaming99, 2023-
 * mmm code-name like thing? oh yes
 * photontext very alpha "anode"
 */

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/plasticgaming99/photon/assets/phfonts"
)

var (
	screenWidth      = 640
	screenHeight     = 480
	mplusNormalFont  font.Face
	mplusBigFont     font.Face
	mplusSmallFont   font.Face
	HackGenFont      font.Face
	smallHackGenFont font.Face

	photonlines = int(1)
	photontext  = []string{}

	cursornowx    = int(1)
	cursornowy    = int(1)
	clickrepeated = false
	returncode    = "\n"

	// options
	hanzenlock     = true
	hanzenlockstat = false
)

func repeatingKeyPressed(key ebiten.Key) bool {
	const (
		delay    = 150
		interval = 50
	)
	d := inpututil.KeyPressDuration(key)
	if d == 1 {
		return true
	}
	if d >= delay && (d-delay)%interval == 0 {
		return true
	}
	return false
}

func checkMixedKanjiLength(kantext string, length int) (int, int) {
	kantext = string([]rune(kantext)[0 : length-1])
	kanji := (len(kantext) - len([]rune(kantext))) / 2
	nonkanji := len([]rune(kantext)) - kanji
	return nonkanji, kanji
}

// func
func init() {
	go ebiten.SetVsyncEnabled(true)
	go ebiten.SetTPS(500)
	go ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	photontextCopy := make([]string, len(photontext))
	copy(photontextCopy, photontext)

	const dpi = 72

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		tt, err := opentype.Parse(phfonts.MPlus1pRegular_ttf)
		if err != nil {
			log.Fatal(err)
		}

		mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
			Size:    24,
			DPI:     dpi,
			Hinting: font.HintingVertical,
		})
		if err != nil {
			log.Fatal(err)
		}
		mplusSmallFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
			Size:    16,
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
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		tt, err := opentype.Parse(phfonts.HackGenRegular_ttf)
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
		wg.Done()
	}()

	wg.Wait()

	// load file
	if len(os.Args) >= 2 {
		phload(os.Args[1])
	}

	// after loaded text to memory, if photontext
	// has not any strings, init photontext with
	// 1-line, 0-column text.
	if len(photontext) == 0 {
		photontext = append(photontext, "")
	}
}

type Game struct {
	/*counter        int
	kanjiText      string
	kanjiTextColor color.RGBA*/
	runeunko []rune
}

func checkcurx(line int) {
	if !(len([]rune(photontext[line-1])) >= cursornowx) {
		if photontext[line-1] == "" {
			cursornowx = 1
		} else {
			fmt.Println([]rune(photontext[line-1]))
			cursornowx = len([]rune(photontext[line-1])) + 1
		}
	}
}

func (g *Game) Update() error {
	// Update Text-info
	photonlines = len(photontext)

	/*\
	 * detect mouse wheel actions.
	\*/
	/*_, dy := ebiten.Wheel()
	if (dy > 0) && (cursornowy > 1) {
		cursornowy--
	}else
	if dy < 0 && (cursornowy < photonlines) {
		cursornowy++
	}*/

	/*\
	 * detect cursor key actions
	\*/
	// Check upper text.
	if (repeatingKeyPressed(ebiten.KeyUp)) && (cursornowy > 1) {
		checkcurx(cursornowy - 1)
		cursornowy--
	} else
	// Check lower text.
	if (repeatingKeyPressed(ebiten.KeyDown)) && (cursornowy < photonlines) {
		checkcurx(cursornowy + 1)
		cursornowy++
	} else if (repeatingKeyPressed(ebiten.KeyLeft)) && (cursornowx > 1) {
		cursornowx--
	} else if (repeatingKeyPressed(ebiten.KeyRight)) && (cursornowx <= len([]rune(photontext[cursornowy-1]))) {
		cursornowx++
	} else if (repeatingKeyPressed(ebiten.KeyControl)) && (repeatingKeyPressed(ebiten.KeyC)) {
		phsave()
	} else if (repeatingKeyPressed(ebiten.KeyBackquote)) && (hanzenlock) {
		if hanzenlockstat == false {
			hanzenlockstat = true
		} else {
			hanzenlockstat = false
		}
	}

	if (repeatingKeyPressed(ebiten.KeyEnter) || repeatingKeyPressed(ebiten.KeyNumpadEnter)) && !hanzenlockstat {
		if cursornowy == len(photontext) {
			photontext = append(photontext, string(""))
			cursornowy++
		} else {
			photontext = append(photontext[:cursornowy], append([]string{""}, photontext[cursornowy:]...)...)
			cursornowy++
			cursornowx = 1
		}
		cursornowx = 1
	} else if repeatingKeyPressed(ebiten.KeyBackspace) && !hanzenlockstat {
		if (photontext[cursornowy-1] == "") && (len(photontext) != 1) {
			photontext[cursornowy-1] = photontext[len(photontext)-1]
			photontext[len(photontext)-1] = os.DevNull
			photontext = photontext[:len(photontext)-1]
			cursornowy--
			cursornowx = len([]rune(photontext[cursornowy-1])) + 1
		} else {
			if !((cursornowy == 1) && (len(photontext[0]) == 0)) {
				// 文字列をruneに変換
				runes := []rune(photontext[cursornowy-1])
				// 最後の文字を削除
				runes = runes[:len(runes)-1]
				// runeを文字列に変換して元のスライスに代入
				photontext[cursornowy-1] = string(runes)
				// Move to left
				cursornowx--
			}
		}
	}

	// Detect text input
	g.runeunko = ebiten.AppendInputChars(g.runeunko[:0])

	// Insert text
	if string(g.runeunko) != "" {
		fmt.Println(string(g.runeunko))
		/*photontext[cursornowy-1] = photontext[cursornowy-1] + string(g.runeunko) (legacy impl) */
		// Detect left side
		if cursornowx == 1 {
			photontext[cursornowy-1] = string(g.runeunko) + photontext[cursornowy-1]
		} else
		// Detect right side
		if cursornowx-1 == len([]rune(photontext[cursornowy-1])) {
			photontext[cursornowy-1] = photontext[cursornowy-1] + string(g.runeunko)
		} else
		// Other
		{
			photontext[cursornowy-1] = string([]rune(photontext[cursornowy-1])[:cursornowx-1]) + string(g.runeunko) + string([]rune(photontext[cursornowy-1])[:len([]rune(photontext[cursornowy-1]))-(cursornowx-1)])
		}
		// Move cursornowx. with cjk support yay!
		cursornowx += len(g.runeunko)
	}

	/*\
	 * Detect touch on buttons.
	\*/
	mousex, mousey := ebiten.CursorPosition()
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		clickrepeated = true
	}
	if clickrepeated && !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if mousex < 80 && mousey < 20 {
			fmt.Println("unko!!!")
		}
		clickrepeated = false
		return nil
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
	/* init cursor image */
	cursorimg := ebiten.NewImage(2, 15)
	cursorimg.Fill(color.RGBA{255, 255, 255, 40})
	/* init top-op-bar image */
	topopbar := ebiten.NewImage(screenWidth, 20)
	topopbar.Fill(color.RGBA{100, 100, 100, 255})
	/* init top-op-bar "files" button */
	filesmenubutton := ebiten.NewImage(80, 20)
	filesmenubutton.Fill(color.RGBA{110, 110, 110, 255})
	/* init top-op-bar separator */
	topopbarsep := ebiten.NewImage(1, 20)
	topopbarsep.Fill(color.RGBA{0, 0, 0, 255})

	screen.Fill(color.RGBA{61, 61, 61, 255})

	sidebarop := &ebiten.DrawImageOptions{}
	sidebarop.GeoM.Translate(float64(0), float64(20))
	screen.DrawImage(sidebar, sidebarop)

	/* Processing Info-Bar Image */
	infobarop := &ebiten.DrawImageOptions{}
	infobarop.GeoM.Translate(float64(0), float64(screenHeight))
	screen.DrawImage(infoBar, infobarop)

	// Draw the text "Photon"
	/*text.Draw(screen, sampleText, mplusNormalFont, x, 80, color.White)*/

	// Draw left information text
	leftinfotxt := "PhotonText alpha "
	if hanzenlockstat == true {
		leftinfotxt += "Hanzenlock "
	}
	text.Draw(screen, leftinfotxt, smallHackGenFont, 5, screenHeight+16, color.White)

	// Draw right infomation text
	text.Draw(screen, strconv.Itoa(cursornowy)+":"+strconv.Itoa(cursornowx), smallHackGenFont, screenWidth-((((len(strconv.Itoa(cursornowx))+len(strconv.Itoa(cursornowy)))+1)*10)+8), screenHeight+16, color.White)

	printext := 0
	for printext < len(photontext) {
		textrepeat := 0
		slicedtext := []rune(photontext[printext])
		x := 60

		for textrepeat < len(slicedtext) {
			if len(string(slicedtext[textrepeat])) != 1 {
				text.Draw(screen, string(slicedtext[textrepeat]), smallHackGenFont, x-1, 20+(printext+1)*18, color.White)
				x += 15
			} else {
				text.Draw(screen, string(slicedtext[textrepeat]), smallHackGenFont, x, 20+(printext+1)*18, color.White)
				x += 9
			}
			textrepeat++
		}
		printext++
	}

	// draw cursor
	cursorop := &ebiten.DrawImageOptions{}
	nonkanj, kanj := checkMixedKanjiLength(photontext[cursornowy-1], cursornowx)
	cursorproceedx := nonkanj*9 + kanj*15

	cursorop.GeoM.Translate(float64(60+cursorproceedx), float64(10+(cursornowy)*18))
	screen.DrawImage(cursorimg, cursorop)

	//// Final render --- Top operation-bar
	screen.DrawImage(topopbar, nil)
	//// Files Button
	screen.DrawImage(filesmenubutton, nil)
	// Label of Files Button
	text.Draw(screen, string("Files"), smallHackGenFont, 10, 15, color.White)
	// Draw separator of files
	topsep1op := &ebiten.DrawImageOptions{}
	topsep1op.GeoM.Translate(float64(80), 0)
	screen.DrawImage(topopbarsep, topsep1op)

	// Draw info
	msg := fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS())
	ebitenutil.DebugPrint(screen, msg)

	// Benchmark
	if len(os.Args) >= 2 {
		if os.Args[1] == "bench" {
			os.Exit(0)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	screenWidth, screenHeight := ebiten.WindowSize()
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("PhotonText(kari)")

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}

	output := strings.Join(photontext, returncode)
	err := os.WriteFile("./output.txt", []byte(output), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func phload(inputpath string) {
	file, err := os.ReadFile(inputpath)
	if err != nil {
		log.Fatal(err)
	}
	ftext := string(file)

	// Check CRLF First, if not, Use LF.
	if strings.Contains(ftext, "\r\n") {
		photontext = strings.Split(ftext, "\r\n")
		returncode = "\r\n"
	} else {
		photontext = strings.Split(ftext, "\n")
		returncode = "\n"
	}
}

func phsave() {
	output := strings.Join(photontext, returncode)
	runeout := []rune(output)
	err := os.WriteFile("./output.txt", []byte(string(runeout)), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

/*for index, runeValue := range "超unko" {
println("位置:", index, "文字:", string([]rune{runeValue}))
}*/
