package agent

import (
	ut "Gophecy/pkg/Utilitaries"

	"github.com/hajimehoshi/ebiten/v2" // Ebiten
)

type InterfaceObjet interface {
	ObjPosition() ut.Position
	ID() IdObjet
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
	Position ut.Position
	Programm Programm
	Img      *ebiten.Image
	Used     bool
}

func (c *Computer) ObjPosition() ut.Position {
	return c.Position
}

func (c *Computer) ID() IdObjet {
	return c.Id
}

func NewComputer(env *Environnement, id IdObjet, pos ut.Position) *Computer {
	return &Computer{Env: env, Id: id, Position: pos, Programm: "None", Used: false}
}
