package pkg

import (
	"math"
	"sort"
	"sync"

	ut "github.com/Tmegaa/The-Gophecy/pkg/Utilitaries"
)

type VisitationMap struct {
	Positions []ut.Position // Liste des positions valides
	Visits    map[int]int   // Carte de comptage des visites (index -> compteur)
	mutex     sync.RWMutex
}

// Création d'une nouvelle carte de positions visitées
func NewVisitationMap(validPositions []ut.Position) *VisitationMap {
	visits := make(map[int]int)
	for i := range validPositions {
		visits[i] = 0
	}

	return &VisitationMap{
		Positions: validPositions,
		Visits:    visits,
	}
}

// Fonction qui augmente le compteur de visites d'une position
func (vm *VisitationMap) IncrementVisit(pos ut.Position) {
	vm.mutex.Lock()
	defer vm.mutex.Unlock()

	// Trouve l'index de la position la plus proche
	idx := vm.findNearestPositionIndex(pos)
	if idx >= 0 {
		vm.Visits[idx]++
	}
}

// Fonction qui renvoie les positions les moins visitées
func (vm *VisitationMap) GetLeastVisitedPositions(currentPos ut.Position, limit int) []ut.Position {
	vm.mutex.RLock()
	defer vm.mutex.RUnlock()

	// Crée un tableau de paires (index, compteur)
	type visitPair struct {
		index int
		count int
		dist  float64
	}
	pairs := make([]visitPair, 0, len(vm.Positions))

	// Remplit le tableau avec les données
	for idx, count := range vm.Visits {
		dist := ut.Distance(currentPos, vm.Positions[idx])
		pairs = append(pairs, visitPair{idx, count, dist})
	}

	// Trie par compteur (les moins visités en premier) et par distance
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].count == pairs[j].count {
			return pairs[i].dist < pairs[j].dist // Si même nombre de visites, préfère le plus proche
		}
		return pairs[i].count < pairs[j].count
	})

	// Retourne les positions les moins visitées
	result := make([]ut.Position, 0, limit)
	for i := 0; i < limit && i < len(pairs); i++ {
		result = append(result, vm.Positions[pairs[i].index])
	}

	return result
}

// Fonction qui renvoie l'index de la position la plus proche
func (vm *VisitationMap) findNearestPositionIndex(pos ut.Position) int {
	minDist := math.MaxFloat64
	nearestIdx := -1

	for i, validPos := range vm.Positions {
		dist := ut.Distance(pos, validPos)
		if dist < minDist {
			minDist = dist
			nearestIdx = i
		}
	}

	return nearestIdx
}
