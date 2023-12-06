package main

import (
	_ "image/gif"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/colornames"
)

var gopher *ebiten.Image

const (
	screenWidth  = 640
	screenHeight = 640
	gopherSpeed  = 12
)

type Game struct {
	gopherX, gopherY float64
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		g.gopherX -= gopherSpeed
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		g.gopherX += gopherSpeed
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		g.gopherY -= gopherSpeed
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		g.gopherY += gopherSpeed
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(colornames.Skyblue)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.gopherX, g.gopherY)
	screen.DrawImage(gopher, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	var err error
	gopher, _, err = ebitenutil.NewImageFromFile("gopher-dance-long-3x.gif")
	if err != nil {
		log.Fatal(err)
	}

	game := &Game{}
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("moving gopher")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
