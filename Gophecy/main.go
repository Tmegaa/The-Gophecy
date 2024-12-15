package main

import (
	sim "Gophecy/pkg/Simulation"
	"time"
)


func main() {
	
	simulation := sim.NewSimulation(
		1000,
		time.Hour,
	)

	simulation.Run()
}

