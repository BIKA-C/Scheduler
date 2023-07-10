package controller

import (
	"scheduler/router"
	"scheduler/router/errors"
	"scheduler/util"
)

func cast(e error) error {
	switch e.(type) {
	case errors.AccumulateError:
		return e
	case errors.Error:
		return e
	}
	return errors.ErrInvalidJSON
}

const (
	uid    = "uid"
	UserID = ":" + uid

	iid    = "iid"
	InstID = ":" + iid

	csid   = "csid"
	CursID = ":" + csid
)

func getUserID(c *router.C) util.UUID {
	return util.UUID(c.Param(uid))
}

func getInstID(c *router.C) util.UUID {
	return util.UUID(c.Param(iid))
}

func getCourseID(c *router.C) util.ID {
	return util.ParseID(c.Param(csid))
}
