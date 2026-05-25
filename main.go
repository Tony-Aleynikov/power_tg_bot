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

func main() {
	http.HandleFunc("GET /handler", handler)
	http.HandleFunc("GET /getme", getMeHandler)
	http.HandleFunc("GET /setWebhook", getWebhook)

	s := &http.Server{
		Addr: ":8080",
	}

	fmt.Println("Listening on port 8080")
	log.Fatal(s.ListenAndServe()) // почему именно log.Fatal
}

func handler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Err")
		return
	}
	defer r.Body.Close()
	fmt.Println(string(body))
}

func getMeHandler(w http.ResponseWriter, r *http.Request) {
	token := "7121104577:AAG9TuCKKVJvCRVp7EzQ8rl-wBIgNwFLUzk"
	ip := "149.154.167.220"

	c := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				ServerName: "api.telegram.org",
			},
		},
	}

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
}

func getWebhook(w http.ResponseWriter, r *http.Request) {
	token := "7121104577:AAG9TuCKKVJvCRVp7EzQ8rl-wBIgNwFLUzk"
	ip := "149.154.167.220"
	webhookURL := "https://umore.serveousercontent.com"

	c := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				ServerName: "api.telegram.org",
			},
		},
	}

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
}
