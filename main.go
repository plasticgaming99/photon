package main

/*
 * PhotonText very alpha (codename anode) by plasticgaming99
 * (c)opyright plasticgaming99, 2023-
 */

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"

	"github.com/plasticgaming99/photon/assets/phfonts"

	"github.com/hugolgst/rich-go/client"
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
	clickrepeated = false
	returncode    = string("\n")

	// options
	hanzenlock     = true
	hanzenlockstat = false
	dbgmode        = false
	editmode       = int(1)

	editorforcused      = true
	commandlineforcused = false

	// Textures
	sideBar         = ebiten.NewImage(60, 3000)
	infoBar         = ebiten.NewImage(4100, 20)
	commandLine     = ebiten.NewImage(4100, 20)
	cursorimg       = ebiten.NewImage(2, 15)
	topopbar        = ebiten.NewImage(4100, 20)
	filesmenubutton = ebiten.NewImage(80, 20)
	topopbarsep     = ebiten.NewImage(1, 20)

	// Texture options
	sideBarop = &ebiten.DrawImageOptions{}
)

var ( /* advanced */
	unneededwg = &sync.WaitGroup{}
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

// func
func init() {
	const dpi = 144

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		ebiten.SetVsyncEnabled(false)

		/*100, 250, 500, 750, 1000 or your monitor's refresh rate*/

		ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
		wg.Done()
	}()

	wg.Add(1)
	// Fill textures
	go func() {
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

		// Init texture options
		sideBarop.GeoM.Translate(float64(0), float64(20))
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
}

type Game struct {
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

func printdbg(text string) {
	if dbgmode {
		fmt.Println(text)
	}
}

func (g *Game) Update() error {
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
		if commandlineforcused {
			goto skiptocommandline
		}
		// Detect text input
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
	skiptocommandline:
		tickwg.Done()
	}()

	tickwg.Add(1)
	go func() {
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
			} else if (repeatingKeyPressed(ebiten.KeyControl)) && (repeatingKeyPressed(ebiten.KeyC)) {
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
					photontext[cursornowy-1] = photontext[len(photontext)-1]
					photontext[len(photontext)-1] = ""
					photontext = photontext[:len(photontext)-1]
					cursornowy--
					cursornowx = len([]rune(photontext[cursornowy-1])) + 1
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
							// Move to left
							cursornowx--
						} else {
							// Convert to rune
							runes := []rune(photontext[cursornowy-1])[:cursornowx-1]
							// Delete last
							runes = runes[:len(runes)-1]
							// Convert to string and insert
							photontext[cursornowy-1] = string(runes) + string([]rune(photontext[cursornowy-1])[cursornowx-1:])
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
			} else if (len([]rune(photoncmd)) >= 1) && (repeatingKeyPressed(ebiten.KeyBackspace)) {
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
			rellines--
		} else if (dy < 0) && (rellines < len(photontext)) {
			rellines++
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
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screenWidth, screenHeight := ebiten.WindowSize()

	screenHeight -= 20

	if commandlineforcused || cmdresult != "" {
		screenHeight -= 20
	}

	screen.Fill(color.RGBA{61, 61, 61, 255})

	screen.DrawImage(sideBar, sideBarop)

	// Draw left information text
	leftinfotxt := ""
	leftinfotxt = "PhotonText alpha "
	if hanzenlockstat {
		leftinfotxt += "Hanzenlock "
	}

	// draw editor text
	Maxtext := math.Ceil(((float64(screenHeight) - 20) / 18)) - 1
	if int(Maxtext) >= len(photontext) {
		textrepeatness = len(photontext) - 1
	} else {
		textrepeatness = int(Maxtext) - 1
	}

	// start line loop
	for printext := 0; printext < len(photontext[rellines:]); {
		if printext > int(Maxtext) {
			break
		}
		slicedtext := []rune(photontext[printext+rellines])
		x := 60
		// start column loop
		for textrepeat := 0; textrepeat < len(slicedtext); {
			if string("	") == string(slicedtext[textrepeat]) {
				x += 30
			} else if len(string(slicedtext[textrepeat])) != 1 {
				// If multi-byte text, print bigger
				text.Draw(screen, string(slicedtext[textrepeat]), smallHackGenFont, x-1, ((printext + 2) * 18), color.White)
				x += 15
			} else {
				// If not, print normally
				text.Draw(screen, string(slicedtext[textrepeat]), smallHackGenFont, x, ((printext + 2) * 18), color.White)
				x += 9
			}
			textrepeat++
		}
		printext++
	}

	// draw cursor
	nonkanj, kanj, tabs := checkMixedKanjiLength(photontext[cursornowy-1], cursornowx)
	cursorproceedx := nonkanj*9 + kanj*15 + tabs*36

	cursorop := &ebiten.DrawImageOptions{}
	cursorop.GeoM.Translate(float64(60+cursorproceedx), float64((cursornowy-(rellines))*18)+5)
	screen.DrawImage(cursorimg, cursorop)

	// Draw info-bar
	infobarop := &ebiten.DrawImageOptions{}
	infobarop.GeoM.Translate(float64(0), float64(screenHeight))
	screen.DrawImage(infoBar, infobarop)

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

	text.Draw(screen, leftinfotxt, smallHackGenFont, 5, screenHeight+16, color.White)

	// Draw right information text
	text.Draw(screen, strconv.Itoa(cursornowy)+":"+strconv.Itoa(cursornowx), smallHackGenFont, screenWidth-((((len(strconv.Itoa(cursornowx))+len(strconv.Itoa(cursornowy)))+1)*10)+8), screenHeight+16, color.White)

	// draw command-line
	if commandlineforcused || cmdresult != "" {
		commandlineop := &ebiten.DrawImageOptions{}
		commandlineop.GeoM.Translate(float64(0), float64(screenHeight+20))
		screen.DrawImage(commandLine, commandlineop)
		if commandlineforcused {
			text.Draw(screen, photoncmd, smallHackGenFont, 5, screenHeight+35, color.White)
		} else {
			text.Draw(screen, cmdresult, smallHackGenFont, 5, screenHeight+35, color.White)
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

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
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

	activityloop:
		err = client.SetActivity(client.Activity{
			Details:    "Coding with PhotonText",
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
		fmt.Println("rich presence active")
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

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

// Proceed command
func proceedcmd(command string) (returnstr string) {
	command2slice := strings.Split(command, " ")
	if len(command2slice) >= 1 {
		// Save override
		if command2slice[0] == "w" {
			if len(command2slice) >= 2 {
				return "Too many arguments for command: w ."
			} else {
				return "dummy: w"
			}
		} else
		// Save with other name.
		if command2slice[0] == "saveas" {
			if len(command2slice) == 1 {
				return "command: saveas Needs more arguments."
			} else if len(command2slice) >= 3 {
				return "Too many arguments for command: saveas ."
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
		} else
		// If not command is avaliable
		{
			return fmt.Sprintf("%s Is not an editor command.", command2slice[0])
		}
	} else {
		return "No command was input."
	}
	return
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

func phsave(dir string) {
	output := strings.Join(photontext, returncode)
	runeout := []rune(output)
	err := os.WriteFile(fmt.Sprintf("%s", dir), []byte(string(runeout)), 0644)
	if err != nil {
		fmt.Println(err, "Save failed")
	}
}

/*for index, runeValue := range "超unko" {
println("位置:", index, "文字:", string([]rune{runeValue}))
}*/
