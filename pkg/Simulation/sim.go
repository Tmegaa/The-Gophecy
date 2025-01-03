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
	"os"
	"sort"
	"sync"
	"time"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/wcharczuk/go-chart/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

const (
	TileSize               = 24
	AgentImageSize         = 16
	WindowWidth            = 1920
	WindowHeight           = 1080
	NumComputers           = 6
	NumStatues             = 1
	AssetsPath             = "assets/images/"
	MapsPath               = "assets/maps/"
	AgentBelieverImageFile = "ninja.png"
	AgentScepticImageFile  = "sceptic.png"
	AgentNeutralImageFile  = "neutre.png"
	TilemapImage           = "img.png"
	TilemapJSONFile        = "spawn.json"
	DiscussionBubbleWidth  = 100
	DiscussionBubbleHeight = 40
	ProbabilityConverter   = 0.2
	ProbabilityPirate      = 0.15
)

type Simulation struct {
	env                ag.Environnement
	agents             []ag.Agent
	objets             []ag.InterfaceObjet
	maxStep            int
	maxDuration        time.Duration
	step               int
	start              time.Time
	syncChans          sync.Map
	carte              carte.Carte
	selected           *ag.Agent
	selectedPC         *ag.Computer
	ctx                context.Context
	cancel             context.CancelFunc
	dialogFont         font.Face
	selectionIndicator *ebiten.Image
	opinionAverages   []float64
}

// NewSimulation initializes a new simulation
// pkg/Simulation/simulation.go

func NewSimulation(config SimulationConfig) *Simulation {
	initializeWindow()
	carte := loadMap()
	env := createEnvironment(*carte, config.NumAgents)
	obj := loadObjects(&env, carte)
	agents := createAgents(&env, carte, config)
	ctx, cancel := context.WithTimeout(context.Background(), config.SimulationTime)
	tt, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal(err)
	}

	selectionIndicator := ebiten.NewImage(TileSize, TileSize)
	selectionIndicator.Fill(color.RGBA{255, 255, 0, 128})

	return &Simulation{
		env:         env,
		agents:      agents,
		objets:      obj,
		maxStep:     10,
		maxDuration: config.SimulationTime,
		start:       time.Now(),
		carte:       *carte,
		ctx:         ctx,
		cancel:      cancel,
		dialogFont: truetype.NewFace(tt, &truetype.Options{
			Size: 12,
			DPI:  72,
		}),
		selectionIndicator: selectionIndicator,
	}
}

func initializeWindow() {
	ebiten.SetWindowSize(WindowWidth, WindowHeight)
	ebiten.SetWindowTitle("Simulation")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
}

func createEnvironment(carte carte.Carte, NumAgents int) ag.Environnement {
	return *ag.NewEnvironment(make([]*ag.Agent, 0), carte, make([]ag.InterfaceObjet, 0))
}

func loadMap() *carte.Carte {
	tilemapImg := loadImage(AssetsPath + TilemapImage)
	tilemapJSON := loadTilemapJSON(MapsPath + TilemapJSONFile)
	tilesets := loadTilesets(tilemapJSON)
	coliders := generateColliders(tilemapJSON, tilesets)
	computers := generateComputers(tilemapJSON, tilesets)
	statues := generateStatues(tilemapJSON, tilesets)

	// générer des ordinateurs.
	return carte.NewCarte(*tilemapJSON, tilesets, tilemapImg, coliders, computers, statues)
}

func loadObjects(env *ag.Environnement, carte *carte.Carte) []ag.InterfaceObjet {
	obj := make([]ag.InterfaceObjet, NumComputers+NumStatues)
	
	for i := 0; i < NumComputers; i++ {
		obj[i] = ag.NewComputer(
			env,
			ag.IdObjet(fmt.Sprintf("Computer%d", i)),
			ut.Position{X: float64(env.Carte.Ordinateurs[i].Min.X), Y: float64(env.Carte.Ordinateurs[i].Min.Y)},
		)
		env.Objs = append(env.Objs, obj[i])
	}

	for i := 0; i < NumStatues; i++ {
		obj[i+NumComputers] = ag.NewStatue(
			env,
			ag.IdObjet(fmt.Sprintf("Statue%d", i)),
			ut.Position{X: float64(env.Carte.Statues[i].Min.X), Y: float64(env.Carte.Statues[i].Min.Y)},
		)
		env.Objs = append(env.Objs, obj[i+NumComputers])
	}
	return obj
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

func createAgents(env *ag.Environnement, carte *carte.Carte, config SimulationConfig) []ag.Agent {
	agents := make([]ag.Agent, config.NumAgents)
	validPositions := getValidSpawnPositions(carte, 1)
	visitationMap := ag.NewVisitationMap(validPositions)

	if len(validPositions) < config.NumAgents {
		log.Fatalf("Not enough valid spawn positions for all agents")
	}

	rand.Shuffle(len(validPositions), func(i, j int) {
		validPositions[i], validPositions[j] = validPositions[j], validPositions[i]
	})

	//default image
	agentsImg := loadImage(AssetsPath + AgentBelieverImageFile)

	for i := 0; i < config.NumAgents; i++ {
		// Gera valores aleatórios
		Opinion := rand.Float64()

		// Determina o tipo base do agente
		var TypeChoosen ag.TypeAgent
		if Opinion < 0.33 {
			TypeChoosen = ag.Sceptic
		} else if Opinion < 0.66 {
			TypeChoosen = ag.Neutral
		} else {
			TypeChoosen = ag.Believer
		}
		id := ag.IdAgent(fmt.Sprintf("Agent%d", i))
		velocite := rand.Float64()
		acuite := 50.0
		position := validPositions[i]
		personalParameter := 0.1 + rand.Float64()*4.0 - 0.1
		// Cria mapa de carisma
		charisme := make(map[ag.IdAgent]float64)
		//charisme[id] = rand.Float64()

		// Cria mapa de relações
		relation := make(map[ag.IdAgent]float64)

		// Define estratégia de movimento
		var strategy ag.MovementStrategy
		switch TypeChoosen {
		case ag.Believer:
			strategy = config.BelieverMovement
			agentsImg = loadImage(AssetsPath + AgentBelieverImageFile)
		case ag.Sceptic:
			strategy = config.ScepticMovement
			agentsImg = loadImage(AssetsPath + AgentScepticImageFile)
		case ag.Neutral:
			strategy = config.NeutralMovement
			agentsImg = loadImage(AssetsPath + AgentNeutralImageFile)
		}

		// Cria o agente usando NewAgent
		agent := ag.NewAgent(
			env,
			id,
			velocite,
			acuite,
			position,
			Opinion,
			charisme,
			relation,
			personalParameter,
			TypeChoosen,
			make(chan ag.Message),
			agentsImg,
		)

		// Configura campos adicionais
		agent.HeatMap = visitationMap
		agent.MovementStrategy = strategy
		agents[i] = *agent
		env.AddAgent(&agents[i])
	}
	env.SetRelations()
	env.SetPoids()
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

func generateStatues(tilemapJSON *tile.TilemapJSON, tilesets []tile.Tileset) []image.Rectangle {
	var computersPositions []image.Rectangle
	for layerIdx, layer := range tilemapJSON.Layers {
		for i, tileID := range layer.Data {
			if tileID == 0 || layerIdx == 0 || layerIdx == 1 || layerIdx == 2 || layerIdx == 4 {
				continue
			}
			
			x, y := (i%layer.Width)*TileSize, (i/layer.Width)*TileSize
			img := tilesets[layerIdx].Img(tileID)
			offsetY := -(img.Bounds().Dy() + TileSize)
			y += offsetY
			computersPositions = append(computersPositions, image.Rect(x, y, x+img.Bounds().Dx(), y+img.Bounds().Dy()))
		}
	}
	return computersPositions
}

func generateComputers(tilemapJSON *tile.TilemapJSON, tilesets []tile.Tileset) []image.Rectangle {
	var computersPositions []image.Rectangle
	for layerIdx, layer := range tilemapJSON.Layers {
		for i, tileID := range layer.Data {
			if tileID == 0 || layerIdx == 0 || layerIdx == 1 || layerIdx == 4 || layerIdx == 3 {
				continue
			}
			
			x, y := (i%layer.Width)*TileSize, (i/layer.Width)*TileSize
			img := tilesets[layerIdx].Img(tileID)
			offsetY := -(img.Bounds().Dy() + TileSize)
			y += offsetY
			computersPositions = append(computersPositions, image.Rect(x, y, x+img.Bounds().Dx(), y+img.Bounds().Dy()))
		}
	}
	
	return computersPositions
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

	screen.Fill(color.RGBA{57, 61, 125, 255})
	sim.drawMap(screen)
	sim.drawAgents(screen)
	sim.drawAcuite(screen)
	sim.drawColliders(screen)
	sim.drawInfoPanel(screen)
	sim.drawSelectionIndicator(screen)
}

func (sim *Simulation) drawSelectionIndicator(screen *ebiten.Image) {
	if sim.selected != nil {
		x := sim.selected.Position.X
		y := sim.selected.Position.Y
		width := float64(AgentImageSize)
		height := float64(AgentImageSize)
		vector.StrokeRect(screen, float32(x), float32(y), float32(width), float32(height), 2, color.RGBA{255, 255, 0, 255}, false)
	} else if sim.selectedPC != nil {
		x := sim.selectedPC.Position.X
		y := sim.selectedPC.Position.Y

		width := float64(sim.carte.Ordinateurs[0].Max.X - sim.carte.Ordinateurs[0].Min.X)
		height := float64(sim.carte.Ordinateurs[0].Max.Y - sim.carte.Ordinateurs[0].Min.Y)

		vector.StrokeRect(screen, float32(x), float32(y), float32(width), float32(height), 2, color.RGBA{255, 255, 0, 255}, false)
	}
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

	// Dessine le panneau de fond
	vector.DrawFilledRect(screen, float32(panelX), float32(panelY), float32(panelWidth), float32(panelHeight), color.RGBA{0, 0, 0, 180}, false)

	// Titre du panneau
	ebitenutil.DebugPrintAt(screen, "Simulation Info", panelX+padding, panelY+padding)

	y := panelY + 30

	// Informations de la simulation
	elapsed := time.Since(sim.start)
	simInfo := fmt.Sprintf("Elapsed: %s", elapsed.Round(time.Second))
	ebitenutil.DebugPrintAt(screen, simInfo, panelX+padding, y)
	y += 40

	// Comptage des agents par type
	ebitenutil.DebugPrintAt(screen, "Agent Count:", panelX+padding, y)
	y += 20
	agentTypes := []ag.TypeAgent{ag.Sceptic, ag.Believer, ag.Neutral}
	for _, agentType := range agentTypes {
		count, _ := sim.env.NbrAgents.Load(agentType)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("  %s: %d", agentType, count), panelX+padding, y)
		y += 20
	}
	y += 20

	// Comptage des ordinateurs par programme None ou Go
	ebitenutil.DebugPrintAt(screen, "Computer Count:", panelX+padding, y)
	y += 20
	computerTypes := []string{"None", "Go"}
	pcs := sim.env.Objs
	for _, computerType := range computerTypes {
		count := 0
		for _, pc := range pcs {
			if pc.GetType() != ag.ComputerType {
				continue
			}
			if string(pc.GetProgramm()) == computerType {
				count++
			}
		}
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("  %s: %d", computerType, count), panelX+padding, y)
		y += 20
	}
	y += 20

	// Informations de l'agent sélectionné
	if sim.selected != nil {
		ebitenutil.DebugPrintAt(screen, "Selected Agent:", panelX+padding, y)
		y += 20
		agentInfo := fmt.Sprintf("  ID: %s\n  Type: %s\n  SubType: %s\n  Personal Param: %.2f\n  Opinion: %.2f\n  Alive: %t\n  DialogTimer : %d\n  CurrentAction : %s\n  Time to change direction : %d \n  Occupied : %t\n  Last Prayer : %d",
			sim.selected.Id,
			sim.selected.TypeAgt,
			sim.selected.SubType,
			sim.selected.PersonalParameter,
			sim.selected.Opinion,
			sim.selected.Vivant,
			sim.selected.DialogTimer,
			sim.selected.CurrentAction,
			0,
			sim.selected.Occupied,
			sim.selected.TimeLastStatue,
		)
		ebitenutil.DebugPrintAt(screen, agentInfo, panelX+padding, y)
		y += 180

		// Informações sobre discussão atual
		if sim.selected.DiscussingWith != nil {
			discussInfo := fmt.Sprintf("  Discussing with:\n  ID: %s\n  Type: %s",
				sim.selected.DiscussingWith.Id,
				sim.selected.DiscussingWith.TypeAgt,
			)
			ebitenutil.DebugPrintAt(screen, discussInfo, panelX+padding, y)
			y += 60
		}

		// Histórico de conversas
		ebitenutil.DebugPrintAt(screen, "  Last conversations with:", panelX+padding, y)
		y += 20
		for i, lastTalked := range sim.selected.LastTalkedTo {
			talkInfo := fmt.Sprintf("    %d. %s (%s)", i+1, lastTalked.Id, lastTalked.TypeAgt)
			ebitenutil.DebugPrintAt(screen, talkInfo, panelX+padding, y)
			y += 15
		}
		y += 20

		// Relações com outros agentes
		ebitenutil.DebugPrintAt(screen, "Relations:", panelX+padding, y)
		y += 20

		// Colete as chaves do mapa
		keys := make([]string, 0, len(sim.selected.Relation))
		for otherId := range sim.selected.Relation {
			keys = append(keys, string(otherId))
		}

		// Ordene as chaves
		sort.Strings(keys)

		// Itere sobre as chaves ordenadas
		for _, otherId := range keys {
			relation := sim.selected.Relation[ag.IdAgent(otherId)]
			relationType := getRelationType(relation)
			ebitenutil.DebugPrintAt(
				screen,
				fmt.Sprintf("  %s: %.2f %s ", otherId, relation, relationType),
				panelX+padding,
				y,
			)
			y += 15
		}
	}

	// Informations de l'ordinateur sélectionné
	if sim.selectedPC != nil {
		ebitenutil.DebugPrintAt(screen, "Selected Computer:", panelX+padding, y)
		y += 20
		pcInfo := fmt.Sprintf("  ID: %s\n  Used: %t\n  Program : %s",
			sim.selectedPC.Id,
			sim.selectedPC.Used,
			sim.selectedPC.Programm,
		)
		ebitenutil.DebugPrintAt(screen, pcInfo, panelX+padding, y)
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

	// Primeiro, desenha as linhas de conexão entre agentes em discussão
	for _, agent := range sim.agents {
		if agent.CurrentAction == "Discussing" && agent.DiscussingWith != nil {
			// Desenha uma linha conectando os agentes que estão discutindo
			startX := agent.Position.X + float64(AgentImageSize)/2
			startY := agent.Position.Y + float64(AgentImageSize)/2
			endX := agent.DiscussingWith.Position.X + float64(AgentImageSize)/2
			endY := agent.DiscussingWith.Position.Y + float64(AgentImageSize)/2

			// Escolhe a cor da linha baseado nos tipos dos agentes
			lineColor := color.RGBA{150, 150, 150, 255}
			vector.StrokeLine(
				screen,
				float32(startX),
				float32(startY),
				float32(endX),
				float32(endY),
				1,
				lineColor,
				false,
			)
		}
	}

	// Depois desenha os agentes
	for _, agent := range sim.agents {
		opts.GeoM.Reset()
		opts.GeoM.Translate(agent.Position.X, agent.Position.Y)

		// Removido o efeito de luz para agentes em discussão
		subImg := agent.Img.SubImage(image.Rect(0, 0, AgentImageSize, AgentImageSize)).(*ebiten.Image)
		screen.DrawImage(subImg, &opts)

		sim.drawDialogBox(screen, agent)
	}
}

func (sim *Simulation) drawDialogBox(screen *ebiten.Image, agent ag.Agent) {
	if agent.CurrentAction == "" || agent.DialogTimer <= 0 {
		return
	}

	dialogWidth := DiscussionBubbleWidth
	dialogHeight := DiscussionBubbleHeight
	x := int(agent.Position.X) - dialogWidth/2 + AgentImageSize/2
	y := int(agent.Position.Y) - dialogHeight - 5

	// Desenha o fundo da caixa de diálogo
	bgColor := color.RGBA{255, 255, 255, 200}

	// Muda a cor do fundo baseado na ação
	switch agent.CurrentAction {
	case "Discussing":
		// Cor diferente para discussão
		switch agent.TypeAgt {
		case ag.Believer:
			bgColor = color.RGBA{200, 230, 255, 200} // Azul claro
		case ag.Sceptic:
			bgColor = color.RGBA{255, 200, 200, 200} // Vermelho claro
		case ag.Neutral:
			bgColor = color.RGBA{200, 255, 200, 200} // Verde claro
		}
	case "Praying":
		bgColor = color.RGBA{255, 255, 200, 200} // Amarelo claro
	case "Using Computer":
		bgColor = color.RGBA{200, 200, 255, 200} // Roxo claro
	}

	// Desenha o fundo da caixa de diálogo
	vector.DrawFilledRect(screen, float32(x), float32(y), float32(dialogWidth), float32(dialogHeight), bgColor, false)

	// Desenha a borda da caixa de diálogo
	vector.StrokeRect(screen, float32(x), float32(y), float32(dialogWidth), float32(dialogHeight), 1, color.Black, false)

	// Prepara o texto baseado na ação
	displayText := agent.CurrentAction
	if agent.CurrentAction == "Discussing" {
		// Adiciona um indicador do tipo de agente na discussão
		displayText = fmt.Sprintf("%s (%s)", agent.CurrentAction, agent.TypeAgt)
	}

	// Desenha o texto da ação
	text.Draw(screen, displayText, sim.dialogFont, x+5, y+20, color.Black)

	// Adiciona uma barra de progresso para o DialogTimer
	if agent.DialogTimer > 0 {
		progressWidth := float32(dialogWidth-10) * (float32(agent.DialogTimer) / 180.0)
		vector.DrawFilledRect(
			screen,
			float32(x+5),
			float32(y+dialogHeight-10),
			progressWidth,
			5,
			color.RGBA{100, 100, 100, 200},
			false,
		)
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
	select {
	case <-sim.ctx.Done():
		return ebiten.Termination
	default:
		// Detecção de cliques
		cursorX, cursorY := ebiten.CursorPosition()

		// Detecta clique do mouse
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			// Verifica clique em agentes
			for i := range sim.agents {
				agent := &sim.agents[i]
				if cursorX >= int(agent.Position.X) &&
					cursorX <= int(agent.Position.X+AgentImageSize) &&
					cursorY >= int(agent.Position.Y) &&
					cursorY <= int(agent.Position.Y+AgentImageSize) {
					sim.selected = sim.env.GetAgentById(agent.Id)
					sim.selectedPC = nil
					sim.selectionIndicator = ebiten.NewImage(AgentImageSize, AgentImageSize)
					sim.selectionIndicator.Fill(color.RGBA{255, 255, 0, 128})
					break
				}
			}

			// Verifica clique em computadores
			for i := range sim.objets {
				if sim.objets[i].GetType() != ag.ComputerType {
					continue
				}

				pc := &sim.objets[i]
				if computer, ok := (*pc).(*ag.Computer); ok {
					if cursorX >= int(computer.ObjPosition().X) &&
						cursorX <= int(computer.ObjPosition().X+TileSize) &&
						cursorY >= int(computer.ObjPosition().Y) &&
						cursorY <= int(computer.ObjPosition().Y+TileSize) {
						sim.selectedPC = computer
						sim.selected = nil
						sim.selectionIndicator = ebiten.NewImage(TileSize, TileSize)
						sim.selectionIndicator.Fill(color.RGBA{255, 255, 0, 128})
						break
					}
				}
			}
		}

		// Atualização dos timers e estados
		for i := range sim.agents {
			if sim.agents[i].TimeLastStatue >= 0 && sim.agents[i].TypeAgt == ag.Believer {
				sim.agents[i].TimeLastStatue++
			}
			if sim.agents[i].DialogTimer > 0 {
				sim.agents[i].DialogTimer--
				if sim.agents[i].DialogTimer == 0 {
					sim.agents[i].ClearAction()
					if sim.agents[i].UseComputer != nil {
						sim.agents[i].UseComputer.Release()
						sim.agents[i].UseComputer = nil
					}
					sim.agents[i].Occupied = false
				}
			}
		}

		// Atualiza o agente selecionado se existir
		if sim.selected != nil {
			// Atualiza a referência para garantir dados atualizados
			sim.selected = sim.env.GetAgentById(sim.selected.Id)
		}

		// Calcular a média de opiniões
		totalOpinion := 0.0
		for _, agent := range sim.agents {
			totalOpinion += agent.Opinion
		}
		averageOpinion := totalOpinion / float64(len(sim.agents))
		sim.opinionAverages = append(sim.opinionAverages, averageOpinion)
	}
	return nil
}

func (sim *Simulation) Run() error {
	defer sim.cancel()

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

	// Gerar e salvar o gráfico
	xValues := make([]float64, len(sim.opinionAverages))
	for i := range xValues {
		xValues[i] = float64(i)
	}

	graph := chart.Chart{
		Series: []chart.Series{
			chart.ContinuousSeries{
				XValues: xValues,
				YValues: sim.opinionAverages,
			},
		},
	}

	file, err := os.Create("opinion_averages.png")
	if err != nil {
		return err
	}
	defer file.Close()
	err = graph.Render(chart.PNG, file)
	if err != nil {
		return err
	}

	return nil
}

func getRelationType(relation float64) string {
	switch {
	case relation == 0.75:
		return "Ennemi"
	case relation == 1.0:
		return "Pas de liens direct"
	case relation == 1.25:
		return "Amis"
	case relation == 1.5:
		return "Famille"
	default:
		return "Inconnu"
	}
}
