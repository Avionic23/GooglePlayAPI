package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	//"time"

	"github.com/PuerkitoBio/goquery"
)

type VersionInfo2 struct {
	AppVersion      string
	UpdateInfo      string
	PublishDateTime string
}

func GooglePlay(language string) (*VersionInfo2, error) {
	// Construct the URL for the app details page in the specified language
	appURL := fmt.Sprintf("https://play.google.com/store/apps/details?id=com.foxwallet.play&hl=%s", language)

	// Send HTTP GET request to the app page
	res, err := http.Get(appURL)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer res.Body.Close()

	// Parse the HTML document using goquery
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML document: %v", err)
	}

	updateInfo := doc.Find("[itemprop='description']").Text()

	// Split the update information into two parts
	updateInfoParts := strings.Split(updateInfo, "- ")

	// Set the first element to be the version and the second - update information
	version := strings.TrimSpace(updateInfoParts[0])
	updateInfo = strings.TrimSpace(updateInfoParts[1])

	publishTimeStr := doc.Find(".xg1aie").First().Text()
	publishTimeStr = strings.TrimSpace(publishTimeStr)
	//layout := "Jan 2, 2006"
	//publishTime, err := time.Parse(layout, publishTimeStr)
	if err != nil {
		fmt.Println(err)
	}

	// Create a new VersionInfo struct with the extracted information
	versionInfo2 := &VersionInfo2{
		AppVersion:      strings.TrimSpace(version),
		UpdateInfo:      strings.TrimSpace(updateInfo),
		PublishDateTime: strings.TrimSpace(publishTimeStr),
	}

	return versionInfo2, nil
}

func main() {
	value, err := GooglePlay("en")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Version: %s\n", value.AppVersion)
	fmt.Printf("Update Information: %s\n", value.UpdateInfo)
	fmt.Printf("Publish Time: %s\n", value.PublishDateTime)
}
