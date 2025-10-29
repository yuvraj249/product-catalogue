package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
)

func main() {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		log.Fatal(err)

	}
	secret := base64.RawStdEncoding.EncodeToString(key)
	fmt.Println(secret)
}
