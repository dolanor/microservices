package server

var (
	// Single secret to do simple token signing
	SymmetricKey = "symmetrickey"
	// password for the cookie store
	Cookiesecret = "cookiesecret"
	// url to the DB Accessor. Need to replace with service discovery.
	DataServiceURL = "http://localhost:8300"
)
