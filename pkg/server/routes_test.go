package server_test

import (
	"fmt"
	"net/http"
	"testing"
)

type TestCase struct {
	Path       string
	StatusCode int
}

const PORT = "8088"

var BASE_URL string = fmt.Sprintf("http://localhost:%s", PORT)

func TestRoutes(t *testing.T) {
	tcs := []TestCase{
		{Path: "healthz", StatusCode: 200},
	}

	for _, tc := range tcs {
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", BASE_URL, tc.Path), nil)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if res.StatusCode != tc.StatusCode {
			t.Errorf("Expected %d but got %d", tc.StatusCode, res.StatusCode)
		}
	}
}
