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
	NbrAgents *sync.Map //key = typeAgent et value = int  -> Compteur d'agents par types
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

func (env *Environnement) NearbyAgents() []Agent {
	env.RLock()
	defer env.RUnlock()
	for _, ag := range env.Ags {

		for _, ag2 := range env.Ags {
			if ag.ID() != ag2.ID() {
				pos2 := ag2.AgtPosition()
				var area ut.Rectangle
				area.PositionDL.X = pos2.X - ag2.Acuite
				area.PositionDL.Y = pos2.Y + ag.Acuite
				area.PositionUR.X = pos2.X + ag2.Acuite
				area.PositionUR.Y = pos2.Y - ag2.Acuite

			}
		}
	}
	return env.Ags
}
