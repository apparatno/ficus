package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var webhook = "https://hooks.slack.com/services/"

type message struct {
	Text   string  `json:"text,omitempty"`
	Blocks []block `json:"blocks"`
}

type block struct {
	Kind string `json:"type"`
	Text text   `json:"text"`
}

type text struct {
	Kind string `json:"type"`
	Text string `json:"text"`
}

func sendMessage(msg message, token string, dryRun bool) error {
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

	if s := res.StatusCode; s > 299 {
		buf = new(bytes.Buffer)
		_, err = buf.ReadFrom(res.Body)
		if err != nil {
			return fmt.Errorf("failed to read response: %v", err)
		}

		switch resp := buf.String(); resp {
		case "no_text":
			return errors.New("no text")
		case "invalid_payload":
			return errors.New("invalid payload")
		default:
			log.Printf("slack response: %v", resp)
			return errors.New(resp)
		}
	}

	return nil
}

func makeMessage(u userFile) message {
	t := makeMessageString(u.username, u.filenames)

	m := block{
		Kind: "section",
		Text: text{
			Kind: "mrkdwn",
			Text: t,
		},
	}
	f := block{
		Kind: "section",
		Text: text{
			Kind: "mrkdwn",
			Text: fmt.Sprintf("<https://drive.google.com/drive/u/0/folders/%s|Gå til utleggsmappen og se alle utleggene til %v>", u.id, u.username),
		},
	}
	return message{
		Blocks: []block{m, f},
	}
}

func makeMessageString(name string, files []string) string {
	msg := strings.Builder{}
	msg.WriteString(name)
	msg.WriteString(" har lastet opp ")
	msg.WriteString(strconv.Itoa(len(files)))
	msg.WriteString(" nye utlegg:\n ")

	for _, f := range files {
		msg.WriteString("* ")
		msg.WriteString(f)
		msg.WriteString("\n")
	}

	return msg.String()
}
