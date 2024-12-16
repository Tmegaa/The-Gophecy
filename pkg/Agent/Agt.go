package pkg

import (
	ut "Gophecy/pkg/Utilitaries"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
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
	Env               *Environnement
	Id                IdAgent
	Velocite          float64
	Acuite            float64
	Position          ut.Position
	Opinion           float64
	Charisme          map[IdAgent]float64 //influence d'un agent sur un autre
	Relation          map[IdAgent]float64
	PersonalParameter float64
	Poid_rel          []float64
	Vivant            bool
	TypeAgt           TypeAgent
	SyncChan          chan int
	Img               *ebiten.Image
}

func NewAgent(env *Environnement, id IdAgent, velocite float64, acuite float64, position ut.Position,
	opinion float64, charisme map[IdAgent]float64, relation map[IdAgent]float64, personalParameter float64,
	agent InterfaceAgent, typeAgt TypeAgent, syncChan chan int, img *ebiten.Image) *Agent {

	//calcul des poids relatif pour chaque agents
	poid_rel := make([]float64, 0)
	personalCharisme := charisme[id]
	for _, v := range charisme {
		char := v / personalCharisme
		poid_rel = append(poid_rel, char)
	}

	return &Agent{Env: env, Id: id, Velocite: velocite, Acuite: acuite,
		Position: position, Opinion: opinion, Charisme: charisme, Relation: relation,
		PersonalParameter: personalParameter, Poid_rel: poid_rel,
		Vivant: true, TypeAgt: typeAgt, SyncChan: syncChan, Img: img}
}

func (ag *Agent) ID() IdAgent {
	return ag.Id
}

func (ag *Agent) AgtPosition() ut.Position {
	return ag.Position
}

func (ag *Agent) Start() {
	log.Printf("%s lancement...\n", ag.Id)

	go func() {
		env := ag.Env
		var step int
		for {
			step = <-ag.SyncChan
			ag.Percept(env)
			ag.Deliberate()
			ag.Act(env)

			ag.SyncChan <- step
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
