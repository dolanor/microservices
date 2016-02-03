package client

import (
	"github.com/dolanor/microservices/api"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

var (
	todos = map[string]*api.Todo{
		"dolanor": {"Make the bed", "Eat", "Change the world", "Sleep"},
		"tanguy":  {"Code", "Prepare food", "Read"},
	}
	users = map[string]*api.User{
		"dolanor": {"dolanor", "Tanguy Herrmann", time.Date(1983, 01, 01, 0, 0, 0, 0, time.UTC)},
	}
)

func createGoodClient() (*MicroserviceClient, error) {
	return NewMicroserviceClient("http://localhost:8080", "dolanor", "test")
}

func clientWithWrongCredentials() (*MicroserviceClient, error) {
	return NewMicroserviceClient("http://localhost:8080", "dolanor", "tes")
}
func clientWithWrongHost() (*MicroserviceClient, error) {
	return NewMicroserviceClient("http://localhost:64321", "dolanor", "test")
}

func TestClientCreation(t *testing.T) {
	Convey("Given a well configured client", t, func() {
		_, err := createGoodClient()
		So(err, ShouldEqual, nil)
	})
	Convey("Given a client with wrong credentials", t, func() {
		_, err := clientWithWrongCredentials()
		So(err, ShouldEqual, api.ErrUnauthorized)
	})
	Convey("Given a client connecting to the wrong host", t, func() {
		_, err := clientWithWrongHost()
		So(err, ShouldNotEqual, nil)
	})
}

func TestGetTodo(t *testing.T) {
	Convey("Given a well configured client", t, func() {
		mc, _ := createGoodClient()
		Convey("Sending a request to get the todo list", func() {
			todo, err := mc.GetTodo()
			So(err, ShouldEqual, nil)
			So(&todo, ShouldResemble, todos["dolanor"])
		})
	})
}

func TestGetUserProfile(t *testing.T) {
	Convey("Given a well configured client", t, func() {
		mc, _ := createGoodClient()
		Convey("Sending a request to get the user profile", func() {
			userprofile, err := mc.GetUserProfile()
			So(err, ShouldEqual, nil)
			So(&user, ShouldResemble, users["dolanor"])
		})
	})
}
