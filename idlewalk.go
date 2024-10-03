package main

import (
	"fmt"

	"github.com/stevefranchak/idlewalk/internal"
	"github.com/stevefranchak/idlewalk/internal/fitbit"
)

const (
	port string = "8443"
)

func main() {
	verifier, err := fitbit.NewCodeVerifier(internal.CryptoRandIntGenerator)
	if err != nil {
		fmt.Println("Error generating the code verifier!")
		return
	}
	fmt.Printf("%v %v %v\n", verifier, verifier.SecretCode(), verifier.ChallengeCode())
}
