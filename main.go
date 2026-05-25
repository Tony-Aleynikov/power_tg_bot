package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	token      = "7121104577:AAG9TuCKKVJvCRVp7EzQ8rl-wBIgNwFLUzk"
	ip         = "149.154.167.220"
	webhookURL = "https://umore.serveousercontent.com"
)

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
	c := getClient()

	req, err := http.NewRequest(http.MethodGet, "https://"+ip+"/bot"+token+"/getMe", nil)
	if err != nil {
		panic(err)
	}
	req.Host = "api.telegram.org"

	resp, err := c.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	fmt.Println(string(b))
	w.WriteHeader(resp.StatusCode)
}

func setWebhook() {
	c := getClient()

	payload := map[string]string{
		"url": webhookURL,
	}
	bodyJSON, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(http.MethodPost, "https://"+ip+"/bot"+token+"/setWebhook", bytes.NewReader(bodyJSON))
	if err != nil {
		panic(err)
	}
	req.Host = "api.telegram.org"
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	fmt.Println(string(b))
	fmt.Println("Webhook set successfully")
}

func getClient() *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				ServerName: "api.telegram.org",
			},
		},
	}
}
