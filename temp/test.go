package main2

// import (
// 	"bytes"
// 	"fmt"
// 	"image"
// 	_ "image/png"
// 	"log"
// 	"math/rand"
// 	"time"

// 	"github.com/hajimehoshi/ebiten/v2"
// 	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
// 	"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
// )

// const (
//     screenWidth  = 240
//     screenHeight = 240
//     tileSize     = 16
//     NumAgents    = 30
//     AgentSize    = 16

//     SeuilCroyance   = 70
//     SeuilScepticisme = 70
// )

// var (
//     tilesImage *ebiten.Image
//     agents     []Agent
// )

// type Agent struct {
//     X, Y        float64
//     VelX, VelY  float64
//     Croyance    int
//     Scepticisme int
//     Type        string
// }

// type Game struct {
//     layers [][]int
// }

// func init() {
//     img, _, err := image.Decode(bytes.NewReader(images.Tiles_png))
//     if err != nil {
//         log.Fatal(err)
//     }
//     tilesImage = ebiten.NewImageFromImage(img)
//     rand.Seed(time.Now().UnixNano())
//     generateAgents()
// }

// func generateAgents() {
//     for i := 0; i < NumAgents; i++ {
//         x := rand.Float64() * (screenWidth - AgentSize)
//         y := rand.Float64() * (screenHeight - AgentSize)

//         agentType := "neutre"
//         switch rand.Intn(3) {
//         case 0:
//             agentType = "croyant"
//         case 1:
//             agentType = "sceptique"
//         }

//         agentCroyance := 0
//         agentScepticisme := 0
//         switch agentType {
//         case "croyant":
//             agentCroyance = rand.Intn(30) + 70
//             agentScepticisme = 100 - agentCroyance
//         case "sceptique":
//             agentScepticisme = rand.Intn(30) + 70
//             agentCroyance = 100 - agentScepticisme
//         default:
//             agentCroyance = rand.Intn(65)
//             agentScepticisme = rand.Intn(65)
//         }

//         agents = append(agents, Agent{
//             X:           x,
//             Y:           y,
//             Croyance:    agentCroyance,
//             Scepticisme: agentScepticisme,
//             Type:        agentType,
//         })
//     }
// }

// func (g *Game) Update() error {
//     for i := range agents {
//         // Update agent positions
//         agents[i].X += agents[i].VelX
//         agents[i].Y += agents[i].VelY

//         // Add random movement
//         agents[i].VelX += rand.Float64()*0.4 - 0.2
//         agents[i].VelY += rand.Float64()*0.4 - 0.2

//         // Keep agents within bounds
//         if agents[i].X < 0 {
//             agents[i].X = 0
//             agents[i].VelX *= -1
//         }
//         if agents[i].X > screenWidth-AgentSize {
//             agents[i].X = screenWidth - AgentSize
//             agents[i].VelX *= -1
//         }
//         if agents[i].Y < 0 {
//             agents[i].Y = 0
//             agents[i].VelY *= -1
//         }
//         if agents[i].Y > screenHeight-AgentSize {
//             agents[i].Y = screenHeight - AgentSize
//             agents[i].VelY *= -1
//         }
//     }
//     return nil
// }

// func (g *Game) Draw(screen *ebiten.Image) {
//     // Draw tile layers
//     w := tilesImage.Bounds().Dx()
//     tileXCount := w / tileSize
//     const xCount = screenWidth / tileSize

//     for _, l := range g.layers {
//         for i, t := range l {
//             op := &ebiten.DrawImageOptions{}
//             op.GeoM.Translate(float64((i%xCount)*tileSize), float64((i/xCount)*tileSize))

//             sx := (t % tileXCount) * tileSize
//             sy := (t / tileXCount) * tileSize
//             screen.DrawImage(tilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image), op)
//         }
//     }

//     // Draw agents
//     for _, agent := range agents {
//         op := &ebiten.DrawImageOptions{}
//         op.GeoM.Translate(agent.X, agent.Y)

//         // Use different tiles for different agent types
//         var tileIndex int
//         switch agent.Type {
//         case "croyant":
//             tileIndex = 26 // Green tile
//         case "sceptique":
//             tileIndex = 51 // Red tile
//         default:
//             tileIndex = 76 // Gray tile
//         }

//         sx := (tileIndex % tileXCount) * tileSize
//         sy := (tileIndex / tileXCount) * tileSize
//         screen.DrawImage(tilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image), op)
//     }

//     ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()))
// }

// func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
//     return screenWidth, screenHeight
// }

// func main() {
// 	g := &Game{
// 		layers: [][]int{
// 			{
// 				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
// 				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
// 				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
// 				243, 218, 243, 243, 243, 243, 243, 243, 243, 243, 243, 218, 243, 244, 243,
// 				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,

// 				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
// 				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
// 				243, 243, 244, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
// 				243, 243, 243, 243, 243, 243, 243, 243, 243, 219, 243, 243, 243, 219, 243,
// 				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,

// 				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
// 				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
// 				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
// 				243, 218, 243, 243, 243, 243, 243, 243, 243, 243, 243, 244, 243, 243, 243,
// 				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
// 			},
// 			{
// 				5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
// 				0, 0, 0, 0, 0, 26, 27, 28, 29, 30, 31, 0, 0, 0, 0,
// 				0, 0, 0, 0, 0, 51, 52, 53, 54, 55, 56, 0, 0, 0, 0,
// 				0, 0, 0, 0, 0, 76, 77, 78, 79, 80, 81, 0, 0, 0, 0,
// 				0, 0, 0, 0, 0, 101, 102, 103, 104, 105, 106, 0, 0, 0, 0,

// 				0, 0, 0, 0, 0, 126, 127, 128, 129, 130, 131, 0, 0, 0, 0,
// 				0, 0, 0, 0, 0, 303, 303, 245, 242, 303, 303, 0, 0, 0, 0,
// 				0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
// 				0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
// 				0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,

// 				0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
// 				0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
// 				0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
// 				0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
// 				0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
// 			},
// 		},
// 	}

// 	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
// 	ebiten.SetWindowTitle("Classroom Simulation")
// 	if err := ebiten.RunGame(g); err != nil {
// 		log.Fatal(err)
// 	}
// }
