package commons

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

func Download(url string) (string, error) {
	client := &http.Client{
		Transport: &http.Transport{},
		Timeout: httpTimeoutSeconds * time.Second,
	}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		err := fmt.Errorf("Failed to create HTTP request (%s): %v", url, err)
		return "", err
	}
	request.Header.Set("User-Agent", userAgent)
	response, err := client.Do(request)
	if err != nil {
		err := fmt.Errorf("Failed to GET data (%s): %v", url, err)
		return "", err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		err := fmt.Errorf("Failed to read response (%s): %v", url, err)
		return "", err
	}
	return string(body), nil
}

func DownloadFile(url string, path string) error {
	data, err := Download(url)
	if err != nil {
		return err
	}
	WriteFile(path, data)
	return nil
}

func DownloadJSON[T any](url string) (T, error) {
	body, err := Download(url)
	if err != nil {
		var empty T
		return empty, err
	}
	var output T
	err = json.Unmarshal([]byte(body), &output)
	if err != nil {
		var empty T
		return empty, err
	}
	return output, nil
}

func BuildURL(base string, parameters map[string]string) string {
	u, err := url.Parse(base)
	if err != nil {
		log.Fatalf("Unable to parse URL (%s): %v", base, err)
	}
	values := url.Values{}
	for key, value := range parameters {
		values.Add(key, value)
	}
	u.RawQuery = values.Encode()
	encoded := u.String()
	return encoded
}