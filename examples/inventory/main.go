package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/LuciusMortified/steam"
	"github.com/joho/godotenv"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Fatal(err)
	}

	timeTip, err := steam.GetTimeTip()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Time tip: %#v\n", timeTip)
	timeDiff := time.Duration(timeTip.Time - time.Now().Unix())

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

	session := steam.NewSession(&http.Client{}, "", true)
	if err := session.Login(username, password, sharedSecret, timeDiff); err != nil {
		log.Fatal(err)
	}
	log.Print("Login successful")

	sid := steam.SteamID(partnerID)
	apps, err := session.GetInventoryAppStats(sid)
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range apps {
		log.Printf("-- AppID total asset count: %d\n", v.AssetCount)
		for _, context := range v.Contexts {
			log.Printf("-- Items on %d %d (count %d)\n", v.AppID, context.ID, context.AssetCount)
			inven, err := session.GetInventory(sid, v.AppID, context.ID, true)
			if err != nil {
				log.Fatal(err)
			}

			for _, item := range inven {
				log.Printf("Item: %s = %d\n", item.Desc.MarketHashName, item.AssetID)
			}
		}
	}

	log.Println("Bye!")
}
