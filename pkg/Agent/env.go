package pkg

import (
	ut "Gophecy/pkg/Utilitaries"
	"sync"
)

var lsType = []TypeAgent{Sceptic, Believer, Neutral}

type Environnement struct {
	sync.RWMutex
	Ags []Agent
	//carte Carte
	//agts []Objet
	NbrAgents      *sync.Map //key = typeAgent et value = int  -> Compteur d'agents par types
	AgentProximity *sync.Map //key = Agent.ID et value = []Agent -> Liste des agents proches
}

func NewEnvironment(ags []Agent) (env *Environnement) {
	counter := &sync.Map{}

	for _, val := range lsType {
		counter.Store(val, 0)
	}

	return &Environnement{Ags: ags, NbrAgents: counter}
}

func (env *Environnement) AddAgent(ag Agent) {
	env.Ags = append(env.Ags, ag)
	nbr, err := env.NbrAgents.Load(ag.TypeAgt)
	if !err {
		nbr = nbr.(int) + 1
		env.NbrAgents.Store(ag.TypeAgt, nbr)
	} else {
		env.NbrAgents.Store(ag.TypeAgt, 1)
	}
}

// NearbyAgents calcule les agents proches de chaque agent
func (env *Environnement) NearbyAgents() {
	env.RLock()
	defer env.RUnlock()
	var nearbyAgents []Agent
	for _, ag := range env.Ags {
		pos := ag.AgtPosition()
		var area ut.Rectangle
		area.PositionDL.X = pos.X - ag.Acuite
		area.PositionDL.Y = pos.Y + ag.Acuite
		area.PositionUR.X = pos.X + ag.Acuite
		area.PositionUR.Y = pos.Y - ag.Acuite

		for _, ag2 := range env.Ags {
			if ag.ID() != ag2.ID() && ut.IsInRectangle(ag2.AgtPosition(), area) {
				nearbyAgents = append(nearbyAgents, ag2)
			}
		}
		env.AgentProximity.Store(ag.Id, nearbyAgents)
	}
}
