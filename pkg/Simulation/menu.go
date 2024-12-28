package simulation

import (
	"fmt"
	"strconv"
	"time"
)

// Type qui gère la configuration de la simulation
type SimulationConfig struct {
	NumAgents      int           // Nombre d'agents
	SimulationTime time.Duration // Durée de la simulation (en minutes)
}

// Fonction qui gère l'initialisation de la simulation avec les valeurs données par l'utilisateur
func ShowMenu() SimulationConfig {
	config := SimulationConfig{}

	fmt.Println("Bienvenue dans la Simulation Gophecie!")

	// Utilisateur donne le nombre d'agents et la durée de la simulation
	config.NumAgents = getIntInput("Nombre d'agents")
	durationMinutes := getIntInput("Durée de la simulation (en minutes)")
	config.SimulationTime = time.Duration(durationMinutes) * time.Minute

	// TODO: mettre plus de paramètres de simulation ici si besoin
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
		if err == nil {
			return value
		}
		fmt.Println("Veuillez donner un entier positif.")
	}
}
