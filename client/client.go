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
	headers := map[string]string{"Content-Type": "application/json"}
	body := makeBody(payload)
	request := makeRequest("/setWebhook", "POST", headers, body)
	response := doRequest(request, c)
	defer response.Body.Close()

	b, _ := io.ReadAll(response.Body)
	fmt.Println(string(b))
	fmt.Println("Webhook set successfully")
}

func (client client) GetMe(w http.ResponseWriter) {
	c := client.client
	headers := map[string]string{"Content-Type": "application/json"}
	request := makeRequest("/getMe", "GET", headers, nil)
	response := doRequest(request, c)
	defer response.Body.Close()

	b, _ := io.ReadAll(response.Body)
	fmt.Println(string(b))
	w.WriteHeader(response.StatusCode)
}

func (client client) SendMessage(chatId int, text string) {
	c := client.client
	requestBody := struct {
		ChatID string `json:"chat_id"`
		Text   string `json:"text"`
	}{
		ChatID: strconv.Itoa(chatId),
		Text:   text,
	}
	headers := map[string]string{"Content-Type": "application/json"}
	body := makeBody(requestBody)
	request := makeRequest("/sendMessage", "POST", headers, body)
	response := doRequest(request, c)
	defer response.Body.Close()

	b, _ := io.ReadAll(response.Body)
	fmt.Println(string(b))
}

func makeBody(body any) []byte {
	b, err := json.Marshal(body)
	if err != nil {
		fmt.Errorf("ошибка создания JSON: %w", err)
		return nil
	}
	return b
}

func makeRequest(url string, method string, headers map[string]string, body []byte) *http.Request {
	baseUrl := "https://" + ip + "/bot" + token
	req, err := http.NewRequest(method, baseUrl+url, bytes.NewReader(body))
	if err != nil {
		fmt.Errorf("ошибка создания запроса: %w", err)
		return nil
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	req.Host = "api.telegram.org"

	return req
}

func doRequest(request *http.Request, client *http.Client) *http.Response {
	resp, err := client.Do(request)
	if err != nil {
		fmt.Errorf("Ошибка в выполнении запроса: %w", err)
		return nil
	}
	return resp
}
