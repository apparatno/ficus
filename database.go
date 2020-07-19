package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type folderID string
type folder struct {
	User      string    `json:"user"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func load(path string) (map[folderID]folder, error) {
	var res map[folderID]folder
	f, err := os.Open(path)
	if err != nil {
		log.Println("no database file found, creating empty database")
		return res, nil
	}

	if err = json.NewDecoder(f).Decode(&res); err != nil {
		return nil, fmt.Errorf("failed to decode database: %v", err)
	}

	return res, nil
}

func updateDB(db map[folderID]folder, changes []userFile) map[folderID]folder {
	for _, f := range changes {
		entry, ok := db[f.id]
		if !ok {
			entry = folder{User: f.username}
		}
		entry.UpdatedAt = time.Now()
		db[f.id] = entry
	}
	return db
}

func save(path string, data map[folderID]folder) error {
	log.Println("saving data")
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("failed to open database file for writing: %v", err)
	}
	if err = json.NewEncoder(f).Encode(&data); err != nil {
		return fmt.Errorf("failed to encode data for writing: %v", err)
	}
	log.Println("data saved")
	return nil
}
