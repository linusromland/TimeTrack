package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	PRODUCTION_CLIENT_ID     string
	PRODUCTION_CLIENT_SECRET string
	tokenDir                 = filepath.Join(userConfigDir(), "TimeTrack")
	tokenFile                = filepath.Join(tokenDir, "token.json")
)

func userConfigDir() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		fmt.Printf("failed to get user config dir: %v", err)
		return ""
	}
	return dir
}

func ensureTokenDirExists() {
	if _, err := os.Stat(tokenDir); os.IsNotExist(err) {
		os.MkdirAll(tokenDir, 0700)
	}
}

func GetClient() *http.Client {
	var config = getOAuthConfig()
	tok, err := tokenFromFile(tokenFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokenFile, tok)
	}
	return config.Client(context.Background(), tok)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	var token *oauth2.Token

	// Generate OAuth URL
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	// Start a local server to receive the OAuth callback
	callback := make(chan string)
	server := http.NewServeMux()
	server.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		queryParts, _ := url.ParseQuery(r.URL.RawQuery)
		code := queryParts["code"][0]
		fmt.Fprintf(w, "Received OAuth callback code. You can now close this window.")
		callback <- code
	})
	go http.ListenAndServe(":8080", server)

	// Open URL in browser
	fmt.Printf("Opening URL in your browser: %s\n", authURL)
	err := openURL(authURL)
	if err != nil {
		fmt.Printf("failed to open URL in browser: %v", err)
	}

	// Wait for the callback to return the authorization code
	authCode := <-callback

	// Use the authorization code to get the access token
	token, err = config.Exchange(context.TODO(), authCode)
	if err != nil {
		fmt.Printf("failed to exchange auth code for token: %v", err)
	}

	return token
}

func saveToken(file string, token *oauth2.Token) {
	ensureTokenDirExists()
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		fmt.Printf("unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

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

func getOAuthConfig() *oauth2.Config {
	CLIENT_ID := os.Getenv("GOOGLE_CLIENT_ID")
	if CLIENT_ID == "" {
		CLIENT_ID = PRODUCTION_CLIENT_ID
	}
	CLIENT_SECRET := os.Getenv("GOOGLE_CLIENT_SECRET")
	if CLIENT_SECRET == "" {
		CLIENT_SECRET = PRODUCTION_CLIENT_SECRET
	}

	config := &oauth2.Config{
		ClientID:     CLIENT_ID,
		ClientSecret: CLIENT_SECRET,
		RedirectURL:  "http://localhost:8080/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/calendar",
		},
		Endpoint: google.Endpoint,
	}
	return config
}

func openURL(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	return err
}
