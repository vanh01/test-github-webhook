package tele

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const BaseUrl = "https://api.telegram.org/bot%s/%s?chat_id=%d&text=%s"

const (
	sendMessage = "sendMessage"
)

type TelegramClient struct {
	ApiKey string
	ChatId int64
}

func (t *TelegramClient) SendMessage(text string) (string, error) {
	genUrl := fmt.Sprintf(BaseUrl, t.ApiKey, sendMessage, t.ChatId, url.QueryEscape(text))
	res, err := SendHttpRequest(http.MethodGet, genUrl, nil)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func SendHttpRequest(method string, url string, payload interface{}) (*http.Response, error) {
	ctx := context.Background()
	body := new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(payload)
	if err != nil {
		return nil, err
	}
	r, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(r)
}
