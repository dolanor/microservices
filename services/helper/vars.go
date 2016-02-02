package helper

var (
	// SymmetricKey is a single secret to do simple token signing
	SymmetricKey = "symmetrickey"
	// Cookiesecret is the password for the cookie store
	Cookiesecret = "cookiesecret"
	// DataServiceURL is the url to the DB Accessor.
	// Should be replaced with service discovery.
	DataServiceURL = "http://localhost:8080/api"
)
