package game

import (
	"asteroids/assets"

	"github.com/hajimehoshi/ebiten/v2"
)

type Background struct {
	sprite *ebiten.Image

	x16 int
	y16 int
}

func NewBackground() *Background {
	sprite := assets.BackgroundSprite

	return &Background{
		sprite: sprite,
	}
}

func (b *Background) Update() {

	s := b.sprite.Bounds().Size()
	maxX16 := s.X * 16
	maxY16 := s.Y * 16

	b.x16 += s.X / 32
	b.y16 += s.Y / 32

	b.x16 %= maxX16
	b.y16 %= maxY16
}

func (b *Background) Draw(screen *ebiten.Image) {
	bounds := b.sprite.Bounds()

	offsetX, offsetY := float64(-b.x16)/16, float64(-b.y16)/16

	w := bounds.Dx()
	h := bounds.Dy()
	hr := ScreenWidth/w + 2
	vr := ScreenHeight/h + 2

	for i := 0; i < hr; i++ {
		for j := 0; j < vr; j++ {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(w*i), float64(h*j))
			op.GeoM.Translate(offsetX, offsetY)
			screen.DrawImage(b.sprite, op)
		}
	}
}
