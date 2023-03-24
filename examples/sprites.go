package examples

import (
	"bytes"
	"fmt"
	"image"
	_ "image/png"
	"log"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
)

const (
	screenWidthSprites  = 320
	screenHeightSprites = 240
	maxAngleSprites     = 256
)

var (
	ebitenImage *ebiten.Image
)

func init() {
	// Decode an image from the image file's byte slice.
	img, _, err := image.Decode(bytes.NewReader(images.Ebiten_png))
	if err != nil {
		log.Fatal(err)
	}
	origEbitenImage := ebiten.NewImageFromImage(img)

	s := origEbitenImage.Bounds().Size()
	ebitenImage = ebiten.NewImage(s.X, s.Y)

	op := &ebiten.DrawImageOptions{}
	op.ColorScale.ScaleAlpha(0.5)
	ebitenImage.DrawImage(origEbitenImage, op)
}

type Sprite struct {
	imageWidth  int
	imageHeight int
	x           int
	y           int
	vx          int
	vy          int
	angle       int
}

func (s *Sprite) Update() {
	s.x += s.vx
	s.y += s.vy
	if s.x < 0 {
		s.x = -s.x
		s.vx = -s.vx
	} else if mx := screenWidthSprites - s.imageWidth; mx <= s.x {
		s.x = 2*mx - s.x
		s.vx = -s.vx
	}
	if s.y < 0 {
		s.y = -s.y
		s.vy = -s.vy
	} else if my := screenHeightSprites - s.imageHeight; my <= s.y {
		s.y = 2*my - s.y
		s.vy = -s.vy
	}
	s.angle++
	if s.angle == maxAngleSprites {
		s.angle = 0
	}
}

type Sprites struct {
	sprites []*Sprite
	num     int
}

func (s *Sprites) Update() {
	for i := 0; i < s.num; i++ {
		s.sprites[i].Update()
	}
}

const (
	MinSprites = 0
	MaxSprites = 50000
)

type SpritesGame struct {
	touchIDs []ebiten.TouchID
	sprites  Sprites
	op       ebiten.DrawImageOptions
	inited   bool
}

func (g *SpritesGame) init() {
	defer func() {
		g.inited = true
	}()

	g.sprites.sprites = make([]*Sprite, MaxSprites)
	g.sprites.num = 500
	for i := range g.sprites.sprites {
		w, h := ebitenImage.Bounds().Dx(), ebitenImage.Bounds().Dy()
		x, y := rand.Intn(screenWidthSprites-w), rand.Intn(screenHeightSprites-h)
		vx, vy := 2*rand.Intn(2)-1, 2*rand.Intn(2)-1
		a := rand.Intn(maxAngleSprites)
		g.sprites.sprites[i] = &Sprite{
			imageWidth:  w,
			imageHeight: h,
			x:           x,
			y:           y,
			vx:          vx,
			vy:          vy,
			angle:       a,
		}
	}
}

func (g *SpritesGame) leftTouched() bool {
	for _, id := range g.touchIDs {
		x, _ := ebiten.TouchPosition(id)
		if x < screenWidthSprites/2 {
			return true
		}
	}
	return false
}

func (g *SpritesGame) rightTouched() bool {
	for _, id := range g.touchIDs {
		x, _ := ebiten.TouchPosition(id)
		if x >= screenWidthSprites/2 {
			return true
		}
	}
	return false
}

func (g *SpritesGame) Update() error {
	if !g.inited {
		g.init()
	}
	g.touchIDs = ebiten.AppendTouchIDs(g.touchIDs[:0])

	// Decrease the number of the sprites.
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || g.leftTouched() {
		g.sprites.num -= 20
		if g.sprites.num < MinSprites {
			g.sprites.num = MinSprites
		}
	}

	// Increase the number of the sprites.
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || g.rightTouched() {
		g.sprites.num += 20
		if MaxSprites < g.sprites.num {
			g.sprites.num = MaxSprites
		}
	}

	g.sprites.Update()
	return nil
}

func (g *SpritesGame) Draw(screen *ebiten.Image) {
	// Draw each sprite.
	// DrawImage can be called many many times, but in the implementation,
	// the actual draw call to GPU is very few since these calls satisfy
	// some conditions e.g. all the rendering sources and targets are same.
	// For more detail, see:
	// https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2#Image.DrawImage
	w, h := ebitenImage.Bounds().Dx(), ebitenImage.Bounds().Dy()
	for i := 0; i < g.sprites.num; i++ {
		s := g.sprites.sprites[i]
		g.op.GeoM.Reset()
		g.op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
		g.op.GeoM.Rotate(2 * math.Pi * float64(s.angle) / maxAngleSprites)
		g.op.GeoM.Translate(float64(w)/2, float64(h)/2)
		g.op.GeoM.Translate(float64(s.x), float64(s.y))
		screen.DrawImage(ebitenImage, &g.op)
	}
	msg := fmt.Sprintf(`TPS: %0.2f
FPS: %0.2f
Num of sprites: %d
Press <- or -> to change the number of sprites`, ebiten.ActualTPS(), ebiten.ActualFPS(), g.sprites.num)
	ebitenutil.DebugPrint(screen, msg)
}

func (g *SpritesGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidthSprites, screenHeightSprites
}

func ExampleSprites() {
	ebiten.SetWindowSize(screenWidthSprites*2, screenHeightSprites*2)
	ebiten.SetWindowTitle("Sprites (Ebitengine Demo)")
	ebiten.SetWindowResizable(true)
	if err := ebiten.RunGame(&SpritesGame{}); err != nil {
		log.Fatal(err)
	}
}
