package game

import (
	"game/assets"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	bulletSpeedPerSecond = 350.0
)

type Bullet struct {
	position Vector
	rotation float64
	sprite   *ebiten.Image
}

func NewBullet(pos Vector, rotation float64) *Bullet {
	sprite := assets.LaserSprite

	return &Bullet{
		position: pos,
		rotation: rotation,
		sprite:   sprite,
	}
}

func (b *Bullet) Update() {
	speed := bulletSpeedPerSecond / float64(ebiten.TPS())

	b.position.X += math.Sin(b.rotation) * speed
	b.position.Y += math.Cos(b.rotation) * -speed
}

func (b *Bullet) Draw(screen *ebiten.Image) {
	mp := MiddlePoint(b.sprite)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-mp.X, -mp.X)
	op.GeoM.Rotate(b.rotation)

	op.GeoM.Translate(b.position.X, b.position.Y)

	screen.DrawImage(b.sprite, op)
}

func (b *Bullet) Collider() Rect {
	bounds := b.sprite.Bounds()

	return NewRect(b.position.X-float64(bounds.Dx())/2, b.position.Y-float64(bounds.Dx())/2, float64(bounds.Dx()), float64(bounds.Dx()))
}

func (b *Bullet) DebugInfo(screen *ebiten.Image) {
	rect := b.Collider()
	vector.StrokeRect(
		screen,
		float32(rect.X),
		float32(rect.Y),
		float32(rect.Width),
		float32(rect.Height),
		1.0,
		color.White,
		false,
	)
}
