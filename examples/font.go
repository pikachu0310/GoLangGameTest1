package examples

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
)

const (
	screenWidthFont  = 640
	screenHeightFont = 480
)

const sampleText = `The quick brown fox jumps over the lazy dog.`

var (
	mplusNormalFont font.Face
	mplusBigFont    font.Face
	jaKanjis        = []rune{}
)

func init() {
	// table is the list of Japanese Kanji characters in a part of JIS X 0208.
	const table = `
あいうえおかきくけこさしすせそたちつてとなにぬねのはひふへほまみむめもやゆよらりるれろわをんがぎぐげござじずぜぞだぢづでどばびぶべぼぱぴぷぺぽ
`
	for _, c := range table {
		if c == '\n' {
			continue
		}
		jaKanjis = append(jaKanjis, c)
	}
}

func init() {
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

func init() {
	rand.Seed(time.Now().UnixNano())
}

type FontGame struct {
	counter        int
	kanjiText      string
	kanjiTextColor color.RGBA
}

func (g *FontGame) Update() error {
	// Change the text color for each second.
	if g.counter%ebiten.TPS() == 0 {
		g.kanjiText = ""
		for j := 0; j < 4; j++ {
			for i := 0; i < 8; i++ {
				g.kanjiText += string(jaKanjis[rand.Intn(len(jaKanjis))])
			}
			g.kanjiText += "\n"
		}

		g.kanjiTextColor.R = 0x80 + uint8(rand.Intn(0x7f))
		g.kanjiTextColor.G = 0x80 + uint8(rand.Intn(0x7f))
		g.kanjiTextColor.B = 0x80 + uint8(rand.Intn(0x7f))
		g.kanjiTextColor.A = 0xff
	}
	g.counter++
	return nil
}

func (g *FontGame) Draw(screen *ebiten.Image) {
	const x = 20

	// Draw info
	msg := fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS())
	text.Draw(screen, msg, mplusNormalFont, x, 40, color.White)

	// Draw the sample text
	text.Draw(screen, sampleText, mplusNormalFont, x, 80, color.White)

	// Draw Kanji text lines
	text.Draw(screen, g.kanjiText, mplusBigFont, x, 160, g.kanjiTextColor)
}

func (g *FontGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidthFont, screenHeightFont
}

func ExampleFont() {
	ebiten.SetWindowSize(screenWidthFont, screenHeightFont)
	ebiten.SetWindowTitle("Font (Ebitengine Demo)")
	if err := ebiten.RunGame(&FontGame{}); err != nil {
		log.Fatal(err)
	}
}
