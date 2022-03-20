package models

import (
	"encoding/csv"
	"net/http"
	"strconv"
	"time"
)

type Messstation struct {
	ID           int
	Url          string
	Station      string
	River        string
	Measurements MeasurementData
	GsSteps      string
}

type MeasurementData struct {
	Readings                   []float64
	TimeStamps                 []string
	LastReading                float64
	LastTimeStamp              time.Time
	MaxReadingLast24h          float64
	MaxReadingTimeStampLast24h time.Time
}

// SetupVirtualStation Setting up the virtual Station RZSON to calculate estimated reading for Niedergösgen
func SetupVirtualStation() Messstation {
	var station01 Messstation
	station01.Station = "RZSON Niedergösgen"
	station01.River = "Aare"
	return station01
}

// SetupStations Hardcoded Station details
func SetupStations() []Messstation {
	var stations []Messstation
	var station01 Messstation
	station01.ID = 2029
	station01.Url = "https://www.hydrodaten.admin.ch/lhg/az/dwh/csv/BAFU_2029_AbflussAFFRA.csv"
	station01.Station = "Brügg, Aegerten"
	station01.River = "Aare"
	var station02 Messstation
	station02.ID = 2155
	station02.Url = "https://www.hydrodaten.admin.ch/lhg/az/dwh/csv/BAFU_2155_AbflussPneumatik.csv"
	station02.Station = "Wiler, Limpachmündung"
	station02.River = "Emme"
	var station03 Messstation
	station03.ID = 2063
	station03.Url = "https://www.hydrodaten.admin.ch/lhg/az/dwh/csv/BAFU_2063_AbflussPneumatikoben.csv"
	station03.Station = "Murgenthal"
	station03.River = "Aare"
	var station04 Messstation
	station04.ID = 2434
	station04.Url = "https://www.hydrodaten.admin.ch/lhg/az/dwh/csv/BAFU_2434_AbflussPneumatik.csv"
	station04.Station = "Olten, Hammermühle"
	station04.River = "Dünnern"
	var station05 Messstation
	station05.ID = 2450
	station05.Url = "https://www.hydrodaten.admin.ch/lhg/az/dwh/csv/BAFU_2450_AbflussPneumatik.csv"
	station05.Station = "Zofingen"
	station05.River = "Wigger"
	stations = append(stations, station01, station02, station03, station04, station05)
	return stations
}

// ReadCSVFromUrl Get .csv from URL
func ReadCSVFromUrl(url string) ([][]string, error) {
	table, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer table.Body.Close()
	reader := csv.NewReader(table.Body)
	reader.Comma = ','
	data, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Parses .csv Table into 2d Slice (Time;value)
func createMeasurementData(data [][]string) MeasurementData {
	var measurementData MeasurementData
	for i, line := range data {
		if i < 1 { // omit header line
			continue
		}
		for j, field := range line {
			if j == 0 {
				measurementData.TimeStamps = append(measurementData.TimeStamps, field)
			} else if j == 1 {
				value, _ := strconv.ParseFloat(field, 64)
				measurementData.Readings = append(measurementData.Readings, value)
			}
		}
	}
	return measurementData
}

// CalculateData Determines the max reading with timestamp for last 24h
func CalculateData(data [][]string, messstation *Messstation) {
	if len(data) == 0 {
		return
	}
	measurementData := createMeasurementData(data)
	messstation.Measurements = measurementData
	columnCount := len(data) - 1
	lastReading := data[columnCount][1]
	lastValue, _ := strconv.ParseFloat(lastReading, 64)
	messstation.Measurements.LastReading = lastValue
	timeValue, _ := time.Parse(time.RFC3339, data[columnCount][0])
	messstation.Measurements.LastTimeStamp = timeValue
	// Determines the first timestamp from 24h ago (one reading every 5min (24*60)/5 = 288 (-1 to skip header)
	last24h := columnCount - 287
	// Determines the highest reading for the past 24h
	maxValue := 0.0
	maxTimestamp := ""
	for i, number := range measurementData.Readings[last24h:] {
		if number > maxValue {
			maxValue = number
			maxTimestamp = measurementData.TimeStamps[i]
		}
	}
	messstation.Measurements.MaxReadingLast24h = maxValue
	messstation.Measurements.MaxReadingTimeStampLast24h, _ = time.Parse(time.RFC3339, maxTimestamp)

}
