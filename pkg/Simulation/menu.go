// pkg/Simulation/menu.go

package simulation

import (
	"fmt"
	"strconv"
	"time"
)

type SimulationConfig struct {
	NumAgents      int
	SimulationTime time.Duration
}

func ShowMenu() SimulationConfig {
	config := SimulationConfig{}

	fmt.Println("Bienvenue dans la Simulation Gophecy!")

	config.NumAgents = getIntInput("Nombre d'agents")
	durationMinutes := getIntInput("Durée de la simulation (en minutes)")
	config.SimulationTime = time.Duration(durationMinutes) * time.Minute

	// Adicione mais opções de configuração aqui
	return config
}

func getIntInput(prompt string) int {
	var input string
	var value int
	var err error

	for {
		fmt.Printf("%s: ", prompt)
		fmt.Scanln(&input)
		value, err = strconv.Atoi(input)
		if err == nil {
			return value
		}
		fmt.Println("Por favor, insira um número válido.")
	}
}
