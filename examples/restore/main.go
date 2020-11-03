package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/LuciusMortified/steam"
	"github.com/joho/godotenv"
)

func checkSession(session *steam.Session) {
	apiKey, err := session.GetWebAPIKey()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Api Key: %s", apiKey)

	resp, err := session.GetTradeOffers(
		steam.TradeFilterSentOffers|steam.TradeFilterRecvOffers,
		time.Now(),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Received offers count: %d", len(resp.ReceivedOffers))
	log.Printf("Sent offers count: %d", len(resp.SentOffers))
}

func NewSession(username, password, sharedSecret string) *steam.Session {
	timeTip, err := steam.GetTimeTip()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Time tip: %#v\n", timeTip)

	timeDiff := time.Duration(timeTip.Time - time.Now().Unix())
	log.Printf("Time diff: %v\n", timeDiff)

	session := steam.NewSession(&http.Client{}, "", true)
	if err := session.Login(username, password, sharedSecret, timeDiff); err != nil {
		log.Fatal(err)
	}
	log.Print("Login successful")

	return session
}

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

	//Initialize new session
	session := NewSession(username, password, sharedSecret)

	checkSession(session)

	//Dump session
	dump, err := session.Dump()
	if err != nil {
		log.Fatal(err)
	}

	dumpBytes, err := json.Marshal(dump)
	if err != nil {
		log.Fatal(err)
	}

	data := &steam.SessionData{}
	err = json.Unmarshal(dumpBytes, data)
	if err != nil {
		log.Fatal(err)
	}

	//Restore session
	session, err = steam.RestoreSession(&http.Client{}, data, true)
	if err != nil {
		log.Fatal(err)
	}

	checkSession(session)
}
