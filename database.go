package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type folder struct {
	ID        string    `json:"id"`
	User      string    `json:"user"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func load() ([]folder, error) {
	f, err := os.Open("db.json")
	if err != nil {
		return nil, fmt.Errorf("failed to open database file for reading: %v", err)
	}

	var res []folder
	if err = json.NewDecoder(f).Decode(&res); err != nil {
		return nil, fmt.Errorf("failed to decode database: %v", err)
	}

	return res, nil
}

func save(data []folder) error {
	log.Println("saving data")
	f, err := os.OpenFile("db.json", os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return fmt.Errorf("failed to open database file for writing: %v", err)
	}
	if err = json.NewEncoder(f).Encode(&data); err != nil {
		return fmt.Errorf("failed to encode data for writing: %v", err)
	}
	log.Println("data saved")
	return nil
}
