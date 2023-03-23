package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path"
	"testing"

	"yatter-backend-go/app/app"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountRegistration(t *testing.T) {
	c := setup(t)
	defer c.Close()

	testCases := []struct {
		name        string
		method      string
		apiPath     string
		payload     string
		expectedRes map[string]interface{}
	}{
		{
			name:    "RegisterAccount",
			method:  "POST",
			apiPath: "/v1/accounts",
			payload: `{"username":"john"}`,
			expectedRes: map[string]interface{}{
				"username": "john",
			},
		},
		{
			name:    "GetAccount",
			method:  "GET",
			apiPath: "/v1/accounts/john",
			expectedRes: map[string]interface{}{
				"username": "john",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var resp *http.Response
			var err error

			if tc.method == "POST" {
				resp, err = c.PostJSON(tc.apiPath, tc.payload)
			} else if tc.method == "GET" {
				resp, err = c.Get(tc.apiPath)
			}

			require.NoError(t, err)
			require.Equal(t, http.StatusOK, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			var res map[string]interface{}
			require.NoError(t, json.Unmarshal(body, &res))

			assert.Equal(t, tc.expectedRes["username"], res["username"])
		})
	}
}

func setup(t *testing.T) *C {
	app, err := app.NewTestApp()
	if err != nil {
		t.Fatal(err)
	}

	if err := app.Dao.InitAll(); err != nil {
		t.Fatal(err)
	}
	if err := app.Dao.SetupTestDB(); err != nil {
		t.Fatal(err)
	}

	server := httptest.NewServer(NewRouter(app))

	return &C{
		App:    app,
		Server: server,
	}
}

// テストで使用するHTTPクライアントとサーバーの情報を保持する
type C struct {
	App    *app.App
	Server *httptest.Server
}

func (c *C) Close() {
	c.Server.Close()
}

func (c *C) PostJSON(apiPath string, payload string) (*http.Response, error) {
	return c.Server.Client().Post(c.asURL(apiPath), "application/json", bytes.NewReader([]byte(payload)))
}

func (c *C) Get(apiPath string) (*http.Response, error) {
	return c.Server.Client().Get(c.asURL(apiPath))
}

func (c *C) asURL(apiPath string) string {
	baseURL, _ := url.Parse(c.Server.URL)
	baseURL.Path = path.Join(baseURL.Path, apiPath)
	return baseURL.String()
}
