package httptestfixture_test

import (
	"encoding/json"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/curlymon/httptestfixture"
)

var svr *httptest.Server

func TestMain(m *testing.M) {
	mocker, err := httptestfixture.New("./fixtures/server.json")
	if err != nil {
		panic(err)
	}
	svr = mocker.Server()
	defer svr.Close()
	os.Exit(m.Run())

}

func TestFixture(t *testing.T) {
	resp, err := svr.Client().Get(svr.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var foo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&foo); err != nil {
		t.Fatal(err)
	}

	_, ok := foo["message"]
	if !ok {
		t.Fatal("resp does not contain 'message'")
	}
}
