package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	token = "7121104577:AAG9TuCKKVJvCRVp7EzQ8rl-wBIgNwFLUzk"
	ip    = "149.154.167.220"
)

type client struct {
	client *http.Client
}

func NewClient() client {
	return client{
		client: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					ServerName: "api.telegram.org",
				},
			},
		},
	}
}

func (client client) SetWebhook(url string) {
	c := client.client
	payload := map[string]string{
		"url": url,
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

func (client client) GetMe(w http.ResponseWriter) {
	c := client.client
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
