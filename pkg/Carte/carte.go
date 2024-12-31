package carte

import (
	tile "Gophecy/pkg/Tile"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// Définition du type Carte
type Carte struct {
	TilemapJSON tile.TilemapJSON
	Tilesets    []tile.Tileset
	TilemapImg  *ebiten.Image
	Coliders    []image.Rectangle
	Ordinateurs []image.Rectangle
	Statues     []image.Rectangle
}

// Fonction de création d'une nouvelle carte
func NewCarte(tilemapJSON tile.TilemapJSON, tilesets []tile.Tileset, tilemapImg *ebiten.Image, coliders []image.Rectangle, ordinateurs []image.Rectangle, statues []image.Rectangle) *Carte {
	return &Carte{
		TilemapJSON: tilemapJSON,
		Tilesets:    tilesets,
		TilemapImg:  tilemapImg,
		Coliders:    coliders,
		Ordinateurs: ordinateurs,
		Statues:     statues,
	}
}
