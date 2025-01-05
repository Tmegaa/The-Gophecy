package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	ag "Gophecy/pkg/Agent"
	simulation "Gophecy/pkg/Simulation"
)

func main() {
	
	agentCount := flag.Int("agents", 10, "Number of agents")
	strategy1 := flag.Int("strategies1", 0, "Believer movement strategy")
	strategy2 := flag.Int("strategies2", 0, "Sceptic movement strategy")
	strategy3 := flag.Int("strategies3", 0, "Neutral movement strategy")
	agentS1 := flag.Int("strategiesA1", 0, "Number of agents with strategy 1")
	agentS2 := flag.Int("strategiesA2", 0, "Number of agents with strategy 2")
	agentS3 := flag.Int("strategiesA3", 0, "Number of agents with strategy 3")
	simulationTime := flag.Int("time", 5, "Simulation time in seconds")
	flag.Parse()

	fmt.Printf("\n--- Running Simulation ---\n")
	fmt.Printf("Agents: %d, Strategies: [%d, %d, %d], Time: %d Minutes\n Agents with strategies: [%d, %d, %d]\n", *agentCount, *strategy1, *strategy2, *strategy3, *simulationTime, *agentS1, *agentS2, *agentS3)
	
	config := simulation.SimulationConfig{
		NumAgents:        *agentCount,
		SimulationTime:   time.Duration(*simulationTime) * time.Minute,
		BelieverMovement: ag.MovementStrategy(*strategy1),
		ScepticMovement:  ag.MovementStrategy(*strategy2),
		NeutralMovement:  ag.MovementStrategy(*strategy3),
		BelieverStrategy: ag.TypeAgentStrategy(*agentS1),
		ScepticStrategy:  ag.TypeAgentStrategy(*agentS2),
		NeutralStrategy:  ag.TypeAgentStrategy(*agentS3),
	}


	sim := simulation.NewSimulation(config)

	
	if err := sim.Run(); err != nil {
		log.Fatalf("Simulation failed: %v", err)
	}
}
