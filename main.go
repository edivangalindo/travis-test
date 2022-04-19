package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	stat, _ := os.Stdin.Stat()

	if (stat.Mode() & os.ModeCharDevice) != 0 {
		fmt.Fprintln(os.Stderr, "No tokens detected. Hint: cat tokens.txt | travis-test")
		os.Exit(1)
	}

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		token := scanner.Text()

		baseUrls := []string{"travis-ci.com", "travis-ci.org"}

		for _, url := range baseUrls {
			err := testTravisToken(token, url)

			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
		}
	}
}

func testTravisToken(token string, baseURL string) error {

	fmt.Printf("Testing token using %v -> %v\n", baseURL, token)

	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.%v/user", baseURL), nil)

	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Travis-API-Version", "3")
	req.Header.Add("Authorization", "token "+token)

	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	var prettyJson bytes.Buffer
	err = json.Indent(&prettyJson, body, "", "\t")

	if err != nil {
		// Travis api don't return a JSON as message, return a plaintext message
		fmt.Println("Error - Status code ->", resp.StatusCode)
		return nil
	}

	fmt.Println(string(prettyJson.Bytes()))

	return nil

}
