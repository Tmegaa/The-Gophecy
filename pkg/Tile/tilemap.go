package tile

import (
	"encoding/json"
	"os"
	"path"
)

// Données que nous voulons pour une couche dans notre liste de couches
type TilemapLayerJSON struct {
	Data   []int  `json:"data"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Name   string `json:"name"`
}

// Toutes les couches dans une carte de tiles
type TilemapJSON struct {
	Layers   []TilemapLayerJSON `json:"layers"`
	Tilesets []map[string]any   `json:"tilesets"`
}

// Fonction de génération des paquets de tiles
func (t *TilemapJSON) GenTilesets() ([]Tileset, error) {

	tilesets := make([]Tileset, 0)

	for _, tilesetData := range t.Tilesets {
		tilesetPath := path.Join("assets/maps", tilesetData["source"].(string))

		tileset, err := NewTileset(tilesetPath, int(tilesetData["firstgid"].(float64)))
		if err != nil {

			return nil, err
		}
		tilesets = append(tilesets, tileset)
	}
	return tilesets, nil
}

// Ouvre le fichier, l'analyse et renvoie l'objet json et erreur potentielle
func NewTilemapJSON(filepath string) (*TilemapJSON, error) {
	contents, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var tilemapJSON TilemapJSON
	err = json.Unmarshal(contents, &tilemapJSON)
	if err != nil {
		return nil, err
	}

	return &tilemapJSON, nil
}
