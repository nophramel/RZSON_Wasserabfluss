package views

import (
	"fmt"
	"github.com/nophramel/RZSON_Wasserstaende/models"
	"os"
	"os/exec"
)

// PrintMenu prints the menu to the console
func PrintMenu() {
	fmt.Println(`
##############################################################################################
#**************************************** Willkommen *****************************************
#***************************** Wählen Sie Ihre Option unten aus ******************************
#
# r.	Aktualisiere die Messwerte.
#
# c.	Leere die Ausgabe und zeige die Übersicht an
	# q.	App beenden
	`)
}

func PrintStations(stations []models.Messstation, virtualStation models.Messstation) {

	fmt.Println("--------------------------------------------------------------------------------------------------")
	fmt.Printf("%7s | %8s | %21s | %15s | %15s | %15s\n", "ID", "Gewässer", "Station", "Abflussmenge", "Messzeitpunkt", "Max 24h")
	fmt.Println("--------------------------------------------------------------------------------------------------")

	for _, currentStation := range stations {

		fmt.Printf("%7v | %8s | %21s | %10v m³/s | %s |  %9v m³/s\n", currentStation.ID, currentStation.River, currentStation.Station, currentStation.Measurements.LastReading, currentStation.Measurements.LastTimeStamp.Format("02.01.06, 15:04"), currentStation.Measurements.MaxReadingLast24h)

	}
	fmt.Println("--------------------------------------------------------------------------------------------------")
	fmt.Printf("%8s  %8s | %21s | %10.f m³/s | %21s\n", "Total:", virtualStation.River, virtualStation.Station, virtualStation.Measurements.LastReading, "Abflussmenge in ca. 1h zu erwarten")
	fmt.Println("--------------------------------------------------------------------------------------------------")

}

// Clear clears the console view
func Clear() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}

// PrintGoodbye prints a goodbye message to the console
func PrintGoodbye() {
	fmt.Println("Auf Wiedersehen!")
}

// ShutDown terminates the application
func ShutDown() {
	os.Exit(0)
}
