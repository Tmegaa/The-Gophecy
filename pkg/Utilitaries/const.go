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
	// Calcula o centro do agente
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
