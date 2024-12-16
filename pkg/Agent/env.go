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

func (env *Environnement) AddAgent(ag Agent) {
	env.ags = append(env.ags, ag)
	nbr, err := env.nbrAgents.Load(ag.TypeAgt)
	if !err {
		nbr = nbr.(int) + 1
		env.nbrAgents.Store(ag.TypeAgt, nbr)
	}else {
	env.nbrAgents.Store(ag.TypeAgt, 1)
	}	
}
