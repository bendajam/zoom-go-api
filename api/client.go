package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Client struct {
	serverURL  string
	httpClient *http.Client
	url        string
	authToken  string
	apiKey     string
	apiSecret  string
	namespace  string
}

func NewClient(
	serverURL string,
	apiKey string,
	apiSecret string) *Client {

	return &Client{
		serverURL:  serverURL,
		httpClient: &http.Client{},
		apiKey:     apiKey,
		apiSecret:  apiSecret,
	}
}

func (client *Client) doRequestJSON(
	path, method string, data bytes.Buffer, resp interface{}) error {

	body, err := client.doRequest(path, method, data)
	if err != nil {
		return fmt.Errorf("request error: %v", err)
	}

	// Check if JSON decoding is required
	if resp != nil {
		err = json.NewDecoder(body).Decode(resp)
	}

	return err
}

func (client *Client) doRequest(
	path, method string, data bytes.Buffer) (io.ReadCloser, error) {

	if client.authToken == "" {
		err := client.requestAuth()
		if err != nil {
			return nil, fmt.Errorf("unable to get auth token: %v", err)
		}
	}

	body, statusCode, err := client.httpRequest(path, method, data)

	// If auth token has expired, then get a new token
	if statusCode == http.StatusUnauthorized {
		if client.requestAuth() != nil {
			return nil, fmt.Errorf("unable to refresh auth token: %v", err)
		}

		// Try with new token, if second error occurs it is returned
		body, _, err = client.httpRequest(path, method, data)
	}

	if err != nil {
		return nil, fmt.Errorf("httpRequest error: %v", err)
	}

	return body, nil
}

func (client *Client) httpRequest(
	path, method string, body bytes.Buffer) (io.ReadCloser, int, error) {

	req, err := http.NewRequest(method, client.serverURL+"/"+path, &body)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", client.authToken)

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("unable to connect to server: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		respBody := new(bytes.Buffer)
		_, err := respBody.ReadFrom(resp.Body)
		if err != nil {
			return nil, resp.StatusCode,
				fmt.Errorf("received non 200 status code: %v",
					resp.StatusCode)
		}
		return nil, resp.StatusCode,
			fmt.Errorf("received non 200 status code: %v - %s",
				resp.StatusCode, respBody.String())
	}
	return resp.Body, 0, nil
}

func (client *Client) requestAuth() error {

	mySigningKey := []byte(client.apiSecret)
	expTime := time.Now().UTC().Add(time.Hour * time.Duration(1))

	// Create the Claims
	claims := &jwt.StandardClaims{
		ExpiresAt: expTime.Unix(),
		Issuer:    client.apiKey,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(mySigningKey)

	client.authToken = fmt.Sprintf("Bearer %s", tokenString)

	return err
}
