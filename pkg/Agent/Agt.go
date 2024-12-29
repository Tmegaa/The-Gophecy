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
	ObjsProximity     []*InterfaceObjet
	UseComputer       *Computer
	LastComputer      *Computer
	LastStatue        *Statue
	TimeLastStatue    int
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
		Vivant: true, TypeAgt: typeAgt, SubType: subTypeAgent, SyncChan: syncChan, Img: img, MoveTimer: 60, CurrentAction: "Praying", DialogTimer: 10, Occupied: false, AgentProximity: make([]Agent, 0), ObjsProximity: make([]*InterfaceObjet, 0), UseComputer: nil, LastComputer: nil, LastStatue: nil, TimeLastStatue: 999}
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
		for {
			
			nearby, obj := ag.Percept(env)
			choice := ag.Deliberate(env, nearby, obj)
			ag.Act(env, choice)
			time.Sleep(20 * time.Millisecond)
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

func (ag *Agent) Percept(env *Environnement) ([]Agent, []*InterfaceObjet) {
	env.RLock()
	defer env.RUnlock()
	msg := Message{Type: "Perception", Agent: ag}
	ag.SendToEnv(msg)
	receive := <-ag.SyncChan
	ag.AgentProximity = receive.NearbyAgents
	if len(ag.AgentProximity) > 0 {
		//log.Printf("Agent %v is perceiving", ag.AgentProximity)
	}

	// percept objs
	ag.ObjsProximity =  env.NearbyObjects(ag)
	
	return nil, ag.ObjsProximity
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

func (ag *Agent) Deliberate(env *Environnement, nearbyAgents []Agent, obj []*InterfaceObjet) string {
    env.Lock()
    defer env.Unlock()

    // Priorizar interação com objetos
    if len(obj) > 0 {
        for _, o := range obj {
            switch concrete := (*o).(type) {
            case *Computer:
                if !concrete.Used {
                    switch concrete.Programm {
                    case "Go":
                        if ag.TypeAgt == Sceptic {
                            concrete.Programm = "None"
                            return ag.useComputer(concrete)
                        }
                    case "None":
                        if ag.TypeAgt == Believer {
                            concrete.Programm = "Go"
                            return ag.useComputer(concrete)
                        }
                    }
                }
            case *Statue:
                switch ag.TypeAgt {
                case Sceptic:
                    continue
                case Believer:
					if ag.LastStatue == nil || ag.LastStatue.ID() != concrete.ID() || ag.TimeLastStatue > 600 {
                    return ag.Prayer(concrete)
					}

                case Neutral:
                    if rand.Float64() < 0.5 && (ag.LastStatue == nil || ag.LastStatue.ID() != concrete.ID() || ag.TimeLastStatue > 350) {
                        return ag.Prayer(concrete)
                    }
                }
            }
        }
    }

    // Interação com outros agentes
    if len(nearbyAgents) > 0 {
        for _, ag2 := range nearbyAgents {
            if !ag2.Occupied {
                if ag.shouldInteract(ag2) {
                    return ag.interactWithAgent(ag2)
                }
            }
        }
    }

    // Movimento padrão ou espera
    if rand.Float64() < 0.8 {
        return "Move"
    }
    return "Move"
}

func (ag *Agent) Prayer(statue *Statue) string {
	ag.LastStatue = statue
	ag.Occupied = true
	ag.TimeLastStatue = 0
	return "Pray"
}

func (ag *Agent) useComputer(computer *Computer) string {
    computer.Used = true
    ag.Occupied = true
    ag.UseComputer = computer
    ag.LastComputer = computer
    return "Computer"
}

func (ag *Agent) shouldInteract(ag2 Agent) bool {
    if ag.TypeAgt == ag2.TypeAgt {
        return ag.Opinion != ag2.Opinion
    }
    return true
}

func (ag *Agent) interactWithAgent(ag2 Agent) string {
    ag.SetAction("Discuss")
    ag2.SetAction("Discuss")
    ag.Occupied = true
    ag2.Occupied = true
    return "Discuss"
}


func (ag *Agent) Act(env *Environnement, choice string) {
	if ag.CurrentAction != "Running" {
		return //l'agent est occupés
	}
	switch choice {
	case "Move":
		//log.Printf("%v", ag.Position)
		ag.SendToEnv(Message{Type: "Move", Agent: ag})
	
	case "Computer":
		ag.SetAction("Using Computer")

	case "Pray":
		ag.SetAction("Praying")

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


func (ag *Agent) Eat() {
	ag.SetAction("Eating")
}

func (ag *Agent) SendToEnv(msg Message) {
	ag.Env.Communication <- msg
}
