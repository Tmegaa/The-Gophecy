package utils

import "math"

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
	// Calculate the center of the agent
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
