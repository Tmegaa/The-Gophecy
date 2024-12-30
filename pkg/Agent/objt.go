package pkg

import (
	pos "Gophecy/pkg/Utilitaries"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

// InterfaceObjet define o comportamento comum para todos os objetos
type InterfaceObjet interface {
    ObjPosition() pos.Position
    ID() IdObjet
    GetUse() bool
    GetProgramm() Programm
    GetType() TypeObjet
	GetImg() *ebiten.Image
	
}

type Programm string

const (
    Go Programm = "Go"
    No Programm = "None"
)

type TypeObjet string

const (
    ComputerType TypeObjet = "Computer"
    StatueType   TypeObjet = "Statue"
)

type IdObjet string

// Objet é a estrutura base para todos os objetos
type Objet struct {
    Env      *Environnement
    Id       IdObjet
    Position pos.Position
    Programm Programm
    Img      *ebiten.Image
    Used     bool
    Type     TypeObjet
}

func (o *Objet) ObjPosition() pos.Position { return o.Position }
func (o *Objet) GetUse() bool              { return o.Used }
func (o *Objet) ID() IdObjet               { return o.Id }
func (o *Objet) GetProgramm() Programm     { return o.Programm }
func (o *Objet) GetType() TypeObjet        { return o.Type }
func (o *Objet) GetImg() *ebiten.Image     { return o.Img }


// Computer é um tipo específico de Objet
type Computer struct {
    Objet
    mutex sync.Mutex // Protege o acesso ao computador
}

func NewComputer(env *Environnement, id IdObjet, pos pos.Position) *Computer {
    return &Computer{
        Objet: Objet{
            Env:      env,
            Id:       id,
            Position: pos,
            Programm: No,
            Used:     false,
            Type:     ComputerType,
        },
    }
}

// Adicione métodos thread-safe para usar o computador
func (c *Computer) TryUse() bool {
    c.mutex.Lock()
    defer c.mutex.Unlock()
    
    if c.Used {
        return false
    }
    c.Used = true
    return true
}

func (c *Computer) Release() {
    c.mutex.Lock()
    defer c.mutex.Unlock()
    c.Used = false
}

func (c *Computer) SetProgramm(p Programm) {
    c.mutex.Lock()
    defer c.mutex.Unlock()
    c.Programm = p
}

// Statue é outro tipo específico de Objet
type Statue struct {
    Objet
    // Campos específicos da Statue, se houver
}

func NewStatue(env *Environnement, id IdObjet, pos pos.Position) *Statue {
    return &Statue{
        Objet: Objet{
            Env:      env,
            Id:       id,
            Position: pos,
            Programm: No,
            Used:     false,
            Type:     StatueType,
        },
    }
}
