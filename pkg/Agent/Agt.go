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

func getRandomSubType(typeAgt TypeAgent) SubTypeAgent {
	if typeAgt == Neutral {
		return None
	}

	// Probabilidade de ter um subtipo (70% de chance)
	if rand.Float64() > 0.7 {
		return None
	}

	switch typeAgt {
	case Believer:
		// Para Believer: 60% Converter, 40% Pirate
		if rand.Float64() < 0.6 {
			return Converter
		}
		return Pirate
	
	case Sceptic:
		// Para Sceptic: 60% Pirate, 40% Converter
		if rand.Float64() < 0.6 {
			return Pirate
		}
		return Converter
	}

	return None
}

func NewAgent(env *Environnement, id IdAgent, velocite float64, acuite float64, position ut.Position,
	opinion float64, charisme map[IdAgent]float64, relation map[IdAgent]float64, personalParameter float64, 
	typeAgt TypeAgent, syncChan chan Message, img *ebiten.Image) *Agent {

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

	// Determina o subtipo baseado no tipo do agente
	subType := getRandomSubType(typeAgt)

	// Log para debug
	log.Printf("New agent created - ID: %s, Type: %s, SubType: %s", id, typeAgt, subType)

	return &Agent{
		Env: env,
		Id: id,
		Velocite: velocite,
		Acuite: acuite,
		Position: position,
		Opinion: opinion,
		Charisme: charisme,
		Relation: relation,
		PersonalParameter: personalParameter,
		Poids_rel: poids_rel,
		Poids_abs: poids_abs,
		Vivant: true,
		TypeAgt: typeAgt,
		SubType: subType,  // Usando o subtipo determinado
		SyncChan: syncChan,
		Img: img,
		MoveTimer: 60,
		CurrentAction: "Running",
		DialogTimer: 10,
		Occupied: false,
		AgentProximity: make([]*Agent, 0),
		ObjsProximity: make([]*InterfaceObjet, 0),
		UseComputer: nil,
		LastComputer: nil,
		LastStatue: nil,
		TimeLastStatue: 999,
		CurrentWaypoint: nil,
		LastTalkedTo: make([]*Agent, 0),
		MaxLastTalked: 3,
	}
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
				if !concrete.Used && (ag.LastComputer == nil || concrete.ID() != ag.LastComputer.ID()) {
					ag.UseComputer = concrete // Apenas guarda a referência do computador
					return "Computer"
				}
			case *Statue:
				switch ag.TypeAgt {
				case Sceptic:
					continue
				case Believer:
					if ag.LastStatue == nil || ag.LastStatue.ID() != concrete.ID() || ag.TimeLastStatue > 600 {
						ag.LastStatue = concrete // Apenas guarda a referência da estátua
						return "Pray"
					}
				case Neutral:
					if rand.Float64() < 0.5 && (ag.LastStatue == nil || ag.LastStatue.ID() != concrete.ID() || ag.TimeLastStatue > 350) {
						ag.LastStatue = concrete // Apenas guarda a referência da estátua
						return "Pray"
					}
				}
			}
		}
	}

	// Interagir com d'autres agents
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
	// Apenas verifica se é possível interagir
	if other.Occupied || other.CurrentAction == "Discussing" {
		return "Wait"
	}
	
	// Retorna "Discuss" com o agente alvo
	ag.DiscussingWith = other
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
		return //l'agent est occupé
	}
	
	switch choice {
	case "Move":
		ag.SendToEnv(Message{Type: "Move", Agent: ag})

	case "Computer":
		if ag.UseComputer != nil && ag.UseComputer.TryUse() {
			ag.SetAction("Using Computer")
			ag.Occupied = true
			ag.LastComputer = ag.UseComputer
		} else {
			ag.UseComputer = nil
		}

	case "Pray":
		if ag.LastStatue != nil {
			ag.SetAction("Praying")
			ag.Occupied = true
			ag.TimeLastStatue = 0
		}

	case "Discuss":
		if ag.DiscussingWith != nil && !ag.DiscussingWith.Occupied {
			ag.SetAction("Discussing")
			ag.DiscussingWith.SetAction("Discussing")
			
			ag.Occupied = true
			ag.DiscussingWith.Occupied = true
			
			ag.DiscussingWith.DiscussingWith = ag
			
			ag.addToTalkHistory(ag.DiscussingWith)
			ag.DiscussingWith.addToTalkHistory(ag)
		} else {
			ag.DiscussingWith = nil
			ag.ClearAction()
		}

	case "Wait":
		ag.ClearAction()
	}
}

func (ag *Agent) SetAction(action string) {
	ag.CurrentAction = action
	ag.DialogTimer = 180 // 2 segundos a 60 FPS
}

func (ag *Agent) ClearAction() {
	switch ag.CurrentAction {
	case "Discussing":
		if ag.DiscussingWith != nil {
			// Limpa também o estado do outro agente
			ag.setOpinion(ag.DiscussingWith)
			ag.DiscussingWith.CurrentAction = "Running"
			ag.DiscussingWith.Occupied = false
			ag.DiscussingWith.DiscussingWith = nil
		}
	
	case "Praying":
		log.Printf("Agent %v finished praying", ag.Id)
		log.Printf("Agent Opinion before: %v", ag.Opinion)
		if ag.TypeAgt == Believer {
			ag.Opinion = math.Min(1.0, ag.Opinion + 0.05)
		} else {
			ag.Opinion = math.Min(1.0, ag.Opinion + 0.1)
		}
		ag.CheckType()
		log.Printf("Agent Opinion after: %v", ag.Opinion)
	
	case "Using Computer":
		log.Printf("Agent %v finished using computer", ag.Id)
		log.Printf("Agent Opinion before: %v", ag.Opinion)
		
		currentProgram := ag.UseComputer.GetProgramm()
		
		switch ag.TypeAgt {
		case Believer:
			if currentProgram == "None" {
				ag.UseComputer.SetProgramm("Go")
				ag.Opinion = math.Min(1.0, ag.Opinion + 0.005)
			} else if currentProgram == "Go" {
				ag.Opinion = math.Min(1.0, ag.Opinion + 0.05)
			}
		
		case Sceptic:
			if currentProgram == "Go" {
				ag.UseComputer.SetProgramm("None")
				ag.Opinion = math.Max(0.0, ag.Opinion - 0.005)
			} else if currentProgram == "None" {
				ag.Opinion = math.Max(0.0, ag.Opinion - 0.05)
			}
		
		case Neutral:
			if currentProgram == "Go" {
				ag.Opinion = math.Min(1.0, ag.Opinion + 0.05)
			} else if currentProgram == "None" {
				ag.Opinion = math.Max(0.0, ag.Opinion - 0.05)
			}
		}
		
		ag.CheckType()
		log.Printf("Agent Opinion after: %v", ag.Opinion)
		log.Printf("Computer program: %v", ag.UseComputer.GetProgramm())
		
		if ag.UseComputer != nil {
			ag.UseComputer.Release()
			ag.UseComputer = nil
		}
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

func (ag *Agent) CheckType() {
	oldType := ag.TypeAgt
	
	// Atualiza o tipo baseado na opinião
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

	// Se o tipo mudou, recalcula o subtipo
	if oldType != ag.TypeAgt {
		// Log para debug
		log.Printf("Agent %v changed type from %v to %v", ag.Id, oldType, ag.TypeAgt)
		log.Printf("Old subtype: %v", ag.SubType)

		// Se virou Neutral, perde o subtipo
		if ag.TypeAgt == Neutral {
			ag.SubType = None
		} else {
			// Se era Neutral e virou outro tipo, ou se mudou entre Believer/Sceptic
			ag.SubType = getRandomSubType(ag.TypeAgt)
		}

		log.Printf("New subtype: %v", ag.SubType)
	}
}

func loadImageAgt(path string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatalf("Failed to load image: %s, error: %v", path, err)
	}
	return img
}
