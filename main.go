package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	listOfAirportsToWatch := []string{"LEPA", "LEBL", "LEIB", "LEAL", "GCTS", "GCLP", "LEBB", "LEMG"}
	networkData := getVatsimNetworkData()
	for _, airport := range listOfAirportsToWatch {
		airportData := getAirportData(airport)
		numberOfPilots := getNumberOfPilotsDepartingAirport(&networkData.Pilots, &airportData)
		fmt.Println()
		fmt.Printf("%v pilots are departing %v at the moment", numberOfPilots, airportData.ICAO)

	}

}

func getVatsimNetworkData() VatsimData {
	const networkStatusUrl string = "https://data.vatsim.net/v3/vatsim-data.json"
	method := "GET"
	client := &http.Client{}

	req, err := http.NewRequest(method, networkStatusUrl, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var vatsimData VatsimData
	json.Unmarshal(body, &vatsimData)
	return vatsimData
}

func getAirportData(airport string) AirportData {
	client := &http.Client{}
	airportApiUrl := "https://api.api-ninjas.com/v1/airports?icao=" + airport
	req, err := http.NewRequest("GET", airportApiUrl, nil)
	if err != nil {
		log.Fatal(err)
	}
	airportApiKey := os.Getenv("AIRPORT_API_KEY")
	req.Header.Add("X-Api-Key", airportApiKey)

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	var airportResponseData []AirportResponse
	if err := json.Unmarshal(body, &airportResponseData); err != nil || len(airportResponseData) == 0 {
		log.Fatal("Failed to parse airport response or no data found")
	}
	return AirportData{
		ICAO:        airportResponseData[0].ICAO,
		Latitude:    airportResponseData[0].Latitude,
		Longitude:   airportResponseData[0].Longitude,
		ElevationFt: airportResponseData[0].ElevationFt,
	}

}
func getNumberOfPilotsDepartingAirport(pilots *[]Pilot, airportData *AirportData) int {
	numberOfPilots := 0
	for _, pilot := range *pilots {
		if pilot.FlightPlan == nil {
			continue
		}
		if pilot.FlightPlan.Departure == airportData.ICAO && determineIfPilotIsOnGround(&pilot, airportData) {
			numberOfPilots++
			// fmt.Println("Number of pilots increased")
			// fmt.Println()
			// fmt.Printf("%v is at coordinates %.8f, %.8f, with altitude %v", pilot.Callsign, pilot.Latitude, pilot.Longitude, pilot.Altitude)
		}

	}
	return numberOfPilots
}
func determineIfPilotIsOnGround(pilot *Pilot, airportData *AirportData) bool {
	const elevationMarginFt = 200
	const maxGroundSpeedKts = 50
	const maxRadiusFromAirportCenterPointKm = 6 // This is waaay to much for most airports, but since the limiting factor will be the speed and the aircraft altitude, this is basically just to check if the aircraft is not at other airport
	distanceBetweenAircraftAndAirportKm := haversineDistanceKm(pilot.Latitude, pilot.Longitude, airportData.Latitude, airportData.Longitude)
	if distanceBetweenAircraftAndAirportKm > maxRadiusFromAirportCenterPointKm {
		return false
	}
	if math.Abs(float64(pilot.Altitude-airportData.ElevationFt)) > float64(elevationMarginFt) {
		// Aircraft is outside the elevation margin
		return false
	}
	if pilot.Groundspeed > maxGroundSpeedKts {
		// On takeoff presumibly
		return false
	}
	return true

}

// haversineDistanceKm calculates the great-circle distance between two points on the Earth
// (specified in decimal degrees) using the Haversine formula.
// Returns distance in kilometers.
func haversineDistanceKm(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371000.0 // Radius of Earth in meters (approximate mean radius)

	// Convert degrees to radians
	phi1 := lat1 * math.Pi / 180
	phi2 := lat2 * math.Pi / 180
	deltaPhi := (lat2 - lat1) * math.Pi / 180
	deltaLambda := (lon2 - lon1) * math.Pi / 180

	// Haversine formula
	a := math.Sin(deltaPhi/2)*math.Sin(deltaPhi/2) +
		math.Cos(phi1)*math.Cos(phi2)*
			math.Sin(deltaLambda/2)*math.Sin(deltaLambda/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := R * c        // Distance in meters
	return distance / 1000.0 // Convert to kilometers
}
