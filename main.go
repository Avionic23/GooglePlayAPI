package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/androidpublisher/v3"
	"google.golang.org/api/option"
)

const (
	keyFile     = "/path_to_json_file"
	packageName = "com.foxwallet.play"
)

type VersionInfo struct {
	AppVersionName string
	AppVersionCode int64
}

func GetGooglePlayVersion(language string) (*VersionInfo, error) {
	// Set up a context and authenticate the client
	ctx := context.Background()

	// Read the contents of the service account key file
	keyBytes, err := ioutil.ReadFile(keyFile)
	if err != nil {
		log.Fatalf("Failed to read service account key file: %v", err)
	}

	fmt.Println(string(keyBytes))

	// Authenticate with Google Play Developer API using client credentials
	creds, err := google.CredentialsFromJSON(ctx, keyBytes, androidpublisher.AndroidpublisherScope)
	if err != nil {
		log.Fatalf("Failed to find valid credentials: %v", err)
	}

	fmt.Println("creds were created succefully")

	// Create a new AndroidPublisher client
	androidPubService, err := androidpublisher.NewService(ctx, option.WithCredentials(creds))
	if err != nil {
		log.Fatalf("Failed to create AndroidPublisher service: %v", err)
	}

	fmt.Println("androidPubService was created succefully")

	// Create a new edit for the app version
	appEdit, err := androidPubService.Edits.Insert(packageName, nil).Do()
	if err != nil {
		log.Fatalf("Failed to create edit: %v", err)
	}
	editId := appEdit.Id

	fmt.Println("appEdit was created succefully")

	// Retrieve the latest tracks for the app
	listTrackResponse, err := androidPubService.Edits.Tracks.List(packageName, editId).Do()
	if err != nil {
		log.Fatalf("Failed to retrieve track list: %v", err)
	}

	fmt.Println("listTrackResponse was created succefully")

	// Loop through the list of tracks to find the latest one
	var latestTrack *androidpublisher.Track
	for _, track := range listTrackResponse.Tracks {
		if latestTrack == nil || track.Releases[0].VersionCodes[0] > latestTrack.Releases[0].VersionCodes[0] {
			latestTrack = track
		}
	}

	fmt.Println("latestTrack was created succefully")

	// Retrieve the details for the latest release in the current track
	trackReleaseResponse, err := androidPubService.Edits.Tracks.Get(packageName, editId, latestTrack.Track).Do()
	if err != nil {
		log.Fatalf("Failed to retrieve app details: %v", err)
	}

	fmt.Println("trackReleaseResponse was created succefully")

	// Find the release that matches the specified language
	var release *androidpublisher.TrackRelease
	for _, r := range trackReleaseResponse.Releases {
		if r.ReleaseNotes[0].Language == language {
			release = r
			break
		}
	}

	fmt.Println(release)

	// Return an error if no release was found for the specified language
	if release == nil {
		return nil, fmt.Errorf("no release found for language %q", language)
	}

	// Retrieve the app version name and code for the release
	appVersionName := release.Name
	appVersionCode := release.VersionCodes[0]

	// Create a new VersionInfo struct with the app version information
	versionInfo := &VersionInfo{
		AppVersionName: appVersionName,
		AppVersionCode: appVersionCode,
	}

	return versionInfo, nil

}
