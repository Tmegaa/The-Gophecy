package main

import (
	sim "Gophecy/pkg/Simulation"
	"log"
)

func main() {
	config := sim.ShowMenu()

	simulation := sim.NewSimulation(config)
	// for _, ag := range simulation.Env.Ags {
	// 	fmt.Println(ag.Id)
	// }

	if err := simulation.Run(); err != nil {
		log.Fatalf("Simulation failed: %v", err)
	}
}
