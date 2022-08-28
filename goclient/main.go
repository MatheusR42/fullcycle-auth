package main

import (
	"context"
	"log"
	"net/http"
	"os"

	oicd "github.com/coreos/go-oidc"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

var (
	clientID     = "myclient"
	clientSecret = "CLIENT_SECRET"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx := context.Background()
	provider, err := oicd.NewProvider(ctx, "http://localhost:8080/realms/myrealm")
	if err != nil {
		log.Fatal(err)
	}

	config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: os.Getenv(clientSecret),
		Endpoint:     provider.Endpoint(),
		RedirectURL:  "http://localhost:8081/auth/callback",
		Scopes:       []string{oicd.ScopeOpenID, "profile", "email", "roles"},
	}

	state := "123"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, config.AuthCodeURL(state), http.StatusFound)
	})

	log.Fatal(http.ListenAndServe(":8081", nil))
}
