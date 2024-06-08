package main

/*
 * PhotonText very alpha (codename anode) by plasticgaming99
 * (c)opyright plasticgaming99, 2023-2024
 */

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"

	"github.com/plasticgaming99/photon/assets/phfonts"
	"github.com/plasticgaming99/photon/assets/phicons"
	"github.com/plasticgaming99/photon/modules/dyntypes"
	"github.com/plasticgaming99/photon/modules/plastk"

	"github.com/hugolgst/rich-go/client"
	"golang.design/x/clipboard"
)

/* basic */
var (
	screenWidth      = 640
	screenHeight     = 480
	mplusNormalFont  font.Face
	mplusBigFont     font.Face
	mplusSmallFont   font.Face
	HackGenFont      font.Face
	smallHackGenFont font.Face

	photontext = []string{}

	photoncmd   = string("")
	cmdresult   = string("")
	clearresult = int(0)
	rellines    = int(0)

	textrepeatness = int(0)

	cursornowx    = int(1)
	cursornowy    = int(1)
	cursorxeffort = int(0)

	closewindow   = false
	clickrepeated = false
	returncode    = string("\n")
	returntype    = string("")

	// options
	hanzenlock      = true
	hanzenlockstat  = false
	limitterenabled = true
	limitterlevel   = int(5)
	dbgmode         = false
	editmode        = int(1)

	editorforcused      = true
	commandlineforcused = false

	editingfile = string("")

	// Textures
	sideBar         *ebiten.Image
	infoBar         *ebiten.Image
	commandLine     *ebiten.Image
	cursorimg       = ebiten.NewImage(2, 15)
	topopbar        *ebiten.Image
	filesmenubutton = ebiten.NewImage(80, 20)
	topopbarsep     = ebiten.NewImage(1, 20)
	linessep        *ebiten.Image
	scrollbar       *ebiten.Image
	scrollbit       *ebiten.Image

	// Texture options
	sideBarop = &ebiten.DrawImageOptions{}
)

/* Texture options */
var (
	topopBarSize    = 20
	infoBarSize     = 20
	commandlineSize = 20
	scrollbarwidth  = 18
)

func repeatingKeyPressed(key ebiten.Key) bool {
	var (
		delay    = ebiten.TPS() / 2
		interval = ebiten.TPS() / 18
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

func checkMixedKanjiLength(kantext string, length int) (int, int, int) {
	kantext = string([]rune(kantext)[0 : length-1])
	kanji := (len(kantext) - len([]rune(kantext))) / 2
	nonkanji := len([]rune(kantext)) - kanji
	tab := strings.Count(kantext, "	")
	return nonkanji, kanji - tab, tab
}

// renew image. if same size, return same image
func renewimg(image *ebiten.Image, targetwidth int, targetheight int, targetcolor color.Color) *ebiten.Image {
	widthnow := image.Bounds().Dx()
	heightnow := image.Bounds().Dy()
	if widthnow != targetwidth || heightnow != targetheight {
		newimage := ebiten.NewImage(targetwidth, targetheight)
		newimage.Fill(targetcolor)
		return newimage
	}
	return image
}

// func
func init() {
	const dpi = 144

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		var iconphoton []image.Image
		ebiten.SetVsyncEnabled(false)
		ebiten.SetScreenClearedEveryFrame(false)

		iconphotonreader := bytes.NewReader(phicons.PhotonIcon16)
		_, iconphoton16, err := ebitenutil.NewImageFromReader(iconphotonreader)
		phloginfo(err)
		iconphotonreader = bytes.NewReader(phicons.PhotonIcon32)
		_, iconphoton32, err := ebitenutil.NewImageFromReader(iconphotonreader)
		phloginfo(err)
		iconphotonreader = bytes.NewReader(phicons.PhotonIcon48)
		_, iconphoton48, err := ebitenutil.NewImageFromReader(iconphotonreader)
		phloginfo(err)
		iconphotonreader = bytes.NewReader(phicons.PhotonIcon128)
		_, iconphoton128, err := ebitenutil.NewImageFromReader(iconphotonreader)
		phloginfo(err)
		iconphoton = append(iconphoton, iconphoton16, iconphoton32, iconphoton48, iconphoton128)

		ebiten.SetWindowIcon(iconphoton)

		/*100, 250, 500, 750, 1000 or your monitor's refresh rate*/
		ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

		err = clipboard.Init()
		if err != nil {
			fmt.Println("**WARN** Clipboard is disabled.", err)
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		tt, err := opentype.Parse(phfonts.MPlus1pRegular_ttf)
		if err != nil {
			log.Fatal(err)
		}

		mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
			Size:    12,
			DPI:     dpi,
			Hinting: font.HintingVertical,
		})
		if err != nil {
			log.Fatal(err)
		}
		mplusSmallFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
			Size:    8,
			DPI:     dpi,
			Hinting: font.HintingVertical,
		})
		if err != nil {
			log.Fatal(err)
		}
		mplusBigFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
			Size:    24,
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
			Size:    12,
			DPI:     dpi,
			Hinting: font.HintingFull,
		})
		if err != nil {
			log.Fatal(err)
		}
		smallHackGenFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
			Size:    8,
			DPI:     dpi,
			Hinting: font.HintingFull,
		})
		if err != nil {
			log.Fatal(err)
		}
		wg.Done()
	}()

	// load file
	wg.Add(1)
	go func() {
		if len(os.Args) >= 2 {
			phload(os.Args[1])
		}
		wg.Done()
	}()

	wg.Wait()

	// after loaded text to memory, if photontext
	// has not any strings, init photontext with
	// 1-line, 0-column text.
	wg.Add(1)
	go func() {
		if len(photontext) == 0 {
			photontext = append(photontext, "")
		}
		wg.Done()
	}()
	wg.Wait()

	// Execute PhotonRC when its avaliable
	{
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
		}
		phrcpath := home + "/.photonrc"

		if err != nil {
			panic(err)
		}
		_, err = os.Stat(phrcpath)
		if err == nil {
			fmt.Println("Using PhotonRC")
			photonRC, err := sliceload(phrcpath)
			if err != nil {
				fmt.Println("PhotonRC initializing failed. Using default")
			}
			for i := 0; i < len(photonRC); i++ {
				proceedcmd(photonRC[i])
			}
			photonRC = nil
		}
	}

	{
		// After executed PhotonRC, Initialize textures
		sideBar = ebiten.NewImage(60, 3000)
		infoBar = ebiten.NewImage(4100, infoBarSize)
		commandLine = ebiten.NewImage(4100, commandlineSize)
		cursorimg = ebiten.NewImage(2, 15)
		topopbar = ebiten.NewImage(4100, topopBarSize)
		filesmenubutton = ebiten.NewImage(80, 20)
		topopbarsep = ebiten.NewImage(1, 20)
		linessep = ebiten.NewImage(2, 3000)
		scrollbar = ebiten.NewImage(scrollbarwidth, 100)
		scrollbit = ebiten.NewImage(scrollbarwidth, 100)
	}

	// Fill textures
	{
		/* init sidebar image. */
		sideBar.Fill(color.RGBA{57, 57, 57, 255})
		/* init information bar image */
		infoBar.Fill(color.RGBA{87, 97, 87, 255})
		/* init commandline image */
		commandLine.Fill(color.RGBA{39, 39, 39, 255})
		/* init cursor image */
		cursorimg.Fill(color.RGBA{255, 255, 255, 5})
		/* init top-op-bar image */
		topopbar.Fill(color.RGBA{100, 100, 100, 255})
		/* init top-op-bar "files" button */
		filesmenubutton.Fill(color.RGBA{110, 110, 110, 255})
		/* init top-op-bar separator */
		topopbarsep.Fill(color.RGBA{0, 0, 0, 255})
		/* init line-bar separator */
		linessep.Fill(color.RGBA{100, 100, 100, 255})
		/* init scroll-bar */
		scrollbar.Fill(color.RGBA{80, 80, 80, 255})
		/* init scroll-bit */
		scrollbit.Fill(color.RGBA{30, 30, 30, 255})

		// Init texture options
		sideBarop.GeoM.Translate(float64(0), float64(20))
	}
}

type Editor struct {
	/*counter        int
	kanjiText      string
	kanjiTextColor color.RGBA*/
	runeunko []rune
}

func checkcurx(line int) {
	if len([]rune(photontext[line-1])) < cursornowx {
		if photontext[line-1] == "" {
			cursornowx = 1
		} else {
			cursornowx = len([]rune(photontext[line-1])) + 1
		}
	}
}

func (g *Editor) Update() error {
	tickwg := &sync.WaitGroup{}
	// Update Text-info
	// photonlines = len(photontext)
	/*\
	 * detect cursor key actions
	\*/
	tickwg.Add(1)
	// Insert text
	go func() {
		/*photontext[cursornowy-1] = photontext[cursornowy-1] + string(g.runeunko) (legacy impl) */
		if editorforcused && !(ebiten.IsKeyPressed(ebiten.KeyControl)) {
			g.runeunko = ebiten.AppendInputChars(g.runeunko[:0])
			// Detect left side
			if cursornowx == 1 {
				photontext[cursornowy-1] = string(g.runeunko) + photontext[cursornowy-1]
			} else
			// Detect right side
			if cursornowx-1 == len([]rune(photontext[cursornowy-1])) {
				photontext[cursornowy-1] = photontext[cursornowy-1] + string(g.runeunko)
			} else
			// Other, Insert
			{
				photontext[cursornowy-1] = string([]rune(photontext[cursornowy-1])[:cursornowx-1]) + string(g.runeunko) + string([]rune(photontext[cursornowy-1])[cursornowx-1:])
			}
			// Move cursornowx. with cjk support yay!
			cursornowx += len(g.runeunko)
		}

		if editorforcused {
			// Check commandline is called
			if (ebiten.IsKeyPressed(ebiten.KeyControl)) && (ebiten.IsKeyPressed(ebiten.KeyShift)) && (ebiten.IsKeyPressed(ebiten.KeyC)) {
				editorforcused = false
				commandlineforcused = true
			} else
			// Check upper text.
			if (repeatingKeyPressed(ebiten.KeyUp)) && (cursornowy > 1) {
				checkcurx(cursornowy - 1)
				cursornowy--
			} else
			// Check lower text.
			if (repeatingKeyPressed(ebiten.KeyDown)) && (cursornowy < len(photontext)) {
				checkcurx(cursornowy + 1)
				cursornowy++
			} else if (repeatingKeyPressed(ebiten.KeyLeft)) && (cursornowx > 1) {
				cursornowx--
			} else if (repeatingKeyPressed(ebiten.KeyRight)) && (cursornowx <= len([]rune(photontext[cursornowy-1]))) {
				cursornowx++
			} else if (ebiten.IsKeyPressed(ebiten.KeyControl)) && (repeatingKeyPressed(ebiten.KeyC)) {
				fmt.Println("c pressed")
			} else if (repeatingKeyPressed(ebiten.KeyBackquote)) && (hanzenlock) {
				if !hanzenlockstat {
					hanzenlockstat = true
				} else {
					hanzenlockstat = false
				}
			} else if repeatingKeyPressed(ebiten.KeyHome) {
				cursornowx = 1
			} else if repeatingKeyPressed(ebiten.KeyEnd) {
				cursornowx = len([]rune(photontext[cursornowy-1])) + 1
			} else if ebiten.IsKeyPressed(ebiten.KeyControl) && repeatingKeyPressed(ebiten.KeyV) {
				testslice := strings.Split(string(clipboard.Read(clipboard.FmtText)), "\n")
				firsttext := string([]rune(photontext[cursornowy-1])[:cursornowx-1])
				lasttext := string([]rune(photontext[cursornowy-1])[cursornowx-1:])
				//{photontext[cursornowy-1] = string(g.runeunko)}

				fmt.Println(testslice)
				if len(testslice) == 1 {
					photontext[cursornowy-1] = firsttext + testslice[0] + lasttext
					cursornowx = len([]rune(firsttext + testslice[0]))
					fmt.Println("one")
				} else {
					for i := 0; i < len(testslice); i++ {
						if i == 0 {
							photontext[cursornowy-1] = firsttext + testslice[i]
						} else {
							{
								photontext = append(photontext[:cursornowy], append([]string{testslice[i]}, photontext[cursornowy:]...)...)
								/*photontext[cursornowy] = string([]rune(photontext[cursornowy-1])[cursornowx-1:])
								photontext[cursornowy-1] = string([]rune(photontext[cursornowy-1])[:cursornowx-1])*/
								cursornowy++
							}
							if i == len(testslice)-1 {
								cursornowx = len([]rune(testslice[i])) + 1
							}
						}
						fmt.Println(i)
					}
				}
			} else if repeatingKeyPressed(ebiten.KeyTab) {
				/*photontext[cursornowy-1] = photontext[cursornowy-1] + string(g.runeunko) (legacy impl) */
				// Detect text input
				// Detect left side
				if cursornowx == 1 {
					photontext[cursornowy-1] = string("	") + photontext[cursornowy-1]
				} else
				// Detect right side
				if cursornowx-1 == len([]rune(photontext[cursornowy-1])) {
					photontext[cursornowy-1] = photontext[cursornowy-1] + string("	")
				} else
				// Other, Insert
				{
					photontext[cursornowy-1] = string([]rune(photontext[cursornowy-1])[:cursornowx-1]) + string("	") + string([]rune(photontext[cursornowy-1])[cursornowx-1:])
				}
				// Move cursornowx. with cjk support yay!
				cursornowx += len("	")
			} else
			// New line
			if (repeatingKeyPressed(ebiten.KeyEnter) || repeatingKeyPressed(ebiten.KeyNumpadEnter)) && !hanzenlockstat {
				{
					photontext = append(photontext[:cursornowy], append([]string{""}, photontext[cursornowy:]...)...)
					photontext[cursornowy] = string([]rune(photontext[cursornowy-1])[cursornowx-1:])
					photontext[cursornowy-1] = string([]rune(photontext[cursornowy-1])[:cursornowx-1])
					cursornowy++
					cursornowx = 1
				}
				cursornowx = 1
			} else
			// Line deletion.
			if repeatingKeyPressed(ebiten.KeyBackspace) && !((len(photontext[0]) == 0) && (cursornowy == 1)) && !hanzenlockstat {
				if (photontext[cursornowy-1] == "") && (len(photontext) != 1) {
					cursornowx = len([]rune(photontext[cursornowy-2])) + 1
					if cursornowy-1 < len(photontext)-1 {
						copy(photontext[cursornowy-1:], photontext[cursornowy:])
					}
					photontext[len(photontext)-1] = ""
					photontext = photontext[:len(photontext)-1]
					cursornowx = len([]rune(photontext[cursornowy-2])) + 1
					cursornowy--
				} else {
					if !((cursornowx == 1) && (cursornowy == 1)) || (cursornowx-1 == len([]rune(photontext[cursornowy-1]))) {
						if cursornowx == 1 {
							cursornowx = len([]rune(photontext[cursornowy-2])) + 1
							photontext[cursornowy-2] = photontext[cursornowy-2] + photontext[cursornowy-1]
							if cursornowy-1 < len(photontext)-1 {
								copy(photontext[cursornowy-1:], photontext[cursornowy:])
							}
							photontext[len(photontext)-1] = ""
							photontext = photontext[:len(photontext)-1]
							cursornowy--
						} else
						//
						if cursornowx-1 == len([]rune(photontext[cursornowy-1])) {
							// 文字列をruneに変換
							runes := []rune(photontext[cursornowy-1])
							// 最後の文字を削除
							runes = runes[:len(runes)-1]
							// runeを文字列に変換して元のスライスに代入
							photontext[cursornowy-1] = string(runes)
							time.Sleep(1 * time.Millisecond)
							// Move to left
							cursornowx--
						} else {
							// Convert to rune
							runes := []rune(photontext[cursornowy-1])[:cursornowx-1]
							// Delete last
							runes = runes[:len(runes)-1]
							// Convert to string and insert
							photontext[cursornowy-1] = string(runes) + string([]rune(photontext[cursornowy-1])[cursornowx-1:])
							time.Sleep(1 * time.Millisecond)
							// Move to left
							cursornowx--
						}
					}
				}
			}
		} else
		// If command-line is forcused
		if commandlineforcused {
			if ebiten.IsKeyPressed(ebiten.KeyEnter) {
				cmdresult = proceedcmd(photoncmd)
				clearresult += 10
				photoncmd = ""
				editorforcused = true
				commandlineforcused = false
			}
			if (len([]rune(photoncmd)) >= 1) && (repeatingKeyPressed(ebiten.KeyBackspace)) {
				cmdrune := []rune(photoncmd)[:len([]rune(photoncmd))-1]
				photoncmd = string(cmdrune)
			} else {
				// detect text input
				g.runeunko = ebiten.AppendInputChars(g.runeunko[:0])

				// insert text
				if string(g.runeunko) != "" {
					photoncmd += string(g.runeunko)
				}
			}
		}
		tickwg.Done()
	}()

	/*\
	 * detect mouse wheel actions.
	\*/
	tickwg.Add(1)
	go func() {
		_, dy := ebiten.Wheel()
		if (dy > 0) && (rellines > 0) {
			rellines -= 3
		} else if (dy < 0) && (rellines+3 < len(photontext)) {
			rellines += 3
		}
		tickwg.Done()
	}()

	/*\
	 * Detect touch on buttons.
	\*/
	tickwg.Add(1)
	go func() {
		mousex, mousey := ebiten.CursorPosition()
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			clickrepeated = true
		}
		if clickrepeated && !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			if mousex < 80 && mousey < 20 {
				fmt.Println("unko!!!")
			}
			clickrepeated = false
		}
		tickwg.Done()
	}()

	/*\
	 * Detect cursor's position, and changes cursor shape
	\*/
	tickwg.Add(1)
	go func() {
		curx, cury := ebiten.CursorPosition()
		windowx, windowy := ebiten.WindowSize()
		if ((60 < curx) && (curx < windowx)) && ((20 < cury) && (cury < windowy-20)) {
			ebiten.SetCursorShape(ebiten.CursorShapeText)
		} else {
			ebiten.SetCursorShape(ebiten.CursorShapeDefault)
		}
		tickwg.Done()
	}()

	tickwg.Add(1)
	go func() {
		if repeatingKeyPressed(ebiten.KeyA) {
			fmt.Println("a pressed")
		}
		tickwg.Done()
	}()

	tickwg.Wait()
	if closewindow {
		return fmt.Errorf("close")
	}
	return nil
}

var (
	prevcurx = int(0)
	prevcury = int(0)
	prevrell = int(0)
	phburst  = int(0)
)

func (g *Editor) Draw(screen *ebiten.Image) {

	if (prevcurx == cursornowx && prevcury == cursornowy && prevrell == rellines) && limitterenabled {
		time.Sleep(time.Duration(phburst) * time.Millisecond)
		if phburst < limitterlevel {
			phburst += 1
		}
	} else {
		if 0 <= phburst {
			phburst = 0
		}
	}

	prevcurx, prevcury, prevrell = cursornowx, cursornowy, rellines

	screenWidth, screenHeight := ebiten.WindowSize()

	screenHeight -= topopBarSize

	if commandlineforcused || cmdresult != "" {
		screenHeight -= 20
	}

	// Init screen
	screen.Fill(color.RGBA{61, 61, 61, 255})

	sideBar = renewimg(sideBar, 60, screenWidth, color.RGBA{57, 57, 57, 255})
	screen.DrawImage(sideBar, sideBarop)

	// Draw left information text
	leftinfotxt := ""
	leftinfotxt = "PhotonText alpha "
	if hanzenlockstat {
		leftinfotxt += "Hanzenlock "
	}

	// Draw right information text
	rightinfotext := " " + returntype + " " + strconv.Itoa(cursornowy) + ":" + strconv.Itoa(cursornowx)

	// draw editor text
	Maxtext := int(math.Ceil(((float64(screenHeight) - 20) / 18)) - 1)
	if int(Maxtext) >= len(photontext) {
		textrepeatness = len(photontext) - 1
	} else {
		textrepeatness = int(Maxtext) - 1
	}

	// start line loop
	var (
		textx        int
		cursorxstart int
	)
	for printext := 0; printext < len(photontext[rellines:]); {
		if printext > int(Maxtext) || (len(photontext)-rellines) == 0 {
			break
		}
		slicedtext := []rune(photontext[printext+rellines])
		textx = 55
		text.Draw(screen, strconv.Itoa(printext+rellines+1), smallHackGenFont, textx+10, ((printext + 2) * 18), color.White)
		textx = textx + (9 * len(strconv.Itoa(int(Maxtext)+rellines)))
		textx += 20
		cursorxstart = textx + 0
		// start column loop
		for textrepeat := 0; textrepeat < len(slicedtext); {
			if string("	") == string(slicedtext[textrepeat]) {
				textx += 30
			} else if len(string(slicedtext[textrepeat])) != 1 {
				// If multi-byte text, print bigger
				text.Draw(screen, string(slicedtext[textrepeat]), smallHackGenFont, textx-1, ((printext + 2) * 18), color.White)
				textx += 15
			} else {
				// If not, print normally
				text.Draw(screen, string(slicedtext[textrepeat]), smallHackGenFont, textx, ((printext + 2) * 18), color.White)
				textx += 9
			}
			textrepeat++
		}
		printext++
	}

	// draw cursor
	nonkanj, kanj, tabs := checkMixedKanjiLength(photontext[cursornowy-1], cursornowx)
	cursorproceedx := (nonkanj*9 + kanj*15 + tabs*36) + cursorxstart

	cursorop := &ebiten.DrawImageOptions{}
	cursorop.GeoM.Translate(float64(cursorproceedx), float64((cursornowy-(rellines))*18)+5)
	screen.DrawImage(cursorimg, cursorop)

	// Draw scroll bar base
	scrollbar = renewimg(scrollbar, scrollbarwidth, screenHeight, color.RGBA{80, 80, 80, 255})
	scrollbarop := &ebiten.DrawImageOptions{}
	scrollbarop.GeoM.Translate(float64(screenWidth)-float64(scrollbarwidth), float64(topopBarSize))
	screen.DrawImage(scrollbar, scrollbarop)

	// Draw scroll bit
	/* init scroll-bit */
	var textsize int
	{
		textsize = len(photontext) + Maxtext
	}
	scrollbartext := float64(screenHeight-20) / float64((float64(textsize) / float64(Maxtext)))
	if scrollbartext < 1 {
		scrollbartext = 1
	}
	scrollbit = renewimg(scrollbit, scrollbarwidth, int(scrollbartext), color.RGBA{30, 30, 30, 255}) //ebiten.NewImage(25, int(scrollbartext))

	scrollbitop := &ebiten.DrawImageOptions{}
	scrollbitop.GeoM.Translate(float64(screenWidth)-float64(scrollbarwidth), float64((float64(screenHeight-20)/float64(textsize))*float64(rellines)+20))
	screen.DrawImage(scrollbit, scrollbitop)

	// Draw lines separator
	linessepop := &ebiten.DrawImageOptions{}
	linessepop.GeoM.Translate(float64(cursorxstart-5), 0)
	screen.DrawImage(linessep, linessepop)

	// Draw info-bar
	infoBarop := &ebiten.DrawImageOptions{}
	infoBarop.GeoM.Translate(0, float64(screenHeight))
	screen.DrawImage(infoBar, infoBarop)

	//// Final render --- Top operation-bar
	screen.DrawImage(topopbar, nil)
	//// Files Button
	screen.DrawImage(filesmenubutton, nil)
	// Label of Files Button
	text.Draw(screen, string("Files"), smallHackGenFont, 10, 15, color.White)
	// Draw separator of files
	topsep1op := &ebiten.DrawImageOptions{}
	topsep1op.GeoM.Translate(80, 0)
	screen.DrawImage(topopbarsep, topsep1op)

	text.Draw(screen, leftinfotxt, smallHackGenFont, 5, screenHeight+infoBarSize-4, color.White)

	text.Draw(screen, rightinfotext, smallHackGenFont, screenWidth-((len(rightinfotext))*10), screenHeight+infoBarSize-4, color.White)

	plastk.DrawMenuBar(screen, color.RGBA{100, 100, 100, 255}, smallHackGenFont, 20, []string{"Files", "Save"}, []string{"Edit", "Undo"}, []string{"View"})

	// draw command-line
	if commandlineforcused || cmdresult != "" {
		commandlineop := &ebiten.DrawImageOptions{}
		commandlineop.GeoM.Translate(float64(0), float64(screenHeight+commandlineSize))
		screen.DrawImage(commandLine, commandlineop)
		if commandlineforcused {
			text.Draw(screen, photoncmd, smallHackGenFont, 5, screenHeight+15+commandlineSize, color.White)
		} else {
			text.Draw(screen, cmdresult, smallHackGenFont, 5, screenHeight+15+commandlineSize, color.White)
		}
	}

	// Draw info
	// ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS()))

	/* Benchmark
	if len(os.Args) >= 2 {
		if os.Args[1] == "bench" {
			os.Exit(0)
		}
	}*/
}

func (g *Editor) Layout(outsideWidth, outsideHeight int) (int, int) {
	screenWidth, screenHeight := ebiten.WindowSize()
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("PhotonText(kari)")

	go func() {
		now := time.Now()
		fmt.Println("photontext will booted with no error(s)")

	loginloop:
		err := client.Login("1199337296307163146")
		if err != nil {
			time.Sleep(20 * time.Second)
			goto loginloop
		}

		success := bool(false)
	activityloop:
		state := strconv.Itoa(cursornowy) + string(":") + strconv.Itoa(cursornowx)
		err = client.SetActivity(client.Activity{
			Details:    "Coding with PhotonText",
			State:      state,
			LargeImage: "photon2",
			LargeText:  "PhotonText Logo",
			Timestamps: &client.Timestamps{
				Start: &now,
			},
		})
		if err != nil {
			fmt.Println(err)
			goto activityloop
		}
		if !success {
			fmt.Println("rich presence active")
			success = true
		}
		time.Sleep(1 * time.Second)
		goto activityloop
	}()

	go func() {
		for {
			looping := false
			if clearresult == 0 {
				if !looping {
					cmdresult = ""
				}
				time.Sleep(1 * time.Second)
				looping = true
			} else if clearresult <= 10 {
				looping = false
				time.Sleep(1 * time.Second)
				clearresult--
			} else if clearresult > 10 {
				looping = false
				clearresult = 10
			}
		}
	}()

	go func() {
		for {
			if editingfile == "" {
				ebiten.SetWindowTitle("PhotonText(kari)")
			} else {
				ebiten.SetWindowTitle(fmt.Sprint(editingfile, " - PhotonText(kari)"))
			}
			time.Sleep(500 * time.Millisecond)
		}
	}()

	if err := ebiten.RunGame(&Editor{}); err != nil {
		log.Fatal(err)
	}
}

// Proceed command
func proceedcmd(command string) (returnstr string) {
	command2slice := strings.Split(command, " ")
	if len(command2slice) >= 1 {
		cmd := command2slice[0]
		// Save override
		if cmd == "w" || cmd == "wr" || cmd == "wri" || cmd == "writ" || cmd == "write" {
			if len(command2slice) >= 2 {
				return "Too many arguments for command: " + cmd
			} else {
				phsave(editingfile)
				return fmt.Sprint("Saved to ", editingfile)
			}
		} else
		//
		if cmd == "q" || cmd == "qu" || cmd == "qui" || cmd == "quit" {
			ebiten.SetWindowClosingHandled(true)
			closewindow = true
		} else if cmd == "wq" {
			proceedcmd("w")
			proceedcmd("q")
		} else if cmd == "version" {
			return "PhotonText rolling " + runtime.Version()
		} else
		// Save with other name.
		if cmd == "sav" || cmd == "save" || cmd == "savea" || cmd == "saveas" {
			if len(command2slice) == 1 {
				return "Too few arguments for command: " + cmd
			} else if len(command2slice) >= 3 {
				return "Too many arguments for command: " + cmd
			} else /* when 2 args */ {
				if strings.HasPrefix(command2slice[1], "~") {
					home, err := os.UserHomeDir()
					if err != nil {
						fmt.Println(err)
					}
					savepath := home + command2slice[1][1:]
					phsave(savepath)
					return fmt.Sprint("Saved to ", savepath)
				} else {
					phsave(command2slice[1])
					return fmt.Sprint("Saved to ", command2slice[1])
				}
			}
		} else
		// Toggle VSync
		if command2slice[0] == "togglevsync" {
			ebiten.SetVsyncEnabled(!ebiten.IsVsyncEnabled())
			return "Toggled VSync"
		} else if command2slice[0] == "set" {
			if len(command2slice) == 1 {
				return "Too few arguments for command: " + cmd
			} else if len(command2slice) >= 3 {
				return "Too many arguments for command: " + cmd
			} else {
				var2set := strings.Split(command2slice[1], "=")[1]
				if 1 < len(var2set) {
					switch strings.Split(command2slice[1], "=")[0] {
					case "vsync":
						if dyntypes.IsDynTypeMatch(var2set, "bool") {
							ebiten.SetVsyncEnabled(dyntypes.DynBool(var2set))
						}
						return strconv.FormatBool(ebiten.IsVsyncEnabled())
					case "rellines":
						if dyntypes.IsDynTypeMatch(var2set, "int") {
							rellines = dyntypes.DynInt(var2set)
						}
					case "topopbarsize":
						if dyntypes.IsDynTypeMatch(var2set, "int") {
							topopBarSize = dyntypes.DynInt(var2set)
						}
					case "infobarsize":
						if dyntypes.IsDynTypeMatch(var2set, "int") {
							infoBarSize = dyntypes.DynInt(var2set)
						}
					case "commandlinesize":
						if dyntypes.IsDynTypeMatch(var2set, "int") {
							commandlineSize = dyntypes.DynInt(var2set)
						}
					case "limitter":
						if dyntypes.IsDynTypeMatch(var2set, "bool") {
							limitterenabled = dyntypes.DynBool(var2set)
						}
					case "limitterlevel":
						if dyntypes.IsDynTypeMatch(var2set, "int") {
							limitterlevel = dyntypes.DynInt(var2set)
						}
					default:
						return "No internal variables named " + (strings.Split(command2slice[1], "="))[0]
					}
				}
			}
		} else
		// If not command is avaliable
		{
			return fmt.Sprintf("Not an editor command: %s", command2slice[0])
		}
	} else {
		return "No command was input."
	}
	return
}

func phloginfo(pherror error) {
	if pherror != nil {
		fmt.Println("error:", pherror)
	}
}

// file load/save
func phload(inputpath string) {
	file, err := os.ReadFile(inputpath)
	if err != nil {
		panic(err)
	}
	editingfile, err = filepath.Abs(inputpath)
	if err != nil {
		panic(err)
	}
	ftext := string(file)

	// Check CRLF(dos) First, if not, Use LF(*nix).
	if strings.Contains(ftext, "\r\n") {
		photontext = strings.Split(ftext, "\r\n")
		returncode = "\r\n"
		returntype = "CRLF"
	} else {
		photontext = strings.Split(ftext, "\n")
		returncode = "\n"
		returntype = "LF"
	}
}

func sliceload(inputpath string) ([]string, error) {
	var slice2load []string
	var sliceerr error
	file, err := os.ReadFile(inputpath)
	if err != nil {
		sliceerr = err
	}
	ftext := string(file)

	// Check CRLF(dos) First, if not, Use LF(*nix).
	if strings.Contains(ftext, "\r\n") {
		slice2load = strings.Split(ftext, "\r\n")
		returncode = "\r\n"
		returntype = "CRLF"
	} else {
		slice2load = strings.Split(ftext, "\n")
		returncode = "\n"
		returntype = "LF"
	}
	return slice2load, sliceerr
}

func phsave(dir string) {
	output := strings.Join(photontext, returncode)
	runeout := []rune(output)
	err := os.WriteFile(dir, []byte(string(runeout)), 0644)
	if err != nil {
		fmt.Println(err, "Save failed")
	}
}

/*for index, runeValue := range "超unko" {
println("位置:", index, "文字:", string([]rune{runeValue}))
}*/
