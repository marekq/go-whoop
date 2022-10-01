go-whoop
========

Download sleep, workout, cycle and recovery data using the [Whoop API](https://developer.whoop.com/api). You need a  [Whoop developer account](https://developer-dashboard.whoop.com) to authenticate. This takes less than 5 minutes to set up and is free for any user.  

## Getting started
 
Follow these steps to download your Whoop data:

- Make sure you have `go` installed locally (on Mac: `brew install go`). 
- Create a new app in the [Whoop Developer portal](https://developer-dashboard.whoop.com). Ensure that it has access to all OAuth2 `read:` scopes. Your app can be in test mode without being published or verified. 
- Create a `.env` file with your `ClientID` and `ClientSecret`. You can view and example in `example.env`. 
- Next, run `go build` to compile the program. 
- Now, run `./go-whoop` to get started with the OAuth2 authentication.
- You will be prompted to sign in using a URL. Click on the URL and login to Whoop. 
- Once you completed authentication, copy/paste the full URL into your terminal and hit enter. This way, the response OAuth2 code is returned. 
- The program should now download your cycle, recovery, sleep and workout data to a log file.
- When you run the program again, it will check if your authentication token is still valid and refresh it if needed.

## Example commands for CLI:

### Download all Whoop logs (cycle, recovery, sleep, workout)
`./go-whoop`