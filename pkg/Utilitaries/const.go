package utils

import "math"

type Position struct {
	X  float64
	Y  float64
	Dx float64
	Dy float64
}

const (
	AssetsPath             = "assets/images/"
	AgentBelieverImageFile = "ninja.png"
	AgentScepticImageFile  = "sceptic.png"
	AgentNeutralImageFile  = "neutre.png"
)

type Pair struct {
	First  float64
	Second float64
}

type Rectangle struct {
	PositionDL Position
	PositionUR Position
}

// Fonction qui vérifie si une position est bien dans un rectangle donné
func IsInRectangle(pos Position, area Rectangle) bool {
	// Calcule le centre d'un agent
	centerX := pos.X + 16/2
	centerY := pos.Y + 16/2

	return centerX >= area.PositionDL.X && centerX <= area.PositionUR.X &&
		centerY >= area.PositionUR.Y && centerY <= area.PositionDL.Y
}

const Maxspeed = 2.0

type UniqueDirection struct {
	Dx float64
	Dy float64
}

type Direction struct {
	direction []UniqueDirection
}

func Distance(p1, p2 Position) float64 {
	return math.Sqrt(math.Pow(p1.X-p2.X, 2) + math.Pow(p1.Y-p2.Y, 2))
}
