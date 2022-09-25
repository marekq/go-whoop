package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/jsonq"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

type OAuthToken struct {
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type"`
	RefreshToken string    `json:"refresh_token"`
	Expiry       time.Time `json:"expiry"`
}
type NextToken struct {
	NextToken string `json:"next_token"`
}

// Check error, exit if error
func check(e error) {
	if e != nil {
		fmt.Println(e)
		panic(e)
	}
}

// Load oauth2 token from local file
func loadToken(ctx context.Context) string {
	f, err := os.Open("token.json")
	check(err)
	defer f.Close()

	var token OAuthToken
	err = json.NewDecoder(f).Decode(&token)
	check(err)

	if token.Expiry.Before(time.Now()) {
		fmt.Println("OAuth token is expired")
		token.AccessToken = oauthRequest(ctx)

	} else {
		fmt.Println("OAuth token not expired")

	}

	return token.AccessToken

}

// Make request to Whoop API
func makeRequest(path string, filename string, access_token string) {

	// Create log file
	f2, err := os.Create(filename)
	check(err)
	defer f2.Close()

	nextToken := "empty"

	for nextToken != "" {

		whoop_url := "https://api.prod.whoop.com/developer/" + path

		if nextToken != "" && nextToken != "empty" {
			whoop_url = whoop_url + "?nextToken=" + nextToken
		}

		fmt.Println("whoop url: " + whoop_url)

		// Request sleep data from Whoop API
		req, err := http.NewRequest("GET", whoop_url, nil)
		req.Header.Add("Authorization", "Bearer "+access_token)
		check(err)

		// Perform request
		client := &http.Client{}
		res, err := client.Do(req)
		check(err)
		defer res.Body.Close()

		// Read response body and write to file
		body, err := io.ReadAll(res.Body)
		check(err)
		fmt.Println(string(body))
		f2.WriteString(string(body))

		// Decode JSON
		data := map[string]interface{}{}
		dec := json.NewDecoder(strings.NewReader(string(body)))
		dec.Decode(&data)

		// Get next token
		jq := jsonq.NewQuery(data)
		nextToken, err = jq.String("next_token")
		check(err)

		// Print response status code and body
		fmt.Println("status code: " + strconv.Itoa(res.StatusCode))
		fmt.Println(nextToken)
	}
}

// OAuth2 request through web browser
// Tokens are valid for 1 hour
func oauthRequest(ctx context.Context) string {

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
			"read:workout",
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

	// Write response body to file
	f1, err := os.Create("token.json")
	check(err)
	defer f1.Close()

	// Marshal JSON
	json, err := json.Marshal(tok)
	check(err)

	// Write JSON to file
	_, err = f1.WriteString(string(json))
	check(err)

	return tok.AccessToken
}

func main() {

	ctx := context.Background()
	access_token := loadToken(ctx)

	// Make requests to Whoop API
	makeRequest("v1/activity/sleep", "sleep.log", access_token)
	makeRequest("v1/activity/workout", "workout.log", access_token)

}
