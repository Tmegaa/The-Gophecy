package pkg

import (
	ut "Gophecy/pkg/Utilitaries"
	"log"
)

type InterfaceAgent interface {
	Avancer()
	Discuter()
	Programmer()
	Perceive(*Environnement)
	Deliberate()
	//Act(*Environnement,typeFunc fun()) A implémenter
}

type TypeAgent string

const (
	Sceptic  TypeAgent = "Sceptic"
	Believer TypeAgent = "Believer"
	Neutral  TypeAgent = "Neutral"
)

type IdAgent string

type Agent struct {
	env               *Environnement
	id                IdAgent
	velocite          float64
	acuite            float64
	position          ut.Position
	opinion           float64
	charisme          map[IdAgent]float64 //influence d'un agent sur un autre
	relation          map[IdAgent]float64
	personalParameter float64
	poid_rel          []float64
	vivant            bool
	typeAgt           TypeAgent
	syncChan          chan int
}

func NewAgent(env *Environnement, id IdAgent, velocite float64, acuite float64, position ut.Position,
	opinion float64, charisme map[IdAgent]float64, relation map[IdAgent]float64, personalParameter float64,
	agent InterfaceAgent, typeAgt TypeAgent, syncChan chan int) *Agent {

	//calcul des poids relatif pour chaque agents
	poid_rel := make([]float64, 0)
	personalCharisme := charisme[id]
	for _, v := range charisme {
		char := v / personalCharisme
		poid_rel = append(poid_rel, char)
	}

	return &Agent{env: env, id: id, velocite: velocite, acuite: acuite,
		position: position, opinion: opinion, charisme: charisme, relation: relation,
		personalParameter: personalParameter, poid_rel: poid_rel,
		vivant: true, typeAgt: typeAgt, syncChan: syncChan}
}

func (agt *Agent) Position() ut.Position {
	return agt.position
}

func (agt *Agent) ID() IdAgent {
	return agt.id
}

func (agt *Agent) Opinion() float64 {
	return agt.opinion
}

func (agt *Agent) Charisme() map[IdAgent]float64 {
	return agt.charisme
}

func (agt *Agent) Relation() map[IdAgent]float64 {
	return agt.relation
}

func (agt *Agent) PersonalParameter() float64 {
	return agt.personalParameter
}

func (agt *Agent) Vivant() bool {
	return agt.vivant
}

func (agt *Agent) TypeAgt() TypeAgent {
	return agt.typeAgt
}

func (ag *Agent) Start() {
	log.Printf("%s lancement...\n", ag.id)

	go func() {
		env := ag.env
		var step int
		for {
			step = <-ag.syncChan
			ag.Percept(env)
			ag.Deliberate()
			ag.Act(env)

			ag.syncChan <- step
		}
	}()
}

func (ag *Agent) Percept(env *Environnement) {
	//ag.position=env.Position() TODO
}

func (ag *Agent) Deliberate() {
	//TODO
}

func (ag *Agent) Act(env *Environnement) {
	//TODO
}
