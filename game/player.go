package game

import (
	"game/assets"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

const (
	shootCooldown     = 350 * time.Millisecond
	rotationPerSecond = math.Pi
	bulletSpawnOffset = 50.0
)

type Player struct {
	game     *Game
	position Vector
	rotation float64

	sprite           *ebiten.Image
	engineFireImages []*ebiten.Image
	engineFireSprite *ebiten.Image

	laserAudio     *audio.Player
	explosionAudio *audio.Player

	shootCooldown *Timer
}

func NewPlayer(g *Game) *Player {
	sprite := assets.PlayerSprite

	// mp := MiddlePoint(sprite)

	pos := Vector{
		X: ScreenWidth / 2,
		Y: ScreenHeight / 2,
		// X: ScreenWidth/2 - mp.X,
		// Y: ScreenHeight/2 - mp.Y,
	}
	return &Player{
		game:     g,
		position: pos,

		sprite:           sprite,
		engineFireImages: assets.PlayerEngineFireSprites,
		engineFireSprite: assets.PlayerEngineFireSprites[0],

		laserAudio:     assets.Laser1SFX,
		explosionAudio: assets.Explosion1SFX,

		shootCooldown: NewTimer(shootCooldown),
	}
}

func (p *Player) Update() {
	p.rotate()
	// p.move()

	p.shootCooldown.Update()
	if p.shootCooldown.IsReady() && ebiten.IsKeyPressed(ebiten.KeySpace) {
		p.shootCooldown.Reset()
		// mp := MiddlePoint(p.sprite)

		spawnPos := Vector{
			// X: p.position.X + mp.X + math.Sin(p.rotation)*bulletSpawnOffset,
			// Y: p.position.Y + mp.Y + math.Cos(p.rotation)*-bulletSpawnOffset,
			X: p.position.X + math.Sin(p.rotation)*bulletSpawnOffset,
			Y: p.position.Y + math.Cos(p.rotation)*-bulletSpawnOffset,
		}

		p.laserAudio.Rewind()
		p.laserAudio.Play()

		bullet := NewBullet(spawnPos, p.rotation)
		p.game.AddBullet(bullet)
	}

	p.engineFireSprite = p.engineFireImages[rand.Intn(len(p.engineFireImages))]
	// p.engineFireSprite = p.engineFireImages[0]
}

func (p *Player) rotate() {
	speed := rotationPerSecond / float64(ebiten.TPS())

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		p.rotation -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.rotation += speed
	}
}

func (p *Player) move() {
	speed := float64(300 / ebiten.TPS()) // 300 pixels per second
	// speed := 0.5

	var delta Vector

	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		delta.Y = speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		delta.Y = -speed
	}
	// if ebiten.IsKeyPressed(ebiten.KeyLeft) {
	// 	delta.X = -speed
	// }
	// if ebiten.IsKeyPressed(ebiten.KeyRight) {
	// 	delta.X = speed
	// }

	// Check for diagonal movement
	if delta.X != 0 && delta.Y != 0 {
		factor := speed / math.Sqrt(delta.X*delta.X+delta.Y*delta.Y)
		delta.X *= factor
		delta.Y *= factor
	}

	p.position.X += delta.X
	p.position.Y += delta.Y

	if p.position.X < 0 {
		p.position.X = 0
	}
	if p.position.Y < 0 {
		p.position.Y = 0
	}
}

func (p *Player) Draw(screen *ebiten.Image) {
	mp := MiddlePoint(p.sprite)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-mp.X, -mp.Y)
	op.GeoM.Rotate(p.rotation)
	// op.GeoM.Translate(mp.X, mp.Y)

	op.GeoM.Translate(p.position.X, p.position.Y)

	screen.DrawImage(p.sprite, op)

	mp2 := MiddlePoint(p.engineFireSprite)
	// b := p.sprite.Bounds()

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-mp2.X, mp.Y)
	op.GeoM.Rotate(p.rotation)
	// op.GeoM.Translate(mp2.X, -float64(b.Dy()))

	op.GeoM.Translate(p.position.X, p.position.Y)
	// op.GeoM.Translate(p.position.X+mp.X, p.position.Y)

	screen.DrawImage(p.engineFireSprite, op)
}

func (p *Player) Collider() Rect {
	bounds := p.sprite.Bounds()

	return NewRect(p.position.X, p.position.Y, float64(bounds.Dx()), float64(bounds.Dy()))
}

func (p *Player) Hit() {
	p.explosionAudio.Rewind()
	p.explosionAudio.Play()
}
