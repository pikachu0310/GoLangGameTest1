package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.design/x/clipboard"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type GameState int

const (
	Playing GameState = iota
	Title
	Inventory
)

var slime struct {
	Hp      int
	Attack  int
	Defense int
}

type Player struct {
	HP        int
	MaxHP     int
	HPRecover int
	Attack    int
	Defense   int
	Inventory []Item
	DaysLeft  int
}

type Enemy struct {
	Name   string
	HP     int
	Attack int
}

type Item struct {
	Name          string
	Category      string
	MaxHp         int
	InstantHeal   int
	SustainedHeal int
	Attack        int
	Defense       int
	Description   string
	*Button
	*CheckBox
}

func generateItem() []*Item {
	items, err := GptGenerateItem()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Item:%+v\n", items)
	}
	return items
}

func combineItem(items []*Item) []*Item {
	items, err := GptCombineItem(items)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Item:%+v\n", items)
	}
	return items
}

var (
	uiImage       *ebiten.Image
	uiFont        font.Face
	uiFontMHeight int
	slimeImage    *ebiten.Image
)

const (
	lineHeight = 16 + 8 // 16
)

func init() {
	// Decode an image from the image file's byte slice.
	img, _, err := image.Decode(bytes.NewReader(images.UI_png))
	if err != nil {
		log.Fatal(err)
	}
	uiImage = ebiten.NewImageFromImage(img)

	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}
	uiFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    12,
		DPI:     72 + 36,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}
	b, _, _ := uiFont.GlyphBounds('M')
	uiFontMHeight = (b.Max.Y - b.Min.Y).Ceil()

	// img, _, err = image.Decode(bytes.NewReader(myimages.Slime_png))
	// if err != nil {
	//	log.Fatal(err)
	// }
	// slimeImage = ebiten.NewImageFromImage(img)
}

type imageType int

const (
	imageTypeButton imageType = iota
	imageTypeButtonPressed
	imageTypeTextBox
	imageTypeVScrollBarBack
	imageTypeVScrollBarFront
	imageTypeCheckBox
	imageTypeCheckBoxPressed
	imageTypeCheckBoxMark
)

var imageSrcRects = map[imageType]image.Rectangle{
	imageTypeButton:          image.Rect(0, 0, 16, 16),
	imageTypeButtonPressed:   image.Rect(16, 0, 32, 16),
	imageTypeTextBox:         image.Rect(0, 16, 16, 32),
	imageTypeVScrollBarBack:  image.Rect(16, 16, 24, 32),
	imageTypeVScrollBarFront: image.Rect(24, 16, 32, 32),
	imageTypeCheckBox:        image.Rect(0, 32, 16, 48),
	imageTypeCheckBoxPressed: image.Rect(16, 32, 32, 48),
	imageTypeCheckBoxMark:    image.Rect(32, 32, 48, 48),
}

const (
	screenWidth  = 1280 + 10
	screenHeight = 720 + 10
)

type Input struct {
	mouseButtonState int
}

func drawNinePatches(dst *ebiten.Image, dstRect image.Rectangle, srcRect image.Rectangle) {
	srcX := srcRect.Min.X
	srcY := srcRect.Min.Y
	srcW := srcRect.Dx()
	srcH := srcRect.Dy()

	dstX := dstRect.Min.X
	dstY := dstRect.Min.Y
	dstW := dstRect.Dx()
	dstH := dstRect.Dy()

	op := &ebiten.DrawImageOptions{}
	for j := 0; j < 3; j++ {
		for i := 0; i < 3; i++ {
			op.GeoM.Reset()

			sx := srcX
			sy := srcY
			sw := srcW / 4
			sh := srcH / 4
			dx := 0
			dy := 0
			dw := sw
			dh := sh
			switch i {
			case 1:
				sx = srcX + srcW/4
				sw = srcW / 2
				dx = srcW / 4
				dw = dstW - 2*srcW/4
			case 2:
				sx = srcX + 3*srcW/4
				dx = dstW - srcW/4
			}
			switch j {
			case 1:
				sy = srcY + srcH/4
				sh = srcH / 2
				dy = srcH / 4
				dh = dstH - 2*srcH/4
			case 2:
				sy = srcY + 3*srcH/4
				dy = dstH - srcH/4
			}

			op.GeoM.Scale(float64(dw)/float64(sw), float64(dh)/float64(sh))
			op.GeoM.Translate(float64(dx), float64(dy))
			op.GeoM.Translate(float64(dstX), float64(dstY))
			dst.DrawImage(uiImage.SubImage(image.Rect(sx, sy, sx+sw, sy+sh)).(*ebiten.Image), op)
		}
	}
}

type Button struct {
	Rect image.Rectangle
	Text string

	mouseDown bool

	onPressed func(b *Button)
	onCursor  func(b *Button)
}

func (b *Button) Update() {
	x, y := ebiten.CursorPosition()
	if b.Rect.Min.X <= x && x < b.Rect.Max.X && b.Rect.Min.Y <= y && y < b.Rect.Max.Y {
		if b.onCursor != nil {
			b.onCursor(b)
		}
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		// x, y := ebiten.CursorPosition()
		if b.Rect.Min.X <= x && x < b.Rect.Max.X && b.Rect.Min.Y <= y && y < b.Rect.Max.Y {
			b.mouseDown = true
		} else {
			b.mouseDown = false
		}
	} else {
		if b.mouseDown {
			if b.onPressed != nil {
				b.onPressed(b)
			}
		}
		b.mouseDown = false
	}
}

func (b *Button) Draw(dst *ebiten.Image) {
	t := imageTypeButton
	if b.mouseDown {
		t = imageTypeButtonPressed
	}
	drawNinePatches(dst, b.Rect, imageSrcRects[t])

	bounds, _ := font.BoundString(uiFont, b.Text)
	w := (bounds.Max.X - bounds.Min.X).Ceil()
	x := b.Rect.Min.X + (b.Rect.Dx()-w)/2
	y := b.Rect.Max.Y - (b.Rect.Dy()-uiFontMHeight)/2
	text.Draw(dst, b.Text, uiFont, x, y, color.Black)
}

func (b *Button) SetOnPressed(f func(b *Button)) {
	b.onPressed = f
}

func (b *Button) SetOnCursor(f func(b *Button)) {
	b.onCursor = f
}

const VScrollBarWidth = 16

type VScrollBar struct {
	X      int
	Y      int
	Height int

	thumbRate           float64
	thumbOffset         int
	dragging            bool
	draggingStartOffset int
	draggingStartY      int
	contentOffset       int
}

func (v *VScrollBar) thumbSize() int {
	const minThumbSize = VScrollBarWidth

	r := v.thumbRate
	if r > 1 {
		r = 1
	}
	s := int(float64(v.Height) * r)
	if s < minThumbSize {
		return minThumbSize
	}
	return s
}

func (v *VScrollBar) thumbRect() image.Rectangle {
	if v.thumbRate >= 1 {
		return image.Rectangle{}
	}

	s := v.thumbSize()
	return image.Rect(v.X, v.Y+v.thumbOffset, v.X+VScrollBarWidth, v.Y+v.thumbOffset+s)
}

func (v *VScrollBar) maxThumbOffset() int {
	return v.Height - v.thumbSize()
}

func (v *VScrollBar) ContentOffset() int {
	return v.contentOffset
}

func (v *VScrollBar) Update(contentHeight int) {
	v.thumbRate = float64(v.Height) / float64(contentHeight)

	if !v.dragging && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		tr := v.thumbRect()
		if tr.Min.X <= x && x < tr.Max.X && tr.Min.Y <= y && y < tr.Max.Y {
			v.dragging = true
			v.draggingStartOffset = v.thumbOffset
			v.draggingStartY = y
		}
	}
	if v.dragging {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			_, y := ebiten.CursorPosition()
			v.thumbOffset = v.draggingStartOffset + (y - v.draggingStartY)
			if v.thumbOffset < 0 {
				v.thumbOffset = 0
			}
			if v.thumbOffset > v.maxThumbOffset() {
				v.thumbOffset = v.maxThumbOffset()
			}
		} else {
			v.dragging = false
		}
	}

	v.contentOffset = 0
	if v.thumbRate < 1 {
		v.contentOffset = int(float64(contentHeight) * float64(v.thumbOffset) / float64(v.Height))
	}
}

func (v *VScrollBar) Draw(dst *ebiten.Image) {
	sd := image.Rect(v.X, v.Y, v.X+VScrollBarWidth, v.Y+v.Height)
	drawNinePatches(dst, sd, imageSrcRects[imageTypeVScrollBarBack])

	if v.thumbRate < 1 {
		drawNinePatches(dst, v.thumbRect(), imageSrcRects[imageTypeVScrollBarFront])
	}
}

const (
	textBoxPaddingLeft = 8
)

type TextBox struct {
	Rect image.Rectangle
	Text string

	contentBuf *ebiten.Image
	vScrollBar *VScrollBar
	offsetX    int
	offsetY    int
}

func (t *TextBox) AppendLine(line string) {
	if t.Text == "" {
		t.Text = line
	} else {
		t.Text += "\n" + line
	}
}

func (t *TextBox) AppendLineToFirst(line string) {
	if t.Text == "" {
		t.Text = line
	} else {
		t.Text = line + "\n" + t.Text
	}
}

func (t *TextBox) Update() {
	if t.vScrollBar == nil {
		t.vScrollBar = &VScrollBar{}
	}
	t.vScrollBar.X = t.Rect.Max.X - VScrollBarWidth
	t.vScrollBar.Y = t.Rect.Min.Y
	t.vScrollBar.Height = t.Rect.Dy()

	_, h := t.contentSize()
	t.vScrollBar.Update(h)

	t.offsetX = 0
	t.offsetY = t.vScrollBar.ContentOffset()
}

func (t *TextBox) contentSize() (int, int) {
	h := len(strings.Split(t.Text, "\n")) * lineHeight
	return t.Rect.Dx(), h
}

func (t *TextBox) viewSize() (int, int) {
	return t.Rect.Dx() - VScrollBarWidth - textBoxPaddingLeft, t.Rect.Dy()
}

func (t *TextBox) contentOffset() (int, int) {
	return t.offsetX, t.offsetY
}

func (t *TextBox) Draw(dst *ebiten.Image) {
	drawNinePatches(dst, t.Rect, imageSrcRects[imageTypeTextBox])

	if t.contentBuf != nil {
		vw, vh := t.viewSize()
		w, h := t.contentBuf.Bounds().Dx(), t.contentBuf.Bounds().Dy()
		if vw > w || vh > h {
			t.contentBuf.Dispose()
			t.contentBuf = nil
		}
	}
	if t.contentBuf == nil {
		w, h := t.viewSize()
		t.contentBuf = ebiten.NewImage(w, h)
	}

	t.contentBuf.Clear()
	for i, line := range strings.Split(t.Text, "\n") {
		x := -t.offsetX + textBoxPaddingLeft
		y := -t.offsetY + i*lineHeight + lineHeight - (lineHeight-uiFontMHeight)/2
		if y < -lineHeight {
			continue
		}
		if _, h := t.viewSize(); y >= h+lineHeight {
			continue
		}
		text.Draw(t.contentBuf, line, uiFont, x, y, color.Black)
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(t.Rect.Min.X), float64(t.Rect.Min.Y))
	dst.DrawImage(t.contentBuf, op)

	t.vScrollBar.Draw(dst)
}

const (
	checkboxWidth       = 16
	checkboxHeight      = 16
	checkboxPaddingLeft = 8
)

type CheckBox struct {
	X    int
	Y    int
	Text string

	checked   bool
	mouseDown bool

	onCheckChanged func(c *CheckBox)
}

func (c *CheckBox) width() int {
	b, _ := font.BoundString(uiFont, c.Text)
	w := (b.Max.X - b.Min.X).Ceil()
	return checkboxWidth + checkboxPaddingLeft + w
}

func (c *CheckBox) Update() {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if c.X <= x && x < c.X+c.width() && c.Y <= y && y < c.Y+checkboxHeight {
			c.mouseDown = true
		} else {
			c.mouseDown = false
		}
	} else {
		if c.mouseDown {
			c.checked = !c.checked
			if c.onCheckChanged != nil {
				c.onCheckChanged(c)
			}
		}
		c.mouseDown = false
	}
}

func (c *CheckBox) Draw(dst *ebiten.Image) {
	t := imageTypeCheckBox
	if c.mouseDown {
		t = imageTypeCheckBoxPressed
	}
	r := image.Rect(c.X, c.Y, c.X+checkboxWidth, c.Y+checkboxHeight)
	drawNinePatches(dst, r, imageSrcRects[t])
	if c.checked {
		drawNinePatches(dst, r, imageSrcRects[imageTypeCheckBoxMark])
	}

	x := c.X + checkboxWidth + checkboxPaddingLeft
	y := (c.Y + 16) - (16-uiFontMHeight)/2
	text.Draw(dst, c.Text, uiFont, x, y, color.Black)
}

func (c *CheckBox) Checked() bool {
	return c.checked
}

func (c *CheckBox) SetOnCheckChanged(f func(c *CheckBox)) {
	c.onCheckChanged = f
}

// My Func Start
func FormatItemsGUI(items []*Item) {
	for i := 0; i < len(items); i++ {
		if i <= 14 {
			items[i].Button.Rect = image.Rect(16*50, 16*(1+3*(i))-4, 16*64, 16*(3+3*(i))+4)
			items[i].CheckBox.X = 16*48 + 12
			items[i].CheckBox.Y = 16*(1+3*(i)) + 8
		} else {
			items[i].Button.Rect = image.Rect(16*66, 16*(1+3*(i-15))-4, 16*80, 16*(3+3*(i-15))+4)
			items[i].CheckBox.X = 16*64 + 12
			items[i].CheckBox.Y = 16*(1+3*(i-15)) + 8
		}
	}
}

func itemStringer(item *Item, interval int) string {
	// return fmt.Sprintf("Item Info\nName: %s\nCategory: %s\nMaxHp: %d\nInstantHeal: %d\nSustainedHeal: %d\nAttck: %d\nDefence: %d\nDescription: \n%s\n", item.Name, item.Category, item.MaxHp, item.InstantHeal, item.SustainedHeal, item.Attack, item.Defense, addNewLineItem(item.Description, 60))
	return fmt.Sprintf("アイテム情報\n名前　　: %s\n種類　　: %s\n最大体力: %d\n即時回復: %d\n持続回復: %d\n攻撃力　: %d\n防御力　: %d\n説明文　:%s\n", item.Name, item.Category, item.MaxHp, item.InstantHeal, item.SustainedHeal, item.Attack, item.Defense, addNewLineItem(item.Description, interval))
}

func intervalStringer(text string, interval int) string {
	// return fmt.Sprintf("Item Info\nName: %s\nCategory: %s\nMaxHp: %d\nInstantHeal: %d\nSustainedHeal: %d\nAttck: %d\nDefence: %d\nDescription: \n%s\n", item.Name, item.Category, item.MaxHp, item.InstantHeal, item.SustainedHeal, item.Attack, item.Defense, addNewLineItem(item.Description, 60))
	return fmt.Sprintf(addNewLine(text, interval))
}

func (g *Game) AddItem(item *Item) {
	if item == nil {
		fmt.Println("item is nil")
		return
	}
	g.items = append(g.items, item)

	item.Button = &Button{
		Text: item.Name,
	}
	item.Button.SetOnPressed(func(b *Button) {
		g.textBoxLog.Text = itemStringer(item, 25)
		item.checked = !item.checked
		g.CheckCheckedItem(item)
	})

	item.Button.SetOnCursor(func(b *Button) {
		g.textBoxLog.Text = itemStringer(item, 25)
	})

	item.CheckBox = &CheckBox{
		Text: "",
	}
	item.CheckBox.SetOnCheckChanged(func(c *CheckBox) {
		g.CheckCheckedItem(item)
	})

	FormatItemsGUI(g.items)
}

func (g *Game) DeleteItem(item *Item) {
	for i := 0; i < len(g.items); i++ {
		if g.items[i] == item {
			g.items = append(g.items[:i], g.items[i+1:]...)
			break
		}
	}
	FormatItemsGUI(g.items)
}

func (g *Game) DeleteItems(items []*Item) {
	for i := 0; i < len(items); i++ {
		g.DeleteItem(items[i])
		g.DeleteCheckedItem(items[i])
	}
}

func (g *Game) ResetItems(items []*Item) {
	g.items = []*Item{}
	for i := 0; i < len(items); i++ {
		g.AddItem(items[i])
	}
	g.checkedItems = []*Item{}
	items2 := getCheckedItems(g.items)
	for _, item := range items2 {
		g.checkedItems = append(g.checkedItems, item)
	}
}

func getCheckedItems(items []*Item) []*Item {
	var checkedItems []*Item
	for i := 0; i < len(items); i++ {
		if items[i].CheckBox.Checked() {
			checkedItems = append(checkedItems, items[i])
		}
	}
	return checkedItems
}

func (g *Game) AddCheckedItem(item *Item) {
	g.checkedItems = append(g.checkedItems, item)
}

func (g *Game) DeleteCheckedItem(item *Item) {
	for i := 0; i < len(g.checkedItems); i++ {
		if g.checkedItems[i] == item {
			g.checkedItems = append(g.checkedItems[:i], g.checkedItems[i+1:]...)
			break
		}
	}
}

func addNewLineItem(str string, interval int) string {
	strSlice := strings.Split(str, "")
	var result string
	for i, s := range strSlice {
		if i%interval == 0 && i != 0 {
			result += "\n　　　　 "
		}
		result += s
	}
	return result
}

func addNewLine(str string, interval int) string {
	strSlice := strings.Split(str, "")
	var result string
	for i, s := range strSlice {
		if i%interval == 0 && i != 0 {
			result += "\n"
		}
		result += s
	}
	return result
}

func (g *Game) SaveItems() {
	err := writeToFile(g.items, "SaveDataAuto.gob")
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to save items: %v", err))
	} else {
		fmt.Println(fmt.Sprintf("Saved %d items Auto", len(g.items)))
	}
}

func (g *Game) LoadItems() {
	items, err := readFromFile("SaveDataAuto.gob")
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to load items: %v", err))
		return
	} else {
		fmt.Println(fmt.Sprintf("Loaded %d items Auto", len(g.items)))
	}
	g.ResetItems(items)
}

func (g *Game) SaveItemsManual() {
	err := writeToFile(g.items, "SaveDataManual.gob")
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to save items: %v", err))
	} else {
		fmt.Println(fmt.Sprintf("Saved %d items Manual", len(g.items)))
		g.textBoxLog5.AppendLineToFirst(intervalStringer(fmt.Sprintf(fmt.Sprintf("Manual Saved %d items", len(g.items))), 12))
	}
}

func (g *Game) LoadItemsManual() {
	items, err := readFromFile("SaveDataManual.gob")
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to load items: %v", err))
		return
	} else {
		fmt.Println(fmt.Sprintf("Loaded %d items Manual", len(g.items)))
		g.textBoxLog5.AppendLineToFirst(intervalStringer(fmt.Sprintf(fmt.Sprintf("Manual Loaded %d items", len(g.items))), 12))
	}
	g.ResetItems(items)
}

func writeToFile(items []*Item, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	encoder := gob.NewEncoder(writer)
	if err := encoder.Encode(items); err != nil {
		return err
	}
	return nil
}

func readFromFile(filePath string) ([]*Item, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	decoder := gob.NewDecoder(reader)

	var items []*Item
	if err := decoder.Decode(&items); err != nil {
		return nil, err
	}
	return items, nil
}

func (g *Game) AddNewButton(x0, y0, x1, y1 int, text string, f func(b *Button)) {
	newButton := &Button{
		Rect:      image.Rect(x0, y0, x1, y1),
		Text:      text,
		onPressed: f,
	}
	g.buttons = append(g.buttons, newButton)
}

func (g *Game) CheckCombineTextBoxLog() {
	fmt.Println(g.checkedItems)
	if len(g.checkedItems) <= 0 {
		return
	}
	g.textBoxLog2.Text = itemStringer(g.checkedItems[0], 7)
	g.textBoxLog3.Text = itemStringer(g.checkedItems[len(g.checkedItems)-1], 7)
}

func (g *Game) CheckCheckedItem(item *Item) {
	if item.checked {
		g.AddCheckedItem(item)
	} else {
		g.DeleteCheckedItem(item)
	}
	g.CheckCombineTextBoxLog()
}

func CopyToClipboard(text string) {
	clipboard.Write(clipboard.FmtText, []byte(text))
}

func (g *Game) GenerateItem() {
	items := generateItem()
	if items == nil {
		fmt.Println("item is nil")
		g.textBoxLog5.AppendLineToFirst(intervalStringer(fmt.Sprintf("Error: item is nil (%d)", g.generating), 12))
		g.generating -= 1
		return
	}
	for _, item := range items {
		if item == nil {
			fmt.Println("item is nil")
			g.textBoxLog5.AppendLineToFirst(intervalStringer(fmt.Sprintf("Error: item is nil (%d)", g.generating), 12))
			g.generating -= 1
			return
		}
		g.AddItem(item)
	}
	g.SaveItems()
	g.textBoxLog5.AppendLineToFirst(intervalStringer(fmt.Sprintf("Item generating is finished. (%d)", g.generating), 12))
	g.generating -= 1
}

func (g *Game) CombineItem() {
	if len(g.checkedItems) <= 1 {
		g.combining -= 1
		g.textBoxLog5.AppendLineToFirst(intervalStringer(fmt.Sprintf("Please select two or more items. (%d)", g.combining), 12))
		return
	}
	items := combineItem(g.checkedItems)

	if items == nil {
		fmt.Println("item is nil")
		g.textBoxLog5.AppendLineToFirst(intervalStringer(fmt.Sprintf("Error: item is nil (%d)", g.combining), 12))
		g.combining -= 1
		return
	}
	for _, item := range items {
		if item == nil {
			fmt.Println("item is nil")
			g.textBoxLog5.AppendLineToFirst(intervalStringer(fmt.Sprintf("Error: item is nil (%d)", g.combining), 12))
			g.combining -= 1
			return
		}
		g.AddItem(item)
		g.textBoxLog4.Text = itemStringer(item, 7)
	}
// My Func End

type Game struct {
	buttons      []*Button
	checkBox     *CheckBox
	textBoxLog   *TextBox
	textBoxLog2  *TextBox
	textBoxLog3  *TextBox
	textBoxLog4  *TextBox
	textBoxLog5  *TextBox
	items        []*Item
	checkedItems []*Item
	slime        *Button
	Player       Player
	Enemy        Enemy
	GameState    GameState
	generating   int
	combining    int
}

func GameMain() *Game {
	g := &Game{}
	g.LoadItems()
	g.AddNewButton(16*40, 16*2, 16*48, 16*6, "Generate", func(b *Button) {
		if g.generating >= 3 {
			g.textBoxLog5.AppendLineToFirst(intervalStringer("Generating is full! (max:3)", 12))
			return
		}
		g.generating += 1
		g.textBoxLog5.AppendLineToFirst(intervalStringer(fmt.Sprintf("Item generating in progress... (%d)", g.generating), 12))
		go g.GenerateItem()
	})
	g.AddNewButton(16*40, 16*7, 16*48, 16*11, "Combine", func(b *Button) {
		if g.combining >= 1 {
			g.textBoxLog5.AppendLineToFirst(intervalStringer("Combining is full! (max:1)", 12))
			return
		}
		g.combining += 1
		g.textBoxLog5.AppendLineToFirst(intervalStringer(fmt.Sprintf("Item combining in progress... (%d)", g.combining), 12))
		go g.CombineItem()
	})
	g.AddNewButton(16*2, 16*1, 16*10, 16*3, "Save", func(b *Button) {
		g.SaveItemsManual()
	})
	g.AddNewButton(16*11, 16*1, 16*19, 16*3, "Load", func(b *Button) {
		g.LoadItemsManual()
	})
	g.AddNewButton(16*20, 16*1, 16*38, 16*3, "左上をクリップボードにコピー", func(b *Button) {
		CopyToClipboard(g.textBoxLog.Text)
		g.textBoxLog5.AppendLineToFirst(intervalStringer(fmt.Sprintf("Copied!"), 12))
	})

	g.checkBox = &CheckBox{
		X:    16,
		Y:    64,
		Text: "Check Box!",
	}
	g.textBoxLog = &TextBox{
		Rect: image.Rect(16*1, 16*4, 16*39, 16*24),
	}
	g.textBoxLog2 = &TextBox{
		Rect: image.Rect(16*1, 16*25, 16*16, 16*45),
	}
	g.textBoxLog3 = &TextBox{
		Rect: image.Rect(16*17, 16*25, 16*32, 16*45),
	}
	g.textBoxLog4 = &TextBox{
		Rect: image.Rect(16*33, 16*25, 16*48, 16*45),
	}
	g.textBoxLog5 = &TextBox{
		Rect: image.Rect(16*39+8, 16*12, 16*48, 16*24),
	}

	g.slime = &Button{
		Rect: image.Rect(16, 480, 144, 512),
	}

	// g.checkBox.SetOnCheckChanged(func(c *CheckBox) {
	// 	msg := "Check box check changed"
	// 	if c.Checked() {
	// 		msg += " (Checked)"
	// 		g.SaveItems()
	// 	} else {
	// 		msg += " (Unchecked)"
	// 		g.LoadItems()
	// 	}
	// 	g.textBoxLog.AppendLine(msg)
	// })
	g.textBoxLog5.AppendLineToFirst(intervalStringer(fmt.Sprintf(fmt.Sprintf("Auto Loaded %d items", len(g.items))), 12))
	return g
}

func (g *Game) Update() error {
	for _, button := range g.buttons {
		button.Update()
	}
	for _, item := range g.items {
		item.Button.Update()
		item.CheckBox.Update()
	}
	g.checkBox.Update()
	g.textBoxLog.Update()
	g.textBoxLog2.Update()
	g.textBoxLog3.Update()
	g.textBoxLog4.Update()
	g.textBoxLog5.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0xeb, 0xeb, 0xeb, 0xff})
	for _, button := range g.buttons {
		button.Draw(screen)
	}
	for _, item := range g.items {
		item.Button.Draw(screen)
		item.CheckBox.Draw(screen)
	}
	// g.checkBox.Draw(screen)
	g.textBoxLog.Draw(screen)
	g.textBoxLog2.Draw(screen)
	g.textBoxLog3.Draw(screen)
	g.textBoxLog4.Draw(screen)
	g.textBoxLog5.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("変幻自在")
	if err := ebiten.RunGame(GameMain()); err != nil {
		log.Fatal(err)
	}
}
