package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/lestrrat-go/jwx/v3/jwk"
)

func createJWKS() error {
	pubKeyData, _ := os.ReadFile("../public.key")
	key, _ := jwk.ParseKey(pubKeyData, jwk.WithPEM(true))
	key.Set(jwk.KeyIDKey, "my-key-id")
	set := jwk.NewSet()
	set.AddKey(key)

	f, _ := os.Create("jwks.json")
	defer f.Close()
	return json.NewEncoder(f).Encode(set)
}

func main() {
	err := createJWKS()
	if err != nil {
		log.Fatal(err)
	}
}
