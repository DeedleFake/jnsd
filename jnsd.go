package jnsd

import "regexp"

var nameRE = regexp.MustCompile(`^[a-zA-Z0-9_]{3,32}$`)

// IsNameValid returns true if the given name is valid.
func IsNameValid(name string) bool {
	return nameRE.MatchString(name)
}
