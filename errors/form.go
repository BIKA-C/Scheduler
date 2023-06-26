package errors

import "net/http"

type FormError struct {
	Status int               `json:"-"`
	Title  string            `json:"title,omitempty"`
	Err    map[string]string `json:"error"`
}

var DefaultFormError = FormError{
	Status: http.StatusBadRequest,
	Title:  "Invalid Input",
}

func (f FormError) Error() string {
	if f.Title == "" {
		return "Invalid Input"
	}
	return f.Title
}

func (f *FormError) Set(key, val string) {
	if f.Err == nil {
		f.Err = make(map[string]string, 5)
	}
	f.Err[key] = val
}

func (f *FormError) Get(key string) string {
	if f.Err == nil {
		f.Err = make(map[string]string, 5)
	}
	return f.Err[key]
}

func (f FormError) Delete(key string) {
	if f.Err == nil {
		return
	}
	delete(f.Err, key)
}

func (f *FormError) Merge(e error) {
	switch e.(type) {
	case FormError:
		m := e.(FormError)
		for k, v := range m.Err {
			f.Set(k, v)
		}
	}
}

func (f FormError) NotOK() bool {
	return len(f.Err) != 0
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
