package types

import "strings"

// Provider is the cloud provider type
type Provider string

// String returns the string representation of the provider
func (p Provider) String() string {
	return string(p)
}

// ToLowerCase returns the lowercase representation of the provider
func (p Provider) ToLowerCase() string {
	return strings.ToLower(string(p))
}
