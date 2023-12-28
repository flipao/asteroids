package assets

import (
	"embed"
	"image"
	_ "image/png"
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

//go:embed *
var assets embed.FS

var PlayerSprite = mustLoadImage("player/player.png")
var PlayerEngineFireSprites = mustLoadImages("player/fire_*.png")
var BigAsteroidSprites = mustLoadImages("asteroids/big/*.png")
var MediumAsteroidSprites = mustLoadImages("asteroids/medium/*.png")
var SmallAsteroidSprites = mustLoadImages("asteroids/small/*.png")
var TinyAsteroidSprites = mustLoadImages("asteroids/tiny/*.png")
var LaserSprite = mustLoadImage("laser.png")
var BackgroundSprite = mustLoadImage("background.png")

var ScoreFont = mustLoadFont("font.ttf")

var Laser1SFX = mustLoadAudio("sfx/sfx_laser1.ogg")
var Explosion1SFX = mustLoadAudio("sfx/sfx_explosion1.ogg")
var Explosion2SFX = mustLoadAudio("sfx/sfx_explosion2.ogg")
var Explosion3SFX = mustLoadAudio("sfx/sfx_explosion3.ogg")

func mustLoadImage(path string) *ebiten.Image {
	f, err := assets.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}
	return ebiten.NewImageFromImage(img)
}

func mustLoadImages(path string) []*ebiten.Image {
	matches, err := fs.Glob(assets, path)
	if err != nil {
		panic(err)
	}

	images := make([]*ebiten.Image, len(matches))
	for i, match := range matches {
		images[i] = mustLoadImage(match)
	}

	return images
}

func mustLoadFont(name string) font.Face {
	f, err := assets.ReadFile(name)
	if err != nil {
		panic(err)
	}

	tt, err := opentype.Parse(f)
	if err != nil {
		panic(err)
	}

	face, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    48,
		DPI:     72,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		panic(err)
	}

	return face
}

func mustLoadAudio(name string) *audio.Player {
	f, err := assets.Open(name)
	if err != nil {
		panic(err)
	}

	s, err := vorbis.DecodeWithoutResampling(f)
	if err != nil {
		panic(err)
	}

	ctx := audio.CurrentContext()
	if ctx == nil {
		ctx = audio.NewContext(48000)
	}

	p, err := ctx.NewPlayer(s)
	if err != nil {
		panic(err)
	}

	return p
}
