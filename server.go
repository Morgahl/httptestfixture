package httptestfixture

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
)

type JSONMock struct {
	fixture map[string]map[string]ResponseMock
}

type ResponseMock struct {
	Headers map[string]string      `json:"headers"`
	Status  int                    `json:"status"`
	Body    map[string]interface{} `json:"body"`
}

func New(path string) (*JSONMock, error) {
	jm := &JSONMock{}
	if err := jm.loadFixture(path); err != nil {
		return nil, err
	}
	return jm, nil
}

func (jm *JSONMock) loadFixture(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	dec := json.NewDecoder(file)
	dec.DisallowUnknownFields()

	return dec.Decode(&jm.fixture)
}

func (jm *JSONMock) fixtureHandler(rw http.ResponseWriter, r *http.Request) {
	resource, exists := jm.fixture[r.URL.Path]
	if !exists {
		panic(fmt.Sprintf("fixture does not exist for: path=%s", r.URL.Path))
	}

	fixture, exists := resource[r.Method]
	if !exists {
		panic(fmt.Sprintf("method does not exist on fixture: method=%s path=%s", r.Method, r.URL.Path))
	}

	for key, value := range fixture.Headers {
		rw.Header().Set(key, value)
	}

	rw.WriteHeader(fixture.Status)

	if fixture.Body != nil {
		json.NewEncoder(rw).Encode(fixture.Body)
	}
}

func (jm *JSONMock) Server() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(jm.fixtureHandler))
}
