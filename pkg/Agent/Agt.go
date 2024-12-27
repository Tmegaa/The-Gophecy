package pkg

import (
	ut "Gophecy/pkg/Utilitaries"
	"log"
	"time"

	"math/rand"

	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type InterfaceAgent interface {
	Move()
	Discuter()
	Programmer()
	Perceive(*Environnement)
	Deliberate()
	Act(*Environnement, string)
}

type SubTypeAgent string

const (
	None      SubTypeAgent = "None"
	Pirate    SubTypeAgent = "Pirate"
	Converter SubTypeAgent = "Converter"
)

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
	SubType           SubTypeAgent
	SyncChan          chan Message
	Img               *ebiten.Image
	MoveTimer         int
	CurrentAction     string
	DialogTimer       int
	Occupied          bool
	AgentProximity    []Agent
}

func NewAgent(env *Environnement, id IdAgent, velocite float64, acuite float64, position ut.Position,
	opinion float64, charisme map[IdAgent]float64, relation map[IdAgent]float64, personalParameter float64, typeAgt TypeAgent, subTypeAgent SubTypeAgent, syncChan chan Message, img *ebiten.Image) *Agent {

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
		Vivant: true, TypeAgt: typeAgt, SubType: subTypeAgent, SyncChan: syncChan, Img: img, MoveTimer: 60, CurrentAction: "Praying", DialogTimer: 10, Occupied: false, AgentProximity: make([]Agent, 0)}
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
		//var step int
		for {
			//step = <-ag.SyncChan
			//time.Sleep(1 * time.Second)
			nearby := ag.Percept(env)
			//log.Printf("Agent %v is perceiving", ag.AgentProximity)
			if len(nearby) > 0 {
				log.Printf("Nearby agents %v", nearby)
			}
			choice := ag.Deliberate(env, nearby)

			if choice != "Move" {
				log.Printf("%s ,choice  %s ", ag.Id, choice)
			}
			ag.Act(env, choice)

			time.Sleep(1 * time.Second)
			//ag.SyncChan <- step
		}

	}()
}

func CheckCollisionHorizontal(x, y float64, coliders []image.Rectangle) bool {
	for _, colider := range coliders {
		if colider.Overlaps(image.Rect(int(x), int(y), int(x)+16, int(y)+16)) {
			if x > 0 {
				return true
			} else if x < 0 {
				return true
			}
		}
	}
	return false
}

func CheckCollisionVertical(x, y float64, coliders []image.Rectangle) bool {
	for _, colider := range coliders {
		if colider.Overlaps(image.Rect(int(x), int(y), int(x)+16, int(y)+16)) {
			if y > 0 {
				return true
			} else if y < 0 {
				return true
			}
		}
	}
	return false
}

/*
func (ag *Agent) Percept(env *Environnement) (nearbyAgents []*Agent) {

		env.RLock()
		defer env.RUnlock()

		log.Printf("Agent %v is perceiving", env.AgentProximity)
		value, _ := env.AgentProximity.Load(ag.Id)
		if value == nil {
			log.Printf("Agent %v has no nearby agents", value)
			nearbyAgents = make([]*Agent, 0)
			return nearbyAgents
		}

		nearby := value.([]*Agent)
		log.Printf("Agent %s has %d nearby agents", ag.Id, len(nearby))
		return nearby
	}
*/
func (ag *Agent) Percept(env *Environnement) []Agent {
	env.RLock()
	defer env.RUnlock()
	msg := Message{Type: "Perception", Agent: ag}
	ag.SendToEnv(msg)
	receive := <-ag.SyncChan
	ag.AgentProximity = receive.NearbyAgents
	if len(ag.AgentProximity) > 0 {
		//log.Printf("Agent %v is perceiving", ag.AgentProximity)
	}
	return ag.AgentProximity
}

func (ag *Agent) SetPriority(nearby []*Agent) []*Agent {
	/*
		switch ag.SubType {
		case Pirate:
			for _,


	*/
	priority := nearby
	return priority
}

func (ag *Agent) Deliberate(env *Environnement, nearbyAgents []Agent) string {
	//TODO GESTION COMPUTER
	env.Lock()
	defer env.Unlock()

	if len(nearbyAgents) > 0 {
		//log.Printf("NNNNNearby agents %v", nearbyAgents)
	}
	//aucun agent à proximité
	if len(nearbyAgents) == 0 {
		//log.Printf("Agent %v has no nearby agents", nearbyAgents)
		ag.ClearAction()
		return "Move"
	}
	//TODO FONCTION SET PRIORITY
	//priority := ag.SetPriority(nearbyAgents)
	priority := nearbyAgents //for testing
	for _, ag2 := range priority {
		if !ag2.Occupied {
			//si l'agent est du même type
			//fmt.Printf("Agent %v, Agent2 %v", ag.TypeAgt, ag2.TypeAgt)
			if ag.TypeAgt == ag2.TypeAgt {
				switch {
				case ag.TypeAgt == Sceptic:
					if ag.Opinion == 0 && ag2.Opinion == 0 {
						ag.ClearAction()
						return "Move"
					}
				case ag.TypeAgt == Believer:
					if ag.Opinion == 1 && ag2.Opinion == 1 {
						ag.ClearAction()
						return "Move"
					}
				case ag.TypeAgt == Neutral:
					if ag.Opinion == 0.5 && ag2.Opinion == 0.5 {
						ag.ClearAction()
						return "Move"
					}
				}
			}
			//si l'agent est d'un autre type
			ag.SetAction("Discuss")
			ag2.SetAction("Discuss")
			ag.Occupied = true
			ag2.Occupied = true
			return "Discuss"
		}
		//gestion aléatoire entre attendre et bouger
		if rand.Intn(2) == 0 {
			return "Wait"
		}

	}
	return "Move"
}

func (ag *Agent) Act(env *Environnement, choice string) {
	switch choice {
	case "Move":
		//log.Printf("%v", ag.Position)
		ag.SendToEnv(Message{Type: "Move", Agent: ag})

	case "Discuss":
		//ag.Discuter()
	case "Wait":
		ag.ClearAction()
	}
}

func (ag *Agent) SetAction(action string) {
	ag.CurrentAction = action
	ag.DialogTimer = 180 // 2 segundos a 60 FPS
}

func (ag *Agent) ClearAction() {
	ag.CurrentAction = "Running"
	ag.DialogTimer = 0
}

func (ag *Agent) Pray() {
	ag.SetAction("Praying")
}

func (ag *Agent) Eat() {
	ag.SetAction("Eating")
}

func (ag *Agent) SendToEnv(msg Message) {
	ag.Env.Communication <- msg
}
