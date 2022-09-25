package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

// Check error, exit if error
func check(e error) {
	if e != nil {
		fmt.Println(e)
		panic(e)
	}
}

func main() {
	ctx := context.Background()

	// Read config file from .env
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	// Set API key and profile
	var ClientID = viper.GetString("ClientID")
	var ClientSecret = viper.GetString("ClientSecret")

	// Set OAuth2 config
	conf := &oauth2.Config{
		ClientID:     ClientID,
		ClientSecret: ClientSecret,
		Scopes: []string{
			"offline",
			//"read:recovery",
			//"read:cycles",
			//"read:workout",
			"read:sleep",
			//"read:profile",
			//"read:body_measurement",
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://api.prod.whoop.com/oauth/oauth2/auth",
			TokenURL: "https://api.prod.whoop.com/oauth/oauth2/token",
		},
	}

	// Redirect user to consent page to ask for permission
	authUrl := conf.AuthCodeURL("stateidentifier", oauth2.AccessTypeOffline)
	fmt.Println("Visit the URL for the auth dialog: \n\n" + authUrl + "\n")
	fmt.Println("Enter the response URL: ")

	// Wait for user to paste in the response URL
	var respUrl string
	if _, err := fmt.Scan(&respUrl); err != nil {
		fmt.Println(respUrl)
		log.Fatal(err)
	}

	// Get code from response URL
	parseUrl, _ := url.Parse(respUrl)
	code := parseUrl.Query().Get("code")

	// Exchange code for token
	tok, err := conf.Exchange(ctx, code)
	check(err)

	// Request sleep data from Whoop API
	req, err := http.NewRequest("GET", "https://api.prod.whoop.com/developer/v1/activity/sleep", nil)
	req.Header.Add("Authorization", "Bearer "+tok.AccessToken)
	check(err)

	// Perform request
	client := &http.Client{}
	res, err := client.Do(req)
	check(err)

	defer res.Body.Close()

	// Read response body
	body, err := io.ReadAll(res.Body)
	check(err)

	// Write response body to file
	f, err := os.Create("sleep.log")

	check(err)
	defer f.Close()

	_, err = f.WriteString(string(body))
	check(err)

	// Print response status code and body
	fmt.Println(strconv.Itoa(res.StatusCode) + string(body))
	fmt.Println(res.Header)
}
