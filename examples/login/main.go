package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/LuciusMortified/steam"
	"github.com/joho/godotenv"
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

	timeTip, err := steam.GetTimeTip()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Time tip: %v\n", timeTip)

	timeDiff := time.Duration(timeTip.Time - time.Now().Unix())
	log.Printf("Time diff: %v\n", timeDiff)

	session := steam.NewSession(&http.Client{}, "", true)
	if err := session.Login(username, password, sharedSecret, timeDiff); err != nil {
		log.Fatal(err)
	}
	log.Print("Login successful")
}
