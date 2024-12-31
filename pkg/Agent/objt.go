package pkg

import (
	pos "Gophecy/pkg/Utilitaries"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

// InterfaceObjet définit un comportement commun pour tous les objets
type InterfaceObjet interface {
	ObjPosition() pos.Position
	ID() IdObjet
	GetUse() bool
	GetProgramm() Programm
	GetType() TypeObjet
	GetImg() *ebiten.Image
}

// On définit les différents langages de programmation possibles: on pourra élargir par la suite
type Programm string

const (
	GoPgm Programm = "Go"
	NoPgm Programm = "None"
)

type TypeObjet string

// Différents types d'objets
const (
	ComputerType TypeObjet = "Computer"
	StatueType   TypeObjet = "Statue"
)

// Chaque objet a un identifiant
type IdObjet string

// Objet est la structure de base de tous les objets
type Objet struct {
	Env      *Environnement
	Id       IdObjet
	Position pos.Position
	Programm Programm
	Img      *ebiten.Image
	Used     bool
	Type     TypeObjet
}

// Fonction qui renvoie la position d'un objet
func (o *Objet) ObjPosition() pos.Position { return o.Position }

// Fonction qui renvoie si un objet est en train d'être utilisé
func (o *Objet) GetUse() bool { return o.Used }

// Fonction qui renvoie l'ID d'un objet
func (o *Objet) ID() IdObjet { return o.Id }

// Fonction qui renvoie le langage de programmation d'un objet
func (o *Objet) GetProgramm() Programm { return o.Programm }

// Fonction qui renvoie le type d'un objet
func (o *Objet) GetType() TypeObjet { return o.Type }

// Fonction qui renvoie l'image d'un objet
func (o *Objet) GetImg() *ebiten.Image { return o.Img }

// L'ordinateur est un type spécifique d'objet
type Computer struct {
	Objet
	mutex sync.Mutex // Protège l'accès à l'ordinateur
}

// Création d'un nouvel ordinateur
func NewComputer(env *Environnement, id IdObjet, pos pos.Position) *Computer {
	return &Computer{
		Objet: Objet{
			Env:      env,
			Id:       id,
			Position: pos,
			Programm: NoPgm,
			Used:     false,
			Type:     ComputerType,
		},
	}
}

// Ajouter des méthodes thread-safe pour utiliser l'ordinateur
func (c *Computer) TryUse() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.Used {
		return false
	}
	c.Used = true
	return true
}

// Fonction qui libère le mutex d'un ordinateur
func (c *Computer) Release() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.Used = false
}

// Fonction qui met à jour le langage de programmation d'un ordinateur
func (c *Computer) SetProgramm(p Programm) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.Programm = p
}

// Le statut est un autre type spécifique d'objet
type Statue struct {
	Objet
	// Champs spécifiques à la statue, le cas échéant
}

// Création d'une nouvelle statue
func NewStatue(env *Environnement, id IdObjet, pos pos.Position) *Statue {
	return &Statue{
		Objet: Objet{
			Env:      env,
			Id:       id,
			Position: pos,
			Programm: NoPgm,
			Used:     false,
			Type:     StatueType,
		},
	}
}
