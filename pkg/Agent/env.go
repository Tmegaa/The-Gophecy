package pkg

import (
	carte "Gophecy/pkg/Carte"
	ut "Gophecy/pkg/Utilitaries"
	"math/rand"
	"sync"
)

var lsType = []TypeAgent{Sceptic, Believer, Neutral}

type Message struct {
	Type         string
	NearbyAgents []Agent
	Agent        *Agent
}

type Environnement struct {
	sync.RWMutex
	Ags             []*Agent
	Carte           carte.Carte
	Objs            []InterfaceObjet
	Communication   chan Message //key = IDAgent et value = []*Message -> Liste des messages reçus par l'agent
	NbrAgents       *sync.Map    //key = typeAgent et value = int  -> Compteur d'agents par types
	AgentProximity  *sync.Map    //key = IDAgent et value = []*Agent -> Liste des agents proches
	ObjectProximity *sync.Map    //key = IDAgent et value = []*Objet -> Liste des objets proches
}

func NewEnvironment(ags []*Agent, carte carte.Carte, objs []InterfaceObjet) (env *Environnement) {
	counter := &sync.Map{}

	for _, val := range lsType {
		counter.Store(val, 0)
	}

	return &Environnement{Ags: ags, Objs: objs, Communication: make(chan Message, 100), NbrAgents: counter, Carte: carte, AgentProximity: &sync.Map{}}
}

func (env *Environnement) AddAgent(ag *Agent) {
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

func (env *Environnement) NearbyAgents(ag *Agent) []Agent {
	nearbyAgents := make([]Agent, 0)
	pos := ag.AgtPosition()
	var area ut.Rectangle
	area.PositionDL.X = pos.X - ag.Acuite
	area.PositionDL.Y = pos.Y + ag.Acuite
	area.PositionUR.X = pos.X + ag.Acuite
	area.PositionUR.Y = pos.Y - ag.Acuite

	for _, ag2 := range env.Ags {
		if ag.ID() != ag2.ID() && ut.IsInRectangle(ag2.AgtPosition(), area) {
			nearbyAgents = append(nearbyAgents, *ag2)
		}
	}
	if len(nearbyAgents) > 0 {
	}
	return nearbyAgents

}

func (env *Environnement) NearbyObjects(ag *Agent) []*Computer {
	nearbyObjects := make([]*Computer, 0)
	pos := ag.AgtPosition()
	var area ut.Rectangle
	area.PositionDL.X = pos.X - ag.Acuite
	area.PositionDL.Y = pos.Y + ag.Acuite
	area.PositionUR.X = pos.X + ag.Acuite
	area.PositionUR.Y = pos.Y - ag.Acuite

	
	for _, pc := range env.Objs {
		PCposition := pc.ObjPosition()

		if ut.IsInRectangle(PCposition, area) {

			if pc.GetUse() && (ag.LastComputer == nil || pc.ID() != ag.LastComputer.ID()) {
				continue
			}

			nearbyObjects = append(nearbyObjects, pc.(*Computer))	
		}
	}

	
	
	return nearbyObjects
}

func (env *Environnement) Listen() {
	go func() {
		for msg := range env.Communication {
			switch {
			case msg.Type == "Perception":
				near := env.NearbyAgents(msg.Agent)
				env.SendToAgent(msg.Agent, Message{Type: "Nearby", NearbyAgents: near})

			case msg.Type == "Move":
				env.Move(msg.Agent)

			}
		}
	}()
}

func (env *Environnement) SendToAgent(agt *Agent, msg Message) {
	agt.SyncChan <- msg
}

func (env *Environnement) Move(ag *Agent) {

	ag.ClearAction()

	if ag.MoveTimer > 0 {

		ag.MoveTimer -= 1
		if CheckCollisionHorizontal((ag.Position.X+ag.Position.Dx), (ag.Position.Y+ag.Position.Dy), ag.Env.Carte.Coliders) || CheckCollisionVertical((ag.Position.X+ag.Position.Dx), (ag.Position.Y+ag.Position.Dy), ag.Env.Carte.Coliders) {
			return
		}
		ag.Position.X += ag.Position.Dx
		ag.Position.Y += ag.Position.Dy
		return
	}

	randIdx := 0
	collision := true
	right := ut.UniqueDirection{Dx: ut.Maxspeed, Dy: 0}
	left := ut.UniqueDirection{Dx: -ut.Maxspeed, Dy: 0}
	down := ut.UniqueDirection{Dx: 0, Dy: ut.Maxspeed}
	up := ut.UniqueDirection{Dx: 0, Dy: -ut.Maxspeed}

	directions := []ut.UniqueDirection{right, left, down, up}

	for collision {
		randIdx = rand.Intn(len(directions))
		tryX := ag.Position.X + directions[randIdx].Dx
		tryY := ag.Position.Y + directions[randIdx].Dy
		if !CheckCollisionHorizontal(tryX, tryY, ag.Env.Carte.Coliders) && !CheckCollisionVertical(tryX, tryY, ag.Env.Carte.Coliders) {
			collision = false
		}
	}
	ag.Position.Dx = directions[randIdx].Dx
	ag.Position.Dy = directions[randIdx].Dy

	ag.MoveTimer = 60

}
