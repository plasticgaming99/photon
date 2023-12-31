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
	"math"
	"os"
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

	photoncmd = ""
	rellines  = int(0)

	textrepeatness = int(0)

	cursornowx    = int(1)
	cursornowy    = int(1)
	clickrepeated = false
	returncode    = "\n"

	// options
	hanzenlock          = true
	hanzenlockstat      = false
	editorforcused      = true
	commandlineforcused = false
	dbgmode             = false

	// Textures
	sidebar         = ebiten.NewImage(60, 3000)
	infoBar         = ebiten.NewImage(4100, 20)
	cursorimg       = ebiten.NewImage(2, 15)
	topopbar        = ebiten.NewImage(4100, 20)
	filesmenubutton = ebiten.NewImage(80, 20)
	topopbarsep     = ebiten.NewImage(1, 20)

	// Texture options
	sidebarop = &ebiten.DrawImageOptions{}
)

func repeatingKeyPressed(key ebiten.Key) bool {
	var (
		delay    = ebiten.TPS() / 4
		interval = ebiten.TPS() / 20
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
	/*kantext = string([]rune(kantext)[0 : length-1])
	kanji := (len(kantext) - len([]rune(kantext))) / 2
	nonkanji := len([]rune(kantext)) - kanji*/
	return len([]rune(string([]rune(kantext)[0:length-1]))) - (len(string([]rune(kantext)[0:length-1]))-len([]rune(string([]rune(kantext)[0:length-1]))))/2, (len(string([]rune(kantext)[0:length-1])) - len([]rune(string([]rune(kantext)[0:length-1])))) / 2
}

// func
func init() {
	go ebiten.SetVsyncEnabled(true)

	/*100, 250, 500, 750, 1000 or your monitor's refresh rate*/
	go ebiten.SetTPS(300)

	go ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	// Fill textures
	/* init sidebar image. */
	go sidebar.Fill(color.RGBA{57, 57, 57, 255})
	/* init information bar image */
	go infoBar.Fill(color.RGBA{87, 97, 87, 255})
	/* init cursor image */
	go cursorimg.Fill(color.RGBA{255, 255, 255, 40})
	/* init top-op-bar image */
	go topopbar.Fill(color.RGBA{100, 100, 100, 255})
	/* init top-op-bar "files" button */
	go filesmenubutton.Fill(color.RGBA{110, 110, 110, 255})
	/* init top-op-bar separator */
	go topopbarsep.Fill(color.RGBA{0, 0, 0, 255})

	// Init texture options
	go sidebarop.GeoM.Translate(float64(0), float64(20))

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
	if dbgmode == true {
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
			}
			// Line deletion.
			if repeatingKeyPressed(ebiten.KeyBackspace) && !((len(photontext[0]) == 0) && (cursornowy == 1)) && !hanzenlockstat {
				if (photontext[cursornowy-1] == "") && (len(photontext) != 1) {
					photontext[cursornowy-1] = photontext[len(photontext)-1]
					photontext[len(photontext)-1] = os.DevNull
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

			// Detect text input
			g.runeunko = ebiten.AppendInputChars(g.runeunko[:0])

			// Insert text
			if len(g.runeunko) > 0 {
				/*photontext[cursornowy-1] = photontext[cursornowy-1] + string(g.runeunko) (legacy impl) */
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
		} else
		// If command-line is forcused
		if commandlineforcused {
			if ebiten.IsKeyPressed(ebiten.KeyEnter) {
				proceedcmd(photoncmd)
				photoncmd = ""
				editorforcused = true
				commandlineforcused = false
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
	tickwg.Wait()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	drawwg := &sync.WaitGroup{}

	screenWidth, screenHeight := ebiten.WindowSize()
	screenHeight -= 20

	drawwg.Add(1)
	go func() {
		screen.Fill(color.RGBA{61, 61, 61, 255})
		screen.DrawImage(sidebar, sidebarop)
		drawwg.Done()
	}()

	// Draw left information text
	drawwg.Add(1)
	leftinfotxt := ""
	go func() {
		leftinfotxt = "PhotonText alpha "
		if hanzenlockstat {
			leftinfotxt += "Hanzenlock "
		}
		drawwg.Done()
	}()

	drawwg.Wait()

	// draw editor text
	Maxtext := math.Ceil(((float64(screenHeight) - 20) / 18)) - 1
	if int(Maxtext) >= len(photontext) {
		textrepeatness = len(photontext) - 1
	} else {
		textrepeatness = int(Maxtext) - 1
	}

	// start line loop
	for printext := 0; printext < len(photontext[rellines:]); {
		slicedtext := []rune(photontext[printext+rellines])
		x := 60
		// start column loop
		for textrepeat := 0; textrepeat < len(slicedtext); {
			if len(string(slicedtext[textrepeat])) != 1 {
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
	// draw comamnd-line
	if commandlineforcused {
		ebitenutil.DebugPrint(screen, photoncmd)
	}

	// draw cursor
	nonkanj, kanj := checkMixedKanjiLength(photontext[cursornowy-1], cursornowx)
	cursorproceedx := nonkanj*9 + kanj*15

	cursorop := &ebiten.DrawImageOptions{}
	cursorop.GeoM.Translate(float64(60+cursorproceedx), float64((cursornowy-(rellines))*18)+5)
	screen.DrawImage(cursorimg, cursorop)

	// Draw info-bar
	infobarop := &ebiten.DrawImageOptions{}
	infobarop.GeoM.Translate(float64(0), float64(screenHeight))
	screen.DrawImage(infoBar, infobarop)

	drawwg.Wait()

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

	// Draw info
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS()))
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
		time.Sleep(5 * time.Second)
		fmt.Println("photontext will loaded")
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
				phsave(".")
			}
		} else
		// Save with other name.
		if command2slice[0] == "saveas" {
			if len(command2slice) == 1 {
				return "command: saveas Needs more arguments."
			} else if len(command2slice) >= 3 {
				return "Too many arguments for command: saveas ."
			}
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
	err := os.WriteFile(fmt.Sprintf("%s/output.txt", dir), []byte(string(runeout)), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

/*for index, runeValue := range "超unko" {
println("位置:", index, "文字:", string([]rune{runeValue}))
}*/
