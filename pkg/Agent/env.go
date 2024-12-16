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

	// Charger le nombre actuel d'agents de ce type
	value, exists := env.nbrAgents.Load(ag.TypeAgt)
	if exists {
		// Si la clé existe, convertir en int et incrémenter
		if nbr, ok := value.(int); ok {
			env.nbrAgents.Store(ag.TypeAgt, nbr+1)
		} else {
			// Si la valeur n'est pas un int, gérer l'erreur ou initialiser à 1
			env.nbrAgents.Store(ag.TypeAgt, 1)
		}
	} else {
		// Si la clé n'existe pas, initialiser à 1
		env.nbrAgents.Store(ag.TypeAgt, 1)
	}
}

