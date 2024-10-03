package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/stevefranchak/idlewalk/internal"
	"github.com/stevefranchak/idlewalk/internal/fitbit"
)

const (
	port string = "8443"
)

func main() {
	var loadedEnv map[string]string
	loadedEnv, err := godotenv.Read()
	if err != nil {
		log.Fatalf("Failed to read env file: %s", err)
	}
	fitbitCreds, err := fitbit.NewFitbitAPICred(loadedEnv)
	if err != nil {
		log.Fatalf("Failed to init Fitbit API Creds: %s", err)
	}
	log.Printf("%+v\n", fitbitCreds)

	verifier, err := fitbit.NewCodeVerifier(internal.CryptoRandIntGenerator)
	if err != nil {
		fmt.Println("Error generating the code verifier!")
		return
	}
	fmt.Printf("%v %v %v\n", verifier, verifier.SecretCode(), verifier.ChallengeCode())
}
