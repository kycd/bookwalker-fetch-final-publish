package p

import (
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHttp(t *testing.T) {
	tests := []struct {
		body string
		want string
	}{
		{body: `{"sn": 1336}`, want: ""},
		//{body: `{"sn": 5530}`, want: ""},
		//{body: `{}`, want: ""},
	}

	for _, test := range tests {
		req := httptest.NewRequest("GET", "/", strings.NewReader(test.body))
		req.Header.Add("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		Fetch(resp, req)

		out, err := ioutil.ReadAll(resp.Result().Body)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		t.Logf("%v\n", string(out))
		// if got := string(out); got != test.want {
		// 	t.Errorf("resp not want: %v", got)
		// }
	}
}
