package game

import (
	"game/assets"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

const (
	rotationSpeedMin = -0.02
	rotationSpeedMax = 0.02
)

type Meteor struct {
	position       Vector
	rotation       float64
	movement       Vector
	rotationSpeed  float64
	sprite         *ebiten.Image
	explosionAudio *audio.Player
}

func NewMeteor(baseVelocity float64) *Meteor {
	sprite := assets.MeteorSprites[rand.Intn(len(assets.MeteorSprites))]

	// Figure out the target position - the screen center, in this case
	target := Vector{
		X: ScreenWidth / 2,
		Y: ScreenHeight / 2,
	}

	// Pick a random angle - 2pi is 360 - so this return 0 to 360
	angle := rand.Float64() * 2 * math.Pi

	// The distance from the center, the meteor should spawn at - half  the width
	r := ScreenWidth/2.0 + float64(sprite.Bounds().Dx())

	// Figure out the spawn position by moving r pixels from the target at the chosen angle
	pos := Vector{
		X: target.X + math.Cos(angle)*r,
		Y: target.Y + math.Sin(angle)*r,
	}

	// Randomize velocity
	velocity := baseVelocity * rand.Float64() * 1.5

	// Direction is the target minus the current position
	direction := Vector{
		X: target.X - pos.X,
		Y: target.Y - pos.Y,
	}

	//Normalize direction - get just the direction without the length
	normalizedDirection := direction.Normalize()

	// Multiply direction by velocity
	movement := Vector{
		X: normalizedDirection.X * velocity,
		Y: normalizedDirection.Y * velocity,
	}

	println("new meteor")

	return &Meteor{
		position:       pos,
		movement:       movement,
		rotationSpeed:  rotationSpeedMin + rand.Float64()*(rotationSpeedMax-rotationSpeedMin),
		sprite:         sprite,
		explosionAudio: assets.Explosion2SFX,
	}
}

func (m *Meteor) Update() {
	m.position.X += m.movement.X
	m.position.Y += m.movement.Y
	m.rotation += m.rotationSpeed
}

func (m *Meteor) Draw(screen *ebiten.Image) {
	mp := MiddlePoint(m.sprite)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-mp.X, -mp.Y)
	op.GeoM.Rotate(m.rotation)
	op.GeoM.Translate(mp.X, mp.Y)

	op.GeoM.Translate(m.position.X, m.position.Y)

	screen.DrawImage(m.sprite, op)
}

func (m *Meteor) Collider() Rect {
	bounds := m.sprite.Bounds()

	return NewRect(m.position.X, m.position.Y, float64(bounds.Dx()), float64(bounds.Dy()))
}

func (m *Meteor) Hit() {
	m.explosionAudio.Rewind()
	m.explosionAudio.Play()
}
