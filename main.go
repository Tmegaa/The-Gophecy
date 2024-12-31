package main

import (
	sim "Gophecy/pkg/Simulation"
	"log"
)

func main() {
	// On récupère les données de configuration
	config := sim.ShowMenu()

	// On génère la simulation
	simulation := sim.NewSimulation(config)

	// On fait tourner la simulation jusqu'à rencontrer une erreur
	if err := simulation.Run(); err != nil {
		log.Fatalf("Simulation failed: %v", err)
	}
}
