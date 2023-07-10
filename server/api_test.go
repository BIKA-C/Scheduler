package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"scheduler/repository"
	"scheduler/repository/database"
	"scheduler/router/errors"

	"github.com/bytedance/sonic"
	"github.com/nsf/jsondiff"
)

func prepareJSONRequest(method, path string, data map[string]any) (*httptest.ResponseRecorder, *http.Request) {
	b, _ := json.Marshal(data)
	r := httptest.NewRequest(method, path, bytes.NewReader(b))
	w := httptest.NewRecorder()
	return w, r
}

func cleanup(b []byte, ignore ...string) []byte {
	if len(ignore) > 1 {
		r, _ := sonic.Get(b)
		path := make([]any, len(ignore)-1)
		for l := 0; l < len(ignore)-1; l++ {
			path[l] = ignore[l]
		}
		parent := r.GetByPath(path...)
		if parent.Valid() {
			parent.Unset(ignore[len(ignore)-1])
			f, _ := r.Raw()
			return []byte(f)
		}
		return b
	} else if len(ignore) == 1 {
		r, _ := sonic.Get(b)
		r.Unset(ignore[0])
		f, _ := r.Raw()
		return []byte(f)
	} else {
		return b
	}
}

func compareJSONBody(expected obj, res io.ReadCloser, match jsondiff.Difference, ignore ...string) (bool, string) {
	buf, _ := io.ReadAll(res)
	defer res.Close()
	buf = cleanup(buf, ignore...)
	b, _ := json.Marshal(expected)
	opt := jsondiff.DefaultJSONOptions()
	opt.SkipMatches = true
	t, _ := jsondiff.Compare(b, buf, &opt)
	if t != jsondiff.FullMatch {
		return false, string(buf)
	}
	return true, ""
}

type obj map[string]any

func TestAPI_UserRegistration(t *testing.T) {
	db := database.NewSQLite(t.TempDir()+"/test.db", 40, database.DefaultPragma)
	api := NewWithDB(db)
	tests := []struct {
		name    string
		request obj
		expCode int
		expBody obj
	}{
		{
			name: "Base Case",
			request: obj{
				"account": obj{
					"email":    "real.email@mail.com",
					"password": "complicatedPassword1234",
				},
				"name": "Real User",
			},
			expCode: http.StatusOK,
			expBody: obj{
				"asset": obj{
					"balance": nil,
					"sum":     0,
				},
				"account": obj{
					"email": "real.email@mail.com",
				},
				"name": "Real User",
			},
		}, {
			name: "Register Already Existed",
			request: obj{
				"account": obj{
					"email":    "real.email@mail.com",
					"password": "complicatedPassword1234",
				},
				"name": "Real User",
			},
			expCode: repository.ErrUserAlreadyExist.Status,
			expBody: obj{
				"error": repository.ErrUserAlreadyExist.Msg,
			},
		}, {
			name: "Invalid Email Request",
			request: obj{
				"account": obj{
					"email":    "real.email-mail.com",
					"password": "complicatedPassword1234",
				},
				"name": "Real User",
			},
			expCode: http.StatusBadRequest,
			expBody: nil,
		}, {
			name: "Invalid Password Request",
			request: obj{
				"account": obj{
					"email": "real.email-mail.com",
				},
				"name": "Real User",
			},
			expCode: http.StatusBadRequest,
			expBody: nil,
		}, {
			name: "Invalid Username Request",
			request: obj{
				"account": obj{
					"email": "real.email-mail.com",
				},
			},
			expCode: http.StatusBadRequest,
			expBody: nil,
		},
	}
	for _, v := range tests {
		w, r := prepareJSONRequest("POST", "http://api.localhost/user/", v.request)
		api.ServeHTTP(w, r)
		if w.Result().StatusCode != v.expCode {
			t.Errorf("%s want %d got %d\n", v.name, v.expCode, w.Result().StatusCode)
		}
		if v.expBody == nil {
			return
		}
		if ok, res := compareJSONBody(v.expBody, w.Result().Body, 0, "account", "id"); !ok {
			t.Errorf("%s: Response not as expected:\n%s", v.name, res)
		}
	}
}

var sampleUser = obj{
	"account": obj{
		"email":    "real.email@mail.com",
		"password": "complicatedPassword1234",
	},
	"name": "Real User",
}

func TestAPI_UserUpdateAccount(t *testing.T) {
	db := database.NewSQLite(t.TempDir()+"/test.db", 40, database.DefaultPragma)
	api := NewWithDB(db)
	w, r := prepareJSONRequest("POST", "http://api.localhost/user/", sampleUser)
	api.ServeHTTP(w, r)
	b, _ := io.ReadAll(w.Result().Body)
	j, _ := sonic.Get(b, "account")
	id, _ := j.Get("id").String()

	tests := []struct {
		name  string
		id    string
		req   obj
		want  obj
		match jsondiff.Difference
		code  int
	}{
		{
			name: "Change Email",
			id:   id,
			req: obj{
				"email":  "new_mail@com.ca",
				"verify": "complicatedPassword1234",
			},
			want: obj{
				"id":    id,
				"email": "new_mail@com.ca",
			},
			code: 200,
		}, {
			name: "Change Nothing",
			id:   id,
			req: obj{
				"verify": "complicatedPassword1234",
			},
			want: obj{
				"error": errors.ErrBadRequest.Msg,
			},
			code: errors.ErrBadRequest.Status,
		}, {
			name: "Change Email Bad Password",
			id:   id,
			req: obj{
				"email":  "new_mail@com.ca",
				"verify": "complicatedPassword",
			},
			want: obj{
				"error": errors.ErrUnauthorized.Msg,
			},
			code: errors.ErrUnauthorized.Status,
		}, {
			name: "Change Password",
			id:   id,
			req: obj{
				"password": "abcdefgh",
				"verify":   "complicatedPassword1234",
			},
			want: obj{
				"id":    id,
				"email": "new_mail@com.ca",
			},
			code: 200,
		}, {
			name: "Change Email Old Password",
			id:   id,
			req: obj{
				"email":  "new_mail@com.ca",
				"verify": "complicatedPassword1234",
			},
			want: obj{
				"error": errors.ErrUnauthorized.Msg,
			},
			code: errors.ErrUnauthorized.Status,
		},
	}
	for _, v := range tests {
		w, r := prepareJSONRequest("PATCH", "http://api.localhost/user/"+v.id+"/account/", v.req)
	api.ServeHTTP(w, r)
		if w.Result().StatusCode != v.code {
			t.Errorf("%s want %d got %d\n", v.name, v.code, w.Result().StatusCode)
		}
		if v.want == nil {
			return
		}
		if ok, res := compareJSONBody(v.want, w.Result().Body, 0); !ok {
			t.Errorf("%s: Response not as expected:\n%s", v.name, res)
		}
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
