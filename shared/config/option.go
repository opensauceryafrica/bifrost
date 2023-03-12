package config

// Options constants.
const (
	// ACL is the option to set the ACL of the file.
	OptACL = "acl"
	// PublicRead is the option to set the ACL of the file to public read.
	ACLPublicRead = "public-read"
	// Private is the option to set the ACL of the file to private.
	ACLPrivate = "private"
	// ContentType is the option to set the content type of the file.
	OptContentType = "content-type"
	// Metadata is the option to set the metadata of the file.
	OptMetadata = "metadata"
)

// Request constants
const (
	// ReqAuth is the authorization header identifier
	ReqAuth = "Authorization"
	// ReqBearer is the bearer token identifier
	ReqBearer = "Bearer %s"
	// MethodGet is the HTTP method for GET requests.
	MethodGet = "GET"
	// MethodPost is the HTTP method for POST requests.
	MethodPost = "POST"
	// MethodPut is the HTTP method for PUT requests.
	MethodPut = "PUT"
	// MethodDelete is the HTTP method for DELETE requests.
	MethodDelete = "DELETE"
)
