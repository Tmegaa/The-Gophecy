package agent

import (
	ut "Gophecy/pkg/Utilitaries"
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
	Velocite          int                 // vitesse à laquelle un agent peut se déplacer (pondère la fréquence de changement de direction)
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
	MoveStepLimit     int                 // limite de nombre de ticks pour une action de mouvement
	DialogStepLimit   int                 // limite de nombre de ticks pour une action de conversation
	WaitStepLimit     int                 // limite de nombre de ticks pour attendre
	CurrentAction     ActionType          // action qui est en train d'être réalisée (soit dernière décision prise)
	StepAction        int                 // nombre de boucles pendant lesquelles une action a été réalisée
	Occupied          bool                // indique si un agent est engagé dans une action bloquante (une conversation par exemple)
	AgentProximity    []Agent             // liste des agents qui sont proches
}

// Création d'un nouvel agent
func NewAgent(env *Environnement, id IdAgent, velocite int, acuite float64, position ut.Position,
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
		Vivant: true, TypeAgt: typeAgt, SubType: subTypeAgent, SyncChan: syncChan, Img: img, MoveStepLimit: 20, CurrentAction: MoveAct, DialogStepLimit: 10, WaitStepLimit: 10, StepAction: 0, Occupied: false, AgentProximity: make([]Agent, 0)}

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

	// Boucle de simulation pour notre agent
	for {
		step := <-ag.SyncChan
		// Perception
		ag.Percept(env)

		// if len(ag.AgentProximity) > 0 {
		// 	log.Printf("Nearby agents %v", ag.AgentProximity)
		// }

		// Délibération
		choice, ag2 := ag.Deliberate(env)

		// Action
		// if choice != MoveAct {
		// 	log.Printf("%s ,choice  %s ", ag.Id, choice)
		// }
		ag.Act(env, choice, ag2)

		ag.SyncChan <- step
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
func (ag *Agent) Deliberate(env *Environnement) (act ActionType, agt2 *Agent) {
	//TODO GESTION COMPUTER
	env.Lock()
	defer env.Unlock()

	// Par défaut, donc même si aucun agent est à proximité, l'agent décide de bouger
	act = MoveAct
	agt2 = nil

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
						act = MoveAct
					}
				case ag.TypeAgt == Believer:
					if ag.Opinion == 1 && ag2.Opinion == 1 {
						act = MoveAct
					}
				// Les deux agents sont neutres
				case ag.TypeAgt == Neutral:
					// Si les deux agents ont une opinion neutre, l'agent décide de bouger
					if ag.Opinion == 0.5 && ag2.Opinion == 0.5 {
						act = MoveAct
					}
					// TODO: autres cas!
				}
			}
			// Si l'autre agent est d'un autre type: l'agent décide de discuter, les deux agents seront désormais occupés
			act = DiscussAct
			agt2 = &ag2
		} else if rand.Intn(2) == 0 {
			// Gestion aléatoire entre attendre et bouger
			act = WaitAct
		}
	}
	// Gestion des limites de temps pour une action: si l'action choisie correspond à l'action en cours on vérifie les limites
	if ag.CurrentAction == act {
		switch ag.CurrentAction {
		case MoveAct:
			// Si on a atteint la limite de mouvement on décide d'attendre et on réinitialise le compteur
			if ag.MoveStepLimit <= ag.StepAction {
				act = WaitAct
				ag.StepAction = 0
			}
		case WaitAct:
			// Si on a atteint la limite d'attente on décide d'attendre et on réinitialise le compteur
			if ag.WaitStepLimit <= ag.StepAction {
				act = MoveAct
				ag.StepAction = 0
			}
		case DiscussAct:
			// Si on a atteint la limite de dialogue on décide de bouger et on réinitialise le compteur
			if ag.DialogStepLimit <= ag.StepAction {
				act = MoveAct
				ag.StepAction = 0
				// On libère les deux agents engagés dans la conversation aussi
				ag.Occupied = false
				agt2.Occupied = false
				agt2.CurrentAction = MoveAct
				agt2.StepAction = 0
				// log.Printf("Agent %s NO LONGER discussing with %s", ag.Id, agt2.Id)
				agt2 = nil
			}
		}
	} else {
		// Sinon, l'action a changé et on réinitialise le compteur
		ag.StepAction = 0
	}
	return
}

// Fonction d'action d'un agent
func (ag *Agent) Act(env *Environnement, choice ActionType, ag2 *Agent) {
	// On incrémente automatiquement le compteur de boucles de l'action
	ag.StepAction += 1
	switch choice {
	case MoveAct:
		ag.SetAction(MoveAct)
		ag.SendToEnv(Message{Type: MoveMsg, Agent: ag})
	case DiscussAct:
		ag.Discuss(ag2)
	case WaitAct:
		ag.SetAction(WaitAct)
	}
}

// Fonction qui met à jour l'action d'un agent
func (ag *Agent) SetAction(action ActionType) {
	ag.CurrentAction = action
}

// Fonction qui réinitialise l'action (en mettant l'agent en attente)
func (ag *Agent) ClearAction() {
	ag.CurrentAction = WaitAct
	ag.DialogStepLimit = 0
}

// Fonction d'action où l'agent engage un autre dans une discussion
func (ag *Agent) Discuss(ag2 *Agent) {
	// Les deux agents sont occupés
	ag.Occupied = true
	ag2.Occupied = true
	// Les deux agents discutent
	ag.SetAction(DiscussAct)
	ag2.SetAction(DiscussAct)
	log.Printf("Agent %s discussing with %s", ag.Id, ag2.Id)
}

// Fonction d'action où l'agent prie pour augmenter sa croyance: pas utilisée pour l'instant
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
