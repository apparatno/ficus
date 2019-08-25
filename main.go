package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func main() {
	token := os.Getenv("FICUS_SLACK_TOKEN")
	if token == "" {
		log.Fatal("missing env var 'FICUS_SLACK_TOKEN'")
	}

	driveID := os.Getenv("FICUS_DRIVE_ID")
	if driveID == "" {
		log.Fatal("missing env var FICUS_DRIVE_ID")
	}

	srv, err := drive.NewService(context.Background(), option.WithCredentialsFile("credentials.json"))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	for {
		if err = run(srv, driveID, token); err != nil {
			log.Fatal(err)
		}
		time.Sleep(10 * time.Minute)
	}
}

func run(srv *drive.Service, driveID, token string) error {
	folders, err := load()
	if err != nil {
		return err
	}

	for i, folder := range folders {
		log.Printf("handling %s", folder.User)
		files, err := listFolder(srv, driveID, folder.ID)
		if err != nil {
			return err
		}
		if len(files) == 0 {
			continue
		}

		newFiles := mapFiles(files, folder)
		if len(newFiles) == 0 {
			log.Println("no new files")
			continue
		}
		log.Printf("found %d files to report", len(newFiles))

		msg := makeMessage(folder.User, newFiles)

		if err = sendMessage(msg, token); err != nil {
			return err
		}

		folders[i].UpdatedAt = time.Now()

		log.Printf("completed %s", folder.User)
	}

	if err = save(folders); err != nil {
		return err
	}
	log.Printf("done")

	return nil
}

func listFolder(srv *drive.Service, driveID, id string) ([]*drive.File, error) {
	r, err := srv.Files.
		List().
		TeamDriveId(driveID).
		SupportsTeamDrives(true).
		IncludeTeamDriveItems(true).
		Corpora("drive").
		Q(fmt.Sprintf("'%s' in parents", id)).
		PageSize(10).
		Fields("nextPageToken, files(id, name, createdTime)").
		Do()
	if err != nil {
		return nil, err
	}

	return r.Files, nil
}

func mapFiles(files []*drive.File, folder folder) []string {
	log.Printf("mapping %d files", len(files))
	newFiles := make([]string, 0, len(files))
	for _, f := range files {
		if f.Name == "Betalt" {
			continue
		}
		createdAt, err := time.Parse(time.RFC3339, f.CreatedTime)
		if err != nil {
			log.Printf("failed to parse created time '%s' from %s as RFC3339: %v", f.Name, f.CreatedTime, err)
			continue
		}
		if createdAt.After(folder.UpdatedAt) {
			newFiles = append(newFiles, f.Name)
		}
	}

	return newFiles
}

func makeMessage(name string, files []string) string {
	msg := strings.Builder{}
	msg.WriteString(name)
	msg.WriteString(" har lastet opp nye utlegg:\n")
	for _, f := range files {
		msg.WriteString(" * ")
		msg.WriteString(f)
		msg.WriteString("\n")
	}
	return msg.String()
}
