package taiwanmap

// Coords is an xy coordinates thing
type Coords struct {
	X         int
	Y         int
	Intensity int
}

// NewCoords is a function that takes lat and long and returns x and y
func NewCoords(long float64, lat float64, intensity int) Coords {
	y := picHeight - int((long-baseLong)*degreeToLong) // - picheight to get y from top
	// y := int((long - baseLong) * degreeToLong)
	x := int((lat - baseLat) * degreeToLat)

	return Coords{
		X:         x,
		Y:         y,
		Intensity: intensity,
	}
}
