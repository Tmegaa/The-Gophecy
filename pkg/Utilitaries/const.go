package utils

type Position struct {
	X  float64
	Y  float64
	Dx float64
	Dy float64
}

type Rectangle struct {
	PositionDL Position
	PositionUR Position
}

func IsInRectangle(pos Position, area Rectangle) bool {
	// Calcule le centre d'un agent
	centerX := pos.X + 16/2
	centerY := pos.Y + 16/2

	return centerX >= area.PositionDL.X && centerX <= area.PositionUR.X &&
		centerY >= area.PositionUR.Y && centerY <= area.PositionDL.Y
}

const Maxspeed = 10.0

type UniqueDirection struct {
	Dx float64
	Dy float64
}

type Direction struct {
	direction []UniqueDirection
}
