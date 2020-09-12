package main

import (
	"log"

	"github.com/LuciusMortified/steam"
)

func main() {
	version, err := steam.NewSessionWithAPIKey("").GetRequiredSteamAppVersion(730)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Current CS:GO version is: ", version)
}
