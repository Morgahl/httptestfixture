package httptestfixture_test

import (
	"context"
	"os"
	"testing"

	"github.com/curlymon/httptestfixture"
	"github.com/curlymon/httptestfixture/testapi"
)

var mocker *httptestfixture.JSONMocker

func TestMain(m *testing.M) {
	var err error
	mocker, err = httptestfixture.New("./fixtures/server.json")
	if err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func Test_GET_Home(t *testing.T) {
	svr := mocker.Server(t, "get home")
	defer svr.Close()

	api := testapi.NewTestAPI(svr.Client(), svr.URL)

	message, err := api.Home.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if message != "Hello World!" {
		t.Fatal("resp does not contain 'message': 'Hello World!'")
	}
}

func Test_POST_Home(t *testing.T) {
	svr := mocker.Server(t, "post home")
	defer svr.Close()

	api := testapi.NewTestAPI(svr.Client(), svr.URL)

	err := api.Home.Create(context.Background(), "Hello World!")
	if err != nil {
		t.Fatal(err)
	}
}
