package controllers

import (
	"RZSON_Wasserstaende/models"
	"RZSON_Wasserstaende/views"
	"bufio"
	"log"
	"os"
	"strings"
)

func Run() {
	views.Clear()
	views.PrintMenu()
	stations, virtualStation := getLatestData()
	views.PrintStations(stations, virtualStation)

	for true {
		executeCommand()
	}
}

func getLatestData() ([]models.Messstation, models.Messstation) {
	virtualStation := models.SetupVirtualStation()
	stations := models.SetupStations()

	for index, currentStation := range stations {
		data, _ := models.ReadCSVFromUrl(currentStation.Url)
		models.CalculateData(data, &stations[index])
		currentStation = stations[index]
		if currentStation.ID != 2029 {
			virtualStation.Measurements.LastReading += currentStation.Measurements.LastReading
		}

	}
	return stations, virtualStation
}

func executeCommand() {
	command := askForInput()
	parseCommand(command)
}

func askForInput() string {
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	response = strings.TrimSpace(response)
	return response
}

func parseCommand(input string) {
	switch {
	case input == "r":
		// Prints Station readings
		stations, virtualStation := getLatestData()
		views.PrintStations(stations, virtualStation)
		break

	case input == "c":
		// Clear view and print menu
		//
		views.Clear()
		views.PrintMenu()
		break
	case input == "q":
		// Terminate application
		//
		views.Clear()
		views.PrintGoodbye()
		views.ShutDown()
		break
	}
}
