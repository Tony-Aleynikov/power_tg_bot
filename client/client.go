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

func (client client) SetWebhook(url string) error {
	bodyData := struct {
		URL                string   `json:"url"`
		AllowedUpdates     []string `json:"allowed_updates"`
		DropPendingUpdates bool     `json:"drop_pending_updates"`
	}{
		URL:                url,
		AllowedUpdates:     []string{"message", "callback_query"},
		DropPendingUpdates: os.Getenv("TG_BOT_FORCE_SET_WEBHOOK") == "1",
	}

	_, err := client.do("/setWebhook", bodyData)
	if err != nil {
		return err
	}
	fmt.Println("Webhook set successfully")
	return nil
}

func (client client) GetMe(w http.ResponseWriter) error {
	response, err := client.do("/getMe", nil)
	if err != nil {
		return err
	}
	w.WriteHeader(response.StatusCode)
	return nil
}

func (client client) SendMessage(chatId int, text string) error {
	bodyData := struct {
		ChatID string `json:"chat_id"`
		Text   string `json:"text"`
	}{
		ChatID: strconv.Itoa(chatId),
		Text:   text,
	}

	_, err := client.do("/sendMessage", bodyData)
	return err
}

func (c client) do(method string, bodyData any) (*http.Response, error) {
	headers := map[string]string{"Content-Type": "application/json"}
	body, err := makeBody(bodyData)
	if err != nil {
		return nil, err
	}
	request, err := makeRequest(method, "POST", headers, body)
	if err != nil {
		return nil, err
	}
	response, err := doRequest(request, c.client)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	b, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(b))
	return response, nil
}

func makeBody(body any) ([]byte, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания JSON: %w", err)
	}
	return b, nil
}

func makeRequest(url string, method string, headers map[string]string, body []byte) (*http.Request, error) {
	baseUrl := "https://" + ip + "/bot" + token
	req, err := http.NewRequest(method, baseUrl+url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %w", err)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	req.Host = "api.telegram.org"

	return req, nil
}

func doRequest(request *http.Request, client *http.Client) (*http.Response, error) {
	resp, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("Ошибка в выполнении запроса: %w", err)
	}
	return resp, nil
}
