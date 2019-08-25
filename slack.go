package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var webhook = "https://hooks.slack.com/services/"

type message struct {
	Text string `json:"text"`
}

func sendMessage(msg, token string) error {
	log.Printf("sending message: %s", msg)
	m := message{Text: msg}
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(&m); err != nil {
		return fmt.Errorf("failed to encode message: %v", err)
	}

	url := fmt.Sprintf("%s%s", webhook, token)
	req, err := http.NewRequest(http.MethodPost, url, buf)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}

	buf = new(bytes.Buffer)
	_, err = buf.ReadFrom(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}
	log.Printf("slack response: %v", buf.String())

	return nil
}
