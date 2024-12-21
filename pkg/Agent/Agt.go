package pkg

import (
	ut "Gophecy/pkg/Utilitaries"
	"log"

	"math/rand"

	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type InterfaceAgent interface {
	Avancer()
	Discuter()
	Programmer()
	Perceive(*Environnement)
	Deliberate()
	//Act(*Environnement,typeFunc fun()) A implémenter
}

type TypeAgent string

const (
	Sceptic  TypeAgent = "Sceptic"
	Believer TypeAgent = "Believer"
	Neutral  TypeAgent = "Neutral"
)

type IdAgent string

type Agent struct {
	Env               *Environnement
	Id                IdAgent
	Velocite          float64
	Acuite            float64
	Position          ut.Position
	Opinion           float64
	Charisme          map[IdAgent]float64 //influence d'un agent sur un autre
	Relation          map[IdAgent]float64
	PersonalParameter float64
	Poid_rel          []float64
	Vivant            bool
	TypeAgt           TypeAgent
	SyncChan          chan int
	Img               *ebiten.Image
	MoveTimer 	      int
	CurrentAction     string
    DialogTimer       int
}

func NewAgent(env *Environnement, id IdAgent, velocite float64, acuite float64, position ut.Position,
	opinion float64, charisme map[IdAgent]float64, relation map[IdAgent]float64, personalParameter float64, typeAgt TypeAgent, syncChan chan int, img *ebiten.Image) *Agent {

	//calcul des poids relatif pour chaque agents
	poid_rel := make([]float64, 0)
	personalCharisme := charisme[id]
	for _, v := range charisme {
		char := v / personalCharisme
		poid_rel = append(poid_rel, char)
	}

	return &Agent{Env: env, Id: id, Velocite: velocite, Acuite: acuite,
		Position: position, Opinion: opinion, Charisme: charisme, Relation: relation,
		PersonalParameter: personalParameter, Poid_rel: poid_rel,
		Vivant: true, TypeAgt: typeAgt, SyncChan: syncChan, Img: img, MoveTimer: 60, CurrentAction: "Praying", DialogTimer: 180}
}

func (ag *Agent) ID() IdAgent {
	return ag.Id
}

func (ag *Agent) AgtPosition() ut.Position {
	return ag.Position
}

func (ag *Agent) Start() {
	log.Printf("%s lancement...\n", ag.Id)

	go func() {
		env := ag.Env
		var step int
		for {
			step = <-ag.SyncChan
			ag.Percept(env)
			ag.Deliberate()
			ag.Act(env)

			ag.SyncChan <- step
		}
	}()
}

func CheckCollisionHorizontal(x, y float64, coliders []image.Rectangle) bool {
	for _, colider := range coliders {
		if colider.Overlaps(image.Rect(int(x), int(y), int(x)+16, int(y)+16)) {
			if x > 0 {
				return true
			} else if x < 0 {
				return true
			}
		}
	}
	return false
}

func CheckCollisionVertical(x, y float64, coliders []image.Rectangle) bool {
	for _, colider := range coliders {
		if colider.Overlaps(image.Rect(int(x), int(y), int(x)+16, int(y)+16)) {
			if y > 0 {
				return true
			} else if y < 0 {
				return true
			}
		}
	}
	return false
}

func (ag *Agent) Move() {
	ag.ClearAction()
	if ag.MoveTimer > 0 {
		
		ag.MoveTimer-=1
		if CheckCollisionHorizontal((ag.Position.X + ag.Position.Dx ),(ag.Position.Y + ag.Position.Dy),ag.Env.Carte.Coliders) || CheckCollisionVertical((ag.Position.X + ag.Position.Dx ),(ag.Position.Y + ag.Position.Dy),ag.Env.Carte.Coliders) {
			return
		}
		ag.Position.X += ag.Position.Dx
		ag.Position.Y += ag.Position.Dy
		return
		
	}
	randIdx := 0
	collision := true
	right := ut.UniqueDirection{Dx: ut.Maxspeed, Dy: 0}
	left := ut.UniqueDirection{Dx: -ut.Maxspeed, Dy: 0}
	down := ut.UniqueDirection{Dx: 0, Dy: ut.Maxspeed}
	up := ut.UniqueDirection{Dx: 0, Dy: -ut.Maxspeed}

	directions := []ut.UniqueDirection{right, left, down, up}

	for collision {

		randIdx = rand.Intn(len(directions))
		tryX := ag.Position.X + directions[randIdx].Dx
		tryY := ag.Position.Y + directions[randIdx].Dy
		if !CheckCollisionHorizontal(tryX,tryY,ag.Env.Carte.Coliders) && !CheckCollisionVertical(tryX,tryY,ag.Env.Carte.Coliders) {
			collision = false
		}
	}

	ag.Position.Dx = directions[randIdx].Dx
	ag.Position.Dy = directions[randIdx].Dy

	

	ag.MoveTimer=60
	
}

func (ag *Agent) Percept(env *Environnement) (nearbyAgents []IdAgent) {

	env.RLock()
	defer env.RUnlock()

	value, _ := env.AgentProximity.Load(ag.Id)
	nearby := value.([]IdAgent)

	return nearby
}

func (ag *Agent) Deliberate() {
	//TODO

	/*


	 */
}

func (ag *Agent) Act(env *Environnement) {
	//TODO
}

func (ag *Agent) SetAction(action string) {
    ag.CurrentAction = action
    ag.DialogTimer = 180 // 2 segundos a 60 FPS
}

func (ag *Agent) ClearAction() {
    ag.CurrentAction = ""
    ag.DialogTimer = 0
}

func (ag *Agent) Pray() {
    ag.SetAction("Praying")
}

func (ag *Agent) Eat() {
    ag.SetAction("Eating")
}