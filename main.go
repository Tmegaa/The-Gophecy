package main

import (
	"front/entities"
	"image"
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/exp/rand"
)

func CheckCollisionHorizontal(sprite *entities.Sprite, coliders []image.Rectangle) {
	for _ , colider := range coliders {
		if colider.Overlaps(image.Rect(int(sprite.X), int(sprite.Y), int(sprite.X) + 16, int(sprite.Y) + 16)) {
			if sprite.Dx > 0 {
				sprite.X = float64(colider.Min.X) - 16
			}else if sprite.Dx < 0 {
				sprite.X = float64(colider.Max.X)
			}
		}
	}
}

func CheckCollisionVertical(sprite *entities.Sprite, coliders []image.Rectangle) {
	for _ , colider := range coliders {
		if colider.Overlaps(image.Rect(int(sprite.X), int(sprite.Y), int(sprite.X) + 16, int(sprite.Y) + 16)) {
			if sprite.Dy > 0 {
				sprite.Y = float64(colider.Min.Y) - 16
			}else if sprite.Dy < 0 {
				sprite.Y = float64(colider.Max.Y)
			}
		}
	}
}

// moveEntity aplica movimento aleatório em um sprite
func moveEntity(entity *entities.Sprite, maxSpeed float64) {
	// Verifica o temporizador de movimento da entidade
	if entity.MoveTimer <= 0 {
		// Gera uma direção aleatória
		directions := []struct{ Dx, Dy float64 }{
			{Dx: maxSpeed, Dy: 0},   // Direita
			{Dx: -maxSpeed, Dy: 0},  // Esquerda
			{Dx: 0, Dy: maxSpeed},   // Para baixo
			{Dx: 0, Dy: -maxSpeed},  // Para cima
		}
		randIdx := rand.Intn(len(directions)) // Escolhe uma direção aleatória
		entity.Dx = directions[randIdx].Dx
		entity.Dy = directions[randIdx].Dy
		entity.MoveTimer = 60 // Define um tempo de 60 frames
	} else {
		entity.MoveTimer--
	}

	// Atualiza a posição da entidade
	entity.X += entity.Dx
	entity.Y += entity.Dy
}



type Game struct {
	// the image and position variables for our player
	players  []*entities.Player
	enemies []*entities.Enemy
	tilemapJSON *TilemapJSON
	tilesets []Tileset
	tilemapImg *ebiten.Image
	coliders []image.Rectangle
	
}

func (g *Game) Update() error {
	const maxSpeed = 2.0
	// Atualiza todos os jogadores
	for _, player := range g.players {
		moveEntity(player.Sprite, maxSpeed)
		CheckCollisionHorizontal(player.Sprite, g.coliders)
		CheckCollisionVertical(player.Sprite, g.coliders)
	}

	// Atualiza todos os inimigos
	for _, enemy := range g.enemies {
		moveEntity(enemy.Sprite, maxSpeed)
		CheckCollisionHorizontal(enemy.Sprite, g.coliders)
		CheckCollisionVertical(enemy.Sprite, g.coliders)
	}
	return nil
	// Verifica o temporizador de movimento do jogador
	// Atualiza todos os jogadores
	// for _, player := range g.players {
	// 	// Verifica o temporizador de movimento do jogador
	// 	if player.MoveTimer <= 0 {
	// 		// Gera uma direção aleatória
	// 		directions := []struct{ Dx, Dy float64 }{
	// 			{Dx: maxSpeed, Dy: 0},   // Direita
	// 			{Dx: -maxSpeed, Dy: 0},  // Esquerda
	// 			{Dx: 0, Dy: maxSpeed},   // Para baixo
	// 			{Dx: 0, Dy: -maxSpeed},  // Para cima
	// 		}
	// 		randIdx := rand.Intn(len(directions)) // Escolhe uma direção aleatória
	// 		player.Dx = directions[randIdx].Dx
	// 		player.Dy = directions[randIdx].Dy
	// 		player.MoveTimer = 60 // Define um tempo de 60 frames
	// 	} else {
	// 		player.MoveTimer--
	// 	}

	// 	// Atualiza a posição do jogador
	// 	player.X += player.Dx
	// 	CheckCollisionHorizontal(player.Sprite, g.coliders)

	// 	player.Y += player.Dy
	// 	CheckCollisionVertical(player.Sprite, g.coliders)
	// }
	
	
	// move the player based on keyboar input (left, right, up down)

	// g.player.Dx = 0.0
	// g.player.Dy = 0.0

	// if ebiten.IsKeyPressed(ebiten.KeyLeft) {
	// 	g.player.Dx = -2
	// }
	// if ebiten.IsKeyPressed(ebiten.KeyRight) {
	// 	g.player.Dx = 2
	// }
	// if ebiten.IsKeyPressed(ebiten.KeyUp) {
	// 	g.player.Dy = -2
	// }
	// if ebiten.IsKeyPressed(ebiten.KeyDown) {
	// 	g.player.Dy = 2
	// }

	// g.player.X += g.player.Dx

	// CheckCollisionHorizontal(g.player.Sprite, g.coliders)

	// g.player.Y += g.player.Dy

	// CheckCollisionVertical(g.player.Sprite, g.coliders)

	// Atualiza os inimigos
	// for _, sprite := range g.enemies {
	// 	sprite.Dx = 0.0
	// 	sprite.Dy = 0.0

	// 	if sprite.FollowsPlayer {
	// 		// Lógica de perseguição do inimigo (pode ser aprimorada para seguir um jogador específico)
	// 		targetPlayer := g.players[0] // Exemplo: inimigo segue o primeiro jogador
	// 		if sprite.X < targetPlayer.X {
	// 			sprite.Dx += 1
	// 		} else if sprite.X > targetPlayer.X {
	// 			sprite.Dx -= 1
	// 		}
	// 		if sprite.Y < targetPlayer.Y {
	// 			sprite.Dy += 1
	// 		} else if sprite.Y > targetPlayer.Y {
	// 			sprite.Dy -= 1
	// 		}
	// 	}

	// 	sprite.X += sprite.Dx
	// 	CheckCollisionHorizontal(sprite.Sprite, g.coliders)

	// 	sprite.Y += sprite.Dy
	// 	CheckCollisionVertical(sprite.Sprite, g.coliders)
	// }


	// return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	// fill the screen with a nice sky color
	screen.Fill(color.RGBA{120, 180, 255, 255})

	opts := ebiten.DrawImageOptions{}

	// draw the tilemap
	for layerIdx, layer := range g.tilemapJSON.Layers {
		for i, tileID := range layer.Data {

			if tileID == 0 {
				continue
			}
			
			x := i % layer.Width
			y := i / layer.Width
			
			y *= 24
			x *= 24


			img := g.tilesets[layerIdx].Img(tileID)
			
			opts.GeoM.Translate(float64(x), float64(y))

			opts.GeoM.Translate(0.0, -(float64(img.Bounds().Dy()) + 24.0))

			screen.DrawImage(
				img,
				&opts,
			)

			
			opts.GeoM.Reset()
		}

		
	}


	
	// Desenha todos os jogadores
	for _, player := range g.players {
		opts.GeoM.Translate(player.X, player.Y)

		screen.DrawImage(
			player.Img.SubImage(
				image.Rect(0, 0, 16, 16),
			).(*ebiten.Image),
			&opts,
		)

		opts.GeoM.Reset()
	}


	for _, sprite := range g.enemies {
		opts.GeoM.Translate(sprite.X, sprite.Y)

		screen.DrawImage(
			sprite.Img.SubImage(
				image.Rect(0, 0, 16, 16),
			).(*ebiten.Image),
			&opts,
		)

		opts.GeoM.Reset()
	}

	opts.GeoM.Reset()

	// for _, sprite := range g.potions {
	// 	opts.GeoM.Translate(sprite.X, sprite.Y)

	// 	screen.DrawImage(
	// 		sprite.Img.SubImage(
	// 			image.Rect(0, 0, 16, 16),
	// 		).(*ebiten.Image),
	// 		&opts,
	// 	)

	// 	opts.GeoM.Reset()
	// }

	for _, colider := range g.coliders {
		vector.StrokeRect(
			screen, 
			float32(colider.Min.X), 
			float32(colider.Min.Y), 
			float32(colider.Dx()), 
			float32(colider.Dy()), 
			float32(1.0),
			color.RGBA{0, 0, 0, 0},true)
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ebiten.WindowSize()
}

func main() {
	
	rand.Seed(uint64(time.Now().UnixNano()))
	
	// set the window size and title whole screen
	ebiten.SetWindowSize(ebiten.WindowSize())
	ebiten.SetWindowTitle("Hello, World!")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	// load the image from file
	playerImg, _, err := ebitenutil.NewImageFromFile("assets/images/ninja.png")
	if err != nil {
		// handle error
		log.Fatal(err)
	}
	// load the image from file
	skeletonImg, _, err := ebitenutil.NewImageFromFile("assets/images/skeleton.png")
	if err != nil {
		// handle error
		log.Fatal(err)
	}

	// potionImg, _, err := ebitenutil.NewImageFromFile("assets/images/potion.png")
	// if err != nil {
	// 	// handle error
	// 	log.Fatal(err)
	// }

	tilemapImg, _, err := ebitenutil.NewImageFromFile("assets/images/img.png")
	if err != nil {
		// handle error
		
		log.Fatal(err)
	}

	tilemapJSON, err := NewTilemapJSON("assets/maps/spawn.json")

	if err != nil {
		
		log.Fatal(err)
	}

	tilesets, err := tilemapJSON.GenTilesets()
	if err != nil {
		log.Fatal(err)
	}

	coliders := []image.Rectangle{}
	for layerIdx, layer := range tilemapJSON.Layers {
		
		for i, tileID := range layer.Data {


			if tileID == 0 {
				continue	
			}else if layerIdx != 0  {
				x := i % layer.Width
				y := i / layer.Width
				
				y *= 24
				x *= 24
				
				img := tilesets[layerIdx].Img(tileID)

				offsetY := -(img.Bounds().Dy() + 24)
				y += offsetY

				yy := img.Bounds().Dy()
				xx := img.Bounds().Dx()

				

				coliders = append(coliders, image.Rect(x, y, x + xx, y + yy))
			}

		}

		
	}

	// for _, layer := range tilemapJSON.Layers {
	// 	for i, tileID := range layer.Data {
			
	// 		if tileID == 385 {
	// 			x := i % layer.Width
	// 			y := i / layer.Width
				
	// 			coliders = append(coliders, image.Rect(x*48, y*56, (x+1)*48, (y+1)*56))
	// 		}
			
	// 	}
	// }
	
	const numPlayers = 5

	// Inicializa jogadores
	playerss := make([]*entities.Player, numPlayers)
	for i := 0; i < numPlayers; i++ {
		playerss[i] = &entities.Player{
			Sprite: &entities.Sprite{
				Img: playerImg,
				X:   float64(200 + i*20), 
				Y:   300.0,
				MoveTimer: rand.Intn(60),
			},
			Health:    3,
			
		}
	}

	const numEnemies = 3
	enemies := make([]*entities.Enemy, numEnemies)
	for i := 0; i < numEnemies; i++ {
		enemies[i] = &entities.Enemy{
			Sprite: &entities.Sprite{
				Img: skeletonImg,
				X:   float64(600 + i*30), 
				Y:   400.0,
			},
			FollowsPlayer: false, 
		}
	}


	game := Game{
		players: playerss,
		enemies: enemies,
		tilemapJSON: tilemapJSON,
		tilemapImg: tilemapImg,
		tilesets: tilesets,
		coliders: coliders,
	}

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}