package main2

// import (
// 	"fmt"
// 	"image/color"
// 	"math/rand"
// 	"time"

// 	"github.com/hajimehoshi/ebiten/v2"
// 	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
// 	"github.com/hajimehoshi/ebiten/v2/vector"
// )

// type Game struct{}

// const (
// 	AgentSize    = 10
// 	ScreenWidth  = 800
// 	ScreenHeight = 600
// 	NumAgents    = 100

// 	// Dimensões dos elementos da sala
// 	DeskWidth   = 40
// 	DeskHeight  = 30
// 	ChairWidth  = 20
// 	ChairHeight = 20

// 	// Layout da sala
// 	NumRows   = 5
// 	NumCols   = 6
// 	MarginX   = 100
// 	MarginY   = 80
// 	SpacingX  = 60
// 	SpacingY  = 50

// 	// Limites de crença
// 	SeuilCroyance   = 70
// 	SeuilScepticisme = 70

// 	maxSpeed = 2.0
// )

// type Agent struct {
// 	X, Y        int
// 	Croyance    int
// 	Scepticisme int
// 	Type        string
// 	VelX, VelY  float64
// }

// type Desk struct {
// 	X, Y int
// }

// type Chair struct {
// 	X, Y int
// }

// type Obstacle struct {
// 	X, Y      int
// 	Width, Height int
// }

// var (
// 	agents        []Agent
// 	desks         []Desk
// 	chairs        []Chair
// 	obstacles     []Obstacle
// 	selectedAgent *Agent
// )

// func init() {
// 	rand.Seed(time.Now().UnixNano())
// 	generateClassroomLayout()
// 	generateAgents(NumAgents)
// }

// func generateClassroomLayout() {
// 	// Paredes
// 	obstacles = append(obstacles, Obstacle{0, 0, ScreenWidth, 20})           // Superior
// 	obstacles = append(obstacles, Obstacle{0, 0, 20, ScreenHeight})          // Esquerda
// 	obstacles = append(obstacles, Obstacle{ScreenWidth - 20, 0, 20, ScreenHeight}) // Direita

// 	// Quadro negro
// 	obstacles = append(obstacles, Obstacle{50, 30, ScreenWidth - 100, 30})

// 	// Mesa do professor
// 	obstacles = append(obstacles, Obstacle{50, 70, 60, 40})

// 	// Gerar carteiras e cadeiras
// 	for row := 0; row < NumRows; row++ {
// 		for col := 0; col < NumCols; col++ {
// 			x := MarginX + col*(DeskWidth+SpacingX)
// 			y := MarginY + row*(DeskHeight+SpacingY)

// 			desks = append(desks, Desk{X: x, Y: y})
// 			chairs = append(chairs, Chair{X: x, Y: y + DeskHeight})

// 			// Adicionar carteiras e cadeiras como obstáculos
// 			obstacles = append(obstacles, Obstacle{x, y, DeskWidth, DeskHeight})
// 			obstacles = append(obstacles, Obstacle{x, y + DeskHeight, ChairWidth, ChairHeight})
// 		}
// 	}
// }

// func generateAgents(count int) {
// 	for i := 0; i < count; i++ {
// 		x := rand.Intn(ScreenWidth - AgentSize)
// 		y := rand.Intn(ScreenHeight - AgentSize)

// 		// Garantir que o agente não seja gerado dentro de um obstáculo
// 		for isCollidingWithAnyObstacle(x, y) {
// 			x = rand.Intn(ScreenWidth - AgentSize)
// 			y = rand.Intn(ScreenHeight - AgentSize)
// 		}

// 		agentType := "neutre"
// 		switch rand.Intn(3) {
// 		case 0:
// 			agentType = "croyant"
// 		case 1:
// 			agentType = "sceptique"
// 		}

// 		agentCroyance := 0
// 		agentScepticisme := 0
// 		switch agentType {
// 		case "croyant":
// 			agentCroyance = rand.Intn(30) + 70
// 			agentScepticisme = 100 - agentCroyance
// 		case "sceptique":
// 			agentScepticisme = rand.Intn(30) + 70
// 			agentCroyance = 100 - agentScepticisme
// 		default:
// 			agentCroyance = rand.Intn(65)
// 			agentScepticisme = rand.Intn(65)
// 		}

// 		agents = append(agents, Agent{
// 			X:           x,
// 			Y:           y,
// 			Croyance:    agentCroyance,
// 			Scepticisme: agentScepticisme,
// 			Type:        agentType,
// 		})
// 	}
// }

// func isCollidingWithAnyObstacle(x, y int) bool {
// 	for _, obs := range obstacles {
// 		if x+AgentSize > obs.X &&
// 			x < obs.X+obs.Width &&
// 			y+AgentSize > obs.Y &&
// 			y < obs.Y+obs.Height {
// 			return true
// 		}
// 	}
// 	return false
// }

// func (g *Game) Update() error {
// 	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
// 		mouseX, mouseY := ebiten.CursorPosition()
// 		for i := range agents {
// 			if checkCollision2(agents[i], mouseX, mouseY) {
// 				selectedAgent = &agents[i]
// 			}
// 		}
// 	}

// 	for i := range agents {
// 		oldX := agents[i].X
// 		oldY := agents[i].Y

// 		agents[i].VelX += rand.Float64()*2 - 1
// 		agents[i].VelY += rand.Float64()*2 - 1

// 		if agents[i].VelX > maxSpeed {
// 			agents[i].VelX = maxSpeed
// 		} else if agents[i].VelX < -maxSpeed {
// 			agents[i].VelX = -maxSpeed
// 		}

// 		if agents[i].VelY > maxSpeed {
// 			agents[i].VelY = maxSpeed
// 		} else if agents[i].VelY < -maxSpeed {
// 			agents[i].VelY = -maxSpeed
// 		}

// 		newX := agents[i].X + int(agents[i].VelX)
// 		newY := agents[i].Y + int(agents[i].VelY)

// 		// Verificar colisões antes de atualizar a posição
// 		if !isCollidingWithAnyObstacle(newX, agents[i].Y) {
// 			agents[i].X = newX
// 		} else {
// 			agents[i].VelX *= -0.5
// 		}

// 		if !isCollidingWithAnyObstacle(agents[i].X, newY) {
// 			agents[i].Y = newY
// 		} else {
// 			agents[i].VelY *= -0.5
// 		}

// 		// class limits
// 		agents[i].X = clamp(agents[i].X, 0, ScreenWidth-AgentSize)
// 		agents[i].Y = clamp(agents[i].Y, 0, ScreenHeight-AgentSize)

// 		for j := range agents {
// 			if i != j && checkCollision(agents[i], agents[j]) {
// 				interact(&agents[i], &agents[j])

// 				agents[i].X = oldX
// 				agents[i].Y = oldY
// 				agents[i].VelX *= -1
// 				agents[i].VelY *= -1
// 			}
// 		}
// 	}
// 	return nil
// }

// func checkCollision(a1, a2 Agent) bool {
// 	return abs(a1.X-a2.X) < AgentSize && abs(a1.Y-a2.Y) < AgentSize
// }

// func checkCollision2(agent Agent, mouseX, mouseY int) bool {
// 	return mouseX >= agent.X && mouseX <= agent.X+AgentSize &&
// 		mouseY >= agent.Y && mouseY <= agent.Y+AgentSize
// }

// func interact(a1, a2 *Agent) {
// 	if a1.Type == "croyant" && a2.Type == "sceptique" {
// 		a1.Croyance -= 5
// 		a2.Croyance += 5
// 	} else if a1.Type == "sceptique" && a2.Type == "croyant" {
// 		a1.Croyance += 5
// 		a2.Croyance -= 5
// 	} else if a1.Type == "croyant" && a2.Type == "neutre" {
// 		a1.Croyance += 5
// 		a2.Croyance += 5
// 	} else if a1.Type == "neutre" && a2.Type == "croyant" {
// 		a1.Croyance += 5
// 		a2.Croyance += 5
// 	}

// 	a1.Croyance = clamp(a1.Croyance, 0, 100)
// 	a2.Croyance = clamp(a2.Croyance, 0, 100)

// 	updateType(a1)
// 	updateType(a2)
// }

// func updateType(agent *Agent) {
// 	if agent.Croyance >= SeuilCroyance {
// 		agent.Type = "croyant"
// 	} else if agent.Scepticisme >= SeuilScepticisme {
// 		agent.Type = "sceptique"
// 	} else {
// 		agent.Type = "neutre"
// 	}
// }

// func abs(x int) int {
// 	if x < 0 {
// 		return -x
// 	}
// 	return x
// }

// func clamp(val, min, max int) int {
// 	if val < min {
// 		return min
// 	} else if val > max {
// 		return max
// 	}
// 	return val
// }

// func (g *Game) Draw(screen *ebiten.Image) {
// 	// walls
// 	vector.DrawFilledRect(screen, 0, 0, float32(ScreenWidth), 20, color.RGBA{180, 180, 180, 255}, false)
// 	vector.DrawFilledRect(screen, 0, 0, 20, float32(ScreenHeight), color.RGBA{180, 180, 180, 255}, false)
// 	vector.DrawFilledRect(screen, float32(ScreenWidth-20), 0, 20, float32(ScreenHeight), color.RGBA{180, 180, 180, 255}, false)

// 	// lousa
// 	vector.DrawFilledRect(screen, 50, 30, float32(ScreenWidth-100), 30, color.RGBA{34, 139, 34, 255}, false)

// 	// prof tables
// 	vector.DrawFilledRect(screen, 50, 70, 60, 40, color.RGBA{139, 69, 19, 255}, false)

// 	for _, desk := range desks {
// 		vector.DrawFilledRect(
// 			screen,
// 			float32(desk.X),
// 			float32(desk.Y),
// 			float32(DeskWidth),
// 			float32(DeskHeight),
// 			color.RGBA{139, 69, 19, 255},
// 			false,
// 		)
// 	}

// 	for _, chair := range chairs {
// 		vector.DrawFilledRect(
// 			screen,
// 			float32(chair.X),
// 			float32(chair.Y),
// 			float32(ChairWidth),
// 			float32(ChairHeight),
// 			color.RGBA{101, 67, 33, 255},
// 			false,
// 		)
// 	}

// 	for _, agent := range agents {
// 		vector.DrawFilledRect(
// 			screen,
// 			float32(agent.X),
// 			float32(agent.Y),
// 			float32(AgentSize),
// 			float32(AgentSize),
// 			agentColor(agent.Type),
// 			false,
// 		)
// 	}

// 	countCroyants, countSceptiques, countNeutres := 0, 0, 0
// 	for _, agent := range agents {
// 		switch agent.Type {
// 		case "croyant":
// 			countCroyants++
// 		case "sceptique":
// 			countSceptiques++
// 		case "neutre":
// 			countNeutres++
// 		}
// 	}

// 	stats := fmt.Sprintf("Croyants: %d | Sceptiques: %d | Neutres: %d",
// 		countCroyants, countSceptiques, countNeutres)
// 	ebitenutil.DebugPrint(screen, stats)

// 	if selectedAgent != nil {
// 		info := fmt.Sprintf("Agent: %s\nCroyance: %d\nScepticisme: %d",
// 			selectedAgent.Type, selectedAgent.Croyance, selectedAgent.Scepticisme)
// 		ebitenutil.DebugPrintAt(screen, info, ScreenWidth-200, ScreenHeight-50)
// 	}
// }

// func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
// 	return ScreenWidth, ScreenHeight
// }

// func agentColor(agentType string) color.Color {
// 	switch agentType {
// 	case "croyant":
// 		return color.RGBA{0, 255, 0, 255}
// 	case "sceptique":
// 		return color.RGBA{255, 0, 0, 255}
// 	default:
// 		return color.RGBA{128, 128, 128, 255}
// 	}
// }

// func main() {
// 	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
// 	ebiten.SetWindowTitle("Simulation avec agents")
// 	if err := ebiten.RunGame(&Game{}); err != nil {
// 		panic(err)
// 	}
// }