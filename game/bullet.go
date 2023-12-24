package game

import (
	"game/assets"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
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

	mp := MiddlePoint(sprite)
	pos.X -= mp.X
	pos.Y -= mp.Y

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
	op.GeoM.Translate(-mp.X, -mp.Y)
	op.GeoM.Rotate(b.rotation)
	op.GeoM.Translate(mp.X, mp.Y)

	op.GeoM.Translate(b.position.X, b.position.Y)

	screen.DrawImage(b.sprite, op)
}

func (b *Bullet) Collider() Rect {
	bounds := b.sprite.Bounds()

	return NewRect(b.position.X, b.position.Y, float64(bounds.Dx()), float64(bounds.Dy()))
}
