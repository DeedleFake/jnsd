package jnsd_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DeedleFake/jnsd"
)

func TestHandler(t *testing.T) {
	tests := []struct {
		name     string
		endpoint string
		status   int
		rsp      map[string]interface{}
	}{
		{
			name:     "Name/Valid",
			endpoint: "/name/valid",
			status:   http.StatusOK,
			rsp: map[string]interface{}{
				"name": "valid",
				"addr": "0xabcdef",
			},
		},
	}

	server := httptest.NewServer(jnsd.HandlerConfig{
		Name: func(name string) (string, error) {
			switch name {
			case "valid":
				return "abcdef", nil
			default:
				return "", jnsd.ErrNotRegistered
			}
		},
		Addr: func(addr string) (string, error) {
			switch addr {
			case "abcdef":
				return "valid", nil
			default:
				return "", jnsd.ErrNotRegistered
			}
		},
	}.Handler())

	client := server.Client()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rsp, err := client.Get(server.URL + test.endpoint)
			if err != nil {
				t.Errorf("get: %v", err)
			}
			defer rsp.Body.Close()

			buf, err := ioutil.ReadAll(rsp.Body)
			if err != nil {
				t.Errorf("read: %v", err)
			}

			var data map[string]interface{}
			err = json.Unmarshal(buf, &data)
			if err != nil {
				t.Errorf("unmarshal: %v", err)
			}

			if len(data) != len(test.rsp) {
				t.Errorf("got length: %v", len(data))
				t.Errorf("expected length: %v", len(test.rsp))
			}

			for k, v := range data {
				if v != test.rsp[k] {
					t.Errorf("got %q: %v", k, v)
					t.Errorf("expected %q: %v", k, test.rsp[k])
				}
			}
		})
	}
}
