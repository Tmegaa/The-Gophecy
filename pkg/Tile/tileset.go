package tile

import (
	"encoding/json"
	"image"
	"os"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Chaque ensemble de tiles doit pouvoir donner une image à partir d'un identifiant
type Tileset interface {
	Img(id int) *ebiten.Image
}

// Des données du jeu de tiles désérialisées à partir d'un jeu de tiles standard à image unique
type UniformTilesetJSON struct {
	Path string `json:"image"`
}

// L'objet d'un jeu de tiles orienté vers l'avant utilisé pour les jeux de tiles à image unique
type UniformTileset struct {
	img *ebiten.Image
	gid int
}

// Fonction qui renvoie une image à partir d'un identifiant
func (u *UniformTileset) Img(id int) *ebiten.Image {
	id -= u.gid

	// Obtenir la position sur l'image où se trouve l'identifiant de la tile
	srcX := id % 16
	srcY := id / 16

	// Convertir la position de la tile source en position source en pixels
	srcX *= 24
	srcY *= 24

	return u.img.SubImage(
		image.Rect(
			srcX, srcY, srcX+24, srcY+24,
		),
	).(*ebiten.Image)
}

// Objet tile à sérialiser en format JSON
type TileJSON struct {
	Id     int    `json:"id"`
	Path   string `json:"image"`
	Width  int    `json:"imagewidth"`
	Height int    `json:"imageheight"`
}

// Objet de jeu de tiles à sérialiser en format JSON
type DynTilesetJSON struct {
	Tiles []*TileJSON `json:"tiles"`
}

// Objet de jeu de tiles avec les images respectives et l'identifient du jeu de tiles
type DynTileset struct {
	imgs []*ebiten.Image
	gid  int
}

// Fonction qui renvoie une image à partir d'un identifient pour un jeu de tiles donné
func (d *DynTileset) Img(id int) *ebiten.Image {
	id -= d.gid

	return d.imgs[id]
}

// Fonction de génération d'un jeu de tuiles
func NewTileset(path string, gid int) (Tileset, error) {
	// Lire le contenu du fichier
	contents, err := os.ReadFile(path)
	if err != nil {

		return nil, err
	}

	if strings.Contains(path, "obj") {
		// Deserialization du contenu
		var dynTilesetJSON DynTilesetJSON
		err = json.Unmarshal(contents, &dynTilesetJSON)
		if err != nil {

			return nil, err
		}

		// Création du jeu de tiles
		dynTileset := DynTileset{}
		dynTileset.gid = gid
		dynTileset.imgs = make([]*ebiten.Image, 0)

		// On itère sur les données des tiles et on charge l'image pour chacune
		for _, tileJSON := range dynTilesetJSON.Tiles {

			// On rajoute un paramètre à l'objet JSON: isCollider (qui est un booléen)
			tileJSONPath := tileJSON.Path
			tileJSONPath = filepath.Clean(tileJSONPath)
			tileJSONPath = strings.ReplaceAll(tileJSONPath, "\\", "/")
			tileJSONPath = strings.TrimPrefix(tileJSONPath, "../")
			tileJSONPath = strings.TrimPrefix(tileJSONPath, "../")
			tileJSONPath = filepath.Join("assets/", tileJSONPath)

			img, _, err := ebitenutil.NewImageFromFile(tileJSONPath)
			if err != nil {

				return nil, err
			}

			dynTileset.imgs = append(dynTileset.imgs, img)
		}

		return &dynTileset, nil
	}
	// On retourne le jeu de tuiles uniforme
	var uniformTilesetJSON UniformTilesetJSON

	// Deserialization du contenu
	err = json.Unmarshal(contents, &uniformTilesetJSON)
	if err != nil {
		return nil, err
	}

	uniformTileset := UniformTileset{}

	// On convertit le chemin relatif du jeu de tuiles en chemin relatif de la racine
	tileJSONPath := uniformTilesetJSON.Path
	tileJSONPath = filepath.Clean(tileJSONPath)
	tileJSONPath = strings.ReplaceAll(tileJSONPath, "\\", "/")
	tileJSONPath = strings.TrimPrefix(tileJSONPath, "../")
	tileJSONPath = strings.TrimPrefix(tileJSONPath, "../")
	tileJSONPath = filepath.Join("assets/", tileJSONPath)

	img, _, err := ebitenutil.NewImageFromFile(tileJSONPath)
	if err != nil {

		return nil, err
	}
	uniformTileset.img = img
	uniformTileset.gid = gid

	return &uniformTileset, nil
}
