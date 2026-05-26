package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/power_tg_bot/client"
)

const webhookURL = "https://umore.serveousercontent.com"

func main() {
	http.HandleFunc("GET /handler", handler)
	http.HandleFunc("POST /", update)

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

func update(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	type update struct {
		Message struct {
			Chat struct {
				ID int `json:"id"`
			} `json:"chat"`
			Text string `json:"text"`
		} `json:"message"`
	}

	decoder := json.NewDecoder(r.Body)
	params := update{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		return
	}

	fmt.Println(params.Message.Text)
	fmt.Println(params.Message.Chat.ID)

	c := client.NewClient()
	c.SendMessage(params.Message.Chat.ID, params.Message.Text)
	fmt.Println("message was sended")
}
