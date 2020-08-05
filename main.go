package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type userFile struct {
	filenames []string
	username  string
	id        folderID
}

func main() {
	var did string
	var rid string
	var dbPath string
	var dryRun bool

	flag.StringVar(&did, "driveid", "", "ID of the Google Drive to use")
	flag.StringVar(&rid, "root", "", "ID of the folder to scan")
	flag.StringVar(&dbPath, "db", "db.json", "path to database JSON file. Defaults to ./db.json")
	flag.BoolVar(&dryRun, "no-slack", false, "don't send Slack messages")

	flag.Parse()

	token := os.Getenv("FICUS_SLACK_TOKEN")
	if token == "" && dryRun == false {
		log.Fatal("missing env var 'FICUS_SLACK_TOKEN'")
	}

	driveID := folderID(did)
	rootID := folderID(rid)

	srv, err := drive.NewService(context.Background(), option.WithCredentialsFile("credentials.json"))
	if err != nil {
		log.Fatalf("unable to retrieve Drive client: %v", err)
	}

	folders, err := getFolders(srv, driveID, rootID)
	if err != nil {
		log.Fatalf("failed to list folders from root folder: %v", err)
	}

	db, err := load(dbPath)
	if err != nil {
		log.Fatalf("failed to read data from database: %v", err)
	}
	log.Printf("database %v", db)

	changes, err := makeFileLists(srv, driveID, folders, db)
	if err != nil {
		log.Fatal(err)
	}

	for _, u := range changes {
		msg := makeMessage(u)
		if err := sendMessage(msg, token, dryRun); err != nil {
			log.Fatalf("failed to send Slack message: %v", err)
		}
	}

	db = updateDB(db, changes)

	if err := save(dbPath, db); err != nil {
		log.Fatalf("failed to save database: %v", err)
	}

	log.Println("update completed")
}

func getFolders(srv *drive.Service, driveID folderID, root folderID) ([]folder, error) {
	var res []folder

	r, err := srv.Files.
		List().
		TeamDriveId(string(driveID)).
		SupportsTeamDrives(true).
		IncludeTeamDriveItems(true).
		Corpora("drive").
		Q(fmt.Sprintf("'%s' in parents", root)).
		PageSize(100).
		Fields("nextPageToken, files(id, name, createdTime)").
		Do()
	if err != nil {
		return nil, err
	}

	for _, f := range r.Files {
		if strings.HasPrefix(f.Name, "xxx") || strings.Contains(f.Name, "README") {
			continue
		}

		log.Printf("%s (%s)", f.Name, f.Id)
		res = append(res, folder{ID: folderID(f.Id), User: f.Name})
	}

	log.Printf("found %d files", len(res))

	return res, nil
}

func makeFileLists(srv *drive.Service, driveID folderID, folders []folder, db map[folderID]folder) ([]userFile, error) {
	var res []userFile

	for _, f := range folders {
		userFolder, ok := db[f.ID]
		if !ok {
			userFolder = f
		}

		log.Printf("handling %s", userFolder.User)
		files, err := listFolder(srv, driveID, f.ID)
		if err != nil {
			log.Printf("failed to list files for user %s (%v)", userFolder.User, err)
			continue
		}
		if len(files) == 0 {
			continue
		}

		newFiles := mapFiles(files, userFolder)
		if len(newFiles) == 0 {
			log.Println("no new files")
			continue
		}
		log.Printf("found %d files to report", len(newFiles))

		res = append(res, userFile{
			id:        f.ID,
			filenames: newFiles,
			username:  userFolder.User,
		})
	}

	return res, nil
}

func listFolder(srv *drive.Service, driveID, id folderID) ([]*drive.File, error) {
	r, err := srv.Files.
		List().
		TeamDriveId(string(driveID)).
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
