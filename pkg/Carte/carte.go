package carte

import (
	tile "Gophecy/pkg/Tile"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Carte struct {
	TilemapJSON tile.TilemapJSON
	Tilesets    []tile.Tileset
	TilemapImg  *ebiten.Image
	Coliders []image.Rectangle
}

func NewCarte(tilemapJSON tile.TilemapJSON, tilesets []tile.Tileset, tilemapImg *ebiten.Image, coliders []image.Rectangle) *Carte {
	return &Carte{
		TilemapJSON: tilemapJSON,
		Tilesets:    tilesets,
		TilemapImg:  tilemapImg,
		Coliders: coliders,
	}
}


