package main

import (
	"bytes"
	"github.com/pikachu0310/GoLangGameTest1.git/myimages"
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
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

type Item struct {
	Name      string
	Category  string
	MaxHP     int
	HPRecover int
	Attack    int
	Defense   int
}

type Enemy struct {
	Name   string
	HP     int
	Attack int
}

var img *ebiten.Image

func generateEnemy() Enemy {
	gpt3 := openai.GetGpt3Model()
	enemyName, _ := gpt3.GenerateText()
	enemy := Enemy{
		Name:   enemyName,
		HP:     rand.Intn(10) + 1,
		Attack: rand.Intn(3) + 1,
	}
	return enemy
}

func generateItem() Item {
	gpt3 := openai.GetGpt3Model()
	itemName, _ := gpt3.GenerateText()
	item := Item{
		Name:      itemName,
		Category:  "Consumable",
		MaxHP:     rand.Intn(5) - 5,
		HPRecover: rand.Intn(5) - 5,
		Attack:    rand.Intn(5) - 5,
		Defense:   rand.Intn(5) - 5,
	}
	if rand.Float64() < 0.1 {
		item.Category = "Weapon"
		item.Attack += rand.Intn(5) + 1
	} else if rand.Float64() < 0.1 {
		item.Category = "Armor"
		item.Defense += rand.Intn(5) + 1
	}
	return item
}

func combineItems(item1, item2 Item) Item {
	gpt3 := openai.GetGpt3Model()
	itemName, _ := gpt3.GenerateText()
	item := Item{
		Name:      itemName,
		Category:  "Consumable",
		MaxHP:     (item1.MaxHP + item2.MaxHP) / 2,
		HPRecover: (item1.HPRecover + item2.HPRecover) / 2,
		Attack:    (item1.Attack + item2.Attack) / 2,
		Defense:   (item1.Defense + item2.Defense) / 2,
	}
	if rand.Float64() < 0.1 {
		item.Category = "Weapon"
		item.Attack += rand.Intn(5) + 1
	} else if rand.Float64() < 0.1 {
		item.Category = "Armor"
		item.Defense += rand.Intn(5) + 1
	}
	return item
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
	var gameState GameState = Playing

	switch gameState {
	case Playing:

	}

	//if gameState == "Exploration" {
	//	if player.HP <= 0 {
	//		gameState = "GameOver"
	//	} else if player.DaysLeft == 0 {
	//		gameState = "BossFight"
	//	} else {
	//		action := getAction()
	//		switch action {
	//		case "Fight":
	//			enemy := generateEnemy()
	//			result := fight(enemy)
	//			if result {
	//				item := generateItem()
	//				player.Inventory = append(player.Inventory, item)
	//			} else {
	//				player.HP -= enemy.Attack - player.Defense
	//			}
	//		case "Rest":
	//			player.HP += player.HPRecover
	//			if player.HP > player.MaxHP {
	//				player.HP = player.MaxHP
	//			}
	//		case "Inventory":
	//			showInventory()
	//		case "Combine":
	//			combineItemsPrompt()
	//		case "Exit":
	//			gameState = "GameOver"
	//		}
	//		player.DaysLeft--
	//	}
	//} else if gameState == "BossFight" {
	//	if player.HP <= 0 {
	//		gameState = "GameOver"
	//	} else {
	//		boss := generateBoss()
	//		result := bossFight(boss)
	//		if result {
	//			gameState = "Win"
	//		} else {
	//			player.HP -= boss.Attack - player.Defense
	//		}
	//	}
	//}
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
