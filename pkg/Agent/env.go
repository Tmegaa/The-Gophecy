package agent

import (
	carte "Gophecy/pkg/Carte"
	ut "Gophecy/pkg/Utilitaries"
	"image"
	"math/rand"
	"sync"
)

var lsType = []TypeAgent{Sceptic, Believer, Neutral}

// Liste des types de messages possibles
type MessageType string

const (
	LoopMsg       MessageType = "StartLoop"
	PerceptionMsg MessageType = "Perception"
	NearbyMsg     MessageType = "Nearby"
	MoveMsg       MessageType = "Move"
)

// Structure d'un message
type Message struct {
	Type         MessageType
	NearbyAgents []Agent
	Agent        *Agent
}

type Environnement struct {
	sync.RWMutex
	Ags             []Agent
	Carte           carte.Carte
	Objs            []InterfaceObjet
	Communication   chan Message //key = IDAgent et value = []*Message -> Liste des messages reçus par l'agent
	NbrAgents       *sync.Map    //key = typeAgent et value = int  -> Compteur d'agents par types
	AgentProximity  *sync.Map    //key = IDAgent et value = []*Agent -> Liste des agents proches
	ObjectProximity *sync.Map    //key = IDAgent et value = []*Objet -> Liste des objets proches
}

// Fonction d'initialisation d'un nouvel environnement
func NewEnvironment(ags []Agent, carte carte.Carte, objs []InterfaceObjet) (env *Environnement) {
	// Initialisation du compteur du nombre d'agents par type
	counter := &sync.Map{}

	for _, val := range lsType {
		counter.Store(val, 0)
	}

	return &Environnement{Ags: ags, Objs: objs, Communication: make(chan Message, 100), NbrAgents: counter, Carte: carte, AgentProximity: &sync.Map{}}
}

// Fonction qui ajoute un nouvel agent dans l'environnement
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

// Fonction qui envoie la liste des agents proches pour un agent donné
func (env *Environnement) NearbyAgents(ag *Agent) []Agent {
	nearbyAgents := make([]Agent, 0)
	pos := ag.AgtPosition()

	// Création du rectangle de perception
	var area ut.Rectangle
	area.PositionDL.X = pos.X - ag.Acuite
	area.PositionDL.Y = pos.Y + ag.Acuite
	area.PositionUR.X = pos.X + ag.Acuite
	area.PositionUR.Y = pos.Y - ag.Acuite

	// On itère sur tous les agents de la carte pour voir si un agent est dans le rectangle de perception
	for _, ag2 := range env.Ags {
		if ag.ID() != ag2.ID() && ut.IsInRectangle(ag2.AgtPosition(), area) {
			nearbyAgents = append(nearbyAgents, ag2)
		}
	}
	return nearbyAgents
}

// Fonction qui envoie la liste des objets proches pour un agent donné
func (env *Environnement) NearbyObjects() {
	env.Lock()
	defer env.Unlock()
	for _, ag := range env.Ags {
		var nearbyObjects []*InterfaceObjet
		pos := ag.AgtPosition()

		// Création du rectangle de perception
		var area ut.Rectangle
		area.PositionDL.X = pos.X - ag.Acuite
		area.PositionDL.Y = pos.Y + ag.Acuite
		area.PositionUR.X = pos.X + ag.Acuite
		area.PositionUR.Y = pos.Y - ag.Acuite

		// On itère sur tous les objets de la carte pour voir si un objet est dans le rectangle de perception
		for _, obj := range env.Objs {
			if ut.IsInRectangle(obj.ObjPosition(), area) {
				nearbyObjects = append(nearbyObjects, &obj)
			}
		}
		env.ObjectProximity.Store(ag.Id, nearbyObjects)
	}
}

// Fonction de l'environnement qui gère la communication avec les agents via les channels
func (env *Environnement) Listen() {
	go func() {
		for msg := range env.Communication {
			//log.Printf("env received a message from %v", msg.Agent.ID())
			switch {
			case msg.Type == PerceptionMsg:
				near := env.NearbyAgents(msg.Agent)
				env.SendToAgent(msg.Agent, Message{Type: NearbyMsg, NearbyAgents: near})

			case msg.Type == MoveMsg:
				env.Move(msg.Agent)
			}
		}
	}()
}

func (env *Environnement) SendToAgent(agt *Agent, msg Message) {
	agt.SyncChan <- msg
}

func (env *Environnement) Move(ag *Agent) {

	// On reste dans la même direction (la vitesse par rapport aux pas écoulés depuis le début de l'action nous permet d'avoir une certaine fréquence de changement de direction )
	if ag.StepAction%ag.Velocite != 0 {
		// Si une collision a lieu on s'arrête
		if CheckCollisionHorizontal((ag.Position.X+ag.Position.Dx), (ag.Position.Y+ag.Position.Dy), ag.Env.Carte.Coliders) || CheckCollisionVertical((ag.Position.X+ag.Position.Dx), (ag.Position.Y+ag.Position.Dy), ag.Env.Carte.Coliders) {
			return
		}
		// Sinon on continue de bouger
		ag.Position.X += ag.Position.Dx
		ag.Position.Y += ag.Position.Dy

		return

	}

	// On change de direction
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
}

// Fonction qui vérifie s'il y a une collision entre un objet à la position x,y et les objets (horizontal)
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

// Fonction qui vérifie s'il y a une collision entre un objet à la position x,y et les objets (vertical)
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
