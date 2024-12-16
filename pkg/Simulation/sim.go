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

type Simulation struct {
	env         ag.Environnement
	agents      []ag.Agent
	maxStep     int
	maxDuration time.Duration
	step        int //Stats
	start       time.Time
	syncChans   sync.Map
	carte       carte.Carte
}


func NewSimulation(maxStep int, maxDuration time.Duration) *Simulation {
	
	const numAgents = 5
	agents := make([]ag.Agent, numAgents)


	// Creation de l'environnement
	env := *ag.NewEnvironment(agents)

	// Creation de la carte
	tilemapImg, _, err := ebitenutil.NewImageFromFile("assets/images/img.png")
	if err != nil {
		// handle error
		log.Fatal(err)
	}

	tilemapJSON, err := tile.NewTilemapJSON("assets/maps/spawn.json")
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


	carte := carte.NewCarte(*tilemapJSON, tilesets, tilemapImg, coliders)

	// Creation des agents
	agentsImg, _, err := ebitenutil.NewImageFromFile("assets/images/ninja.png")
	if err != nil {
		// handle error
		log.Fatal(err)
	}

	
	for i := 0; i < numAgents; i++ {

		agents[i] = ag.Agent{
			Env:               &env,
			Id:                ag.IdAgent(fmt.Sprintf("Agent%d", i)),
			Velocite:          rand.Float64(),
			Acuite:            rand.Float64(),
			Position:          ut.Position{X: rand.Float64(), Y: rand.Float64()},
			Opinion:           rand.Float64(),
			Charisme:          make(map[ag.IdAgent]float64),
			Relation:          make(map[ag.IdAgent]float64),
			PersonalParameter: rand.Float64(),
			Poid_rel:          []float64{rand.Float64(), rand.Float64()},
			Vivant:            true,
			TypeAgt:           ag.TypeAgent(rand.Intn(3)),
			SyncChan:          make(chan int),
			Img:               agentsImg,
		}
		env.AddAgent(agents[i])
	}
	
	return &Simulation{env: env, agents: agents, maxStep: maxStep, maxDuration: maxDuration, step: 0, start: time.Now(), carte: *carte}
}


func (g *Simulation) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{120, 180, 255, 255})

	opts := ebiten.DrawImageOptions{}

	// draw the tilemap
	for layerIdx, layer := range g.carte.TilemapJSON.Layers {
		for i, tileID := range layer.Data {
			if tileID == 0 {
				continue
			}
			
			x := i % layer.Width
			y := i / layer.Width
			
			y *= 24
			x *= 24


			img := g.carte.Tilesets[layerIdx].Img(tileID)
			
			opts.GeoM.Translate(float64(x), float64(y))

			opts.GeoM.Translate(0.0, -(float64(img.Bounds().Dy()) + 24.0))

			screen.DrawImage(
				img,
				&opts,
			)

			opts.GeoM.Reset()
		}
	}

	// draw the agents
	for _, agent := range g.agents {
		opts.GeoM.Translate(agent.Position.X, agent.Position.Y)
		screen.DrawImage(
			agent.Img.SubImage(
				image.Rect(0, 0, 16, 16),
			).(*ebiten.Image),
			&opts,
		)

		opts.GeoM.Reset()
	}

	//draw coliders
	for _, colider := range g.carte.Coliders {
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

func (g *Simulation) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func (sim *Simulation) Update() error {	
	return nil
}




func (sim *Simulation) Run(){

	
	for _, ag := range sim.agents {
		ag.Start()
	}
	
	sim.start = time.Now()

	
	for _, agent := range sim.agents {
		go func(agent ag.Agent) {
			step := 0
			for {
				step++
				c , _ := sim.syncChans.Load(agent.Id)
				c.(chan int) <- step
				time.Sleep(1 * time.Millisecond)
				<-c.(chan int)
			}
		}(agent)
		time.Sleep(sim.maxDuration)
	}	
}
