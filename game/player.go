package game

import (
	"game/assets"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

const (
	shootCooldown     = 250 * time.Millisecond
	rotationPerSecond = math.Pi
	bulletSpawnOffset = 50.0
)

type Player struct {
	game           *Game
	position       Vector
	rotation       float64
	sprite         *ebiten.Image
	laserAudio     *audio.Player
	explosionAudio *audio.Player

	shootCooldown *Timer
}

func NewPlayer(g *Game) *Player {
	sprite := assets.PlayerSprite
	laserAudio := assets.Laser1SFX

	mp := MiddlePoint(sprite)

	pos := Vector{
		X: ScreenWidth/2 - mp.X,
		Y: ScreenHeight/2 - mp.Y,
	}
	return &Player{
		game:           g,
		position:       pos,
		sprite:         sprite,
		laserAudio:     laserAudio,
		explosionAudio: assets.Explosion1SFX,

		shootCooldown: NewTimer(shootCooldown),
	}
}

func (p *Player) Update() {
	p.rotate()

	p.shootCooldown.Update()
	if p.shootCooldown.IsReady() && ebiten.IsKeyPressed(ebiten.KeySpace) {
		p.shootCooldown.Reset()
		mp := MiddlePoint(p.sprite)

		spawnPos := Vector{
			X: p.position.X + mp.X + math.Sin(p.rotation)*bulletSpawnOffset,
			Y: p.position.Y + mp.Y + math.Cos(p.rotation)*-bulletSpawnOffset,
		}

		p.laserAudio.Rewind()
		p.laserAudio.Play()

		bullet := NewBullet(spawnPos, p.rotation)
		p.game.AddBullet(bullet)

	}
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
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		delta.X = -speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		delta.X = speed
	}

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
	op.GeoM.Translate(mp.X, mp.Y)

	op.GeoM.Translate(p.position.X, p.position.Y)

	screen.DrawImage(p.sprite, op)
}

func (p *Player) Collider() Rect {
	bounds := p.sprite.Bounds()

	return NewRect(p.position.X, p.position.Y, float64(bounds.Dx()), float64(bounds.Dy()))
}

func (p *Player) Hit() {
	p.explosionAudio.Rewind()
	p.explosionAudio.Play()
}
