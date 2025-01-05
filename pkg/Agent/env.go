package pkg

import (
	carte "Gophecy/pkg/Carte"
	ut "Gophecy/pkg/Utilitaries"
	"image"
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

// Liste des types de messages possibles
type MessageType string

const (
	LoopMsg       MessageType = "StartLoop"
	PerceptionMsg MessageType = "Perception"
	NearbyMsg     MessageType = "Nearby"
	MoveMsg       MessageType = "Move"
)

type Message struct {
	Type         MessageType
	NearbyAgents []*Agent
	Agent        *Agent
}

type Environnement struct {
	sync.RWMutex
	Ags             []*Agent
	Carte           *carte.Carte
	Objs            []InterfaceObjet
	Communication   chan Message //key = IDAgent et value = []*Message -> Liste des messages reçus par l'agent
	NbrAgents       *sync.Map    //key = typeAgent et value = int  -> Compteur d'agents par types
	AgentProximity  *sync.Map    //key = IDAgent et value = []*Agent -> Liste des agents proches
	ObjectProximity *sync.Map    //key = IDAgent et value = []*Objet -> Liste des objets proches
}

// Fonction d'initialisation d'un nouvel environnement
func NewEnvironment(ags []*Agent, carte *carte.Carte, objs []InterfaceObjet) *Environnement {
	// Initialisation du compteur du nombre d'agents par type
	counter := &sync.Map{}

	for _, val := range lsType {
		counter.Store(val, 0)
	}

	return &Environnement{Ags: ags, Objs: objs, Communication: make(chan Message, 100), NbrAgents: counter, Carte: carte, AgentProximity: &sync.Map{}}
}

// Fonction qui ajoute un nouvel agent dans l'environnement
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

// Fonction qui envoie la liste des agents proches pour un agent donné
func (env *Environnement) NearbyAgents(ag *Agent) []*Agent {
	nearbyAgents := make([]*Agent, 0)
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
func (env *Environnement) NearbyObjects(ag *Agent) []*InterfaceObjet {
	nearbyObjects := make([]*InterfaceObjet, 0)
	pos := ag.AgtPosition()

	// Création du rectangle de perception
	var area ut.Rectangle
	area.PositionDL.X = pos.X - ag.Acuite
	area.PositionDL.Y = pos.Y + ag.Acuite
	area.PositionUR.X = pos.X + ag.Acuite
	area.PositionUR.Y = pos.Y - ag.Acuite

	// Pour les ordinateurs: On itère sur tous les objets de la carte pour voir si un objet est dans le rectangle de perception
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

// Fonction de l'environnement qui gère la communication avec les agents via les channels
func (env *Environnement) Listen() {
	go func() {
		for msg := range env.Communication {
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

// Fonction d'envoi d'un message à un agent via son channel
func (env *Environnement) SendToAgent(agt *Agent, msg Message) {
	agt.SyncChan <- msg
}

// Fonction de movement d'un agent dans l'environnement (soit mise à jour de sa position)
func (env *Environnement) Move(ag *Agent) {
	ag.ClearAction()

	// On reste dans la même direction
	if ag.MoveTimer > 0 {
		ag.MoveTimer -= 1
		// Si une collision a lieu on s'arrête
		if CheckCollisionHorizontal((ag.Position.X+ag.Position.Dx), (ag.Position.Y+ag.Position.Dy), ag.Env.Carte.Coliders) ||
			CheckCollisionVertical((ag.Position.X+ag.Position.Dx), (ag.Position.Y+ag.Position.Dy), ag.Env.Carte.Coliders) {
			return
		}

		// Sinon on continue de bouger
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
		env.moveRandom(ag) // Fallback pour le mouvement aléatoire
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
	// Possibilité de changer de direction même si vous n'avez pas atteint le waypoint
	if ag.CurrentWaypoint != nil && rand.Float64() < 0.02 {
		ag.CurrentWaypoint = nil
	}

	// S'il n'y a pas de waypoint actuel ou s'il est proche du waypoint actuel
	if ag.CurrentWaypoint == nil || ut.Distance(ag.Position, *ag.CurrentWaypoint) < 5.0 {
		if len(ag.HeatMap.Positions) > 0 {
			// Essaye de trouver un waypoint valide
			maxAttempts := 10 // Limiter les tentatives pour éviter une boucle infinie
			attempts := 0

			for attempts < maxAttempts {
				// Sélectionne 3 points aléatoires
				numChoices := 3
				choices := make([]ut.Position, 0, numChoices)

				for i := 0; i < numChoices; i++ {
					randomIdx := rand.Intn(len(ag.HeatMap.Positions))
					pos := ag.HeatMap.Positions[randomIdx]

					// Vérifie si le chemin vers le point est dégagé
					if isPathClear(ag.Position, pos, env.Carte.Coliders) {
						choices = append(choices, pos)
					}
				}

				if len(choices) > 0 {
					var bestChoice ut.Position
					bestScore := -1.0

					for _, pos := range choices {
						dist := ut.Distance(ag.Position, pos)

						// Ajuste les critères de distance
						distScore := 0.0
						if dist < 50 { // Réduit la distance minimale
							distScore = dist / 50
						} else if dist > 200 { // Réduit la distance maximale
							distScore = 2 - (dist-200)/200
						} else {
							distScore = 1.0
						}

						// Ajoute un facteur de déviation des objets
						obstacleScore := getObstacleAvoidanceScore(pos, env.Carte.Coliders)

						// Combine les scores
						randomFactor := 0.5 + rand.Float64()
						finalScore := (distScore*0.4 + obstacleScore*0.4) * randomFactor

						if finalScore > bestScore {
							bestScore = finalScore
							bestChoice = pos
						}
					}

					if bestScore > 0 {
						ag.CurrentWaypoint = &bestChoice
						break
					}
				}

				attempts++
			}

			// S'il n'a pas trouvé de waypoint valide, utilisez un mouvement aléatoire
			if ag.CurrentWaypoint == nil {
				env.moveRandom(ag)
				return
			}
		}
	}

	// Déplacement vers le waypoint
	if ag.CurrentWaypoint != nil {
		dx := ag.CurrentWaypoint.X - ag.Position.X
		dy := ag.CurrentWaypoint.Y - ag.Position.Y

		// Réduit la variation aléatoire
		dx += (rand.Float64()*2 - 1) * 2 // Réduit à ±2 pixels
		dy += (rand.Float64()*2 - 1) * 2

		length := math.Sqrt(dx*dx + dy*dy)
		if length > 0 {
			speed := ut.Maxspeed * (0.9 + rand.Float64()*0.2) // Vitesse plus constante
			ag.Position.Dx = (dx / length) * speed
			ag.Position.Dy = (dy / length) * speed
		}
	}
}

// Fonction auxiliaire pour vérifier si le chemin est clair
func isPathClear(start, end ut.Position, coliders []image.Rectangle) bool {
	// Vérifie quelques points en cours de route
	steps := 10
	dx := (end.X - start.X) / float64(steps)
	dy := (end.Y - start.Y) / float64(steps)

	for i := 0; i <= steps; i++ {
		x := start.X + dx*float64(i)
		y := start.Y + dy*float64(i)

		// Vérifier la collision au point
		for _, colider := range coliders {
			if colider.Overlaps(image.Rect(int(x)-8, int(y)-8, int(x)+8, int(y)+8)) {
				return false
			}
		}
	}
	return true
}

// Fonction pour calculer le score d'évitement d'obstacles
func getObstacleAvoidanceScore(pos ut.Position, coliders []image.Rectangle) float64 {
	minDistance := math.MaxFloat64

	// Trouve la distance jusqu'à l'obstacle le plus proche
	for _, colider := range coliders {
		centerX := float64(colider.Min.X+colider.Max.X) / 2
		centerY := float64(colider.Min.Y+colider.Max.Y) / 2

		dist := math.Sqrt(math.Pow(pos.X-centerX, 2) + math.Pow(pos.Y-centerY, 2))
		if dist < minDistance {
			minDistance = dist
		}
	}

	// Normalise le score (plus on s'éloigne des obstacles, mieux c'est)
	if minDistance < 30 {
		return 0
	}
	if minDistance > 100 {
		return 1
	}
	return (minDistance - 30) / 70
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

// Fonction qui définit les poids des agents
func (env *Environnement) SetPoids() {
	minWeight := 0.01
	maxWeight := 1.0
	for _, ag := range env.Ags {
		sum := 0.0
		for _, ag2 := range env.Ags {
			// Pour chaque agent déjà existant de l'environnement,
			// on affect un poids absolu aléatoire et impacté par la relation entre agents
			// et on calcule le poids relatif
			ag.Poids_abs[ag2.ID()] = minWeight + rand.Float64()*(maxWeight-minWeight)*ag.Relation[ag2.ID()]
			sum += ag.Poids_abs[ag2.Id]
		}
		// On applique la propriété de normalisation des poids absolus
		for _, ag2 := range env.Ags {
			pairAg := ut.Pair{}
			if sum > 0 {
				ag.Poids_abs[ag2.Id] = ag.Poids_abs[ag2.Id] / sum
			} else {
				ag.Poids_abs[ag2.Id] = 0
			}
			if ag.Id != ag2.Id {
				pairAg.Second = ag.Poids_abs[ag2.Id] / (ag.Poids_abs[ag2.Id] + ag.Poids_abs[ag.Id])
				pairAg.First = ag.Poids_abs[ag.Id] / (ag.Poids_abs[ag2.Id] + ag.Poids_abs[ag.Id])
			}
			ag.Poids_rel[ag2.Id] = pairAg
		}

	}
}

// Fonction qui définie les relations entre agents
func (env *Environnement) SetRelations() {
	for _, ag := range env.Ags {
		for _, ag2 := range env.Ags {
			if ag.ID() != ag2.ID() {
				close := rand.Float64()
				switch {
				case close < 0.25: //ennemi
					ag.Relation[ag2.ID()] = 0.75
				case close < 0.5: //pas de lien direct
					ag.Relation[ag2.ID()] = 1
				case close < 0.75: //amis
					ag.Relation[ag2.ID()] = 1.25
				case close < 1: //famille
					ag.Relation[ag2.ID()] = 1.5
				}
			} else {
				ag.Relation[ag2.ID()] = 1
			}

		}
	}
}
