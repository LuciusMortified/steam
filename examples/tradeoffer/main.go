package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/LuciusMortified/steam"
)

func processOffer(session *steam.Session, offer *steam.TradeOffer) {
	var sid steam.SteamID
	sid.ParseDefaults(offer.Partner)

	log.Printf("Offer id: %d, Receipt ID: %d, State: %d", offer.ID, offer.ReceiptID, offer.State)
	log.Printf("Offer partner SteamID 64: %d", uint64(sid))
	if offer.State == steam.TradeStateAccepted {
		items, err := session.GetTradeReceivedItems(offer.ReceiptID)
		if err != nil {
			log.Printf("error getting items: %v", err)
		} else {
			for _, item := range items {
				log.Printf("Item: %#v", item)
			}
		}
	}
	if offer.State == steam.TradeStateActive && !offer.IsOurOffer {
		err := offer.Accept(session)
		if err != nil {
			log.Printf("error accept trade: %v", err)
		}
	}
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

	key, err := session.GetWebAPIKey()
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Key: ", key)

	resp, err := session.GetTradeOffers(
		steam.TradeFilterSentOffers|steam.TradeFilterRecvOffers,
		time.Now(),
	)
	if err != nil {
		log.Fatal(err)
	}

	for _, offer := range resp.SentOffers {
		processOffer(session, offer)
	}
	for _, offer := range resp.ReceivedOffers {
		processOffer(session, offer)
	}

	log.Println("Bye!")
}
