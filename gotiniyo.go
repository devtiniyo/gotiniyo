// Package gotiniyo is a library for interacting with http://www.tiniyo.com/ API.
package gotiniyo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

const (
	baseURL       = "https://api.tiniyo.com/v1/"
	clientTimeout = time.Second * 30
)

// The default http.Client that is used if none is specified
var defaultClient = &http.Client{
	Timeout: time.Second * 30,
}

// Tiniyo stores basic information important for connecting to the
// tiniyo.com REST api such as AccountSid and AuthToken.
type Tiniyo struct {
	AuthID string
	AuthToken  string
	BaseUrl    string
	HTTPClient *http.Client
}

// Exception is a representation of a tiniyo exception.
type Exception struct {
	Status   int           `json:"status"`    // HTTP specific error code
	Message  string        `json:"message"`   // HTTP error message
	Code     ExceptionCode `json:"code"`      // Tiniyo specific error code
	MoreInfo string        `json:"more_info"` // Additional info from Tiniyo
}

// Print the RESTException in a human-readable form.
func (r Exception) Error() string {
	var errorCode ExceptionCode
	var status int
	if r.Code != errorCode {
		return fmt.Sprintf("Code %d: %s", r.Code, r.Message)
	} else if r.Status != status {
		return fmt.Sprintf("Status %d: %s", r.Status, r.Message)
	}
	return r.Message
}

// Create a new Tiniyo struct.
func NewTiniyoClient(accountSid, authToken string) *Tiniyo {
	return NewTiniyoClientCustomHTTP(accountSid, authToken, nil)
}

// Create a new Tiniyo client, optionally using a custom http.Client
func NewTiniyoClientCustomHTTP(authID, authToken string, HTTPClient *http.Client) *Tiniyo {
	if HTTPClient == nil {
		HTTPClient = defaultClient
	}

	return &Tiniyo{
		AuthID: authID,
		AuthToken:  authToken,
		BaseUrl:    baseURL,
		HTTPClient: HTTPClient,
	}
}

func (tiniyo *Tiniyo) getJSON(url string, result interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.SetBasicAuth(tiniyo.getBasicAuthCredentials())
	resp, err := tiniyo.do(req)
	if err != nil {
		return fmt.Errorf("failed to submit HTTP request: %v", err)
	}

	if resp.StatusCode != 200 {
		re := Exception{}
		json.NewDecoder(resp.Body).Decode(&re)
		return re
	}
	return json.NewDecoder(resp.Body).Decode(&result)
}

func (tiniyo *Tiniyo) getBasicAuthCredentials() (string, string) {

	return tiniyo.AuthID, tiniyo.AuthToken
}

func (tiniyo *Tiniyo) post(formValues url.Values, tiniyoUrl string) (*http.Response, error) {
	req, err := http.NewRequest("POST", tiniyoUrl, strings.NewReader(formValues.Encode()))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(tiniyo.getBasicAuthCredentials())
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return tiniyo.do(req)
}

func (tiniyo *Tiniyo) get(tiniyoUrl string) (*http.Response, error) {
	req, err := http.NewRequest("GET", tiniyoUrl, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(tiniyo.getBasicAuthCredentials())

	return tiniyo.do(req)
}

func (tiniyo *Tiniyo) delete(tiniyoUrl string) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", tiniyoUrl, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(tiniyo.getBasicAuthCredentials())

	return tiniyo.do(req)
}

func (tiniyo *Tiniyo) do(req *http.Request) (*http.Response, error) {
	client := tiniyo.HTTPClient
	if client == nil {
		client = defaultClient
	}

	return client.Do(req)
}

// Build path to a resource within the Tiniyo account
func (tiniyo *Tiniyo) buildUrl(resourcePath string) string {
	return tiniyo.BaseUrl + "/" + path.Join("Accounts", tiniyo.AuthID, resourcePath)
}
