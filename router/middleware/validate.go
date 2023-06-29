package middleware

import "scheduler/router"

func Validate(f func(c *router.C) bool) router.HandlerFunc {
	return func(c *router.C) error {
		if f(c) {
			return c.Next()
		}

		if c.GetValidation() {
			return c.Next()
		}

		c.SetValidation(true)
		e := c.Next()
		c.SetValidation(false)
		return e
	}
}
