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
		log.Fatal("Error while reading key", err)
	}

	secret_key := base64.RawStdEncoding.EncodeToString(key)
	fmt.Println("Secret key Generated: ", secret_key)

}
