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
	return pos.X >= area.PositionDL.X && pos.X <= area.PositionUR.X && pos.Y >= area.PositionDL.Y && pos.Y <= area.PositionUR.Y
}

const Maxspeed = 2.0

type UniqueDirection struct {
	Dx float64
	Dy float64
}

type Direction struct {
	direction []UniqueDirection
}
