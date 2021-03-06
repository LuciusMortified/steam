package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"

	"github.com/LuciusMortified/steam"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Fatal(err)
	}

	username := os.Getenv("USERNAME")
	if username == "" {
		log.Fatal(errors.New("specify USERNAME env"))
	}

	password := os.Getenv("PASSWORD")
	if password == "" {
		log.Fatal(errors.New("specify PASSWORD env"))
	}

	sharedSecret := os.Getenv("SHARED_SECRET")
	if sharedSecret == "" {
		log.Fatal(errors.New("specify SHARED_SECRET env"))
	}

	partnerIDStr := os.Getenv("PARTNER_ID")
	if partnerIDStr == "" {
		log.Fatal(errors.New("specify PARTNER_ID env"))
	}

	partnerID, err := strconv.ParseUint(partnerIDStr, 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	timeTip, err := steam.GetTimeTip()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Time tip: %#v\n", timeTip)

	timeDiff := time.Duration(timeTip.Time - time.Now().Unix())
	session := steam.NewSession(&http.Client{}, "", true)
	if err := session.Login(username, password, sharedSecret, timeDiff); err != nil {
		log.Fatal(err)
	}
	log.Print("Login successful")

	err = session.RevokeWebAPIKey()
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Revoked API Key")

	key, err := session.RegisterWebAPIKey("test.org")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Registered new API Key: %s", key)

	ownedGames, err := session.GetOwnedGames(steam.SteamID(partnerID), false, true)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Games count: %d\n", ownedGames.Count)
	for _, game := range ownedGames.Games {
		log.Printf("Game: %d 2 weeks play time: %d\n", game.AppID, game.Playtime2Weeks)
	}
}
