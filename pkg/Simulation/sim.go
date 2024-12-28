package simulation

import (
	ag "Gophecy/pkg/Agent"
	carte "Gophecy/pkg/Carte"
	tile "Gophecy/pkg/Tile"
	ut "Gophecy/pkg/Utilitaries"
	"context"
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

// Constantes de la simulation
const (
	TileSize        = 24
	AgentImageSize  = 16
	WindowWidth     = 1920
	WindowHeight    = 1080
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
	ctx         context.Context
	cancel      context.CancelFunc
	dialogFont  font.Face
}

// Fonction qui lance une nouvelle simulation
func NewSimulation(config SimulationConfig) *Simulation {
	initializeWindow()
	carte := loadMap()
	env := createEnvironment(*carte, config.NumAgents)
	agents := createAgents(&env, carte, config.NumAgents)
	ctx, cancel := context.WithTimeout(context.Background(), config.SimulationTime)
	tt, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal(err)
	}

	return &Simulation{
		env:         env,
		agents:      agents,
		maxStep:     10, // Possible de l'ajouter à la configuration
		maxDuration: config.SimulationTime,
		start:       time.Now(),
		carte:       *carte,
		ctx:         ctx,
		cancel:      cancel,
		dialogFont: truetype.NewFace(tt, &truetype.Options{
			Size: 12,
			DPI:  72,
		}),
	}
}

// Fonction qui initialise la fenêtre d'affichage de la simulation
func initializeWindow() {
	ebiten.SetWindowSize(WindowWidth, WindowHeight)
	ebiten.SetWindowTitle("Simulation")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
}

// Fonction qui crée et retourne un nouvel environnement
func createEnvironment(carte carte.Carte, NumAgents int) ag.Environnement {
	return *ag.NewEnvironment(make([]ag.Agent, 0), carte, make([]ag.InterfaceObjet, 0))
}

// Fonction qui charge la carte
func loadMap() *carte.Carte {
	tilemapImg := loadImage(AssetsPath + TilemapImage)
	tilemapJSON := loadTilemapJSON(MapsPath + TilemapJSONFile)
	tilesets := loadTilesets(tilemapJSON)
	coliders := generateColliders(tilemapJSON, tilesets)
	return carte.NewCarte(*tilemapJSON, tilesets, tilemapImg, coliders)
}

// Fonction qui renvoie une liste des positions possibles que peuvent prendre les objets ou les agents
func getValidSpawnPositions(carte *carte.Carte, tilesetID int) []ut.Position {
	validPositions := []ut.Position{}
	for layerIdx, layer := range carte.TilemapJSON.Layers {
		for i, tileID := range layer.Data {
			if tileID == 5 || tileID == 6 || tileID == 21 || tileID == 22 { // En supposant que 0 représente une tile vide
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

// Fonction qui crée et rajoute à la carte les nouveaux agents
func createAgents(env *ag.Environnement, carte *carte.Carte, NumAgents int) []ag.Agent {
	agentsImg := loadImage(AssetsPath + AgentImageFile)
	agents := make([]ag.Agent, NumAgents)

	validPositions := getValidSpawnPositions(carte, 1)

	if len(validPositions) < NumAgents {
		log.Fatalf("Not enough valid spawn positions for all agents")
	}

	rand.Shuffle(len(validPositions), func(i, j int) {
		validPositions[i], validPositions[j] = validPositions[j], validPositions[i]
	})

	// for i := 0; i < NumAgents; i++ {
	// 	agents[i] = *ag.NewAgent(&env, ag.IdAgent(fmt.Sprintf("Agent%d", i)), rand.Float64(), rand.Float64(), validPositions[i], rand.Float64(), make(map[ag.IdAgent]float64), make(map[ag.IdAgent]float64), rand.Float64(), []float64{rand.Float64(), rand.Float64()}, []ag.TypeAgent{ag.Sceptic, ag.Believer, ag.Neutral}[rand.Intn(3)], make(chan int), agentsImg)
	// 	env.AddAgent(agents[i])
	// }

	//i dont why we are creating agents like this and not using the function NewAgent
	for i := 0; i < NumAgents; i++ {
		agents[i] = ag.Agent{
			Env:               env,
			Id:                ag.IdAgent(fmt.Sprintf("Agent%d", i)),
			Velocite:          rand.Float64(),
			Acuite:            30.0, //float64(rand.Intn(10)),
			Position:          validPositions[i],
			Opinion:           rand.Float64(),
			Charisme:          make(map[ag.IdAgent]float64),
			Relation:          make(map[ag.IdAgent]float64),
			PersonalParameter: rand.Float64(),
			Poid_rel:          []float64{rand.Float64(), rand.Float64()},
			Vivant:            true,
			TypeAgt:           []ag.TypeAgent{ag.Sceptic, ag.Believer, ag.Neutral}[rand.Intn(3)],
			SyncChan:          make(chan ag.Message),
			Img:               agentsImg,
			MoveTimer:         2,
			CurrentAction:     "Praying",
			DialogTimer:       2,
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
	// Desenha o fundo e os agentesrgba(57,61,125,255)
	screen.Fill(color.RGBA{57, 61, 125, 255})
	sim.drawMap(screen)
	sim.drawAgents(screen)
	sim.drawAcuite(screen)
	sim.drawColliders(screen)
	sim.drawInfoPanel(screen)
}

func (sim Simulation) drawAcuite(screen *ebiten.Image) {
	for _, agent := range sim.agents {

		centerX := agent.Position.X + float64(AgentImageSize)/2
		centerY := agent.Position.Y + float64(AgentImageSize)/2

		area := ut.Rectangle{
			PositionDL: ut.Position{
				X: centerX - agent.Acuite,
				Y: centerY + agent.Acuite,
			},
			PositionUR: ut.Position{
				X: centerX + agent.Acuite,
				Y: centerY - agent.Acuite,
			},
		}

		vector.StrokeRect(
			screen,
			float32(area.PositionDL.X),
			float32(area.PositionUR.Y),
			float32(area.PositionUR.X-area.PositionDL.X),
			float32(area.PositionDL.Y-area.PositionUR.Y),
			1,
			color.RGBA{255, 0, 0, 128},
			false,
		)
	}
}

func (sim *Simulation) drawInfoPanel(screen *ebiten.Image) {
	panelX, panelY := 0, 0
	panelWidth, panelHeight := 240, WindowHeight-20
	padding := 10

	// Desenha o painel de fundo
	vector.DrawFilledRect(screen, float32(panelX), float32(panelY), float32(panelWidth), float32(panelHeight), color.RGBA{0, 0, 0, 180}, false)

	// Título do painel
	ebitenutil.DebugPrintAt(screen, "Simulation Info", panelX+padding, panelY+padding)

	y := panelY + 30

	// Informações da simulação
	elapsed := time.Since(sim.start)
	simInfo := fmt.Sprintf("Elapsed: %s", elapsed.Round(time.Second))
	ebitenutil.DebugPrintAt(screen, simInfo, panelX+padding, y)
	y += 40

	// Contagem de agentes por tipo
	ebitenutil.DebugPrintAt(screen, "Agent Count:", panelX+padding, y)
	y += 20
	agentTypes := []ag.TypeAgent{ag.Sceptic, ag.Believer, ag.Neutral}
	for _, agentType := range agentTypes {
		count, _ := sim.env.NbrAgents.Load(agentType)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("  %s: %d", agentType, count), panelX+padding, y)
		y += 20
	}
	y += 20

	// Informações do agente selecionado
	if sim.selected != nil {
		ebitenutil.DebugPrintAt(screen, "Selected Agent:", panelX+padding, y)
		y += 20
		agentInfo := fmt.Sprintf("  ID: %s\n  Type: %s\n  Position: (%.2f, %.2f)\n  Personal Param: %.2f\n  Alive: %t\n  DialogTimer : %d\n  CurrentAction : %s\n  Time to change direction : %d \n  Occupied : %t",
			sim.selected.Id,
			sim.selected.TypeAgt,
			sim.selected.Position.X,
			sim.selected.Position.Y,
			sim.selected.PersonalParameter,
			sim.selected.Vivant,
			sim.selected.DialogTimer,
			sim.selected.CurrentAction,
			sim.selected.MoveTimer,
			sim.selected.Occupied,
		)
		ebitenutil.DebugPrintAt(screen, agentInfo, panelX+padding, y)
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

		sim.drawDialogBox(screen, agent)
	}
}

func (sim *Simulation) drawDialogBox(screen *ebiten.Image, agent ag.Agent) {
	if agent.CurrentAction == "" || agent.DialogTimer <= 0 {
		return
	}

	dialogWidth := 100
	dialogHeight := 30
	x := int(agent.Position.X) - dialogWidth/2 + AgentImageSize/2
	y := int(agent.Position.Y) - dialogHeight - 5

	// Desenha o fundo da caixa de diálogo
	vector.DrawFilledRect(screen, float32(x), float32(y), float32(dialogWidth), float32(dialogHeight), color.RGBA{255, 255, 255, 200}, false)

	// Desenha a borda da caixa de diálogo
	vector.StrokeRect(screen, float32(x), float32(y), float32(dialogWidth), float32(dialogHeight), 1, color.Black, false)

	// Escreve o texto da ação
	text.Draw(screen, string(agent.CurrentAction), sim.dialogFont, x+5, y+20, color.Black)
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
	select {
	case <-sim.ctx.Done():
		return ebiten.Termination
	default:
		//Detection des nearby agents

		//log.Printf("Agents proches détectés: %v", sim.env.AgentProximity)
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

		for i := range sim.agents {

			if sim.agents[i].DialogTimer > 0 {
				sim.agents[i].DialogTimer--
				if sim.agents[i].DialogTimer == 0 {
					sim.agents[i].ClearAction()
				}
			} else {
				/*
					if nearbyAgents, ok := sim.env.AgentProximity.Load(sim.agents[i].Id); ok {
						log.Printf("Nearby agents of %s : %v", sim.agents[i].Id, nearbyAgents)
					}

				*/
				/*
					//sim.env.NearbyAgents(&ag)
					switch sim.agents[i].CurrentAction {

					case "Praying":
						sim.agents[i].Pray()
					case "Running":
						sim.agents[i].Move()
					case "Discuss":
						//log.Printf("Agent %s is discussing", sim.agents[i].Id)
						//sim.agents[i].Discuss()
					}
				*/
			}
		}
	}
	return nil
}

func (sim *Simulation) Run() error {
	defer sim.cancel() // Ensure context is canceled when Run() exits

	go sim.env.Listen()
	go func() {
		for i := range sim.agents {
			go sim.agents[i].Start()
		}
		sim.start = time.Now()
	}()

	if err := ebiten.RunGame(sim); err != nil && err != ebiten.Termination {
		return err
	}
	// Affichages de finalisation
	fmt.Println("\n--- Simulation Terminée ---")
	fmt.Printf("Durée totale : %s\n", time.Since(sim.start).Round(time.Second))

	// Comptage des agents par type
	agentCounts := make(map[ag.TypeAgent]int)
	for _, agent := range sim.agents {
		agentCounts[agent.TypeAgt]++
	}
	fmt.Println("\nNombre final d'agents par type :")
	for agentType, count := range agentCounts {
		fmt.Printf("- %s : %d\n", agentType, count)
	}

	// Statistiques supplémentaires
	totalOpinion := 0.0
	for _, agent := range sim.agents {
		totalOpinion += agent.Opinion
	}
	averageOpinion := totalOpinion / float64(len(sim.agents))
	fmt.Printf("\nOpinion moyenne des agents : %.2f\n", averageOpinion)

	return nil

}
