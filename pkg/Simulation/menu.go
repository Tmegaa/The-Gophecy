// pkg/Simulation/menu.go

package simulation

import (
	ag "Gophecy/pkg/Agent"
	"fmt"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)


type SimulationConfig struct {
	NumAgents        int
	SimulationTime   time.Duration
	BelieverMovement ag.MovementStrategy
	ScepticMovement  ag.MovementStrategy
	NeutralMovement  ag.MovementStrategy
	BelieverStrategy ag.TypeAgentStrategy
	ScepticStrategy  ag.TypeAgentStrategy
	NeutralStrategy  ag.TypeAgentStrategy
}

func ShowMenu() {
	fmt.Println("\nBienvenue dans la Simulation Gophecy!")
	fmt.Println("----------------------------------------")
	fmt.Println("1 - Configurer une simulation manuellement")
	fmt.Println("2 - Exécuter des simulations automatiques")
	fmt.Println("----------------------------------------")

	choice := getIntInput("Choisissez une option")

	switch choice {
	case 1:
		config := configureManualSimulation()
		sim := NewSimulation(config)
		if err := ebiten.RunGame(sim); err != nil && err != ebiten.Termination {
			fmt.Printf("Error running simulation: %v\n", err)
		}
	case 2:
		RunAutomatedSimulations()
	default:
		fmt.Println("Option invalide. Veuillez choisir 1 ou 2.")
	}
}

func configureManualSimulation() SimulationConfig {
	config := SimulationConfig{}

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

func RunAutomatedSimulations() {
	agentCounts := []int{10, 20, 30, 40, 50, 60}
	movementStrategies := [][]int{
		{0, 2, 3},
		{0, 1, 2},
		{2, 3, 2},
		{0, 1, 3},
		{1, 2, 3},
		{3, 0, 1},
		{2, 1, 0},
	}
	simulationTime := 5 * time.Minute

	for _, count := range agentCounts {
		for _, strategies := range movementStrategies {
			fmt.Printf("\n--- Running Automated Simulation ---\n")
			fmt.Printf("Agents: %d, Strategies: %v\n", count, strategies)

			config := SimulationConfig{
				NumAgents:        count,
				SimulationTime:   simulationTime,
				BelieverMovement: ag.MovementStrategy(strategies[0]),
				ScepticMovement:  ag.MovementStrategy(strategies[1]),
				NeutralMovement:  ag.MovementStrategy(strategies[2]),
			}

			sim := NewSimulation(config)
			if err := ebiten.RunGame(sim); err != nil && err != ebiten.Termination {
				fmt.Printf("Error running simulation: %v\n", err)
			}

			time.Sleep(1 * time.Second)
		}
	}
}
