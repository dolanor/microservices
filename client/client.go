package client

import (
	"github.com/dolanor/microservices/api"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

// MicroserviceClient is a client library to connect directly to our microservices
// without people having to do the http connection work themselves.
type MicroserviceClient struct {
	// Host is the host for the microservices
	Host string
	// Client is the http.Client needed to keep track of jwt in cookies
	Client *http.Client
	// Username is the username of the user connecting to the service
	Username string
}

// NewMicroserviceClient create a client that connects to the host, given the username
// and password. It will handle the connection to the service and the persistence of the
// jwt token.
func NewMicroserviceClient(host, username, password string) (*MicroserviceClient, error) {
	cj, _ := cookiejar.New(nil)
	client := &http.Client{Jar: cj}
	microserviceclient := &MicroserviceClient{Host: host, Client: client, Username: username}
	resp, err := microserviceclient.authenticateCredentials(username, password)
	if err != nil {
		return nil, err
	}

	// Check for HTTP StatusCode
	switch resp.StatusCode {
	case http.StatusNotFound:
		return nil, api.ErrDataNotFound
	case http.StatusUnauthorized:
		return nil, api.ErrUnauthorized
	case http.StatusServiceUnavailable:
		return nil, api.ErrConnectingEndpoint
	case http.StatusInternalServerError:
		return nil, api.ErrInternalServer
	}

	return microserviceclient, nil
}

func (mc *MicroserviceClient) authenticateCredentials(username, password string) (*http.Response, error) {
	form := url.Values{}
	form.Add("username", username)
	form.Add("password", password)
	resp, err := QueryService(mc, "POST", mc.Host+"/login", TypeForm, form)
	return resp, err
}

// GetTodo returns the TODO list of the current user
func (mc *MicroserviceClient) GetTodo() (api.Todo, error) {
	var respTodo api.Todo

	url := mc.Host + "/todo"

	r, err := QueryService(mc, "GET", url, TypeJSON, nil)
	if err != nil {
		return respTodo, err
	}

	// Check for HTTP StatusCode
	switch r.StatusCode {
	case http.StatusNotFound:
		return respTodo, api.ErrDataNotFound
	case http.StatusUnauthorized:
		return respTodo, api.ErrUnauthorized
	case http.StatusServiceUnavailable:
		return respTodo, api.ErrConnectingEndpoint
	case http.StatusInternalServerError:
		return respTodo, api.ErrInternalServer
	}
	err = processResponseEntity(r, &respTodo, http.StatusOK)
	return respTodo, err
}

// GetUserProfile returns the user profile of the current user
func (mc *MicroserviceClient) GetUserProfile() (api.User, error) {
	var respUser api.User

	url := mc.Host + "/user/" + mc.Username

	r, err := QueryService(mc, "GET", url, TypeJSON, nil)
	if err != nil {
		return respUser, err
	}

	// Check for HTTP StatusCode
	switch r.StatusCode {
	case http.StatusNotFound:
		return respUser, api.ErrDataNotFound
	case http.StatusUnauthorized:
		return respUser, api.ErrUnauthorized
	case http.StatusServiceUnavailable:
		return respUser, api.ErrConnectingEndpoint
	case http.StatusInternalServerError:
		return respUser, api.ErrInternalServer
	}
	err = processResponseEntity(r, &respUser, http.StatusOK)
	return respUser, err
}
