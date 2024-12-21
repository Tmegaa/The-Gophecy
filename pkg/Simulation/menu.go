// pkg/Simulation/menu.go

package simulation

import (
	"fmt"
	"strconv"
	"time"
)

type SimulationConfig struct {
	NumAgents       int
	SimulationTime  time.Duration
	// Adicione outros hiperparâmetros aqui
}

func ShowMenu() SimulationConfig {
	config := SimulationConfig{}

	fmt.Println("Bem-vindo à Simulação Gophecy!")
	
	config.NumAgents = getIntInput("Número de agentes")
	durationMinutes := getIntInput("Duração da simulação (em minutos)")
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
