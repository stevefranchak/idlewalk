package main

import (
	"context"
	"embed"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	core "github.com/stevefranchak/idlewalk/internal"
	"github.com/stevefranchak/idlewalk/internal/fitbit"
	"golang.org/x/oauth2"

	"github.com/gin-gonic/gin"
)

const (
	port string = "8443"
)

var currentVerifier string

//go:embed migrations/*.sql
var migrationFiles embed.FS

func loadEnvFile() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, relying on environment variables only:", err)
	}
}

func fitbitAuthHandler(fitbitOauthConfig *oauth2.Config) func(*gin.Context) {
	return func(c *gin.Context) {
		currentVerifier = oauth2.GenerateVerifier()
		s256ChallengeOption := oauth2.S256ChallengeOption(currentVerifier)

		authUrl := fitbitOauthConfig.AuthCodeURL("", s256ChallengeOption)
		c.Redirect(http.StatusSeeOther, authUrl)
	}
}

func fitbitCallbackHandler(fitbitOauthConfig *oauth2.Config) func(*gin.Context) {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		code := c.Query("code")
		log.Println("Received code:", code)

		token, err := fitbitOauthConfig.Exchange(ctx, code, oauth2.VerifierOption(currentVerifier))
		if err != nil {
			log.Println("Failed to exchange for token:", err)
			c.String(http.StatusInternalServerError, "Could not exchange for token")
			return
		}

		client := fitbitOauthConfig.Client(ctx, token)
		resp, err := client.Get("https://api.fitbit.com/1/user/-/profile.json")
		if err != nil {
			log.Println("Failed to get user's Fitbit profile:", err)
			c.String(http.StatusInternalServerError, "Could not get user's Fitbit profile")
			return
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		log.Println(string(body))
	}
}

func main() {
	ctx := context.Background()
	loadEnvFile()
	fitbitOauthConfig, err := fitbit.NewFitbitOauthConfig()
	if err != nil {
		log.Fatalln("Could not create Fitbit Oauth2 config:", err)
	}

	db, err := core.SetupDb(ctx, migrationFiles)
	if err != nil {
		log.Fatalln("Could not setup db:", err)
	}
	defer db.Close()

	r := gin.Default()
	// https://github.com/gin-gonic/gin/blob/master/docs/doc.md#dont-trust-all-proxies
	// Not worried about proxies to this service quite yet - but want to remove the startup warning
	r.SetTrustedProxies(nil)

	r.GET("/fitbit-auth", fitbitAuthHandler(fitbitOauthConfig))
	r.GET("/fitbit-callback", fitbitCallbackHandler(fitbitOauthConfig))

	listenPort := fmt.Sprintf(":%s", port)
	log.Printf("Listening on port %s", listenPort)
	err = r.Run(listenPort)
	if err != nil {
		log.Fatalln("Failed to start http server:", err)
	}
}
