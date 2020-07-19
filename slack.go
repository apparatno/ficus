package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var webhook = "https://hooks.slack.com/services/"

const markdown = "mrkdwn"

type payload []interface{}

type section struct {
	Text        text   `json:"text"`
	MessageType string `json:"type"`
}

type text struct {
	Text           string `json:"text"`
	FormattingType string
}

type footer struct {
	MessageType string `json:"type"`
	Elements    []text `json:"elements"`
}

func sendMessage(msg payload, token string, dryRun bool) error {
	log.Printf("sending message")

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(&msg); err != nil {
		return fmt.Errorf("failed to encode message: %v", err)
	}

	url := fmt.Sprintf("%s%s", webhook, token)
	req, err := http.NewRequest(http.MethodPost, url, buf)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")

	if dryRun {
		log.Println("skipping Slack message")
		return nil
	}

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

func makeMessage(u userFile) payload {
	msg := makeMessageString(u.username, u.filenames)
	body := section{
		MessageType: "section",
		Text: text{
			Text:           msg,
			FormattingType: markdown,
		},
	}
	f := footer{
		MessageType: "context",
		Elements: []text{
			{
				FormattingType: markdown,
				Text:           fmt.Sprintf("<https://drive.google.com/drive/u/0/folders/%s|Se utleggene til %v i utleggsmappen>", u.username, u.id),
			},
		},
	}
	return payload{body, f}
}

func makeMessageString(name string, files []string) string {
	msg := strings.Builder{}
	msg.WriteString(name)
	msg.WriteString(" har lastet opp ")
	msg.WriteString(strconv.Itoa(len(files)))
	msg.WriteString(" nye utlegg:\n")
	for _, f := range files {
		msg.WriteString(" * ")
		msg.WriteString(f)
		msg.WriteString("\n")
	}
	return msg.String()
}
