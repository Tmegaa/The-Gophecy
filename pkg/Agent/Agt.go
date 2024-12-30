package agent

import (
	ut "Gophecy/pkg/Utilitaries"
	"image"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

// Interface qui regroupe les méthodes de tous les types d'agents
type InterfaceAgent interface {
	Move()
	Discuter()
	Programmer()
	Perceive(*Environnement)
	Deliberate()
	Act(*Environnement, string)
}

// Différents types d'agents neutres: ce sous-type affecte la prise de décision
type SubTypeAgent string

const (
	None      SubTypeAgent = "None"
	Pirate    SubTypeAgent = "Pirate"
	Converter SubTypeAgent = "Converter"
)

// Il existe trois types d'agents dans la simulation, on définie ici le type TypeAgent et les valeurs possibles
type TypeAgent string

const (
	Sceptic  TypeAgent = "Sceptic"
	Believer TypeAgent = "Believer"
	Neutral  TypeAgent = "Neutral"
)

// Chaque agent va avoir un ID
type IdAgent string

// On définie toutes les actions possibles
type ActionType string

const (
	MoveAct    ActionType = "Move"
	DiscussAct ActionType = "Discuss"
	WaitAct    ActionType = "Wait"
	PrayAct    ActionType = "Pray"
	EatAct     ActionType = "Eat"
)

type Agent struct {
	Env               *Environnement      // pointeur vers l'environnement
	Id                IdAgent             // identifiant agent
	Velocite          float64             // vitesse à laquelle un agent peut se déplacer
	Acuite            float64             // définie la perception de cet agent par rapport à l'environnement
	Position          ut.Position         // position d'un agent dans la carte
	Opinion           float64             // définie son degré de croyance ou de scepticisme
	Charisme          map[IdAgent]float64 // influence exercée par les autres agents sur celui-ci
	Relation          map[IdAgent]float64 // relation que cet agent a aux autres agents
	PersonalParameter float64             // paramètre personnel influençant l’attraction ou la répulsion des opinions
	Poid_rel          []float64           // paramètre qui donne le poids relatif des opinions des autres en fonction du charisme, du paramètre personnel et des relations entre agents
	Vivant            bool                // booléen indiquant si l'agent est vivant: inutile pour l'instant
	TypeAgt           TypeAgent           // agent sceptique, neutre ou croyant
	SubType           SubTypeAgent        // agent de sous-type pirate, évangéliste ou sans sous-type
	SyncChan          chan Message        // channel propre à l'agent pour communiquer avec l'environnement
	Img               *ebiten.Image       // image à afficher sur l'interface graphique
	MoveTimer         int                 // limite de nombre de ticks pour une action de mouvement
	CurrentAction     ActionType          // action qui est en train d'être réalisée (soit dernière décision prise)
	DialogTimer       int                 // limite de nombre de ticks pour une action de conversation
	Occupied          bool                // indique si un agent est engagé dans une action bloquante (une conversation par exemple)
	AgentProximity    []Agent             // liste des agents qui sont proches
}

// Création d'un nouvel agent
func NewAgent(env *Environnement, id IdAgent, velocite float64, acuite float64, position ut.Position,
	opinion float64, charisme map[IdAgent]float64, relation map[IdAgent]float64, personalParameter float64, typeAgt TypeAgent, subTypeAgent SubTypeAgent, syncChan chan Message, img *ebiten.Image) *Agent {

	//calcul des poids relatifs du nouvel agent par rapport à chaque autre agent
	poid_rel := make([]float64, 0)
	personalCharisme := charisme[id]
	for _, v := range charisme {
		char := v / personalCharisme
		poid_rel = append(poid_rel, char)
	}

	return &Agent{Env: env, Id: id, Velocite: velocite, Acuite: acuite,
		Position: position, Opinion: opinion, Charisme: charisme, Relation: relation,
		PersonalParameter: personalParameter, Poid_rel: poid_rel,
		Vivant: true, TypeAgt: typeAgt, SubType: subTypeAgent, SyncChan: syncChan, Img: img, MoveTimer: 60, CurrentAction: MoveAct, DialogTimer: 10, Occupied: false, AgentProximity: make([]Agent, 0)}
}

// Fonction qui renvoie d'ID d'un agent
func (ag *Agent) ID() IdAgent {
	return ag.Id
}

// Fonction qui renvoie la position d'un agent
func (ag *Agent) AgtPosition() ut.Position {
	return ag.Position
}

// Fonction qui lance la boucle de perception, délibération et action pour chaque agent
func (ag *Agent) Start() {
	log.Printf("%s lancement...\n", ag.Id)
	env := ag.Env
	//var step int
	// Boucle de simulation pour notre agent
	for {
		// Perception
		ag.Percept(env)
		//step = <-ag.SyncChan
		//time.Sleep(1 * time.Second)
		if len(ag.AgentProximity) > 0 {
			log.Printf("Nearby agents %v", ag.AgentProximity)
		}

		// Délibération
		choice := ag.Deliberate(env)

		// Action
		if choice != MoveAct {
			log.Printf("%s ,choice  %s ", ag.Id, choice)
		}
		ag.Act(env, choice)

		// time.Sleep(15 * time.Millisecond)
		//ag.SyncChan <- step
	}

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

// Fonction de perception d'un agent
func (ag *Agent) Percept(env *Environnement) {
	// On utilise un mutex en Read pour avoir la certitude qu'il n'y aura pas d'accès concurrent aux données
	env.RLock()
	defer env.RUnlock()

	// Message à envoyer à l'environnement
	msg := Message{Type: PerceptionMsg, Agent: ag}
	ag.SendToEnv(msg)

	// Réception des agents à proximité dans le channel de l'agent
	receive := <-ag.SyncChan
	// Message d'erreur en cas de mauvais type de message
	if receive.Type != NearbyMsg {
		log.Println("Error: received message of wrong type. Expected type:", NearbyMsg, ". Type received:", receive.Type)
	}

	// Mise à jour de la liste des agents à proximité de l'agent
	ag.AgentProximity = receive.NearbyAgents
	// if len(ag.AgentProximity) > 0 {
	// 	log.Printf("Agent %v is perceiving", ag.AgentProximity)
	// }
}

// Fonction qui détermine la priorité d'un agent par rapport à son sous-type
func (ag *Agent) SetPriority(nearby []*Agent) []*Agent {
	/*
		switch ag.SubType {
		case Pirate:
			for _,

	*/
	priority := nearby
	return priority
}

// Fonction de délibération d'un agent
func (ag *Agent) Deliberate(env *Environnement) ActionType {
	//TODO GESTION COMPUTER
	env.Lock()
	defer env.Unlock()

	// if len(ag.AgentProximity) > 0 {
	//log.Printf("NNNNNearby agents %v", nearbyAgents)
	// }

	// Si aucun agent est à proximité, l'agent décide de bouger
	if len(ag.AgentProximity) == 0 {
		//log.Printf("Agent %v has no nearby agents", nearbyAgents)
		ag.ClearAction()
		return MoveAct
	}

	//TODO FONCTION SET PRIORITY
	//priority := ag.SetPriority(nearbyAgents)
	priority := ag.AgentProximity //for testing
	// On itère sur les agents à proximité
	for _, ag2 := range priority {
		// Si l'autre agent n'est pas occupé
		if !ag2.Occupied {
			// Si l'agent est du même type
			if ag.TypeAgt == ag2.TypeAgt {
				switch {
				// Les deux agents sont sceptiques pu croyants: l'agent décide de bouger
				case ag.TypeAgt == Sceptic:
					if ag.Opinion == 0 && ag2.Opinion == 0 {
						ag.ClearAction()
						return MoveAct
					}
				case ag.TypeAgt == Believer:
					if ag.Opinion == 1 && ag2.Opinion == 1 {
						ag.ClearAction()
						return MoveAct
					}
				// Les deux agents sont neutres
				case ag.TypeAgt == Neutral:
					// Si les deux agents ont une opinion neutre, l'agent décide de bouger
					if ag.Opinion == 0.5 && ag2.Opinion == 0.5 {
						ag.ClearAction()
						return MoveAct
					}
					// TODO: autres cas!
				}
			}
			// Si l'autre agent est d'un autre type: l'agent décide de discuter, les deux agents sont désormais occupés
			ag.SetAction(DiscussAct)
			ag2.SetAction(DiscussAct)
			ag.Occupied = true
			ag2.Occupied = true
			return DiscussAct
		}
		// Gestion aléatoire entre attendre et bouger
		if rand.Intn(2) == 0 {
			return WaitAct
		}

	}
	return MoveAct
}

// Fonction d'action d'un agent
func (ag *Agent) Act(env *Environnement, choice ActionType) {
	switch choice {
	case MoveAct:
		//log.Printf("%v", ag.Position)
		ag.SendToEnv(Message{Type: MoveMsg, Agent: ag})

	case DiscussAct:
		//ag.Discuter()
	case WaitAct:
		ag.ClearAction()
	}
}

// Fonction qui met à jour l'action d'un agent
func (ag *Agent) SetAction(action ActionType) {
	ag.CurrentAction = action
	// TODO: gérer les timers
	ag.DialogTimer = 180 // 2 secondes à 60 FPS
}

// Fonction qui réinitialise l'action (en mettant l'agent en attente)
func (ag *Agent) ClearAction() {
	ag.CurrentAction = WaitAct
	ag.DialogTimer = 0
}

// Fonction d'action où l'agent prie pour augmenter sa croyance
func (ag *Agent) Pray() {
	ag.SetAction(PrayAct)
}

// Fonction d'action où l'agent mange: pas utilisée pour l'instant
func (ag *Agent) Eat() {
	ag.SetAction(EatAct)
}

// Fonction qui envoie un message à l'environnement via le channel de communication l'environnement
func (ag *Agent) SendToEnv(msg Message) {
	ag.Env.Communication <- msg
}

// Fontion qui vérifie s'il y a une collision entre un objet à la position x,y et les objets (horizontal)
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

// Fontion qui vérifie s'il y a une collision entre un objet à la position x,y et les objets (vertical)
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
