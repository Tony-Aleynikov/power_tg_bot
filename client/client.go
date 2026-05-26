package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
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

	payload := struct {
		URL                string   `json:"url"`
		AllowedUpdates     []string `json:"allowed_updates"`
		DropPendingUpdates bool     `json:"drop_pending_updates"`
	}{
		URL:                url,
		AllowedUpdates:     []string{"message", "callback_query"},
		DropPendingUpdates: os.Getenv("TG_BOT_FORCE_SET_WEBHOOK") == "1",
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

func (client client) SendMessage(chaiId int, text string) {
	c := client.client
	requestBody := struct {
		ChatID string `json:"chat_id"`
		Text   string `json:"text"`
	}{
		ChatID: strconv.Itoa(chaiId),
		Text:   text,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Errorf("ошибка создания JSON: %w", err)
		return
	}

	req, err := http.NewRequest(http.MethodPost, "https://"+ip+"/bot"+token+"/sendMessage", bytes.NewReader(jsonData))
	if err != nil {
		fmt.Errorf("ошибка создания запроса: %w", err)
		return
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
