package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func MiddlePoint(s *ebiten.Image) Vector {
	bounds := s.Bounds()
	return Vector{
		X: float64(bounds.Dx()) / 2,
		Y: float64(bounds.Dy()) / 2,
	}
}
