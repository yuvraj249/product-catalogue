package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	fmt.Println("Enter Password for hash: ")
	reader := bufio.NewReader(os.Stdin)
	pwd, _ := reader.ReadString('\n')
	for len(pwd) > 0 && (pwd[len(pwd)-1] == '\n' || pwd[len(pwd)-1] == '\r') {
		pwd = pwd[:len(pwd)-1]
	}
	if pwd == "" {
		fmt.Println("Please enter a valid Password!!")
		return
	}
	hashed_pwd, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Error while generating hash", err)
	}
	fmt.Println("Hash for password generated")
	fmt.Println(string(hashed_pwd))
}
