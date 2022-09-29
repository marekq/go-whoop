package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

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

// Load oauth2 token from local file
func loadToken(ctx context.Context) *http.Client {

	var token *oauth2.Token

	// Read config file from .env
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	// Set API key and profile
	var ClientID = viper.GetString("ClientID")
	var ClientSecret = viper.GetString("ClientSecret")

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
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://api.prod.whoop.com/oauth/oauth2/auth",
			TokenURL: "https://api.prod.whoop.com/oauth/oauth2/token",
		},
	}

	if _, err := os.Stat("token.json"); err == nil {

		f, err := os.Open("token.json")
		check(err)
		defer f.Close()

		// Read token from file
		err = json.NewDecoder(f).Decode(&token)
		check(err)
		fmt.Println("Token loaded from file" + token.TokenType + " " + token.Expiry.String())

		// Create client
		cclient := conf.Client(ctx, token)
		return cclient

	} else {

		// Request new client
		cclient := oauthRequest(ctx, conf)
		return cclient
	}
}

// Make request to Whoop API
func makeRequest(path string, filename string, cclient *http.Client) {

	// Create log file
	f2, err := os.Create(filename)
	check(err)
	defer f2.Close()

	// Set empty next token
	nextToken := "empty"

	// Loop through all next tokens
	for nextToken != "" {

		whoop_url := "https://api.prod.whoop.com/developer/" + path

		// If next token is not empty, add it to the get URL
		if nextToken != "" && nextToken != "empty" {
			whoop_url = whoop_url + "?nextToken=" + nextToken
		}

		// Request sleep data from Whoop API using client
		res, err := cclient.Get(whoop_url)
		check(err)
		fmt.Println("Status: " + res.Status)

		// Decode JSON to get nextToken
		var decodeStruct models.Sleep
		err = json.NewDecoder(res.Body).Decode(&decodeStruct)
		check(err)

		// Iterate through all structs
		for _, record := range decodeStruct.Records {

			fmt.Println(res.Status, record)
			json, err := json.Marshal(record)
			check(err)

			// Write JSON to file
			f2.WriteString(string(json) + ",\n")
		}

		// Get nextToken
		nextToken = decodeStruct.NextToken
	}
}

// OAuth2 request through web browser
// Tokens are valid for 1 hour
func oauthRequest(ctx context.Context, conf *oauth2.Config) *http.Client {

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
	accessToken, err := conf.Exchange(ctx, code)
	check(err)

	// Write response body to file
	f1, err := os.Create("token.json")
	check(err)
	defer f1.Close()

	// Marshal JSON
	json, err := json.Marshal(accessToken)
	check(err)

	// Write JSON to file
	_, err = f1.WriteString(string(json))
	check(err)

	// Create client
	cclient := conf.Client(ctx, accessToken)
	return cclient
}

// Main function
func main() {

	ctx := context.Background()
	cclient := loadToken(ctx)

	// Make requests to Whoop API
	makeRequest("v1/activity/sleep", "sleep.log", cclient)

}
