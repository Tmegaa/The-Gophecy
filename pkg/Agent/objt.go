package pkg

import (
	pos "Gophecy/pkg/Utilitaries"

	"github.com/hajimehoshi/ebiten/v2" // Ebiten
)

type InterfaceObjet interface {
	ObjPosition() pos.Position
	ID() IdObjet
	GetUse() bool
}

type Programm string

const (
	Go Programm = "Go"
	No Programm = "None"
)

type IdObjet string

type Computer struct {
	Env      *Environnement
	Id       IdObjet
	Position pos.Position
	Programm Programm
	Img      *ebiten.Image
	Used     bool
}

func (c *Computer) ObjPosition() pos.Position {
	return c.Position
}

func (c *Computer) GetUse() bool {
	return c.Used
}

func (c *Computer) ID() IdObjet {
	return c.Id
}

func NewComputer(env *Environnement, id IdObjet, pos pos.Position) *Computer {
	return &Computer{Env: env, Id: id, Position: pos, Programm: "None", Used: false}
}
