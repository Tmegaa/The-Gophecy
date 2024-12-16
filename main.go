package main

import (
	sim "Gophecy/pkg/Simulation"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)


func main() {
	
	simulation := sim.NewSimulation(
		1000,
		time.Hour,
	)

	if err := ebiten.RunGame(simulation); err != nil {
		log.Fatal(err)
	}
}

