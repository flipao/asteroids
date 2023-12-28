package game

import (
	"fmt"
	"game/assets"
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 600

	meteorSpawnTime     = 2 * time.Second
	baseMeteorVelocity  = 0.35
	meteorSpeedUpAmount = 0.1
	meteorSpeedUpTime   = 5 * time.Second
)

type Game struct {
	background         *Background
	player             *Player
	asteroidSpawnTimer *Timer
	asteroids          []*Asteroid
	bullets            []*Bullet

	score int

	baseVelocity  float64
	velocityTimer *Timer
}

func NewGame() *Game {
	g := &Game{
		background:         NewBackground(),
		asteroidSpawnTimer: NewTimer(meteorSpawnTime),
		baseVelocity:       baseMeteorVelocity,
		velocityTimer:      NewTimer(meteorSpeedUpTime),
	}

	g.player = NewPlayer(g)

	return g
}

func (g *Game) Update() error {
	g.velocityTimer.Update()
	if g.velocityTimer.IsReady() {
		g.velocityTimer.Reset()
		g.baseVelocity += meteorSpeedUpAmount
	}

	g.background.Update()

	g.player.Update()

	g.asteroidSpawnTimer.Update()
	if g.asteroidSpawnTimer.IsReady() {
		g.asteroidSpawnTimer.Reset()

		sizes := []AsteroidSize{SizeBig, SizeMedium, SizeSmall}

		g.AddAsteroid(sizes[rand.Intn(len(sizes))], nil)
	}
	for _, m := range g.asteroids {
		m.Update()
	}
	for _, b := range g.bullets {
		b.Update()
	}
	// Check meteors hits by bullets
	for i, m := range g.asteroids {
		for j, b := range g.bullets {
			if m.Collider().Intersects(b.Collider()) {
				println("Meteor", i, "hit by Bullet", j)
				m.Hit()
				g.score++
				g.asteroids = append(g.asteroids[:i], g.asteroids[i+1:]...)
				g.bullets = append(g.bullets[:j], g.bullets[j+1:]...)
			}
		}
	}
	for i, m := range g.asteroids {
		if m.Collider().Intersects(g.player.Collider()) {
			println("Player hit by Meteor", i)
			g.player.Hit()
			g.Reset()
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.background.Draw(screen)

	g.player.Draw(screen)

	for _, m := range g.asteroids {
		m.Draw(screen)
	}
	for _, b := range g.bullets {
		b.Draw(screen)
	}

	// Draw score
	text.Draw(screen, fmt.Sprintf("%06d", g.score), assets.ScoreFont, ScreenWidth/2-100, 50, color.White)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) AddBullet(b *Bullet) {
	g.bullets = append(g.bullets, b)
}

func (g *Game) AddAsteroid(size AsteroidSize, fromAsteroid *Asteroid) {
	a := NewAsteroid(g, size, g.baseVelocity, fromAsteroid)
	g.asteroids = append(g.asteroids, a)
}

func (g *Game) Reset() {
	g.player = NewPlayer(g)
	g.asteroids = nil
	g.bullets = nil
	g.score = 0
	g.asteroidSpawnTimer.Reset()
	g.baseVelocity = baseMeteorVelocity
	g.velocityTimer.Reset()
}
