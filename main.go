package main

import (
	sim "Gophecy/pkg/Simulation"
	"log"
)

func main() {
	config := sim.ShowMenu()
	
	simulation := sim.NewSimulation(config)

	if err := simulation.Run(); err != nil {
		log.Fatalf("Simulation failed: %v", err)
	}
}
