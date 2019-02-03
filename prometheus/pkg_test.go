package prometheus

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	SetupFunction()
	code := m.Run()
	TeardownFunction(code)

}

func SetupFunction() {
	a = App{}
	a.Initilize()
}

func TeardownFunction(code int) {
	os.Exit(code)
}

func TestRoot(t *testing.T) {

	req, _ := http.NewRequest("GET", "/", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body !=
		`{repo:"https://github.com/mchirico/go_script"}` {
		t.Errorf("Expected an array. Got %s", body)
	}
}

func TestMetrics(t *testing.T) {

	req, _ := http.NewRequest("GET", "/metrics", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); strings.Contains(body,
		"A summary of the GC invocation durations") != true {
		t.Errorf("Expected an array. Got %s", body)
	}

}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
