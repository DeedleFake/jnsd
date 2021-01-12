package jnsd

import "regexp"

var nameRE = regexp.MustCompile(`^[a-zA-Z0-9_]{3,32}$`)

func IsValidName(name string) bool {
	return nameRE.MatchString(name)
}
