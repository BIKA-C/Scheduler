package errors

type Error struct {
	Err    error  `json:"-"`
	Status int    `json:"-"`
	Title  string `json:"error"`
}

func (e Error) Error() string {
	if e.Err == nil {
		return e.Title
	}
	return e.Err.Error()
}

func (e Error) Unwrap() error {
	return e.Err
}

func (e Error) String() string {
	return e.Title
}
