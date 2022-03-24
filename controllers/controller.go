package controllers

import (
	"bufio"
	"github.com/nophramel/RZSON_Wasserabfluss/models"
	"github.com/nophramel/RZSON_Wasserabfluss/views"
	"log"
	"os"
	"strings"
	"time"
)

//Run Starts the application
func Run() {
	views.Clear()
	views.PrintHeader()
	stations, virtualStation := getLatestData()
	views.PrintStations(stations, virtualStation)
	views.PrintMenu()

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
		views.PrintHeader()
		views.Clear()
		views.PrintMenu()
		break
	case input == "q":
		// Terminate application with a 5 sec delay
		views.Clear()
		views.PrintGoodbye()
		time.Sleep(3 * time.Second)
		views.ShutDown()
		break
	}
}
