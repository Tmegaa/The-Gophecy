package simulation

import (
	ag "Gophecy/pkg/Agent"
	carte "Gophecy/pkg/Carte"
	tile "Gophecy/pkg/Tile"
	ut "Gophecy/pkg/Utilitaries"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"sort"
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

// Constantes de la simulation
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
	env                *ag.Environnement
	agents             []*ag.Agent
	objets             []ag.InterfaceObjet
	maxDuration        time.Duration
	start              time.Time
	carte              *carte.Carte
	selected           *ag.Agent
	selectedPC         *ag.Computer
	ctx                context.Context
	cancel             context.CancelFunc
	dialogFont         font.Face
	selectionIndicator *ebiten.Image
	opinionAverages    []float64
}

// Fonction qui initialize une nouvelle simulation
func NewSimulation(config SimulationConfig) *Simulation {
	initializeWindow()
	carte := loadMap()
	env := createEnvironment(carte)
	var agents []*ag.Agent
	var err error

	if config.AgentsFilePath != "" {
		// Load agents from file
		agents, err = createAgentsFromFile(env, config.AgentsFilePath)
		if err != nil {
			log.Fatalf("Failed to create agents from file: %v", err)
		}
	} else {
		// Create agents normally
		agents = createAgents(env, carte, config)
	}

	obj := loadObjects(env)
	ctx, cancel := context.WithTimeout(context.Background(), config.SimulationTime)
	tt, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal(err)
	}

	selectionIndicator := ebiten.NewImage(TileSize, TileSize)
	selectionIndicator.Fill(color.RGBA{255, 255, 0, 128})

	return &Simulation{
		env:    env,
		agents: agents,
		objets: obj,
		// maxStep:     10,
		maxDuration: config.SimulationTime,
		start:       time.Now(),
		carte:       carte,
		ctx:         ctx,
		cancel:      cancel,
		dialogFont: truetype.NewFace(tt, &truetype.Options{
			Size: 12,
			DPI:  72,
		}),
		selectionIndicator: selectionIndicator,
	}
}

// Fonction qui initialise la fenêtre d'affichage de la simulation
func initializeWindow() {
	ebiten.SetWindowSize(WindowWidth, WindowHeight)
	ebiten.SetWindowTitle("Simulation")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
}

// Fonction qui crée et retourne un nouvel environnement
func createEnvironment(carte *carte.Carte) *ag.Environnement {
	return ag.NewEnvironment(make([]*ag.Agent, 0), carte, make([]ag.InterfaceObjet, 0))
}

// Fonction qui charge la carte
func loadMap() *carte.Carte {
	tilemapImg := loadImage(AssetsPath + TilemapImage)
	tilemapJSON := loadTilemapJSON(MapsPath + TilemapJSONFile)
	tilesets := loadTilesets(tilemapJSON)
	coliders := generateColliders(tilemapJSON, tilesets)
	computers := generateComputers(tilemapJSON, tilesets)
	statues := generateStatues(tilemapJSON, tilesets)

	return carte.NewCarte(*tilemapJSON, tilesets, tilemapImg, coliders, computers, statues)
}

// Fonction qui charge les objets dans la carte
func loadObjects(env *ag.Environnement) []ag.InterfaceObjet {
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

// Fonction qui renvoie une liste des positions possibles que peuvent prendre les objets ou les agents
func getValidSpawnPositions(carte *carte.Carte) []ut.Position {
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

// Fonction qui crée et rajoute à la carte les nouveaux agents
func createAgents(env *ag.Environnement, carte *carte.Carte, config SimulationConfig) []*ag.Agent {
	agents := make([]*ag.Agent, config.NumAgents)
	validPositions := getValidSpawnPositions(carte)
	visitationMap := ag.NewVisitationMap(validPositions)

	if len(validPositions) < config.NumAgents {
		log.Fatalf("Not enough valid spawn positions for all agents")
	}

	rand.Shuffle(len(validPositions), func(i, j int) {
		validPositions[i], validPositions[j] = validPositions[j], validPositions[i]
	})

	// Image par défaut
	agentsImg := loadImage(AssetsPath + AgentBelieverImageFile)

	for i := 0; i < config.NumAgents; i++ {
		// Génère des valeurs aléatoires en rescpectant les contraintes de type s'il y en a
		
		var Opinion float64
		if i < config.NumBelievers {
			Opinion = rand.Float64()*(1./3.) + 2./3.
		} else if i < config.NumBelievers + config.NumSceptics {
			Opinion = rand.Float64()*(1./3.)
		} else if i < config.NumBelievers + config.NumSceptics + config.NumNeutrals {
			Opinion = rand.Float64()*(1./3.) + 1./3.
		} else {
			Opinion = rand.Float64()
		}

		// Détermine le type de base de l'agent
		var TypeChoosen ag.TypeAgent
		if Opinion > 2./3. {
			TypeChoosen = ag.Believer
		} else if Opinion > 1./3. {
			TypeChoosen = ag.Neutral
		} else {
			TypeChoosen = ag.Sceptic
		}
		id := ag.IdAgent(fmt.Sprintf("Agent%d", i))
		velocite := rand.Float64()
		acuite := 50.0
		position := validPositions[i]
		personalParameter := 0.1 + rand.Float64()*4.0 - 0.1
		// Crée une carte de charisme
		charisme := make(map[ag.IdAgent]float64)


		// Crée une carte des relations entre agents
		relation := make(map[ag.IdAgent]float64)

		// Définit la stratégie de mouvement
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

		// Créez l'agent à l'aide de NewAgent
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

		// Configure les champs supplémentaires
		agent.HeatMap = visitationMap
		agent.MovementStrategy = strategy
		agents[i] = agent
		env.AddAgent(agent)
	}
	env.SetRelations()
	env.SetPoids()
	return env.Ags
}

// Structure pour contenir les données des agents à partir d'un fichier
type AgentData struct {
	Id                string  `json:"id"`
	Opinion           float64 `json:"opinion"`
	Charisme          map[string]float64 `json:"charisme"`
	Relation          map[string]float64 `json:"relation"`
	PersonalParameter float64 `json:"personalParameter"`
	SubType           string  `json:"subType"`
}

// Fonction pour créer des agents à partir d'un fichier
func createAgentsFromFile(env *ag.Environnement, filePath string) ([]*ag.Agent, error) {
	// Lire le contenu du fichier
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("échec de la lecture du fichier: %v", err)
	}

	// Analyser les données JSON
	var agentsData []AgentData
	err = json.Unmarshal(data, &agentsData)
	if err != nil {
		return nil, fmt.Errorf("échec de l'analyse des données JSON: %v", err)
	}

	// Générer des positions valides
	carte := env.Carte
	validPositions := getValidSpawnPositions(carte)
	if len(validPositions) < len(agentsData) {
		return nil, fmt.Errorf("pas assez de positions de spawn valides pour tous les agents")
	}

	rand.Shuffle(len(validPositions), func(i, j int) {
		validPositions[i], validPositions[j] = validPositions[j], validPositions[i]
	})

	// Créer des agents à partir des données analysées
	agents := make([]*ag.Agent, len(agentsData))
	for i, agentData := range agentsData {
		id := ag.IdAgent(agentData.Id)
		charisme := make(map[ag.IdAgent]float64)
		for k, v := range agentData.Charisme {
			charisme[ag.IdAgent(k)] = v
		}
		relation := make(map[ag.IdAgent]float64)
		for k, v := range agentData.Relation {
			relation[ag.IdAgent(k)] = v
		}

		// Déterminer le type en fonction de l'opinion
		var typeAgt ag.TypeAgent
		if agentData.Opinion > 2.0/3.0 {
			typeAgt = ag.Believer
		} else if agentData.Opinion > 1.0/3.0 {
			typeAgt = ag.Neutral
		} else {
			typeAgt = ag.Sceptic
		}

		subType := ag.SubTypeAgent(agentData.SubType)

		// Générer la vélocité, l'acuité et la position
		velocite := rand.Float64()
		acuite := 50.0
		position := validPositions[i]

		// Charger l'image appropriée en fonction du type d'agent
		var agentsImg *ebiten.Image
		switch typeAgt {
		case ag.Believer:
			agentsImg = loadImage(AssetsPath + AgentBelieverImageFile)
		case ag.Sceptic:
			agentsImg = loadImage(AssetsPath + AgentScepticImageFile)
		case ag.Neutral:
			agentsImg = loadImage(AssetsPath + AgentNeutralImageFile)
		default:
			return nil, fmt.Errorf("type d'agent inconnu: %s", typeAgt)
		}

		if agentsImg == nil {
			return nil, fmt.Errorf("échec du chargement de l'image pour le type d'agent: %s", typeAgt)
		}

		// Créer l'agent
		agent := ag.NewAgent(
			env,
			id,
			velocite,
			acuite,
			position,
			agentData.Opinion,
			charisme,
			relation,
			agentData.PersonalParameter,
			typeAgt,
			make(chan ag.Message),
			agentsImg,
		)
		agent.SubType = subType

		agents[i] = agent
		env.AddAgent(agent)
	}
	env.SetPoids()
	return agents, nil
}

// Fonction qui affiche une image sur la fenêtre d'affichage
func loadImage(path string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatalf("Failed to load image: %s, error: %v", path, err)
	}
	return img
}

// Fonction qui charge la carte des tiles
func loadTilemapJSON(path string) *tile.TilemapJSON {
	tilemap, err := tile.NewTilemapJSON(path)
	if err != nil {
		log.Fatalf("Failed to load tilemap JSON: %s, error: %v", path, err)
	}
	return tilemap
}

// Fonction qui charge les jeux de tiles
func loadTilesets(tilemapJSON *tile.TilemapJSON) []tile.Tileset {
	tilesets, err := tilemapJSON.GenTilesets()
	if err != nil {
		log.Fatalf("Failed to generate tilesets, error: %v", err)
	}
	return tilesets
}

// Fonction de génération des statues
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

// Fonction de génération des ordinateurs
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

// Fonction de génération des bornes d'une zone de collision (où deux éléments ne peuvent pas coexister)
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

// Fonction qui affiche les éléments dans la fenêtre d'affichage
func (sim *Simulation) Draw(screen *ebiten.Image) {
	// Dessine l'arrière-plan et les agents rgba(57,61,125,255)
	screen.Fill(color.RGBA{57, 61, 125, 255})
	sim.drawMap(screen)
	sim.drawAgents(screen)
	sim.drawAcuite(screen)
	sim.drawColliders(screen)
	sim.drawInfoPanel(screen)
	sim.drawSelectionIndicator(screen)
}

// Fonction qui affiche dans la fenêtre d'affichage un indicateur de la séléction de l'utilisateur
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

// Fonction qui affiche la zone où un agent peut percevoir d'autres agents ou des objets (rectangle)
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

// Fonction qui affiche les informations de la simulation (nombre d'agents, temps écoulé...) dans un cadre de la fenêtre d'affichage
func (sim *Simulation) drawInfoPanel(screen *ebiten.Image) {
	panelX, panelY := 0, 0
	panelWidth, panelHeight := 240, WindowHeight-20
	padding := 10

	// Dessine le panneau de fond
	vector.DrawFilledRect(screen, float32(panelX), float32(panelY), float32(panelWidth), float32(panelHeight), color.RGBA{0, 0, 0, 180}, false)

	// Titre du panneau
	ebitenutil.DebugPrintAt(screen, "Informations de la simulation", panelX+padding, panelY+padding)

	y := panelY + 30

	// Informations de la simulation
	elapsed := time.Since(sim.start)
	simInfo := fmt.Sprintf("Temps écoulé: %s", elapsed.Round(time.Second))
	ebitenutil.DebugPrintAt(screen, simInfo, panelX+padding, y)
	y += 40

	// Nombre d'agents par type
	ebitenutil.DebugPrintAt(screen, "Nombre d'agents:", panelX+padding, y)
	y += 20
	agentTypes := []ag.TypeAgent{ag.Sceptic, ag.Believer, ag.Neutral}
	for _, agentType := range agentTypes {
		count, _ := sim.env.NbrAgents.Load(agentType)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("  %s: %d", agentType, count), panelX+padding, y)
		y += 20
	}
	y += 20

	// Comptage des ordinateurs par programme None ou Go
	ebitenutil.DebugPrintAt(screen, "Nombre d'ordinateurs:", panelX+padding, y)
	y += 20
	computerTypes := []ag.Programm{ag.NoPgm, ag.GoPgm}
	pcs := sim.env.Objs
	for _, computerType := range computerTypes {
		count := 0
		for _, pc := range pcs {
			if pc.GetType() != ag.ComputerType {
				continue
			}
			if pc.GetProgramm() == computerType {
				count++ 
			}
		}
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("  %s: %d", computerType, count), panelX+padding, y)
		y += 20
	}
	y += 20

	// Informations de l'agent sélectionné
	if sim.selected != nil {
		ebitenutil.DebugPrintAt(screen, "Agent sélectionné:", panelX+padding, y)
		y += 20
		agentInfo := fmt.Sprintf("  ID: %s\n  Type: %s\n  Sous-Type: %s\n  Paramètre Personnel: %.2f\n  Opinion: %.2f\n  Vivant: %t\n  Temps de Dialogue: %d\n  Action: %s\n  Stratégie de mouvement: %s\n  Occupé : %t\n  Dernière prière : %.2f",
			sim.selected.Id,
			sim.selected.TypeAgt,
			sim.selected.SubType,
			sim.selected.PersonalParameter,
			sim.selected.Opinion,
			sim.selected.Vivant,
			sim.selected.DialogTimer,
			sim.selected.CurrentAction,
			sim.selected.MovementStrategy,
			sim.selected.Occupied,
			time.Since(sim.selected.TimeLastStatue).Seconds(),
		)
		ebitenutil.DebugPrintAt(screen, agentInfo, panelX+padding, y)
		y += 180

		// Informations sur la discussion actuelle
		if sim.selected.DiscussingWith != nil {
			discussInfo := fmt.Sprintf("  En discussion avec:\n  ID: %s\n  Type: %s",
				sim.selected.DiscussingWith.Id,
				sim.selected.DiscussingWith.TypeAgt,
			)
			ebitenutil.DebugPrintAt(screen, discussInfo, panelX+padding, y)
			y += 60
		}

		// Historique des conversations
		ebitenutil.DebugPrintAt(screen, "  Dernières conversations avec:", panelX+padding, y)
		y += 20
		for i, lastTalked := range sim.selected.LastTalkedTo {
			talkInfo := fmt.Sprintf("    %d. %s (%s)", i+1, lastTalked.Id, lastTalked.TypeAgt)
			ebitenutil.DebugPrintAt(screen, talkInfo, panelX+padding, y)
			y += 15
		}
		y += 20

		// Relations avec les autres agents
		ebitenutil.DebugPrintAt(screen, "Relations:", panelX+padding, y)
		y += 20

		// Récupére les clés de la carte
		keys := make([]string, 0, len(sim.selected.Relation))
		for otherId := range sim.selected.Relation {
			keys = append(keys, string(otherId))
		}

		// Trie les clés
		sort.Strings(keys)

		// Itére sur les clés triées
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
		ebitenutil.DebugPrintAt(screen, "Ordinateur sélectionné:", panelX+padding, y)
		y += 20
		pcInfo := fmt.Sprintf("  ID: %s\n  En utilisation: %t\n  Langage de programmation: %s",
			sim.selectedPC.Id,
			sim.selectedPC.Used,
			sim.selectedPC.Programm,
		)
		ebitenutil.DebugPrintAt(screen, pcInfo, panelX+padding, y)
	}
}

// Fonction d'affichage de la carte dans la fenêtre d'affichage
func (sim *Simulation) drawMap(screen *ebiten.Image) {
	opts := ebiten.DrawImageOptions{}
	// Gestion par couche
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

// Fonction d'affichage des agents dans la fenêtre d'affichage
func (sim *Simulation) drawAgents(screen *ebiten.Image) {
	opts := ebiten.DrawImageOptions{}

	// Premièrement, il trace les lignes de connexion entre les agents en discussion.
	for _, agent := range sim.agents {
		if agent.CurrentAction == ag.DiscussAct && agent.DiscussingWith != nil {
			// Trace une ligne reliant les agents qui se discutent
			startX := agent.Position.X + float64(AgentImageSize)/2
			startY := agent.Position.Y + float64(AgentImageSize)/2
			endX := agent.DiscussingWith.Position.X + float64(AgentImageSize)/2
			endY := agent.DiscussingWith.Position.Y + float64(AgentImageSize)/2

			// Choisissez la couleur de la ligne en fonction des types d'agents
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

	// Puis affiche les agents
	for _, agent := range sim.agents {
		opts.GeoM.Reset()
		opts.GeoM.Translate(agent.Position.X, agent.Position.Y)

		// Suppression de l'effet d'éclairage pour les agents en discussion
		subImg := agent.Img.SubImage(image.Rect(0, 0, AgentImageSize, AgentImageSize)).(*ebiten.Image)
		screen.DrawImage(subImg, &opts)

		sim.drawDialogBox(screen, *agent)
	}
}

// Fonction d'affichage des boîtes de dialogue dans la fenêtre d'affichage
func (sim *Simulation) drawDialogBox(screen *ebiten.Image, agent ag.Agent) {
	if agent.CurrentAction == "" || agent.DialogTimer <= 0 {
		return
	}

	dialogWidth := DiscussionBubbleWidth
	dialogHeight := DiscussionBubbleHeight
	x := int(agent.Position.X) - dialogWidth/2 + AgentImageSize/2
	y := int(agent.Position.Y) - dialogHeight - 5

	// Dessiner l'arrière-plan de la boîte de dialogue
	bgColor := color.RGBA{255, 255, 255, 200}

	// Changer la couleur d'arrière-plan en fonction de l'action
	switch agent.CurrentAction {
	case ag.DiscussAct:
		// Couleur différente pour chaque type de discussion
		switch agent.TypeAgt {
		case ag.Believer:
			bgColor = color.RGBA{200, 230, 255, 200} // Bleu clair
		case ag.Sceptic:
			bgColor = color.RGBA{255, 200, 200, 200} // Rouge clair
		case ag.Neutral:
			bgColor = color.RGBA{200, 255, 200, 200} // Vert clair
		}
	case ag.PrayAct:
		bgColor = color.RGBA{255, 255, 200, 200} // Jaune clair
	case ag.ComputerAct:
		bgColor = color.RGBA{200, 200, 255, 200} // Violet clair
	}

	// Dessine l'arrière-plan de la boîte de dialogue
	vector.DrawFilledRect(screen, float32(x), float32(y), float32(dialogWidth), float32(dialogHeight), bgColor, false)

	// Dessine la bordure de la boîte de dialogue
	vector.StrokeRect(screen, float32(x), float32(y), float32(dialogWidth), float32(dialogHeight), 1, color.Black, false)

	// Prépare un texte basé sur l'action
	displayText := string(agent.CurrentAction)
	if agent.CurrentAction == ag.DiscussAct {
		// Ajoute un indicateur de type d'agent à la discussion
		displayText = fmt.Sprintf("%s (%s)", agent.CurrentAction, agent.TypeAgt)
	}

	// Écrit le texte de l'action
	text.Draw(screen, displayText, sim.dialogFont, x+5, y+20, color.Black)

	// Ajouter une barre de progression pour le DialogTimer
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

// Fonction d'affichage des bornes d'une zone de collision dans la fenêtre d'affichage
func (sim *Simulation) drawColliders(screen *ebiten.Image) {
	for _, colider := range sim.carte.Coliders {
		vector.StrokeRect(screen, float32(colider.Min.X), float32(colider.Min.Y), float32(colider.Dx()), float32(colider.Dy()), 1.0, color.RGBA{0, 0, 0, 0}, true)
	}
}

// Fonction qui retourne les dimentions de la fenêtre d'affichage
func (sim *Simulation) Layout(outsideWidth, outsideHeight int) (int, int) {
	return WindowWidth, WindowHeight
}

// Fonction de mise à jour de la simulation
func (sim *Simulation) Update() error {
	select {
	case <-sim.ctx.Done():
		return ebiten.Termination
	default:
		// Position du curseur
		cursorX, cursorY := ebiten.CursorPosition()

		// Détection d'un clic
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			// Vérifie si le clic se situe dans la zone agent
			for i := range sim.agents {
				agent := sim.agents[i]
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

			// Vérifie si le clic se situe sur un ordinateur
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

		// Mettre à jour du timer et les états
		for i := range sim.agents {
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

		// Met à jour l'agent sélectionné s'il existe
		if sim.selected != nil {
			// Met à jour la référence pour garantir que les données soient à jour
			sim.selected = sim.env.GetAgentById(sim.selected.Id)
		}

		// Calcule l'avis moyen
		totalOpinion := 0.0
		for _, agent := range sim.agents {
			totalOpinion += agent.Opinion
		}
		averageOpinion := totalOpinion / float64(len(sim.agents))
		sim.opinionAverages = append(sim.opinionAverages, averageOpinion)
	}
	return nil
}

// Fonction qui fait tourner la simulation
func (sim *Simulation) Run() error {
	defer sim.cancel() // On s'assure que le contexte est annulé lorsque Run() se termine

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
	fmt.Printf("Durée totale: %s\n", time.Since(sim.start).Round(time.Second))

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
	fmt.Printf("\nOpinion moyenne des agents: %.2f\n", averageOpinion)

	// Générer et enregistrer le graphique
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

// Fonction qui retourne le type de relation en fonction d'un nombre qui les identifie
func getRelationType(relation float64) string {
	switch {
	case relation == 0.75:
		return "Ennemi"
	case relation == 1.0:
		return "Pas de lien direct"
	case relation == 1.25:
		return "Amis"
	case relation == 1.5:
		return "Famille"
	default:
		return "Inconnu"
	}
}
