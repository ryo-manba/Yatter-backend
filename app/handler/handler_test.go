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
	"yatter-backend-go/app/domain/object"

	"github.com/stretchr/testify/assert"
)

/// Accounts
func TestAccount_Registration(t *testing.T) {
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
			name:         "異常系：不正なパラメータ",
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

func TestAccount_Get(t *testing.T) {
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
			name:         "異常系：ユーザーが存在しない",
			pathParam:    "notfound",
			expectedCode: http.StatusNotFound,
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

/// status
func TestStatus_Create(t *testing.T) {
	c := setup(t)
	defer c.Close()

	const apiPath = "/v1/statuses"
	testCases := []struct {
		name         string
		username     string
		payload      string
		expectedCode int
		expectedRes  map[string]interface{}
	}{
		{
			name:         "正常系：ステータスを作成できる",
			username:     "test-user1",
			payload:      `{"status": "ピタ ゴラ スイッチ♪", "media_ids": [0]}`,
			expectedCode: http.StatusOK,
			expectedRes: map[string]interface{}{
				"content": "ピタ ゴラ スイッチ♪",
			},
		},
		{
			name:         "異常系：認証できない",
			username:     "",
			payload:      `{"status": "ピタ ゴラ スイッチ♪", "media_ids": [0]}`,
			expectedCode: http.StatusUnauthorized,
			expectedRes:  map[string]interface{}{},
		},
		{
			name:         "異常系：不正なパラメータ",
			username:     "test-user1",
			payload:      `{"data":"ピタ ゴラ スイッチ♪", "media_ids": [0]}`,
			expectedCode: http.StatusBadRequest,
			expectedRes:  map[string]interface{}{},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var resp *http.Response
			var err error

			// 認証する
			resp, err = c.PostJSONWithAuth(apiPath, tc.payload, tc.username)

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCode, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			if tc.expectedCode == http.StatusOK {
				var res map[string]interface{}
				assert.NoError(t, json.Unmarshal(body, &res))
				assert.Equal(t, tc.expectedRes["content"], res["content"])
			}
		})
	}
}

func TestStatus_Get(t *testing.T) {
	c := setup(t)
	defer c.Close()

	const apiPath = "/v1/statuses/"
	testCases := []struct {
		name         string
		pathParam    string
		expectedCode int
		expectedRes  map[string]interface{}
	}{
		{
			name:         "正常系：ステータスが存在する",
			pathParam:    "1",
			expectedCode: http.StatusOK,
			expectedRes: map[string]interface{}{
				"content": "Test content for user 1",
			},
		},
		{
			name:         "異常系：ステータスが存在しない",
			pathParam:    "10000",
			expectedCode: http.StatusNotFound,
			expectedRes:  map[string]interface{}{},
		},
		{
			name:         "異常系：パラメータが数値ではない",
			pathParam:    "hoge",
			expectedCode: http.StatusBadRequest,
			expectedRes:  map[string]interface{}{},
		},
		{
			name:         "異常系：パラメータなし",
			pathParam:    "",
			expectedCode: http.StatusUnauthorized,
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
				assert.Equal(t, tc.expectedRes["content"], res["content"])
			}
		})
	}
}

func TestStatus_Delete(t *testing.T) {
	c := setup(t)
	defer c.Close()

	const apiPath = "/v1/statuses/"
	testCases := []struct {
		name         string
		username     string
		pathParam    string
		expectedCode int
		expectedRes  map[string]interface{}
	}{
		{
			name:         "正常系：ステータスを削除できる",
			username:     "test-user1",
			pathParam:    "1",
			expectedCode: http.StatusOK,
			expectedRes:  map[string]interface{}{},
		},
		{
			name:         "異常系：ステータスが存在しない",
			username:     "test-user1",
			pathParam:    "10000",
			expectedCode: http.StatusNotFound,
			expectedRes:  map[string]interface{}{},
		},
		{
			name:         "異常系：ステータス作成者とユーザが一致しない",
			username:     "test-user1",
			pathParam:    "5",
			expectedCode: http.StatusBadRequest,
			expectedRes:  map[string]interface{}{},
		},
		{
			name:         "異常系：認証できない",
			username:     "",
			pathParam:    "10000",
			expectedCode: http.StatusUnauthorized,
			expectedRes:  map[string]interface{}{},
		},
		{
			name:         "異常系：パラメータが数値ではない",
			username:     "test-user1",
			pathParam:    "hoge",
			expectedCode: http.StatusBadRequest,
			expectedRes:  map[string]interface{}{},
		},
		{
			name:         "異常系：パラメータなし",
			username:     "test-user1",
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

			// 認証する
			resp, err = c.DeleteJSONWithAuth(apiPath+tc.pathParam, tc.username)

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCode, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			if tc.expectedCode == http.StatusOK {
				var res map[string]interface{}
				assert.NoError(t, json.Unmarshal(body, &res))
				assert.Equal(t, tc.expectedRes["content"], res["content"])
			}
		})
	}
}

func TestTimeline_PublicGet(t *testing.T) {
	c := setup(t)
	defer c.Close()

	const apiPath = "/v1/timelines/public"
	testCases := []struct {
		name         string
		query        string
		limit        int
		expectedCode int
	}{
		{
			name:         "正常系：タイムラインが存在する",
			query:        "?only_media=1&since_id=1&max_id=10&limit=5",
			expectedCode: http.StatusOK,
		},
		{
			name:         "正常系：タイムラインが存在しない",
			query:        "?only_media=1&since_id=450&max_id=500&limit=5",
			expectedCode: http.StatusOK,
		},
		{
			name:         "正常系：パラメータの指定なし",
			query:        "",
			expectedCode: http.StatusOK,
		},
		{
			name:         "異常系：max_idが負の値",
			query:        "?only_media=1&max_id=-1&since_id=1&limit=5",
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var resp *http.Response
			var err error

			resp, err = c.GetWithQuery(apiPath, tc.query)

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCode, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			if tc.expectedCode == http.StatusOK {
				var res []*object.Status
				assert.NoError(t, json.Unmarshal(body, &res))
				length := len(res)
				if length > 0 {
					assert.NotEqual(t, nil, res[0].Account)
					assert.NotEqual(t, nil, res[0].Content)
				}
			}
		})
	}
}

/// utils
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

func (c *C) GetWithQuery(apiPath, query string) (*http.Response, error) {
	reqURL := c.asURLWithQuery(apiPath, query)
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	return c.Server.Client().Do(req)
}

func (c *C) asURLWithQuery(apiPath, query string) string {
	baseURL, _ := url.Parse(c.Server.URL)
	baseURL.Path = path.Join(baseURL.Path, apiPath)
	baseURL.RawQuery = query
	return baseURL.String()
}

func (c *C) PostJSONWithAuth(apiPath string, payload string, username string) (*http.Response, error) {
	req, err := http.NewRequest("POST", c.asURL(apiPath), bytes.NewReader([]byte(payload)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authentication", "username "+username)
	return c.Server.Client().Do(req)
}

func (c *C) DeleteJSONWithAuth(apiPath string, username string) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", c.asURL(apiPath), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("Authentication", "username "+username)
	return c.Server.Client().Do(req)
}
