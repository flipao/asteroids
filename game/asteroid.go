package game

import (
	"asteroids/assets"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	rotationSpeedMin = -0.02
	rotationSpeedMax = 0.02
)

type AsteroidSize int

const (
	AsteroidSizeBig AsteroidSize = iota
	AsteroidSizeMedium
	AsteroidSizeSmall
	AsteroidSizeTiny
)

func (s AsteroidSize) String() string {
	switch s {
	case AsteroidSizeBig:
		return "Big"
	case AsteroidSizeMedium:
		return "Medium"
	case AsteroidSizeSmall:
		return "Small"
	case AsteroidSizeTiny:
		return "Tiny"
	default:
		panic("not reached")
	}
}

type Asteroid struct {
	game           *Game
	position       Vector
	movement       Vector
	rotation       float64
	velocity       float64
	rotationSpeed  float64
	sprite         *ebiten.Image
	explosionAudio *audio.Player
	size           AsteroidSize
}

func NewAsteroid(game *Game, size AsteroidSize, baseVelocity float64, fromAsteroid *Asteroid) *Asteroid {
	var sprites []*ebiten.Image
	switch size {
	case AsteroidSizeBig:
		sprites = assets.BigAsteroidSprites
	case AsteroidSizeMedium:
		sprites = assets.MediumAsteroidSprites
	case AsteroidSizeSmall:
		sprites = assets.SmallAsteroidSprites
	default:
		sprites = assets.TinyAsteroidSprites
	}

	sprite := sprites[rand.Intn(len(sprites))]

	// Figure out the target position - the screen center, in this case
	target := Vector{
		X: ScreenWidth / 2,
		Y: ScreenHeight / 2,
	}

	// Pick a random angle - 2pi is 360 - so this return 0 to 360
	angle := rand.Float64() * 2 * math.Pi

	// The distance from the center, the meteor should spawn at - half  the width
	r := ScreenWidth / 2.0 // + float64(sprite.Bounds().Dx())

	var pos Vector
	var velocity float64
	if fromAsteroid != nil {
		// Debris from a bigger asteroid will spawn from the original's asteroid position with less speed
		pos = fromAsteroid.position
		pos.X += rand.Float64() * 10
		pos.Y += rand.Float64() * 10
		velocity = fromAsteroid.velocity / 3
	} else {
		// Figure out the spawn position by moving r pixels from the target at the chosen angle
		pos = Vector{
			X: target.X + math.Cos(angle)*r,
			Y: target.Y + math.Sin(angle)*r,
		}
		// Randomize velocity
		velocity = baseVelocity * rand.Float64() * 1.5
	}

	// Direction is the target minus the current position
	direction := Vector{
		X: target.X - pos.X,
		Y: target.Y - pos.Y,
	}

	// Normalize direction - get just the direction without the length
	normalizedDirection := direction.Normalize()

	// Multiply direction by velocity
	movement := Vector{
		X: normalizedDirection.X * velocity,
		Y: normalizedDirection.Y * velocity,
	}

	println("new meteor", velocity)

	return &Asteroid{
		game:           game,
		position:       pos,
		movement:       movement,
		velocity:       velocity,
		rotationSpeed:  rotationSpeedMin + rand.Float64()*(rotationSpeedMax-rotationSpeedMin),
		sprite:         sprite,
		explosionAudio: assets.Explosion2SFX,
		size:           size,
	}
}

func (a *Asteroid) Update() {
	a.position.X += a.movement.X
	a.position.Y += a.movement.Y
	a.rotation += a.rotationSpeed
}

func (a *Asteroid) Draw(screen *ebiten.Image) {
	mp := MiddlePoint(a.sprite)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-mp.X, -mp.Y)
	op.GeoM.Rotate(a.rotation)
	op.GeoM.Translate(mp.X, mp.Y)

	op.GeoM.Translate(a.position.X, a.position.Y)

	screen.DrawImage(a.sprite, op)
}

func (a *Asteroid) Collider() Rect {
	bounds := a.sprite.Bounds()

	return NewRect(a.position.X, a.position.Y, float64(bounds.Dx()), float64(bounds.Dy()))
}

func (a *Asteroid) Hit() {
	a.explosionAudio.Rewind()
	a.explosionAudio.Play()

	if a.size == AsteroidSizeBig {
		a.game.AddAsteroid(AsteroidSizeSmall, a)
		a.game.AddAsteroid(AsteroidSizeSmall, a)
	}
	if a.size == AsteroidSizeMedium {
		a.game.AddAsteroid(AsteroidSizeSmall, a)
		a.game.AddAsteroid(AsteroidSizeTiny, a)
	}
}

func (a *Asteroid) DebugInfo(screen *ebiten.Image) {
	rect := a.Collider()
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
