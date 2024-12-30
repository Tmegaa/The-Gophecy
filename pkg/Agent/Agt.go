package pkg

import (
	ut "Gophecy/pkg/Utilitaries"
	"log"
	"time"

	"image"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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
	Charisme          map[IdAgent]float64 // influence d'un agent sur un autre
	Relation          map[IdAgent]float64
	PersonalParameter float64
	Poids_rel         map[IdAgent]float64
	Poids_abs         map[IdAgent]float64
	Vivant            bool
	TypeAgt           TypeAgent
	SubType           SubTypeAgent
	SyncChan          chan Message
	Img               *ebiten.Image
	MoveTimer         int
	CurrentAction     string
	DialogTimer       int
	Occupied          bool
	AgentProximity    []*Agent
	ObjsProximity     []*InterfaceObjet
	UseComputer       *Computer
	LastComputer      *Computer
	LastStatue        *Statue
	TimeLastStatue    int
	HeatMap           *VisitationMap
	CurrentWaypoint   *ut.Position // Point actuel de patrouille pour les agents Neutral
	MovementStrategy  MovementStrategy
	DiscussingWith    *Agent   // Référence à l'agent avec qui il discute
	LastTalkedTo      []*Agent // Liste des derniers agents avec qui il a conversé
	MaxLastTalked     int      // Taille maximale de la liste des derniers agents
}

func NewAgent(env *Environnement, id IdAgent, velocite float64, acuite float64, position ut.Position,
	opinion float64, charisme map[IdAgent]float64, relation map[IdAgent]float64, personalParameter float64, typeAgt TypeAgent, subTypeAgent SubTypeAgent, syncChan chan Message, img *ebiten.Image) *Agent {

	//calcul des poids relatif pour chaque agents
	poids_rel := make(map[IdAgent]float64, 0)
	poids_abs := make(map[IdAgent]float64, 0)
	min := 0.01
	max := 1.0

	if len(env.Ags) == 0 {
		poids_abs[id] = min + rand.Float64()*(max-min)
		poids_rel[id] = poids_abs[id] / (poids_abs[id] + poids_abs[id])
	} else {

		for _, v := range env.Ags {
			// pour chaque agent déja exisatant de l'environnement, on affect un poids absolu aléatoire et on calcule le poids relatif
			poids_abs[v.Id] = min + rand.Float64()*(max-min)
			poids_rel[v.Id] = poids_abs[v.Id] / (poids_abs[id] + poids_abs[v.Id])

			// on affecte les poids absolu et relatif de l'agent que l'on vient de créer à chaque agent déja existant
			v.Poids_abs[id] = min + rand.Float64()*(max-min)
			v.Poids_rel[id] = v.Poids_abs[id] / (v.Poids_abs[id] + poids_abs[id])
		}
	}
	personalCharisme := charisme[id]
	for i, v := range charisme {
		char := v / personalCharisme
		poids_rel[i] = char
	}

	return &Agent{Env: env, Id: id, Velocite: velocite, Acuite: acuite,
		Position: position, Opinion: opinion, Charisme: charisme, Relation: relation,
		PersonalParameter: personalParameter, Poids_rel: poids_rel, Poids_abs: poids_abs,
		Vivant: true, TypeAgt: typeAgt, SubType: subTypeAgent, SyncChan: syncChan, Img: img, MoveTimer: 60, CurrentAction: "Praying", DialogTimer: 10, Occupied: false, AgentProximity: make([]*Agent, 0), ObjsProximity: make([]*InterfaceObjet, 0), UseComputer: nil, LastComputer: nil, LastStatue: nil, TimeLastStatue: 999, CurrentWaypoint: nil, LastTalkedTo: make([]*Agent, 0), MaxLastTalked: 3}
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

func (ag *Agent) Percept(env *Environnement) ([]*Agent, []*InterfaceObjet) {
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
	ag.ObjsProximity = env.NearbyObjects(ag)

	return ag.AgentProximity, ag.ObjsProximity
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

func (ag *Agent) Deliberate(env *Environnement, nearbyAgents []*Agent, obj []*InterfaceObjet) string {
	env.Lock()
	defer env.Unlock()

	// Prioriser l'interaction avec les objets
	if len(obj) > 0 {
		for _, o := range obj {
			switch concrete := (*o).(type) {
			case *Computer:
				if !concrete.Used {
					switch concrete.Programm {
					case "Go":
						if ag.TypeAgt == Sceptic {
							concrete.SetProgramm("None")
							return ag.useComputer(concrete)
						}
					case "None":
						if ag.TypeAgt == Believer {
							concrete.SetProgramm("Go")
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

	// Interagir avec d'autres agents
	if len(nearbyAgents) > 0 {
		for i := range nearbyAgents {
			// Obtém o ponteiro para o agente original do ambiente
			otherAgent := env.GetAgentById(nearbyAgents[i].Id)
			if otherAgent != nil && !otherAgent.Occupied {
				if ag.shouldInteract(otherAgent) {
					return ag.interactWithAgent(otherAgent)
				}
			}
		}
	}

	// Mouvement par défaut ou attente
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
	if !computer.TryUse() {
		return "Wait"
	}

	ag.Occupied = true
	ag.UseComputer = computer
	ag.LastComputer = computer
	return "Computer"
}

func (ag *Agent) shouldInteract(other *Agent) bool {
	// Se algum dos agentes está ocupado, não deve interagir
	if ag.Occupied || other.Occupied {
		return false
	}

	// Verifica se o outro agente já está em discussão com alguém
	if other.CurrentAction == "Discussing" {
		return false
	}

	// Verifica se já conversou recentemente com este agente
	for _, lastTalked := range ag.LastTalkedTo {
		if lastTalked.Id == other.Id {
			return false
		}
	}

	// Verifica se o tipo de agente influencia a interação
	if ag.TypeAgt == other.TypeAgt {
		return ag.Opinion != other.Opinion
	}
	return true
}
func (ag *Agent) setOpinion(ag2 *Agent) {
	if ag.TypeAgt == "Sceptic" && ag2.TypeAgt == "Believer" {
		//faire proba avec charisme
		ag.Opinion = ag.Opinion - 0.05
		ag2.Opinion = ag2.Opinion + 0.05
	} else if ag.TypeAgt == "Believer" && ag2.TypeAgt == "Sceptic" {
		ag.Opinion = ag.Opinion + 0.05
		ag2.Opinion = ag2.Opinion - 0.05
	} else {
		ag.Opinion = ag.Poids_rel[ag.Id]*ag.PersonalParameter*ag.Opinion*(1-ag.Opinion) + ag.Poids_rel[ag2.Id]*ag2.Opinion    //calcul du nouvel opinion
		ag2.Opinion = ag.Poids_rel[ag.Id]*ag.Opinion + ag.Poids_rel[ag2.Id]*ag2.PersonalParameter*ag2.Opinion*(1-ag2.Opinion) //calcul du nouvel opinion
	}
	ag.Opinion = math.Max(0, math.Min(1, ag.Opinion))
	ag2.Opinion = math.Max(0, math.Min(1, ag2.Opinion))
	ag.CheckType()
	ag2.CheckType()

}
func (ag *Agent) interactWithAgent(other *Agent) string {
	// Dupla verificação de segurança
	if other.Occupied || other.CurrentAction == "Discussing" {
		return "Wait"
	}

	// Configura ambos os agentes para discussão
	ag.SetAction("Discussing")
	other.SetAction("Discussing")
	ag.Occupied = true
	other.Occupied = true

	// Guarda referência do outro agente para visualização
	ag.DiscussingWith = other
	other.DiscussingWith = ag

	// Adiciona ambos os agentes ao histórico um do outro
	ag.addToTalkHistory(other)
	other.addToTalkHistory(ag)

	return "Discuss"
}

func (ag *Agent) addToTalkHistory(other *Agent) {
	// Verifica se o agente já está no histórico
	for _, a := range ag.LastTalkedTo {
		if a.Id == other.Id {
			return // Se já está no histórico, não adiciona novamente
		}
	}

	// Adiciona o novo agente ao início do histórico
	ag.LastTalkedTo = append([]*Agent{other}, ag.LastTalkedTo...)

	// Mantém apenas os últimos MaxLastTalked agentes
	if len(ag.LastTalkedTo) > ag.MaxLastTalked {
		ag.LastTalkedTo = ag.LastTalkedTo[:ag.MaxLastTalked]
	}
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
	if ag.CurrentAction == "Discussing" && ag.DiscussingWith != nil {
		// Limpa também o estado do outro agente
		ag.setOpinion(ag.DiscussingWith)
		ag.DiscussingWith.CurrentAction = "Running"
		ag.DiscussingWith.Occupied = false
		ag.DiscussingWith.DiscussingWith = nil

	}
	ag.CurrentAction = "Running"
	ag.DialogTimer = 0
	ag.Occupied = false
	ag.DiscussingWith = nil
}

func (ag *Agent) Eat() {
	ag.SetAction("Eating")
}

func (ag *Agent) SendToEnv(msg Message) {
	ag.Env.Communication <- msg
}

func (env *Environnement) GetAgentById(id IdAgent) *Agent {
	for _, agent := range env.Ags {
		if agent.Id == id {
			return agent
		}
	}
	return nil
}

func (ag *Agent) CheckType() { //fonction qui permet de définir le type de l'agent en fonction de son opinion; à verifier tous les X top d'horloges
	{
		if ag.Opinion > 0.66 {
			ag.TypeAgt = Believer
			ag.Img = loadImageAgt(ut.AssetsPath + ut.AgentBelieverImageFile)
		} else if ag.Opinion < 0.33 {
			ag.TypeAgt = Sceptic
			ag.Img = loadImageAgt(ut.AssetsPath + ut.AgentScepticImageFile)
		} else {
			ag.TypeAgt = Neutral
			ag.Img = loadImageAgt(ut.AssetsPath + ut.AgentNeutralImageFile)
		}
	}
}

func loadImageAgt(path string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatalf("Failed to load image: %s, error: %v", path, err)
	}
	return img
}
