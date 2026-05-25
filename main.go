package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/power_tg_bot/client"
)

const webhookURL = "https://umore.serveousercontent.com"

func main() {
	http.HandleFunc("GET /handler", handler)

	s := &http.Server{
		Addr: ":8080",
	}
	setWebhook()
	fmt.Println("Listening on port 8080")
	log.Fatal(s.ListenAndServe()) // почему именно log.Fatal
}

func handler(w http.ResponseWriter, r *http.Request) {
	c := client.NewClient()
	c.GetMe(w)
}

func setWebhook() {
	c := client.NewClient()
	c.SetWebhook(webhookURL)
}
