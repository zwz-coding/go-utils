package restful

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

const succeed = "\u2713"
const failed = "\u2717"

func TestReadJSON(t *testing.T) {

	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	tests := []struct {
		name string
		json string
		err  bool
		p    Person
	}{
		{"Empty json", ``, true, Person{}},
		{"json data without matching", `{"abc":"123"}`, false, Person{}},
		{"Invalid json data", `{"abc":"123"}{"dog": 123}`, true, Person{}},
		{"Valid json", `{"name":"kevin","age":30}`, false, Person{"kevin", 30}},
	}

	t.Log("Given the need to test ReadJSON function.")
	{
		for index, test := range tests {
			t.Logf("\tTest %d[%s]:\tWhen checking %#q for error %v and person %v", index, test.name, test.json, test.err, test.p)
			{
				r, _ := http.NewRequest("GET", "http://localhost", bytes.NewBufferString(test.json))
				var p Person
				w := httptest.NewRecorder()
				err := ReadJSON(w, r, &p)
				if (err == nil) == test.err {
					t.Fatalf("\t%s\tTest %d[%s]:\tShould be a %v error : %v", failed, index, test.name, test.err, err)
				}
				t.Logf("\t%s\tTest %d[%s]:\tShould be a %v error", succeed, index, test.name, err)
				if err != nil {
					continue
				}
				if p != test.p {
					t.Fatalf("\t%s\tTest %d[%s]:\tShould be a person object %v: %v", failed, index, test.name, test.p, p)
				}
				t.Logf("\t%s\tTest %d[%s]:\tShould be a person object %v", succeed, index, test.name, p)
			}
		}
	}
}

func TestWriteJSON(t *testing.T) {

	tests := []struct {
		name   string
		status int
		data   any
		header http.Header
		err    error
	}{
		{"Empty body and empty header", 200, ``, http.Header{}, nil},
		{"Body and header", 200, struct {
			A string `json:"a"`
			B int    `json:"b"`
		}{"ok", 10}, http.Header{"test": []string{"1", "2"}, "hello": []string{"world"}}, nil},
		{"Body and empty header", 202, struct {
			A string `json:"a"`
			B int    `json:"b"`
		}{"ok", 10}, http.Header{}, nil},
	}

	t.Log("Given the need to test WriteJSON function.")
	{
		for index, test := range tests {
			t.Logf("\tTest %d[%s]:\tWhen input are status %d, data %v, header %v for error %v", index, test.name, test.status, test.data, test.header, test.err)
			{

				w := httptest.NewRecorder()
				err := WriteJSON(w, test.status, test.data, test.header)
				if err != test.err {
					t.Fatalf("\t%s\tTest %d[%s]:\tShould be a %v error : %v", failed, index, test.name, test.err, err)
				}
				t.Logf("\t%s\tTest %d[%s]:\tShould be a %v error", succeed, index, test.name, err)
				if err != nil {
					continue
				}

				// check json
				exp, _ := json.Marshal(test.data)
				bs, _ := io.ReadAll(w.Body)
				if string(exp) != string(bs) {
					t.Fatalf("\t%s\tTest %d[%s]:\tShould be json %#q : %#q", failed, index, test.name, string(exp), string(bs))
				}
				t.Logf("\t%s\tTest %d[%s]:\tShould be json %#q", succeed, index, test.name, string(exp))

				// check headers
				ct := w.Header().Get("Content-Type")
				if ct != "application/json" {
					t.Fatalf("\t%s\tTest %d[%s]:\tShould be included header Content-Type %s : %s", failed, index, test.name, "application/json", ct)
				}
				t.Logf("\t%s\tTest %d[%s]:\tShould be be included header Content-Type application/json", succeed, index, test.name)

				for key, value := range test.header {
					v, ok := w.Header()[key]
					if !ok || !same(value, v) {
						t.Fatalf("\t%s\tTest %d[%s]:\tShould be included header %s=%v : %s=%v", failed, index, test.name, key, value, key, v)
					}
					t.Logf("\t%s\tTest %d[%s]:\tShould be be included header %s=%v", succeed, index, test.name, key, value)
				}

				// check status
				if w.Code != test.status {
					t.Fatalf("\t%s\tTest %d[%s]:\tShould be status code %d : %d", failed, index, test.name, test.status, w.Code)
				}
				t.Logf("\t%s\tTest %d[%s]:\tShould be status code %d", succeed, index, test.name, w.Code)
			}
		}
	}
}

func TestErrorJSON(t *testing.T) {

	tests := []struct {
		name   string
		status int
		err    error
		result error
	}{
		{"With status code", 500, errors.New("test"), nil},
		{"Without status code", -1, errors.New("test"), nil},
	}

	t.Log("Given the need to test WriteJSON function.")
	{
		for index, test := range tests {
			t.Logf("\tTest %d[%s]:\tWhen input are status %d, err %v for error %v", index, test.name, test.status, test.err, test.result)
			{

				var err error
				expCode := http.StatusBadRequest
				w := httptest.NewRecorder()
				if test.status == -1 {
					err = ErrorJSON(w, test.err)
				} else {
					err = ErrorJSON(w, test.err, test.status)
					expCode = test.status
				}
				if err != test.result {
					t.Fatalf("\t%s\tTest %d[%s]:\tShould be a %v error : %v", failed, index, test.name, test.err, err)
				}
				t.Logf("\t%s\tTest %d[%s]:\tShould be a %v error", succeed, index, test.name, err)
				if err != nil {
					continue
				}

				// check headers
				ct := w.Header().Get("Content-Type")
				if ct != "application/json" {
					t.Fatalf("\t%s\tTest %d[%s]:\tShould be included header Content-Type %s : %s", failed, index, test.name, "application/json", ct)
				}
				t.Logf("\t%s\tTest %d[%s]:\tShould be be included header Content-Type application/json", succeed, index, test.name)

				// check status
				if w.Code != expCode {
					t.Fatalf("\t%s\tTest %d[%s]:\tShould be status code %d : %d", failed, index, test.name, expCode, w.Code)
				}
				t.Logf("\t%s\tTest %d[%s]:\tShould be status code %d", succeed, index, test.name, w.Code)
			}
		}
	}
}

func same(in, out []string) bool {
	if len(in) != len(out) {
		return false
	}
	for i, v := range in {
		if out[i] != v {
			return false
		}
	}
	return true
}
