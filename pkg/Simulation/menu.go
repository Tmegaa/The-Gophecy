package simulation

import (
	ag "Gophecy/pkg/Agent"
	"fmt"
	"strconv"
	"time"
)

// Type qui gère la configuration de la simulation
type SimulationConfig struct {
	NumAgents        int                 // Nombre d'agents
	SimulationTime   time.Duration       // Durée de la simulation (en minutes)
	BelieverMovement ag.MovementStrategy // Stratégie de mouvement des croyants
	ScepticMovement  ag.MovementStrategy // Stratégie de mouvement des sceptiques
	NeutralMovement  ag.MovementStrategy // Stratégie de mouvement des agents neutres
}

// Fonction qui gère l'initialisation de la simulation avec les valeurs données par l'utilisateur
func ShowMenu() SimulationConfig {
	config := SimulationConfig{}

	fmt.Println("\nBienvenue dans la Simulation Gophecy!")
	fmt.Println("----------------------------------------")

	// Utilisateur donne le nombre d'agents et la durée de la simulation
	config.NumAgents = getIntInput("Nombre d'agents")
	durationMinutes := getIntInput("Durée de la simulation (en minutes)")
	config.SimulationTime = time.Duration(durationMinutes) * time.Minute

	// Utilisateur choisit les stratégies de mouvement
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

// Fonction qui affiche un message et récupère un entier rentré par l'utilisateur
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

// Fonction qui affiche le type d'un agent et récupère un entier rentré par l'utilisateur
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
