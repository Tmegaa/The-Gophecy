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

// Interface qui regroupe les méthodes de tous les types d'agents
type InterfaceAgent interface {
	Move()
	Discuter()
	Programmer()
	Perceive(*Environnement)
	Deliberate()
	Act(*Environnement, string)
}

// Différents types d'agents : ce sous-type affecte la prise de décision
type SubTypeAgent string

const (
	None      SubTypeAgent = "None"
	Pirate    SubTypeAgent = "Pirate"
	Converter SubTypeAgent = "Converter"
)

// Il existe trois types d'agents dans la simulation, on définie ici le type TypeAgent et les valeurs possibles
type TypeAgent string

const (
	Sceptic  TypeAgent = "Sceptique"
	Believer TypeAgent = "Croyant"
	Neutral  TypeAgent = "Neutre"
)

// Chaque agent va avoir un ID
type IdAgent string

// On définie toutes les actions possibles
type ActionType string

const (
	MoveAct     ActionType = "Bouge"
	ComputerAct ActionType = "Utilise un ordinateur"
	RunAct      ActionType = "Cours"
	DiscussAct  ActionType = "Discute"
	WaitAct     ActionType = "Attend"
	PrayAct     ActionType = "Prie"
	EatAct      ActionType = "Mange"
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
	Poids_rel         map[IdAgent]ut.Pair // paramètre qui donne le poids relatif des opinions des autres en fonction du charisme, du paramètre personnel et des relations entre agents
	Poids_abs         map[IdAgent]float64 // paramètre de poids absolu
	Vivant            bool                // booléen indiquant si l'agent est vivant: inutile pour l'instant
	TypeAgt           TypeAgent           // agent sceptique, neutre ou croyant
	SubType           SubTypeAgent        // agent de sous-type pirate, évangéliste ou sans sous-type
	SyncChan          chan Message        // channel propre à l'agent pour communiquer avec l'environnement
	Img               *ebiten.Image       // image à afficher sur l'interface graphique
	MoveTimer         int                 // Temps de mouvement d'un agent
	CurrentAction     ActionType          // action qui est en train d'être réalisée (soit dernière décision prise)
	DialogTimer       int                 // Temps de dialogue d'un agent
	Occupied          bool                // indique si un agent est engagé dans une action bloquante (une conversation par exemple)
	AgentProximity    []*Agent            // liste des agents qui sont proches
	ObjsProximity     []*InterfaceObjet   // liste des objets qui sont proches
	UseComputer       *Computer           // Ordinateur en cours d'utilisation
	LastComputer      *Computer           // Dernier ordinateur utilisé
	LastStatue        *Statue             // Dernière statue utilisée
	TimeLastStatue    int                 // Temps écoulé depuis la dernière utilisation d'une statue
	HeatMap           *VisitationMap      // Carte des endroits visités
	CurrentWaypoint   *ut.Position        // Point actuel de patrouille pour les agents Neutral
	MovementStrategy  MovementStrategy    // Stratégie de mouvement de l'agent
	DiscussingWith    *Agent              // Référence à l'agent avec qui il discute
	LastTalkedTo      []*Agent            // Liste des derniers agents avec qui il a conversé
	MaxLastTalked     int                 // Taille maximale de la liste des derniers agents
}

// Fonction qui renvoie un sous-type par rapport au type de l'agent
func getRandomSubType(typeAgt TypeAgent) SubTypeAgent {
	// Un agent neutre n'a pas de sous-type
	if typeAgt == Neutral {
		return None
	}

	// Probabilité d'avoir un sous-type (70 % de chance)
	if rand.Float64() > 0.7 {
		return None
	}

	switch typeAgt {
	case Believer:
		// Pour les croyants : 60 % convertisseur, 40 % pirate
		if rand.Float64() < 0.6 {
			return Converter
		}
		return Pirate

	case Sceptic:
		// Pour les sceptiques : 60% pirate, 40% convertisseur
		if rand.Float64() < 0.6 {
			return Pirate
		}
		return Converter
	}

	return None
}

// Création d'un nouvel agent
func NewAgent(env *Environnement, id IdAgent, velocite float64, acuite float64, position ut.Position,
	opinion float64, charisme map[IdAgent]float64, relation map[IdAgent]float64, personalParameter float64,
	typeAgt TypeAgent, syncChan chan Message, img *ebiten.Image) *Agent {

	// Calcul des poids relatifs du nouvel agent par rapport à chaque autre agent
	poids_rel := make(map[IdAgent]ut.Pair, 0)
	poids_abs := make(map[IdAgent]float64, 0)

	// Détermine le sous-type en fonction du type d'agent
	subType := getRandomSubType(typeAgt)

	// Log pour debug
	log.Printf("New agent created - ID: %s, Type: %s, SubType: %s", id, typeAgt, subType)

	return &Agent{
		Env:               env,
		Id:                id,
		Velocite:          velocite,
		Acuite:            acuite,
		Position:          position,
		Opinion:           opinion,
		Charisme:          charisme,
		Relation:          relation,
		PersonalParameter: personalParameter,
		Poids_rel:         poids_rel,
		Poids_abs:         poids_abs,
		Vivant:            true,
		TypeAgt:           typeAgt,
		SubType:           subType, // En utilisant le sous-type donné
		SyncChan:          syncChan,
		Img:               img,
		MoveTimer:         60,
		CurrentAction:     RunAct,
		DialogTimer:       10,
		Occupied:          false,
		AgentProximity:    make([]*Agent, 0),
		ObjsProximity:     make([]*InterfaceObjet, 0),
		UseComputer:       nil,
		LastComputer:      nil,
		LastStatue:        nil,
		TimeLastStatue:    999,
		CurrentWaypoint:   nil,
		LastTalkedTo:      make([]*Agent, 0),
		MaxLastTalked:     3,
	}
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

	go func() {
		env := ag.Env
		for {
			// Perception
			nearby, obj := ag.Percept(env)
			// Délibération
			choice := ag.Deliberate(env, nearby, obj)
			// Action
			ag.Act(env, choice)
			// Temps d'attente entre actions
			time.Sleep(20 * time.Millisecond)
		}
	}()
}

// Fonction de perception d'un agent
func (ag *Agent) Percept(env *Environnement) ([]*Agent, []*InterfaceObjet) {
	// On utilise un mutex en Read pour avoir la certitude qu'il n'y aura pas d'accès concurrent aux données
	env.RLock()
	defer env.RUnlock()

	// Message à envoyer à l'environnement
	msg := Message{Type: "Perception", Agent: ag}
	ag.SendToEnv(msg)

	// Réception des agents à proximité dans le channel de l'agent
	receive := <-ag.SyncChan

	// Message d'erreur en cas de mauvais type de message
	if receive.Type != NearbyMsg {
		log.Println("Error: received message of wrong type. Expected type:", NearbyMsg, ". Type received:", receive.Type)
	}

	// Mise à jour de la liste des agents à proximité de l'agent
	ag.AgentProximity = receive.NearbyAgents

	// Mise à jour de la liste des objets à proximité de l'agent
	ag.ObjsProximity = env.NearbyObjects(ag)

	return ag.AgentProximity, ag.ObjsProximity
}

// Fonction de délibération d'un agent
func (ag *Agent) Deliberate(env *Environnement, nearbyAgents []*Agent, obj []*InterfaceObjet) ActionType {
	env.Lock()
	defer env.Unlock()

	// Vérifie s'il y a des objets ou des agents à proximité
	hasObjects := len(obj) > 0
	hasAgents := len(nearbyAgents) > 0

	if !hasObjects && !hasAgents {
		return MoveAct
	}

	// Define prioridade baseada no subtipo
	switch ag.SubType {
	case Pirate:
		// Les pirates donnent la priorité aux ordinateurs
		if hasObjects {
			for _, o := range obj {
				if computer, ok := (*o).(*Computer); ok {
					if !computer.Used && (ag.LastComputer == nil || computer.ID() != ag.LastComputer.ID()) {
						ag.UseComputer = computer
						return ComputerAct
					}
				}
			}
		}
		// Si on ne trouve pas d'ordinateur, on essaye d'interagir avec des agents
		if hasAgents {
			return ag.tryInteractWithAgents(env, nearbyAgents)
		}

	case Converter:
		// Les convertisseurs donnent la priorité à l’interaction avec d’autres agents
		if hasAgents {
			result := ag.tryInteractWithAgents(env, nearbyAgents)
			if result != MoveAct {
				return result
			}
		}
		// Si on ne peut pas interagir, on essaye d'utiliser les objets
		if hasObjects {
			return ag.tryUseObjects(obj)
		}

	default: // Aucun ou autres sous-types
		// Comportement par défaut : choisit aléatoirement entre les objets et les agents
		if rand.Float64() < 0.5 && hasObjects {
			return ag.tryUseObjects(obj)
		} else if hasAgents {
			return ag.tryInteractWithAgents(env, nearbyAgents)
		}
	}

	return MoveAct
}

// Fonction auxiliaire pour tenter d'interagir avec les agents à proximité
func (ag *Agent) tryInteractWithAgents(env *Environnement, nearbyAgents []*Agent) ActionType {
	for i := range nearbyAgents {
		otherAgent := env.GetAgentById(nearbyAgents[i].Id)
		if otherAgent != nil && !otherAgent.Occupied {
			if ag.shouldInteract(otherAgent) {
				return ag.interactWithAgent(otherAgent)
			}
		}
	}
	return MoveAct
}

// Fonction d'assistance pour essayer d'utiliser des objets à proximité
func (ag *Agent) tryUseObjects(obj []*InterfaceObjet) ActionType {
	for _, o := range obj {
		switch concrete := (*o).(type) {
		case *Computer:
			if !concrete.Used && (ag.LastComputer == nil || concrete.ID() != ag.LastComputer.ID()) {
				ag.UseComputer = concrete
				return ComputerAct
			}
		case *Statue:
			switch ag.TypeAgt {
			case Sceptic:
				continue
			case Believer:
				if ag.LastStatue == nil || ag.LastStatue.ID() != concrete.ID() || ag.TimeLastStatue > 600 {
					ag.LastStatue = concrete
					return PrayAct
				}
			case Neutral:
				if rand.Float64() < 0.5 && (ag.LastStatue == nil || ag.LastStatue.ID() != concrete.ID() || ag.TimeLastStatue > 350) {
					ag.LastStatue = concrete
					return PrayAct
				}
			}
		}
	}
	return MoveAct
}

// Fonction d'action où l'agent prie pour augmenter sa croyance
func (ag *Agent) Prayer(statue *Statue) ActionType {
	ag.LastStatue = statue
	ag.Occupied = true
	ag.TimeLastStatue = 0
	return PrayAct
}

// Fonction d'action où l'agent tente d'utiliser un ordinateur
func (ag *Agent) useComputer(computer *Computer) ActionType {
	if !computer.TryUse() {
		return WaitAct
	}

	ag.Occupied = true
	ag.UseComputer = computer
	ag.LastComputer = computer
	return ComputerAct
}

// Fonction auxiliaire pour vérifier si deux agents peuvent interagir
func (ag *Agent) shouldInteract(other *Agent) bool {
	// Si l'un des agents est occupé, il ne doit pas interagir
	if ag.Occupied || other.Occupied {
		return false
	}

	// Vérifie si l'autre agent est déjà en discussion avec quelqu'un
	if other.CurrentAction == DiscussAct {
		return false
	}

	// Vérifiez si vous avez récemment parlé à cet agent
	for _, lastTalked := range ag.LastTalkedTo {
		if lastTalked.Id == other.Id {
			return false
		}
	}

	// Vérifie si le type d'agent influence l'interaction
	if ag.TypeAgt == other.TypeAgt {
		return ag.Opinion != other.Opinion
	}
	return true
}

// Fonction qui fait évoluer les opinions de deux agents en discussion
func (ag *Agent) setOpinion(ag2 *Agent) {
	// Un agent scéptique qui parle à un croyant va augmenter sa croyance et à l'inverse
	if ag.TypeAgt == Sceptic && ag2.TypeAgt == Believer {
		ag.Opinion = ag.Opinion - 0.05
		ag2.Opinion = ag2.Opinion + 0.05
	} else if ag.TypeAgt == Believer && ag2.TypeAgt == Sceptic {
		ag.Opinion = ag.Opinion + 0.05
		ag2.Opinion = ag2.Opinion - 0.05
	} else {
		// Pour deux agents neutres, on utilise les équations faisant rentrer en compte les paramètres
		newOpinionAg := ag.Poids_rel[ag2.Id].First*ag.PersonalParameter*ag.Opinion*(1.0-ag.Opinion) + ag.Poids_rel[ag2.Id].Second*ag2.Opinion
		newOpinionAg2 := ag2.Poids_rel[ag.Id].Second*ag.Opinion + ag2.Poids_rel[ag.Id].First*ag2.PersonalParameter*ag2.Opinion*(1.0-ag2.Opinion)
		ag.Opinion = newOpinionAg
		ag2.Opinion = newOpinionAg2
	}
	ag.Opinion = math.Max(0, math.Min(1, ag.Opinion))
	ag2.Opinion = math.Max(0, math.Min(1, ag2.Opinion))
	// On met à jour les types des agents
	ag.CheckType()
	ag2.CheckType()
}

// Fonction d'action où l'agent s'engage dans une discussion
func (ag *Agent) interactWithAgent(other *Agent) ActionType {
	// Il vérifie simplement s'il est possible d'interagir
	if other.Occupied || other.CurrentAction == DiscussAct {
		return WaitAct
	}

	// Renvoie "Discuter" avec l'agent cible
	ag.DiscussingWith = other
	return DiscussAct
}

// Fonction de mise à jour de l'historique de discussion d'un agent
func (ag *Agent) addToTalkHistory(other *Agent) {
	// Vérifie si l'agent est déjà dans l'historique
	for _, a := range ag.LastTalkedTo {
		if a.Id == other.Id {
			// S'il est déjà dans l'historique, on ne l'ajoute pas à nouveau
			return
		}
	}

	// Ajoute le nouvel agent au début de l'historique
	ag.LastTalkedTo = append([]*Agent{other}, ag.LastTalkedTo...)

	// Conserve uniquement les MaxLastTalked dernier agents
	if len(ag.LastTalkedTo) > ag.MaxLastTalked {
		ag.LastTalkedTo = ag.LastTalkedTo[:ag.MaxLastTalked]
	}
}

// Fonction d'action d'un agent
func (ag *Agent) Act(env *Environnement, choice ActionType) {
	if ag.CurrentAction != RunAct {
		// L'agent est occupé
		return
	}

	switch choice {
	case MoveAct:
		ag.SendToEnv(Message{Type: MoveMsg, Agent: ag})

	case ComputerAct:
		if ag.UseComputer != nil && ag.UseComputer.TryUse() {
			ag.SetAction(ComputerAct)
			ag.Occupied = true
			ag.LastComputer = ag.UseComputer
		} else {
			ag.UseComputer = nil
		}

	case PrayAct:
		if ag.LastStatue != nil {
			ag.SetAction(PrayAct)
			ag.Occupied = true
			ag.TimeLastStatue = 0
		}

	case DiscussAct:
		if ag.DiscussingWith != nil && !ag.DiscussingWith.Occupied {
			ag.SetAction(DiscussAct)
			ag.DiscussingWith.SetAction(DiscussAct)

			ag.Occupied = true
			ag.DiscussingWith.Occupied = true

			ag.DiscussingWith.DiscussingWith = ag
			ag.addToTalkHistory(ag.DiscussingWith)
			ag.DiscussingWith.addToTalkHistory(ag)
		} else {
			ag.DiscussingWith = nil
			ag.ClearAction()
		}

	case WaitAct:
		ag.ClearAction()
	}
}

// Fonction qui met à jour l'action d'un agent
func (ag *Agent) SetAction(action ActionType) {
	ag.CurrentAction = action
	ag.DialogTimer = 180 // 2 secondes à 60 FPS
}

// Fonction qui réinitialise l'action d'un agent
func (ag *Agent) ClearAction() {
	switch ag.CurrentAction {
	case DiscussAct:
		if ag.DiscussingWith != nil {
			// Efface également l'état de l'autre agent
			ag.setOpinion(ag.DiscussingWith)
			ag.DiscussingWith.CurrentAction = RunAct
			ag.DiscussingWith.Occupied = false
			ag.DiscussingWith.DiscussingWith = nil
		}

	case PrayAct:
		log.Printf("Agent %v finished praying", ag.Id)
		log.Printf("Agent Opinion before: %v", ag.Opinion)
		if ag.TypeAgt == Believer {
			ag.Opinion = math.Min(1.0, ag.Opinion+0.05)
		} else {
			ag.Opinion = math.Min(1.0, ag.Opinion+0.1)
		}
		ag.CheckType()
		log.Printf("Agent Opinion after: %v", ag.Opinion)

	case ComputerAct:
		log.Printf("Agent %v finished using computer", ag.Id)
		log.Printf("Agent Opinion before: %v", ag.Opinion)

		currentProgram := ag.UseComputer.GetProgramm()

		switch ag.TypeAgt {
		case Believer:
			if currentProgram == NoPgm {
				ag.UseComputer.SetProgramm(GoPgm)
				ag.Opinion = math.Min(1.0, ag.Opinion+0.005)
			} else if currentProgram == GoPgm {
				ag.Opinion = math.Min(1.0, ag.Opinion+0.05)
			}

		case Sceptic:
			if currentProgram == GoPgm {
				ag.UseComputer.SetProgramm(NoPgm)
				ag.Opinion = math.Max(0.0, ag.Opinion-0.005)
			} else if currentProgram == NoPgm {
				ag.Opinion = math.Max(0.0, ag.Opinion-0.05)
			}

		case Neutral:
			if currentProgram == GoPgm {
				ag.Opinion = math.Min(1.0, ag.Opinion+0.05)
			} else if currentProgram == NoPgm {
				ag.Opinion = math.Max(0.0, ag.Opinion-0.05)
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

	ag.CurrentAction = RunAct
	ag.DialogTimer = 0
	ag.Occupied = false
	ag.DiscussingWith = nil
}

// Fonction d'action où l'agent mange: pas utilisée pour l'instant
func (ag *Agent) Eat() {
	ag.SetAction(EatAct)
}

// Fonction qui envoie un message à l'environnement via le channel de communication l'environnement
func (ag *Agent) SendToEnv(msg Message) {
	ag.Env.Communication <- msg
}

// Fonction qui renvoie un pointeur vers un agent donné à partir de son identifiant
func (env *Environnement) GetAgentById(id IdAgent) *Agent {
	for _, agent := range env.Ags {
		if agent.Id == id {
			return agent
		}
	}
	return nil
}

// Fonction qui vérifie la cohérence entre le type d'un agent et sa croyance et la met à jour si besoin
func (ag *Agent) CheckType() {
	oldType := ag.TypeAgt

	// Mise à jour du type de l'agent par rapport à son opinion
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

	// Si le type a changé, on recalcule le sous-type
	if oldType != ag.TypeAgt {
		// Log pour debug
		log.Printf("Agent %v changed type from %v to %v", ag.Id, oldType, ag.TypeAgt)
		log.Printf("Old subtype: %v", ag.SubType)

		// Si l'agent devient neutre, il perd son sous-type
		if ag.TypeAgt == Neutral {
			ag.SubType = None
		} else {
			// S'il était Neutre et est devenu un autre type, ou s'il a changé entre Croyant/Sceptique
			ag.SubType = getRandomSubType(ag.TypeAgt)
		}

		log.Printf("New subtype: %v", ag.SubType)
	}
}

// Fonction qui affiche une image sur la fenêtre d'affichage
func loadImageAgt(path string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatalf("Failed to load image: %s, error: %v", path, err)
	}
	return img
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
