package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"golang.org/x/term"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"syscall"
	"time"
)

type ValidationResult struct {
	ID int `json:"id"`
}

type ValidationResultResponse struct {
	ValidationResults []ValidationResult `json:"validationResults"`
}

func createClient() *http.Client {
	return &http.Client{Timeout: 10 * time.Second}
}

func fetchAndDeleteValidationResults(client *http.Client, baseURL, auth string) error {
	for {
		results, err := fetchValidationResults(client, baseURL, auth)
		if err != nil {
			return err
		}

		if len(results) == 0 {
			fmt.Println("✅ No more validation results found. Stopping.")
			break
		}

		fmt.Printf("%d validation results found. Deleting...\n", len(results))

		deleteValidationResults(client, baseURL, auth, results)
	}

	return nil
}

func fetchValidationResults(client *http.Client, baseURL, auth string) ([]ValidationResult, error) {
	url := fmt.Sprintf("%s/api/validationResults?page=1&pageSize=50", baseURL)
	fmt.Println("ℹ️ Fetching:", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("❌ Failed to fetch validation results, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result ValidationResultResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result.ValidationResults, nil
}

func deleteValidationResults(client *http.Client, baseURL, auth string, results []ValidationResult) {
	for _, validation := range results {
		err := deleteValidationResult(client, baseURL, auth, validation.ID)
		if err != nil {
			fmt.Printf("❌ Error deleting validation result %d: %v\n", validation.ID, err)
		} else {
			fmt.Printf("✅ Deleted validation result: %d\n", validation.ID)
		}
	}
}


func deleteValidationResult(client *http.Client, baseURL, auth string, id int) error {
	req, err := http.NewRequest("DELETE", baseURL+"/api/validationResults/"+url.PathEscape(strconv.Itoa(id)), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("❌ Failed to delete validation result %d, status code: %d", id, resp.StatusCode)
	}

	return nil
}

func getAuthToken(username, password string) string {
	credentials := fmt.Sprintf("%s:%s", username, password)
	return base64.StdEncoding.EncodeToString([]byte(credentials))
}

func main() {
	fmt.Print("Enter DHIS2 base URL: ")
	var baseURL string
	fmt.Scanln(&baseURL)

	fmt.Print("Enter your DHIS2 username: ")
	var username string
	fmt.Scanln(&username)

	fmt.Print("Enter your DHIS2 password: ")
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Println("\n❌ Failed to read password")
		return
	}
	password := string(passwordBytes)
	fmt.Println()

	auth := getAuthToken(username, password)

	client := createClient()

	err = fetchAndDeleteValidationResults(client, baseURL, auth)
	if err != nil {
		fmt.Println("❌ Error:", err)
	}
}