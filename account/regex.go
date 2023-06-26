package account

import "regexp"

var phone = regexp.MustCompile(`^\+1\d{10}$`)
var email = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{1,}$`)
