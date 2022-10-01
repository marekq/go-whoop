package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
	"golang.org/x/oauth2"

	models "github.com/marekq/go-whoop/model"
)

// Check error, exit if error
func check(e error) {
	if e != nil {
		fmt.Println(e)
		panic(e)
	}
}

func readLocalToken() oauth2.Token {

	f, err := os.Open("token.json")
	check(err)
	defer f.Close()

	var token oauth2.Token
	json.NewDecoder(f).Decode(&token)

	return token
}

func writeLocalToken(token *oauth2.Token) {

	f, err := os.Create("token.json")
	check(err)
	defer f.Close()

	json, err := json.Marshal(token)
	check(err)
	f.WriteString(string(json))

}

func getOauthConfig() (*oauth2.Config, string, string) {

	// Read config file from .env
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	// Set API key and profile from config file
	ClientID := viper.GetString("ClientID")
	ClientSecret := viper.GetString("ClientSecret")

	// Check ClientID and ClientSecret values exist
	if ClientID == "" || ClientSecret == "" {

		fmt.Println("ClientID and ClientSecret must be set in .env file")
		os.Exit(1)
	}

	// Set OAuth2 config
	conf := &oauth2.Config{
		ClientID:     ClientID,
		ClientSecret: ClientSecret,
		Scopes: []string{
			"offline",
			"read:recovery",
			"read:cycles",
			"read:workout",
			"read:sleep",
			"read:profile",
			"read:body_measurement",
		},
		RedirectURL: "https://coldstart.dev/",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://api.prod.whoop.com/oauth/oauth2/auth",
			TokenURL: "https://api.prod.whoop.com/oauth/oauth2/token",
		},
	}

	return conf, ClientID, ClientSecret
}

// Load oauth2 token from local file
func loadToken() string {

	// Set accessToken variable
	accessToken := ""

	// Set OAuth2 config
	conf, ClientID, ClientSecret := getOauthConfig()

	// Check if token.json file exists
	if _, err := os.Stat("token.json"); err == nil {

		localToken := readLocalToken()

		if !localToken.Valid() {

			fmt.Println("Local token expired at " + localToken.Expiry.String() + " , refreshing...")

			form := url.Values{}
			form.Add("grant_type", "refresh_token")
			form.Add("refresh_token", localToken.RefreshToken)
			form.Add("client_id", ClientID)
			form.Add("client_secret", ClientSecret)
			form.Add("scope", "offline")

			body := strings.NewReader(form.Encode())
			req, err := http.NewRequest("POST", "https://api.prod.whoop.com/oauth/oauth2/token", body)
			check(err)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			client := &http.Client{}
			resp, err := client.Do(req)
			check(err)

			// Decode JSON
			var tokenResponse models.TokenLocalFile
			err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
			check(err)

			// Marshal JSON
			newToken := &oauth2.Token{
				AccessToken:  tokenResponse.AccessToken,
				TokenType:    tokenResponse.TokenType,
				RefreshToken: tokenResponse.RefreshToken,
				Expiry:       time.Now().Local().Add(time.Second * time.Duration(tokenResponse.ExpiresIn)),
			}

			// Write token to file
			writeLocalToken(newToken)

			accessToken = tokenResponse.AccessToken

		} else {

			// Token is valid, use it without refresh
			fmt.Println("Local token valid till " + localToken.Expiry.String() + ", reused without refresh")
			accessToken = localToken.AccessToken

		}

	} else {

		// If token.json not present, start browser authentication flow
		fmt.Println("No token.json found, starting OAuth2 flow")

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

		// Get response code from response URL string
		parseUrl, _ := url.Parse(respUrl)
		code := parseUrl.Query().Get("code")

		// Exchange response code for token
		accessToken, err := conf.Exchange(context.Background(), code)
		check(err)

		// Write token to file
		writeLocalToken(accessToken)

	}

	// Return access token and newline
	fmt.Println("")
	return accessToken

}

// Make request to Whoop API
func makeRequest(path string, filename string, accessToken string) {

	fmt.Println("Making requests to " + path)

	// Create log file
	f2, err := os.Create(filename)
	check(err)
	defer f2.Close()

	// Set empty next token
	nextToken := "empty"
	count := 0

	// Loop through all next tokens
	for nextToken != "" {

		whoop_url := "https://api.prod.whoop.com/developer/" + path

		// If next token is not empty, add it to the get URL
		if nextToken != "" && nextToken != "empty" {
			whoop_url = whoop_url + "?nextToken=" + nextToken
		}

		// Request sleep data from Whoop API using client
		req, err := http.NewRequest("GET", whoop_url, nil)
		check(err)

		// Add authorization and content header
		req.Header.Add("Authorization", "Bearer "+accessToken)
		req.Header.Add("Content-Type", "application/json")

		// Make request
		client := &http.Client{}
		resp, err := client.Do(req)
		check(err)

		// Decode JSON to get nextToken
		var decodeStruct models.All
		err = json.NewDecoder(resp.Body).Decode(&decodeStruct)
		check(err)

		// Iterate through all structs
		for _, record := range decodeStruct.Records {

			// Write JSON to file
			json, err := json.Marshal(record)
			check(err)
			f2.WriteString(string(json) + ",\n")

			// Increment count
			count++
		}

		// Print status message per 100 records
		xrate_str := resp.Header.Get("X-RateLimit-Remaining")
		xrate_int, err := strconv.Atoi(xrate_str)
		check(err)

		if count%100 == 0 {
			fmt.Println("Processed " + strconv.Itoa(count) + " " + path + ", X-RateLimit remaining: " + xrate_str)
		}

		if xrate_int < 25 {
			fmt.Println("X-RateLimit low: " + xrate_str + ", waiting 5 seconds...")
			time.Sleep(5 * time.Second)
		}

		// Get nextToken
		nextToken = decodeStruct.NextToken
	}

	fmt.Println("Completed " + strconv.Itoa(count) + " " + path + " records\n")

}

// Main function
func main() {

	// Create client
	accessToken := loadToken()

	// Make requests to Whoop Sleep API
	makeRequest("v1/activity/sleep", "sleep.log", accessToken)
	makeRequest("v1/recovery", "recovery.log", accessToken)
	makeRequest("v1/cycle", "cycle.log", accessToken)
	makeRequest("v1/activity/workout", "workout.log", accessToken)

}
