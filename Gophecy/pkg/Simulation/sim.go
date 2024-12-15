package simulation

import (
	ag "Gophecy/pkg/Agent"
	ut "Gophecy/pkg/Utilitaries"
	"fmt"
	"log"
	"sync"
	"time"

	"math/rand"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Simulation struct {
	env         ag.Environnement
	agents      []ag.Agent
	maxStep     int
	maxDuration time.Duration
	step        int //Stats
	start       time.Time
	syncChans   sync.Map
}


func NewSimulation(maxStep int, maxDuration time.Duration) *Simulation {
	
	const numAgents = 5
	agents := make([]ag.Agent, numAgents)


	// Creation de l'environnement
	env := *ag.NewEnvironment(agents)



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
	
	return &Simulation{env: env, agents: agents, maxStep: maxStep, maxDuration: maxDuration, step: 0, start: time.Now()}
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
