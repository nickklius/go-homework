package handler

import (
	"bytes"
	"go-homework/internal/checker"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	logger *zap.Logger
	ch     *checker.YandexSpellChecker
)

func TestMain(m *testing.M) {
	logger, _ = zap.NewDevelopment()
	defer logger.Sync()

	ch = checker.New()

	os.Exit(m.Run())
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)

	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp, string(respBody)
}

func TestHandler_checkSpellsBatchHandler(t *testing.T) {

	type want struct {
		statusCode   int
		contentType  string
		responseBody string
	}

	tests := []struct {
		name string
		body string
		path string
		want want
	}{
		{
			name: "success: checker return 200",
			body: "{\"texts\":[\"На лисной опушки распускаюца колоколчики, незабутки, шыповник.\"]}",
			path: "/",
			want: want{
				statusCode:   http.StatusOK,
				contentType:  "application/json; charset=utf-8",
				responseBody: "[\"На лесной опушке распускаются колокольчики, незабудки, шиповник.\"]",
			},
		},
		{
			name: "fail: wrong body format",
			body: "{\"texts\":}",
			path: "/",
			want: want{
				statusCode:   http.StatusInternalServerError,
				contentType:  "application/json; charset=utf-8",
				responseBody: "{\"error\":\"wrong format\"}",
			},
		},
		{
			name: "fail: wrong body format",
			body: "{\"text\":[\"На лисной опушки распускаюца колоколчики, незабутки, шыповник.\"]}",
			path: "/",
			want: want{
				statusCode:   http.StatusInternalServerError,
				contentType:  "application/json; charset=utf-8",
				responseBody: "{\"error\":\"wrong format\"}",
			},
		},
	}
	for _, tt := range tests {
		h := NewHandler(logger, ch)

		gin.SetMode(gin.ReleaseMode)

		router := gin.New()
		router.POST(tt.path, h.checkSpellsBatchHandler)

		ts := httptest.NewServer(router)
		defer ts.Close()

		resp, resultBody := testRequest(t, ts, http.MethodPost, tt.path, bytes.NewBuffer([]byte(tt.body)))
		defer resp.Body.Close()

		assert.Equal(t, tt.want.statusCode, resp.StatusCode)
		assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
		assert.Equal(t, tt.want.responseBody, resultBody)
	}
}
