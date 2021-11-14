package httptestfixture

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type JSONMocker struct {
	file     string
	fixtures map[string]struct {
		Request  `json:"request"`
		Response `json:"response"`
	}
}

type Request struct {
	Path    string            `json:"path"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
}

type Response struct {
	Headers map[string]string      `json:"headers"`
	Status  int                    `json:"status"`
	Body    map[string]interface{} `json:"body"`
}

func New(path string) (*JSONMocker, error) {
	jm := &JSONMocker{file: path}
	if err := jm.loadFixture(path); err != nil {
		return nil, err
	}
	return jm, nil
}

func (jm *JSONMocker) loadFixture(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	dec := json.NewDecoder(file)
	dec.DisallowUnknownFields()

	return dec.Decode(&jm.fixtures)
}

func (jm *JSONMocker) fixtureHandler(t *testing.T, fixtureName string) func(http.ResponseWriter, *http.Request) {
	fixture, ok := jm.fixtures[fixtureName]
	if !ok {
		t.Fatalf("fixture does not exist: file=%s fixture=%s", jm.file, fixtureName)
		return nil
	}

	return func(rw http.ResponseWriter, r *http.Request) {
		if fixture.Request.Path != "" && fixture.Request.Path != r.URL.Path {
			t.Fatalf(
				"requested path incorrect: file=%s fixture=%s fixture_path=%s request_path=%s",
				jm.file, fixtureName, fixture.Request.Path, r.URL.Path,
			)
			return
		}

		if fixture.Request.Method != "" && fixture.Request.Method != r.Method {
			t.Fatalf(
				"request path incorrect: file=%s fixture=%s fixture_method=%s request_method=%s",
				jm.file, fixtureName, fixture.Request.Method, r.Method,
			)
			return
		}

		for key, value := range fixture.Request.Headers {
			if reqValue := r.Header.Get(key); reqValue != value {
				t.Fatalf(
					"request header incorrect: file=%s fixture=%s fixture_header=%s request_header=%s",
					jm.file, fixtureName, value, reqValue,
				)
				return
			}
		}

		for key, value := range fixture.Response.Headers {
			rw.Header().Set(key, value)
		}

		rw.WriteHeader(fixture.Response.Status)

		if fixture.Response.Body != nil {
			if err := json.NewEncoder(rw).Encode(fixture.Response.Body); err != nil {
				panic(fmt.Sprintf(
					"error encoding body for: file=%s fixture=%s error=%s",
					jm.file, fixtureName, err,
				))
			}
		}
	}
}

func (jm *JSONMocker) Server(t *testing.T, fixtureName string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(jm.fixtureHandler(t, fixtureName)))
}
