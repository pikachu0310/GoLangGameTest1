package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

var (
	screenWidth  = 320
	screenHeight = 240
)

func main() {
	a := app.New()
	w := a.NewWindow("My Game")

	img, _, err := ebitenutil.NewImageFromFile("image.png", ebiten.FilterDefault)
	if err != nil {
		panic(err)
	}

	imgWidget := widget.NewIcon(ebiten.NewImageFromImage(img))
	content := fyne.NewContainerWithLayout(layout.NewGridLayout(1), imgWidget)
	content.Refresh()

	w.Resize(fyne.NewSize(screenWidth, screenHeight))
	w.SetPadded(false)

	innerContent := fyne.NewContainerWithLayout(layout.NewCenterLayout(), content)
	innerContent.Resize(content.MinSize())
	w.SetContent(innerContent)

	w.ShowAndRun()
}
