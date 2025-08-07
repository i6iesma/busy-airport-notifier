package main

import (
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

func GetBusyAirportData() map[string]int {
	const dataFileName = "airport_data.gob"
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	listOfAirportsToWatch := []string{"LEMD", "LEPA", "LEBL", "LEIB", "LEAL", "GCTS", "GCLP", "LEBB", "LEMG"}
	airportDataMap := make(map[string]AirportData)
	file, err := os.Open(dataFileName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// fmt.Println("file does not exist, creating it...")
			airportDataMap = refetchAllAirportData(listOfAirportsToWatch)
			saveAirportData(airportDataMap, dataFileName)
		} else {
			log.Fatal(err)
		}

	} else {
		defer file.Close()
		// fmt.Println("File found, decoding data...")
		decoder := gob.NewDecoder(file)
		err := decoder.Decode(&airportDataMap)
		if err != nil {
			log.Fatalf("error parsing the airport data info. %v", err)
		}
	}

	networkData := getVatsimNetworkData()
	departingPilotsData := make(map[string]int)
	for _, airport := range listOfAirportsToWatch {
		numberOfPilots := getNumberOfPilotsDepartingAirport(networkData.Pilots, airportDataMap[airport])
		departingPilotsData[airport] = numberOfPilots

	}
	return departingPilotsData

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
func getNumberOfPilotsDepartingAirport(pilots []Pilot, airportData AirportData) int {
	numberOfPilots := 0
	for _, pilot := range pilots {
		// Skip pilots arriving at this airport
		if pilot.FlightPlan != nil && pilot.FlightPlan.Arrival == airportData.ICAO {
			continue
		}
		// Count pilots departing from this airport and on ground
		if determineIfPilotIsOnGround(pilot, airportData) {
			if pilot.FlightPlan == nil || pilot.FlightPlan.Departure == airportData.ICAO {
				numberOfPilots++
			}
		}

	}
	return numberOfPilots
}
func determineIfPilotIsOnGround(pilot Pilot, airportData AirportData) bool {
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

func refetchAllAirportData(listOfAirportsToWatch []string) map[string]AirportData {
	airportDataMap := make(map[string]AirportData)
	var wg sync.WaitGroup
	for _, airportIcao := range listOfAirportsToWatch {
		wg.Add(1)
		go func(icao string) {
			defer wg.Done()
			airportData := getAirportData(icao)
			airportDataMap[icao] = airportData
		}(airportIcao)

	}
	wg.Wait()
	return airportDataMap
}

func saveAirportData(airportData map[string]AirportData, fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(airportData); err != nil {
		fmt.Printf("Error encoding map to Gob: %v\n", err)
		return err
	}
	return nil
}
