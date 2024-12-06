package simulation

import (
	ag "Gophecy/pkg/Agent"
	"sync"
	"time"
)

type Simulation struct {
	env         ag.Environnement
	agents      []ag.Agent
	maxStep     int
	maxDuration time.Duration
	step        int //Stats
	start       time.Time
	syncChans   sync.Map
}
