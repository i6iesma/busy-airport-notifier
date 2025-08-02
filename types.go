package main

import "time"

// The relevant data we need from the airport api
type AirportData struct {
	ICAO        string
	Latitude    float64
	Longitude   float64
	ElevationFt int
}

// JSON format for the response obtained from the Airport API
type AirportResponse struct {
	ICAO        string  `json:"icao"`
	IATA        string  `json:"iata"`
	Name        string  `json:"name"`
	City        string  `json:"city"`
	Region      string  `json:"region"`
	Country     string  `json:"country"`
	ElevationFt int     `json:"elevation_ft"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Timezone    string  `json:"timezone"`
}

// Main top-level struct to hold all the data from the JSON
type VatsimData struct {
	General         GeneralInfo      `json:"general"`
	Pilots          []Pilot          `json:"pilots"`
	Controllers     []Controller     `json:"controllers"`
	Atis            []ATIS           `json:"atis"`
	Servers         []Server         `json:"servers"`
	Prefiles        []Prefile        `json:"prefiles"`
	Facilities      []Facility       `json:"facilities"`
	Ratings         []Rating         `json:"ratings"`
	PilotRatings    []PilotRating    `json:"pilot_ratings"`
	MilitaryRatings []MilitaryRating `json:"military_ratings"`
}

// Struct for the "general" object
type GeneralInfo struct {
	Version          int       `json:"version"`
	UpdateTimestamp  time.Time `json:"update_timestamp"`
	ConnectedClients int       `json:"connected_clients"`
	UniqueUsers      int       `json:"unique_users"`
}

// Struct for a single pilot entry
type Pilot struct {
	CID            int         `json:"cid"`
	Name           string      `json:"name"`
	Callsign       string      `json:"callsign"`
	Server         string      `json:"server"`
	PilotRating    int         `json:"pilot_rating"`
	MilitaryRating int         `json:"military_rating"`
	Latitude       float64     `json:"latitude"`
	Longitude      float64     `json:"longitude"`
	Altitude       int         `json:"altitude"`
	Groundspeed    int         `json:"groundspeed"`
	Transponder    string      `json:"transponder"`
	Heading        int         `json:"heading"`
	QnhIHG         float64     `json:"qnh_i_hg"`
	QnhMB          int         `json:"qnh_mb"`
	FlightPlan     *FlightPlan `json:"flight_plan"` // Using a pointer since the field can be null
	LogonTime      time.Time   `json:"logon_time"`
	LastUpdated    time.Time   `json:"last_updated"`
}

// Struct for a pilot's flight plan
type FlightPlan struct {
	FlightRules         string `json:"flight_rules"`
	Aircraft            string `json:"aircraft"`
	AircraftFAA         string `json:"aircraft_faa"`
	AircraftShort       string `json:"aircraft_short"`
	Departure           string `json:"departure"`
	Arrival             string `json:"arrival"`
	Alternate           string `json:"alternate"`
	Deptime             string `json:"deptime"`
	EnrouteTime         string `json:"enroute_time"`
	FuelTime            string `json:"fuel_time"`
	Remarks             string `json:"remarks"`
	Route               string `json:"route"`
	RevisionID          int    `json:"revision_id"`
	AssignedTransponder string `json:"assigned_transponder"`
}

// Struct for a single controller entry
type Controller struct {
	CID         int       `json:"cid"`
	Name        string    `json:"name"`
	Callsign    string    `json:"callsign"`
	Frequency   string    `json:"frequency"`
	Facility    int       `json:"facility"`
	Rating      int       `json:"rating"`
	Server      string    `json:"server"`
	VisualRange float64   `json:"visual_range"`
	TextATIS    []string  `json:"text_atis"`
	LastUpdated time.Time `json:"last_updated"`
	LogonTime   time.Time `json:"logon_time"`
}

// Struct for a single ATIS entry
type ATIS struct {
	CID         int       `json:"cid"`
	Name        string    `json:"name"`
	Callsign    string    `json:"callsign"`
	Frequency   string    `json:"frequency"`
	Facility    int       `json:"facility"`
	Rating      int       `json:"rating"`
	Server      string    `json:"server"`
	VisualRange int       `json:"visual_range"`
	AtisCode    string    `json:"atis_code"`
	TextATIS    []string  `json:"text_atis"`
	LastUpdated time.Time `json:"last_updated"`
	LogonTime   time.Time `json:"logon_time"`
}

// Struct for a single server entry
type Server struct {
	Ident                    string `json:"ident"`
	HostnameOrIP             string `json:"hostname_or_ip"`
	Location                 string `json:"location"`
	Name                     string `json:"name"`
	ClientConnectionsAllowed bool   `json:"client_connections_allowed"`
	IsSweatbox               bool   `json:"is_sweatbox"`
}

// Struct for a single prefile entry
type Prefile struct {
	CID         int        `json:"cid"`
	Name        string     `json:"name"`
	Callsign    string     `json:"callsign"`
	FlightPlan  FlightPlan `json:"flight_plan"`
	LastUpdated time.Time  `json:"last_updated"`
}

// Struct for a single facility entry
type Facility struct {
	ID       int    `json:"id"`
	Short    string `json:"short"`
	LongName string `json:"long_name"`
}

// Struct for a single rating entry
type Rating struct {
	ID        int    `json:"id"`
	ShortName string `json:"short_name"`
	LongName  string `json:"long_name"`
}

// Struct for a single pilot rating entry
type PilotRating struct {
	ID        int    `json:"id"`
	ShortName string `json:"short_name"`
	LongName  string `json:"long_name"`
}

// Struct for a single military rating entry
type MilitaryRating struct {
	ID        int    `json:"id"`
	ShortName string `json:"short_name"`
	LongName  string `json:"long_name"`
}
