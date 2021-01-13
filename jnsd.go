package jnsd

import "regexp"

var nameRE = regexp.MustCompile(`^[a-zA-Z0-9_]{3,32}$`)

// IsValidName returns true if the given name is valid.
func IsValidName(name string) bool {
	return nameRE.MatchString(name)
}
