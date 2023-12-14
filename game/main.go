package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	debug          = false
	screenX        = 1000
	screenY        = 600
	baseX          = 190
	groundY        = 110
	speed          = 12
	jumpingPower   = 8.5
	gravity        = 0.5
	interval       = 10
	minEbifryDist  = 50
	maxEbifryCount = 3
	modeFontSize   = 20
	msgFontSize    = 12

	// game modes
	modeTitle    = 0
	modeGame     = 1
	modeGameover = 2

	// image sizes
	gopherHeight = 25
	gopherWidth  = 50
	groundHeight = 50
	groundWidth  = 50

	gopherFrameOX     = 324
	gopherFrameWidth  = 12
	gopherFrameHeight = 14
	gopherFrameCount  = 27

	tileFrameOX     = 0
	tileFrameOY     = 32
	tileFrameWidth  = 32
	tileFrameHeight = 14
	tileFrameCount  = 27
)

//go:embed assets/sheet.png
var gopherImage []byte

//go:embed assets/ebiten.png
var ebitenImage []byte

//go:embed assets/tiles.png
var byteGroundImg []byte

var (
	gopherImg    *ebiten.Image
	ebitenImg    *ebiten.Image
	groundImg    *ebiten.Image
	groundSubImg *ebiten.Image
	modeFont     font.Face
	messageFont  font.Face
)

func init() {
	img, _, err := image.Decode(bytes.NewReader(gopherImage))
	if err != nil {
		log.Fatal(err)
	}
	gopherImg = ebiten.NewImageFromImage(img)

	img, _, err = image.Decode(bytes.NewReader(ebitenImage))
	if err != nil {
		log.Fatal(err)
	}
	ebitenImg = ebiten.NewImageFromImage(img)

	img, _, err = image.Decode(bytes.NewReader(byteGroundImg))
	if err != nil {
		log.Fatal(err)
	}
	groundImg = ebiten.NewImageFromImage(img)

	tt, err := opentype.Parse(fonts.PressStart2P_ttf)
	if err != nil {
		log.Fatal(err)
	}
	const dpi = 72
	modeFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    modeFontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	messageFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    msgFontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
}

type ebifry struct {
	count   int
	x       int
	y       float64
	visible bool
}

func (t *ebifry) move(speed int) {
	t.x += speed
}

func (t *ebifry) show() {
	t.x = 0
	t.y = groundY
	t.visible = true
}

func (t *ebifry) hide() {
	t.visible = false
}

func (t *ebifry) isOutOfScreen() bool {
	return t.x > 3000
}

type ground struct {
	count int
	x     int
	y     int
}

func (g *ground) move(speed int) {
	g.x += speed
	if g.x > -groundWidth {
		g.x = g.x - groundWidth
	}
}

type Game struct {
	mode        int
	count       int
	score       int
	hiscore     int
	gopherX     int
	gopherY     float64
	gy          float64
	jumpFlg     bool
	ebifrys     [maxEbifryCount]*ebifry
	lastEbifryX int
	ground      *ground
}

func NewGame() *Game {
	g := &Game{}
	g.init()
	return g
}

func (g *Game) init() {
	g.hiscore = g.score
	g.count = 0
	g.score = 0
	g.lastEbifryX = 0
	g.gy = 0
	g.gopherX = baseX
	g.gopherY = groundY - gopherHeight
	for i := 0; i < maxEbifryCount; i++ {
		g.ebifrys[i] = &ebifry{}
	}
	g.ground = &ground{y: groundY - 30}
}

func (g *Game) Update() error {
	switch g.mode {
	case modeTitle:
		if g.isKeyJustPressed() {
			g.mode = modeGame
		}
	case modeGame:
		g.count++
		g.score = g.count / 5

		g.ground.count++

		if !g.jumpFlg && g.isKeyJustPressed() {
			g.jumpFlg = true
			g.gy = -jumpingPower
		}

		if g.jumpFlg {
			g.gopherY += g.gy
			g.gy += gravity
		}

		if g.gopherY >= groundY-gopherHeight {
			g.jumpFlg = false
		}

		for _, t := range g.ebifrys {
			if t.visible {
				t.move(speed)
				if t.isOutOfScreen() {
					t.hide()
				}
			} else {
				if g.ground.count-g.lastEbifryX > minEbifryDist && g.ground.count%interval == 0 && rand.Intn(10) == 0 {
					g.lastEbifryX = g.ground.count
					t.show()
					break
				}
			}
		}

		g.ground.move(speed)

		if g.hit() {
			g.mode = modeGameover
		}
	case modeGameover:
		if g.isKeyJustPressed() {
			g.init()
			g.mode = modeGame
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
	var xs [3]int
	var ys [3]float64

	if len(g.ebifrys) > 0 {
		for i, t := range g.ebifrys {
			xs[i] = t.x
			ys[i] = t.y
		}
	}

	if debug {
		ebitenutil.DebugPrint(screen, fmt.Sprintf(
			"groundY: %d\ngopherWidth: %d, gopherHeight: %d\ng.gopherX: %d, g.gopherY: %d\nEbiten1 x:%d, y:%d\nEbiten2 x:%d, y:%d\nEbiten3 x:%d, y:%d",
			groundY,
			gopherWidth,
			gopherHeight,
			g.gopherX,
			g.gopherY,
			xs[0],
			ys[0],
			xs[1],
			ys[1],
			xs[2],
			ys[2],
		))
	}

	// g.drawGround(screen)
	g.drawEbifrys(screen)
	g.drawGopher(screen)

	switch g.mode {
	case modeTitle:
		text.Draw(screen, "PRESS SPACE KEY", modeFont, 360, 240, color.RGBA{0, 100, 0, 100})
	case modeGameover:
		text.Draw(screen, "GAME OVER", modeFont, 410, 240, color.RGBA{255, 0, 0, 100})
		text.Draw(screen, "To play again, press the space key.", messageFont, 300, 260, color.RGBA{0, 0, 0, 100})
	}
}

func (g *Game) drawGopher(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(baseX, float64(g.gopherY))
	op.GeoM.Scale(4.0, 4.0)
	op.Filter = ebiten.FilterNearest
	i := (g.count / 5) % gopherFrameCount
	sx := gopherFrameOX - i*gopherFrameWidth
	if sx <= 96 {
		g.count = 0
	}
	screen.DrawImage(gopherImg.SubImage(image.Rect(sx, 0, sx-gopherFrameWidth, 14)).(*ebiten.Image), op)
}

func (g *Game) drawEbifrys(screen *ebiten.Image) {
	for _, t := range g.ebifrys {
		if t.visible {
			op := &ebiten.DrawImageOptions{}
			t.y = 900
			op.GeoM.Translate(float64(t.x), float64(t.y))
			op.GeoM.Scale(0.35, 0.35)
			op.Filter = ebiten.FilterLinear
			screen.DrawImage(ebitenImg, op)
		}
	}
}

func (g *Game) drawGround(screen *ebiten.Image) {
	for i := 0; i < 14; i++ {
		x := float64(groundWidth * i)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(x, float64(g.ground.y))
		op.GeoM.Translate(float64(g.ground.x), 0.0)
		op.Filter = ebiten.FilterLinear

		i := (g.ground.count / 5) % tileFrameCount
		sx := tileFrameOX + i*tileFrameWidth
		if sx > 32 {
			g.ground.count = 0
		}
		subImgRect := image.Rect(sx, 0, sx+tileFrameWidth, 32)
		groundSubImg = groundImg.SubImage(subImgRect).(*ebiten.Image)
		screen.DrawImage(groundImg, op)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return screenX, screenY
}

func (g *Game) isKeyJustPressed() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		return true
	}
	return false
}

func (g *Game) hit() bool {
	for _, t := range g.ebifrys {
		if t.visible {
			if t.x > 2000 && t.x < 2200 && 340 == g.gopherY*4 {
				return true
			}
		}
	}
	return false
}

func main() {
	ebiten.SetWindowSize(screenX, screenY)
	ebiten.SetWindowTitle("Gopher Jump")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
