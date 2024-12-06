package pkg

import (
	"sync"
)

var lsType = []TypeAgent{Sceptic, Believer, Neutral}

type Environnement struct {
	sync.RWMutex
	ags []Agent
	//carte Carte
	//agts []Objet
	nbrAgents *sync.Map //key = typeAgent et value = int  -> Compteur d'agents par types
}

func NewEnvironment(ags []Agent) (env *Environnement) {
	counter := &sync.Map{}

	for _, val := range lsType {
		counter.Store(val, 0)
	}

	return &Environnement{ags: ags, nbrAgents: counter}
}
