package pkg

import (
	ut "Gophecy/pkg/Utilitaries"
)

type InterfaceObjet interface {
	AgtPosition() ut.Position
	ID() IdAgent
}
