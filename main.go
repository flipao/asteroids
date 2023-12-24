package main

import (
	"game/game"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(game.ScreenWidth, game.ScreenHeight)
	ebiten.SetWindowTitle("Asteroids (learning games in Go)")
	g := game.NewGame()
	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}
}
