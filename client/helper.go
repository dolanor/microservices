package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// PostType indicates the type of post data (url-encoded or json).
type PostType int

const (
	// TypeForm is the urlencoded version for POST request.
	TypeForm PostType = iota
	// TypeJSON is the JSON version for POST request.
	TypeJSON
)

// QueryService is a helper function that queries the API and handles all the errors in the way.
func QueryService(client *MicroserviceClient, method, url string, postType PostType, entity interface{}) (*http.Response, error) {
	req, err := buildRequest(method, url, postType, entity)
	if err != nil {
		return nil, err
	}
	return client.Client.Do(req)
}

func buildRequest(method, url string, postType PostType, entity interface{}) (*http.Request, error) {
	body, err := marshalEntity(postType, entity)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return req, err
	}

	// Makes the server reply only json to our client
	req.Header.Set("Accept", "application/json")

	if postType == TypeJSON {
		req.Header.Set("content-type", "application/json")
	} else {
		req.Header.Set("content-type", "application/x-www-form-urlencoded")
	}
	return req, err
}

func marshalEntity(postType PostType, entity interface{}) (io.Reader, error) {
	if entity == nil {
		return nil, nil
	} else if postType == TypeJSON {
		b, err := json.Marshal(entity)
		if err != nil {
			return nil, err
		}
		return bytes.NewBuffer(b), nil
	} else if postType == TypeForm {
		return strings.NewReader(entity.(url.Values).Encode()), nil
	}
	return nil, nil
}

func processResponseEntity(r *http.Response, entity interface{}, expectedStatus int) error {
	if err := processResponse(r, expectedStatus); err != nil {
		return err
	}

	respBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(respBody, entity); err != nil {
		return err
	}

	return nil
}

func processResponse(r *http.Response, expectedStatus int) error {

	if r.StatusCode != expectedStatus {
		return errors.New("response status of " + r.Status)
	}
	return nil
}
