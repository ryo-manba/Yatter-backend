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
)

func TestAccountRegistration(t *testing.T) {
	c := setup(t)
	defer c.Close()

	const apiPath = "/v1/accounts"
	testCases := []struct {
		name         string
		payload      string
		expectedCode int
		expectedRes  map[string]interface{}
	}{
		{
			name:         "正常系：作成できる",
			payload:      `{"username":"john"}`,
			expectedCode: http.StatusOK,
			expectedRes: map[string]interface{}{
				"username": "john",
			},
		},
		{
			name:         "異常系：すでにユーザーが存在する",
			payload:      `{"username":"john"}`,
			expectedCode: http.StatusInternalServerError, // TODO: 409 Conflictを返すようにする？
			expectedRes:  map[string]interface{}{},
		},
		{
			name:         "異常系：パラメータが不正",
			payload:      `{"name":"john"}`,
			expectedCode: http.StatusBadRequest,
			expectedRes:  map[string]interface{}{},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var resp *http.Response
			var err error

			resp, err = c.PostJSON(apiPath, tc.payload)

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCode, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			if tc.expectedCode == http.StatusOK {
				var res map[string]interface{}
				assert.NoError(t, json.Unmarshal(body, &res))
				assert.Equal(t, tc.expectedRes["username"], res["username"])
			}
		})
	}
}

func TestAccountGet(t *testing.T) {
	c := setup(t)
	defer c.Close()

	const apiPath = "/v1/accounts/"
	testCases := []struct {
		name         string
		pathParam    string
		expectedCode int
		expectedRes  map[string]interface{}
	}{
		{
			name:         "正常系：ユーザーが存在する",
			pathParam:    "test-user1",
			expectedCode: http.StatusOK,
			expectedRes: map[string]interface{}{
				"username": "test-user1",
			},
		},
		{
			name:         "正常系：ユーザーが存在しない",
			pathParam:    "notfound",
			expectedCode: http.StatusOK,
			expectedRes:  map[string]interface{}{},
		},
		{
			name:         "異常系：パラメータなし",
			pathParam:    "",
			expectedCode: http.StatusMethodNotAllowed,
			expectedRes:  map[string]interface{}{},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var resp *http.Response
			var err error

			resp, err = c.Get(apiPath + tc.pathParam)

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCode, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			if tc.expectedCode == http.StatusOK {
				var res map[string]interface{}
				assert.NoError(t, json.Unmarshal(body, &res))
				assert.Equal(t, tc.expectedRes["username"], res["username"])
			}
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
