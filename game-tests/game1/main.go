package main

import (
	"bytes"
	"github.com/pikachu0310/GoLangGameTest1.git/myimages"
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

var img *ebiten.Image
var slime struct {
	Hp      int
	Attack  int
	Defense int
}

func init() {
	var err error
	imga, _, err := image.Decode(bytes.NewReader(myimages.Slime_png))
	if err != nil {
		log.Fatal(err)
	}
	img = ebiten.NewImageFromImage(imga)
}

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(50, 50)
	op.GeoM.Scale(1, 1)
	screen.DrawImage(img, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Geometry Matrix")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
