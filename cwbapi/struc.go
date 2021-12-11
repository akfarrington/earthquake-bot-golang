package cwbapi

// Response is the top level of json response from server
type Response struct {
	Records Records `json:"records"`
}

// Records holds the earthquake info in an array
type Records struct {
	Earthquake []Earthquake `json:"earthquake"`
}

// Earthquake holds the information I want to collect
type Earthquake struct {
	ReportContent  string         `json:"reportContent"`
	ReportImageURI string         `json:"reportImageURI"`
	Web            string         `json:"web"`
	EarthquakeInfo EarthquakeInfo `json:"earthquakeInfo"`
	Intensity      Intensity      `json:"intensity"`
}

// EarthquakeInfo includes some extra info about the eq, such as
// epicenter, lat, long, etc.
type EarthquakeInfo struct {
	OriginTime string    `json:"originTime"`
	Epicenter  Epicenter `json:"epiCenter"`
}

// Intensity is a field with eq intensity
type Intensity struct {
	ShakingArea []ShakingArea `json:"shakingArea"`
}

// ShakingArea stores the station info
type ShakingArea struct {
	EqStation []EqStation `json:"eqStation"`
}

// EqStation stores station info
type EqStation struct {
	StationIntensity StationIntensity `json:"stationIntensity"`
	StationLat       LongLat          `json:"stationLat"`
	StationLong      LongLat          `json:"stationLon"`
}

// StationIntensity is the intensity of the eq at the point
type StationIntensity struct {
	Value int `json:"value"`
}

// LongLat is the location of the station
type LongLat struct {
	Value float64 `json:"value"`
}

// Epicenter location of the epicenter
type Epicenter struct {
	Long LongLat `json:"epiCenterLat"`
	Lat  LongLat `json:"epiCenterLon"`
}
