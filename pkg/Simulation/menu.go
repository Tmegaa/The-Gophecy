// pkg/Simulation/menu.go

package simulation

import (
	ag "Gophecy/pkg/Agent"
	"fmt"
	"strconv"
	"time"
)

type SimulationConfig struct {
	NumAgents           int
	SimulationTime      time.Duration
	BelieverMovement    ag.MovementStrategy
	ScepticMovement     ag.MovementStrategy
	NeutralMovement     ag.MovementStrategy
}

func ShowMenu() SimulationConfig {
	config := SimulationConfig{}

	fmt.Println("\nBienvenue dans la Simulation Gophecy!")
	fmt.Println("----------------------------------------")

	config.NumAgents = getIntInput("Nombre d'agents")
	durationMinutes := getIntInput("Durée de la simulation (en minutes)")
	config.SimulationTime = time.Duration(durationMinutes) * time.Minute

	fmt.Println("\nChoisissez la stratégie de mouvement pour chaque type d'agent:")
	fmt.Println("0 - Random")
	fmt.Println("1 - Patrol")
	fmt.Println("2 - HeatMap")
	fmt.Println("3 - Center of Mass")
	fmt.Println("----------------------------------------")

	config.BelieverMovement = ag.MovementStrategy(getStrategyInput("Believer"))
	config.ScepticMovement = ag.MovementStrategy(getStrategyInput("Sceptic"))
	config.NeutralMovement = ag.MovementStrategy(getStrategyInput("Neutral"))

	fmt.Println("\nRésumé de la configuration:")
	fmt.Printf("Nombre d'agents: %d\n", config.NumAgents)
	fmt.Printf("Durée: %v\n", config.SimulationTime)
	fmt.Printf("Stratégie Believer: %s\n", config.BelieverMovement)
	fmt.Printf("Stratégie Sceptic: %s\n", config.ScepticMovement)
	fmt.Printf("Stratégie Neutral: %s\n", config.NeutralMovement)
	fmt.Println("----------------------------------------")

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
		if err == nil && value > 0 {
			return value
		}
		fmt.Println("Veuillez entrer un nombre valide supérieur à 0.")
	}
}

func getStrategyInput(agentType string) int {
	var input string
	var value int
	var err error

	for {
		fmt.Printf("Stratégie pour %s (0-3): ", agentType)
		fmt.Scanln(&input)
		value, err = strconv.Atoi(input)
		if err == nil && value >= 0 && value <= 3 {
			return value
		}
		fmt.Println("Veuillez entrer un nombre entre 0 et 3.")
	}
}
