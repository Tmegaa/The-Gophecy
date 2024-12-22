package pkg

import (
	carte "Gophecy/pkg/Carte"
	ut "Gophecy/pkg/Utilitaries"
	"sync"
)

var lsType = []TypeAgent{Sceptic, Believer, Neutral}

type Environnement struct {
	sync.RWMutex
	Ags   []Agent
	Carte carte.Carte
	//agts []Objet
	NbrAgents      *sync.Map //key = typeAgent et value = int  -> Compteur d'agents par types
	AgentProximity *sync.Map //key = IDAgent et value = []Agent -> Liste des agents proches
}

func NewEnvironment(ags []Agent, carte carte.Carte) (env *Environnement) {
	counter := &sync.Map{}

	for _, val := range lsType {
		counter.Store(val, 0)
	}

	return &Environnement{Ags: ags, NbrAgents: counter, Carte: carte, AgentProximity: &sync.Map{}}
}

func (env *Environnement) AddAgent(ag Agent) {
	env.Ags = append(env.Ags, ag)

	// Charger le nombre actuel d'agents de ce type
	value, exists := env.NbrAgents.Load(ag.TypeAgt)
	if exists {
		// Si la clé existe, convertir en int et incrémenter
		if nbr, ok := value.(int); ok {
			env.NbrAgents.Store(ag.TypeAgt, nbr+1)
		} else {
			// Si la valeur n'est pas un int, gérer l'erreur ou initialiser à 1
			env.NbrAgents.Store(ag.TypeAgt, 1)
		}
	} else {
		// Si la clé n'existe pas, initialiser à 1
		env.NbrAgents.Store(ag.TypeAgt, 1)
	}
}

// NearbyAgents calcule les agents proches de chaque agent
func (env *Environnement) NearbyAgents() {
	env.Lock()
	defer env.Unlock()
	for _, ag := range env.Ags {
		var nearbyAgents []*Agent
		pos := ag.AgtPosition()
		var area ut.Rectangle
		area.PositionDL.X = pos.X - ag.Acuite
		area.PositionDL.Y = pos.Y + ag.Acuite
		area.PositionUR.X = pos.X + ag.Acuite
		area.PositionUR.Y = pos.Y - ag.Acuite

		for _, ag2 := range env.Ags {
			if ag.ID() != ag2.ID() && ut.IsInRectangle(ag2.AgtPosition(), area) {
				nearbyAgents = append(nearbyAgents, &ag2)
			}
		}
		env.AgentProximity.Store(ag.Id, nearbyAgents)
	}
}
