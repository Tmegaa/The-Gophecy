package entities

import (
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
)

type Sprite struct {
	Img  *ebiten.Image
	X, Y, Dx, Dy float64
    MoveTimer int
}

func (sprite *Sprite) RandomMovement(maxSpeed float64) {
    // Atualize as velocidades com um fator aleatório
    sprite.Dx += rand.Float64()*2 - 1 // Aleatório no intervalo [-1, 1]
    sprite.Dy += rand.Float64()*2 - 1

    // Limite a velocidade no eixo X
    if sprite.Dx > maxSpeed {
        sprite.Dx = maxSpeed
    } else if sprite.Dx < -maxSpeed {
        sprite.Dx = -maxSpeed
    }

    // Limite a velocidade no eixo Y
    if sprite.Dy > maxSpeed {
        sprite.Dy = maxSpeed
    } else if sprite.Dy < -maxSpeed {
        sprite.Dy = -maxSpeed
    }

    // Atualize a posição com base nas velocidades
    sprite.X += sprite.Dx
    sprite.Y += sprite.Dy
}