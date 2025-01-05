package simulation

import (
	ag "Gophecy/pkg/Agent"
	"fmt"
	"strconv"
	"time"
)

// Type qui gère la configuration de la simulation
type SimulationConfig struct {
	NumAgents        int                 // Nombre total d'agents
	NumBelievers     int                 // Nombre d'agents croyants
	NumSceptics      int                 // Nombre d'agents sceptiques
	NumNeutrals      int                 // Nombre d'agents neutres
	SimulationTime   time.Duration       // Durée de la simulation (en minutes)
	BelieverMovement ag.MovementStrategy // Stratégie de mouvement des croyants
	ScepticMovement  ag.MovementStrategy // Stratégie de mouvement des sceptiques
	NeutralMovement  ag.MovementStrategy // Stratégie de mouvement des agents neutres
	AgentsFilePath   string              // Chemin du fichier JSON contenant les agents
}

// Fonction qui gère l'initialisation de la simulation avec les valeurs données par l'utilisateur
func ShowMenu() SimulationConfig {
	config := SimulationConfig{}

	fmt.Println("\nBienvenue dans la Simulation Gophecy!")
	fmt.Println("----------------------------------------")

	// Choix du mode de configuration
	fmt.Println("Choisissez le mode de configuration:")
	fmt.Println("1 - Nombre total d'agents")
	fmt.Println("2 - Quantité de chaque type d'agent")
	fmt.Println("3 - Charger les agents depuis un fichier")
	mode := getChoiceInput("Mode de configuration (1, 2 ou 3)")

	if mode == 1 {
		// Utilisateur donne le nombre total d'agents et la durée de la simulation
		config.NumAgents = getNumAgentsInput("Nombre total d'agents")
	} else if mode == 2 {
		// Utilisateur donne le nombre d'agents de chaque type et la durée de la simulation
		for {
			config.NumBelievers = getNumAgentsInput("Nombre d'agents croyants")
			config.NumSceptics = getNumAgentsInput("Nombre d'agents sceptiques")
			config.NumNeutrals = getNumAgentsInput("Nombre d'agents neutres")
			config.NumAgents = config.NumBelievers + config.NumSceptics + config.NumNeutrals
			if config.NumAgents <= 1911 {
				break
			} else {
				fmt.Println("La somme des agents ne doit pas dépasser 1911.")
			}
		}
	} else if mode == 3 {
		// Utilisateur donne le chemin du fichier JSON contenant les agents
		config.AgentsFilePath = getFilePathInput("Chemin du fichier JSON contenant les agents")
	}

	durationMinutes := getDurationInput("Durée de la simulation (en minutes)")
	config.SimulationTime = time.Duration(durationMinutes) * time.Minute

	// Utilisateur choisit les stratégies de mouvement
	fmt.Println("\nChoisissez la stratégie de mouvement pour chaque type d'agent:")
	fmt.Println("0 - Random")
	fmt.Println("1 - Patrol")
	fmt.Println("2 - HeatMap")
	fmt.Println("3 - Center of Mass")
	fmt.Println("----------------------------------------")

	config.BelieverMovement = ag.MovementStrategy(getStrategyInput(ag.Believer))
	config.ScepticMovement = ag.MovementStrategy(getStrategyInput(ag.Sceptic))
	config.NeutralMovement = ag.MovementStrategy(getStrategyInput(ag.Neutral))

	fmt.Println("\nRésumé de la configuration:")
	if mode == 1 {
		fmt.Printf("Nombre total d'agents: %d\n", config.NumAgents)
	} else if mode == 2 {
		fmt.Printf("Nombre d'agents croyants: %d\n", config.NumBelievers)
		fmt.Printf("Nombre d'agents sceptiques: %d\n", config.NumSceptics)
		fmt.Printf("Nombre d'agents neutres: %d\n", config.NumNeutrals)
	} else if mode == 3 {
		fmt.Printf("Chemin du fichier JSON contenant les agents: %s\n", config.AgentsFilePath)
	}
	fmt.Printf("Durée: %v\n", config.SimulationTime)
	fmt.Printf("Stratégie %ss: %s\n", ag.Believer, config.BelieverMovement)
	fmt.Printf("Stratégie %ss: %s\n", ag.Sceptic, config.ScepticMovement)
	fmt.Printf("Stratégie %ss: %s\n", ag.Neutral, config.NeutralMovement)
	fmt.Println("----------------------------------------")

	return config
}

// Fonction qui affiche un message et récupère un entier rentré par l'utilisateur
func getNumAgentsInput(prompt string) int {
	var input string
	var value int
	var err error

	for {
		fmt.Printf("%s: ", prompt)
		fmt.Scanln(&input)
		value, err = strconv.Atoi(input)
		if err == nil && value >= 0 && value < 1912 {
			return value
		}
		fmt.Println("Veuillez entrer un nombre valide supérieur ou égal à 0 et inférieur à 1912.")
	}
}

// Fonction qui affiche un message et récupère un entier rentré par l'utilisateur
func getDurationInput(prompt string) int {
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

// Fonction qui affiche un message et récupère un entier rentré par l'utilisateur
func getChoiceInput(prompt string) int {
	var input string
	var value int
	var err error

	for {
		fmt.Printf("%s: ", prompt)
		fmt.Scanln(&input)
		value, err = strconv.Atoi(input)
		if err == nil && (value == 1 || value == 2 || value == 3) {
			return value
		}
		fmt.Println("Veuillez entrer 1, 2 ou 3.")
	}
}

// Fonction qui affiche le type d'un agent et récupère un entier rentré par l'utilisateur
func getStrategyInput(agentType ag.TypeAgent) int {
	var input string
	var value int
	var err error

	for {
		fmt.Printf("Stratégie pour %ss (0-3): ", agentType)
		fmt.Scanln(&input)
		value, err = strconv.Atoi(input)
		if err == nil && value >= 0 && value <= 3 {
			return value
		}
		fmt.Println("Veuillez entrer un nombre entre 0 et 3.")
	}
}

// Fonction qui affiche un message et récupère le chemin du fichier rentré par l'utilisateur
func getFilePathInput(prompt string) string {
	var input string

	for {
		fmt.Printf("%s: ", prompt)
		fmt.Scanln(&input)
		if input != "" {
			return input
		}
		fmt.Println("Veuillez entrer un chemin de fichier valide.")
	}
}
