package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

func GetDriveService() (*drive.Service, error) {
	b, err := ioutil.ReadFile("oauth2.json")
	if err != nil {
		fmt.Printf("Unable to read credentials.json file. Err: %v\n", err)
		return nil, err
	}

	// If you want to modifyt this scope, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, drive.DriveScope)

	if err != nil {
		return nil, err
	}

	client := getClient(config)

	service, err := drive.New(client)

	if err != nil {
		fmt.Printf("Cannot create the Google Drive service: %v\n", err)
		return nil, err
	}

	return service, err
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	fmt.Println("Paste Authrization code here :")
	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func CreateFile(service *drive.Service, name string, mimeType string, content multipart.File, parentId string) (*drive.File, error) {
	f := &drive.File{
		MimeType: mimeType,
		Name:     name,
		Parents:  []string{parentId},
	}

	file, err := service.Files.Create(f).Media(content).Do()

	if err != nil {
		log.Println("Could not create file: " + err.Error())
		return nil, err
	}

	return file, nil
}

func SetFilePermissions(service *drive.Service, fileID string) error{
	permission := &drive.Permission{
		Type: "anyone",
		Role: "reader",
		AllowFileDiscovery: false,
	}

	_,err := service.Permissions.Create(fileID, permission).Do()

	return err
}

func DeleteFile(service *drive.Service, fileID string) error{
	err := service.Files.Delete(fileID).Do()

	return err
}

func GetFileIDFromURL(fileURL string) (string, error) {
    u, err := url.Parse(fileURL)
    if err != nil {
        return "", err
    }

    queryValues, err := url.ParseQuery(u.RawQuery)
    if err != nil {
        return "", err
    }

    fileID := queryValues.Get("id")
    if fileID == "" {
        return "", SetError(http.StatusInternalServerError, "File ID not found in the URL")
    }

    return fileID, nil
}