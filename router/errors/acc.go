package errors

import (
	"net/http"
)

var ErrInvalidJSON = BadRequest("Invalid JSON")

// AccumulateError stores errors in the process
type AccumulateError struct {
	Status int               `json:"-"`
	Title  string            `json:"title,omitempty"`
	Err    map[string]string `json:"error"`
}

var DefaultAccError = AccumulateError{
	Status: http.StatusBadRequest,
	Title:  "Invalid Input",
}

func (f AccumulateError) Error() string {
	if f.Title == "" {
		return "Invalid Input"
	}
	return f.Title
}

func (f *AccumulateError) Set(key, val string) {
	if f.Err == nil {
		f.Err = make(map[string]string, 5)
	}
	f.Err[key] = val
}

func (f AccumulateError) Get(key string) string {
	return f.Err[key]
}

func (f AccumulateError) Delete(key string) {
	if f.Err == nil {
		return
	}
	delete(f.Err, key)
}

func (f *AccumulateError) Merge(e error) {
	switch e.(type) {
	case AccumulateError:
		m := e.(AccumulateError)
		for k, v := range m.Err {
			f.Set(k, v)
		}
	}
}

func (f AccumulateError) NotOK() bool {
	return len(f.Err) != 0
}

func (f AccumulateError) OK() bool {
	return len(f.Err) == 0
}


func (f AccumulateError) Build() error {
	if f.NotOK() {
		return f
	}
	return nil
}

/*
	1xx		Form Errors
		100:	Name
		101:	Phone
		102:	Email
		103:	Password

		104:	Street
		105:	Province
		106:	Country
		107:	PostCode
	4xx		Other
		400		JSON
*/

// const (
// 	NameErr = 100 + iota
// 	PhoneErr
// 	EmailErr
// 	PasswordErr
// 	StreetErr
// 	ProvinceErr
// 	CountryErr
// 	PostcodeErr

// 	JSONErr = 400 + iota - 8
// )

// func Error(c int) error {

// }
