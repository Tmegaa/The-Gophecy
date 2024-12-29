package pkg

import (
	carte "Gophecy/pkg/Carte"
	ut "Gophecy/pkg/Utilitaries"
	"math"
	"math/rand"
	"sync"
)

var lsType = []TypeAgent{Sceptic, Believer, Neutral}

type MovementStrategy int

const (
	RandomMovement MovementStrategy = iota
	PatrolMovement
	HeatMapMovement
	CenterOfMassMovement
)

func (m MovementStrategy) String() string {
	return [...]string{"Random", "Patrol", "HeatMap", "CenterOfMass"}[m]
}

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

func (env *Environnement) NearbyObjects(ag *Agent) []*InterfaceObjet {
	nearbyObjects := make([]*InterfaceObjet, 0)
	pos := ag.AgtPosition()
	var area ut.Rectangle
	area.PositionDL.X = pos.X - ag.Acuite
	area.PositionDL.Y = pos.Y + ag.Acuite
	area.PositionUR.X = pos.X + ag.Acuite
	area.PositionUR.Y = pos.Y - ag.Acuite


	//for PC
	for _, pc := range env.Objs {
		
		PCposition := pc.ObjPosition()

		if ut.IsInRectangle(PCposition, area) {
			if pc.GetUse() && (ag.LastComputer == nil || pc.ID() != ag.LastComputer.ID()) {
				continue
			}
			nearbyObjects = append(nearbyObjects, &pc)
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
		if CheckCollisionHorizontal((ag.Position.X+ag.Position.Dx), (ag.Position.Y+ag.Position.Dy), ag.Env.Carte.Coliders) || 
		   CheckCollisionVertical((ag.Position.X+ag.Position.Dx), (ag.Position.Y+ag.Position.Dy), ag.Env.Carte.Coliders) {
			return
		}
		
		ag.Position.X += ag.Position.Dx
		ag.Position.Y += ag.Position.Dy
		return
	}

	// Utilise la stratégie définie pour chaque type d'agent
	switch ag.MovementStrategy {
	case RandomMovement:
		env.moveRandom(ag)
	case PatrolMovement:
		env.movePatrol(ag)
	case HeatMapMovement:
		env.moveWithHeatMap(ag)
	case CenterOfMassMovement:
		env.moveToCenterOfMass(ag)
	default:
		env.moveRandom(ag) // Fallback pour mouvement aléatoire
	}

	ag.MoveTimer = 60
}

// Mouvement basé sur la carte de chaleur pour les Believers et Sceptics
func (env *Environnement) moveWithHeatMap(ag *Agent) {
	// Obtient les 3 positions les moins visitées proches de l'agent
	leastVisited := ag.HeatMap.GetLeastVisitedPositions(ag.Position, 3)
	
	// 70% de chance d'aller vers une position moins visitée
	if len(leastVisited) > 0 && rand.Float64() < 0.7 {
		// Choisit aléatoirement une des positions les moins visitées
		targetPos := leastVisited[rand.Intn(len(leastVisited))]
		// Calcule la direction vers la position choisie
		dx := targetPos.X - ag.Position.X
		dy := targetPos.Y - ag.Position.Y
		// Normalise la direction
		length := math.Sqrt(dx*dx + dy*dy)
		if length > 0 {
			ag.Position.Dx = (dx / length) * ut.Maxspeed
			ag.Position.Dy = (dy / length) * ut.Maxspeed
		}
	} else {
		env.moveRandom(ag)
	}
}

// Mouvement de patrouille pour les Neutrals
func (env *Environnement) movePatrol(ag *Agent) {
	// Chance de changer de direction même s'il n'a pas atteint le waypoint
	if ag.CurrentWaypoint != nil && rand.Float64() < 0.02 { // 2% de chance par tick de changer de direction
		ag.CurrentWaypoint = nil
	}

	// Si pas de waypoint actuel ou proche du waypoint actuel
	if ag.CurrentWaypoint == nil || ut.Distance(ag.Position, *ag.CurrentWaypoint) < 5.0 {
		// Choisit un nouveau waypoint
		if len(ag.HeatMap.Positions) > 0 {
			// Prend 3 points aléatoires et en choisit un
			numChoices := 3
			choices := make([]ut.Position, 0, numChoices)
			
			for i := 0; i < numChoices; i++ {
				randomIdx := rand.Intn(len(ag.HeatMap.Positions))
				choices = append(choices, ag.HeatMap.Positions[randomIdx])
			}

			// Choisit le point basé sur une combinaison de :
			// - Distance (préfère les points ni trop proches ni trop éloignés)
			// - Aléatoire
			var bestChoice ut.Position
			bestScore := -1.0

			for _, pos := range choices {
				dist := ut.Distance(ag.Position, pos)
				
				// Score basé sur la distance (préfère les distances moyennes, entre 100 et 300 pixels)
				distScore := 0.0
				if dist < 100 {
					distScore = dist / 100 // Score augmente jusqu'à 100
				} else if dist > 300 {
					distScore = 2 - (dist-300)/300 // Score diminue après 300
				} else {
					distScore = 1.0 // Distance idéale
				}

				// Ajoute de l'aléatoire au score
				randomFactor := 0.5 + rand.Float64()
				finalScore := distScore * randomFactor

				if finalScore > bestScore {
					bestScore = finalScore
					bestChoice = pos
				}
			}

			ag.CurrentWaypoint = &bestChoice
		}
	}

	// Si a un waypoint, se déplace vers lui avec une certaine variation
	if ag.CurrentWaypoint != nil {
		dx := ag.CurrentWaypoint.X - ag.Position.X
		dy := ag.CurrentWaypoint.Y - ag.Position.Y
		
		// Ajoute une petite variation aléatoire au mouvement
		dx += (rand.Float64()*2 - 1) * 5 // Variation de ±5 pixels
		dy += (rand.Float64()*2 - 1) * 5

		// Normalise la direction
		length := math.Sqrt(dx*dx + dy*dy)
		if length > 0 {
			speed := ut.Maxspeed * (0.8 + rand.Float64()*0.4) // Vitesse varie entre 80% et 120% de la maximale
			ag.Position.Dx = (dx / length) * speed
			ag.Position.Dy = (dy / length) * speed
		}
	} else {
		env.moveRandom(ag)
	}
}

// Mouvement aléatoire (utilisé comme fallback)
func (env *Environnement) moveRandom(ag *Agent) {
	directions := []ut.UniqueDirection{
		{Dx: ut.Maxspeed, Dy: 0},
		{Dx: -ut.Maxspeed, Dy: 0},
		{Dx: 0, Dy: ut.Maxspeed},
		{Dx: 0, Dy: -ut.Maxspeed},
	}
	
	randIdx := rand.Intn(len(directions))
	ag.Position.Dx = directions[randIdx].Dx
	ag.Position.Dy = directions[randIdx].Dy
}

// Mouvement basé sur le centre de masse des agents proches
func (env *Environnement) moveToCenterOfMass(ag *Agent) {
	// Obtient tous les agents dans un rayon plus grand que la normale pour considérer des groupes distants
	searchRadius := ag.Acuite * 2 // Double le rayon de recherche pour détecter des groupes plus grands
	nearbyAgents := make([]Agent, 0)
	
	// Recherche des agents proches dans un rayon plus grand
	pos := ag.Position
	var searchArea ut.Rectangle
	searchArea.PositionDL.X = pos.X - searchRadius
	searchArea.PositionDL.Y = pos.Y + searchRadius
	searchArea.PositionUR.X = pos.X + searchRadius
	searchArea.PositionUR.Y = pos.Y - searchRadius

	// Collecte tous les agents dans la zone
	for _, other := range env.Ags {
		if ag.ID() != other.ID() && ut.IsInRectangle(other.Position, searchArea) {
			nearbyAgents = append(nearbyAgents, *other)
		}
	}

	// S'il n'y a pas d'agents proches, se déplace aléatoirement
	if len(nearbyAgents) == 0 {
		env.moveRandom(ag)
		return
	}

	// Calcule le centre de masse
	var centerX, centerY float64
	totalWeight := 0.0

	for _, other := range nearbyAgents {
		// Poids basé sur la distance (les agents plus proches ont plus d'influence)
		dist := ut.Distance(ag.Position, other.Position)
		weight := 1.0 / (dist + 1) // Évite la division par zéro
		
		// Ajoute un bonus de poids pour les agents du même type
		if other.TypeAgt == ag.TypeAgt {
			weight *= 1.5
		}

		centerX += other.Position.X * weight
		centerY += other.Position.Y * weight
		totalWeight += weight
	}

	// Normalise le centre de masse
	if totalWeight > 0 {
		centerX /= totalWeight
		centerY /= totalWeight
	}

	// Ajoute un élément de hasard pour éviter un regroupement parfait
	centerX += (rand.Float64()*2 - 1) * 10
	centerY += (rand.Float64()*2 - 1) * 10

	// Calcule la direction vers le centre de masse
	dx := centerX - ag.Position.X
	dy := centerY - ag.Position.Y

	// Normalise la direction et applique la vitesse
	length := math.Sqrt(dx*dx + dy*dy)
	if length > 0 {
		// Varie la vitesse en fonction de la distance au centre de masse
		speedFactor := 0.8 + rand.Float64()*0.4 // Vitesse entre 80% et 120% de la vitesse maximale
		ag.Position.Dx = (dx / length) * ut.Maxspeed * speedFactor
		ag.Position.Dy = (dy / length) * ut.Maxspeed * speedFactor
	}

	// Petite chance de passer à un mouvement aléatoire pour éviter un regroupement excessif
	if rand.Float64() < 0.05 { // 5% de chance
		env.moveRandom(ag)
	}
}
