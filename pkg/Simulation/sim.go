package simulation

import (
	ag "Gophecy/pkg/Agent"
	carte "Gophecy/pkg/Carte"
	tile "Gophecy/pkg/Tile"
	ut "Gophecy/pkg/Utilitaries"
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	TileSize        = 24
	AgentImageSize  = 16
	WindowWidth     = 1920
	WindowHeight    = 1080
	NumAgents       = 20
	AssetsPath      = "assets/images/"
	MapsPath        = "assets/maps/"
	AgentImageFile  = "ninja.png"
	TilemapImage    = "img.png"
	TilemapJSONFile = "spawn.json"
	
)

type Simulation struct {
	env         ag.Environnement
	agents      []ag.Agent
	maxStep     int
	maxDuration time.Duration
	step        int
	start       time.Time
	syncChans   sync.Map
	carte       carte.Carte
	selected    *ag.Agent
}

// NewSimulation initializes a new simulation
func NewSimulation(maxStep int, maxDuration time.Duration) *Simulation {
	initializeWindow()
	env := createEnvironment()
	carte := loadMap()
	agents := createAgents(env,carte)

	return &Simulation{
		env:         env,
		agents:      agents,
		maxStep:     maxStep,
		maxDuration: maxDuration,
		start:       time.Now(),
		carte:       *carte,
	}
}

func initializeWindow() {
	ebiten.SetWindowSize(WindowWidth, WindowHeight)
	ebiten.SetWindowTitle("Simulation")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
}

func createEnvironment() ag.Environnement {
	return *ag.NewEnvironment(make([]ag.Agent, NumAgents))
}

func loadMap() *carte.Carte {
	tilemapImg := loadImage(AssetsPath + TilemapImage)
	tilemapJSON := loadTilemapJSON(MapsPath + TilemapJSONFile)
	tilesets := loadTilesets(tilemapJSON)
	coliders := generateColliders(tilemapJSON, tilesets)
	return carte.NewCarte(*tilemapJSON, tilesets, tilemapImg, coliders)
}
func getValidSpawnPositions(carte *carte.Carte, tilesetID int) []ut.Position {
    validPositions := []ut.Position{}
    for layerIdx, layer := range carte.TilemapJSON.Layers {
        for i, tileID := range layer.Data {
			if tileID == 5 || tileID == 6 || tileID == 21 || tileID == 22 { // Assumindo que 0 representa um tile vazio
                x := float64((i % layer.Width) * TileSize)
                y := float64((i / layer.Width) * TileSize)
				img := carte.Tilesets[layerIdx].Img(tileID)
				offsetY := -(img.Bounds().Dy() + TileSize)
				y += float64(offsetY)

                validPositions = append(validPositions, ut.Position{X: x, Y: y})
            }
        }
    }
    return validPositions
}


func createAgents(env ag.Environnement, carte *carte.Carte) []ag.Agent {
	agentsImg := loadImage(AssetsPath + AgentImageFile)
	agents := make([]ag.Agent, NumAgents)

	validPositions := getValidSpawnPositions(carte, 1)

	if len(validPositions) < NumAgents {
        log.Fatalf("Not enough valid spawn positions for all agents")
    }

	rand.Shuffle(len(validPositions), func(i, j int) {
        validPositions[i], validPositions[j] = validPositions[j], validPositions[i]
    })


	for i := 0; i < NumAgents; i++ {
		agents[i] = ag.Agent{
			Env:               &env,
			Id:                ag.IdAgent(fmt.Sprintf("Agent%d", i)),
			Velocite:          rand.Float64(),
			Acuite:            rand.Float64(),
			Position:          validPositions[i],
			Opinion:           rand.Float64(),
			Charisme:          make(map[ag.IdAgent]float64),
			Relation:          make(map[ag.IdAgent]float64),
			PersonalParameter: rand.Float64(),
			Poid_rel:          []float64{rand.Float64(), rand.Float64()},
			Vivant:            true,
			TypeAgt:           []ag.TypeAgent{ag.Sceptic, ag.Believer, ag.Neutral}[rand.Intn(3)],
			SyncChan:          make(chan int),
			Img:               agentsImg,
		}
		env.AddAgent(agents[i])
	}
	return agents
}

func loadImage(path string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatalf("Failed to load image: %s, error: %v", path, err)
	}
	return img
}

func loadTilemapJSON(path string) *tile.TilemapJSON {
	tilemap, err := tile.NewTilemapJSON(path)
	if err != nil {
		log.Fatalf("Failed to load tilemap JSON: %s, error: %v", path, err)
	}
	return tilemap
}

func loadTilesets(tilemapJSON *tile.TilemapJSON) []tile.Tileset {
	tilesets, err := tilemapJSON.GenTilesets()
	if err != nil {
		log.Fatalf("Failed to generate tilesets, error: %v", err)
	}
	return tilesets
}

func generateColliders(tilemapJSON *tile.TilemapJSON, tilesets []tile.Tileset) []image.Rectangle {
	var coliders []image.Rectangle
	for layerIdx, layer := range tilemapJSON.Layers {
		for i, tileID := range layer.Data {
			if tileID == 0 || layerIdx == 0 {
				continue
			}
			x, y := (i%layer.Width)*TileSize, (i/layer.Width)*TileSize
			img := tilesets[layerIdx].Img(tileID)
			offsetY := -(img.Bounds().Dy() + TileSize)
			y += offsetY
			coliders = append(coliders, image.Rect(x, y, x+img.Bounds().Dx(), y+img.Bounds().Dy()))
		}
	}
	return coliders
}

func (sim *Simulation) Draw(screen *ebiten.Image) {
	// Fundo e agentes já desenhados
	screen.Fill(color.RGBA{120, 180, 255, 255})
	sim.drawMap(screen)
	sim.drawAgents(screen)
	sim.drawColliders(screen)

	// Exibe atributos do agente selecionado
	if sim.selected != nil {
		// Desenha um retângulo para exibir informações
		infoBoxX, infoBoxY := 10, 10
		infoBoxWidth, infoBoxHeight := 200, 100
		vector.StrokeRect(screen, float32(infoBoxX), float32(infoBoxY), float32(infoBoxWidth), float32(infoBoxHeight), 2, color.RGBA{0, 0, 0, 255}, true)

		// Renderiza informações do agente selecionado
		text := fmt.Sprintf("ID: %s\nType Agent: %s\nPos: (%.2f, %.2f)\nPersonalParameter: %.2f \nVivant: %t",
			sim.selected.Id,
			sim.selected.TypeAgt,
			sim.selected.Position.X,
			sim.selected.Position.Y,
			sim.selected.PersonalParameter,
			sim.selected.Vivant,
		)

		// Adiciona texto dentro do retângulo
		ebitenutil.DebugPrintAt(screen, text, infoBoxX+10, infoBoxY+10)
	}
}


func (sim *Simulation) drawMap(screen *ebiten.Image) {
	opts := ebiten.DrawImageOptions{}
	for layerIdx, layer := range sim.carte.TilemapJSON.Layers {
		for i, tileID := range layer.Data {
			if tileID == 0 {
				continue
			}
			x, y := (i%layer.Width)*TileSize, (i/layer.Width)*TileSize
			img := sim.carte.Tilesets[layerIdx].Img(tileID)
			opts.GeoM.Translate(float64(x), float64(y-img.Bounds().Dy()-TileSize))
			screen.DrawImage(img, &opts)
			opts.GeoM.Reset()
		}
	}
}

func (sim *Simulation) drawAgents(screen *ebiten.Image) {
	opts := ebiten.DrawImageOptions{}
	for _, agent := range sim.agents {
		opts.GeoM.Translate(agent.Position.X, agent.Position.Y)
		subImg := agent.Img.SubImage(image.Rect(0, 0, AgentImageSize, AgentImageSize)).(*ebiten.Image)
		screen.DrawImage(subImg, &opts)
		opts.GeoM.Reset()
	}
}

func (sim *Simulation) drawColliders(screen *ebiten.Image) {
	for _, colider := range sim.carte.Coliders {
		vector.StrokeRect(screen, float32(colider.Min.X), float32(colider.Min.Y), float32(colider.Dx()), float32(colider.Dy()), 1.0, color.RGBA{0, 0, 0, 0}, true)
	}
}

func (sim *Simulation) Layout(outsideWidth, outsideHeight int) (int, int) {
	return WindowWidth, WindowHeight
}

func (sim *Simulation) Update() error {
	// Posição do cursor
	cursorX, cursorY := ebiten.CursorPosition()

	// Detecta clique
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		for i := range sim.agents {
			agent := &sim.agents[i]
			// Verifica se o clique está dentro da área do agente
			if cursorX >= int(agent.Position.X) && cursorX <= int(agent.Position.X+AgentImageSize) &&
				cursorY >= int(agent.Position.Y) && cursorY <= int(agent.Position.Y+AgentImageSize) {
				sim.selected = agent // Define o agente selecionado
				break
			}
		}
	}

	return nil
}


func (sim *Simulation) Run() {	
	for _, ag := range sim.agents {
		go ag.Start()
	}
	sim.start = time.Now()
	time.Sleep(sim.maxDuration)
}
