package util

import "regexp"

var PhoneRegex = regexp.MustCompile(`^\+1\d{10}$`)
var EmailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{1,}$`)
