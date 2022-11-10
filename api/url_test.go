package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	mockdb "shortURL/db/mock"
	db "shortURL/db/sqlc"
	"shortURL/util"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}

func TestServer_createShortURL(t *testing.T) {
	url := db.Url{
		OriginUrl: "https://" + util.RandomLongURL(),
	}

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockQuerier, redis *mockdb.MockRedisQuerier)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Success case",
			body: gin.H{
				"originUrl": url.OriginUrl,
			},
			buildStubs: func(store *mockdb.MockQuerier, redis *mockdb.MockRedisQuerier) {
				store.EXPECT().
					CreateURL(gomock.Any(), gomock.Any()).
					Times(1).
					Return(url, nil)

				redis.EXPECT().
					SetBloom(gomock.Any(), gomock.Any()).
					Times(1).
					Return(true, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "OriginUrl empty",
			body: gin.H{
				"originUrl": "",
			},
			buildStubs: func(store *mockdb.MockQuerier, redis *mockdb.MockRedisQuerier) {
				store.EXPECT().
					CreateURL(gomock.Any(), gomock.Any()).
					Times(0)

				redis.EXPECT().
					SetBloom(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "OriginUrl not acceptable Url",
			body: gin.H{
				"originUrl": "www.google.com",
			},
			buildStubs: func(store *mockdb.MockQuerier, redis *mockdb.MockRedisQuerier) {
				store.EXPECT().
					CreateURL(gomock.Any(), gomock.Any()).
					Times(0)

				redis.EXPECT().
					SetBloom(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"originUrl": url.OriginUrl,
			},
			buildStubs: func(store *mockdb.MockQuerier, redis *mockdb.MockRedisQuerier) {
				store.EXPECT().
					CreateURL(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Url{}, sql.ErrConnDone)

				redis.EXPECT().
					SetBloom(gomock.Any(), gomock.Any()).
					Times(1).
					Return(true, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// go mock
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockQueries := mockdb.NewMockQuerier(ctrl)
			mockRedis := mockdb.NewMockRedisQuerier(ctrl)
			tc.buildStubs(mockQueries, mockRedis)

			// start test server and send request
			server := NewServer(mockQueries, nil)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			apiUrl := "/short"
			request, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestServer_getRedirect(t *testing.T) {
	url := db.Url{
		OriginUrl: util.RandomLongURL(),
		ShortUrl:  util.RandomString(6),
	}

	testCases := []struct {
		name          string
		shortUrl      string
		buildStubs    func(store *mockdb.MockQuerier, redis *mockdb.MockRedisQuerier)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:     "Success case",
			shortUrl: url.ShortUrl,
			buildStubs: func(store *mockdb.MockQuerier, redis *mockdb.MockRedisQuerier) {
				store.EXPECT().
					GetURL(gomock.Any(), gomock.Eq(url.ShortUrl)).
					Times(1).
					Return(url, nil)

				redis.EXPECT().
					ExistBloom(gomock.Any(), gomock.Eq(url.ShortUrl)).
					Times(1).
					Return(true, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusMovedPermanently, recorder.Code)
			},
		},
		{
			name:     "shortURL too short",
			shortUrl: "000",
			buildStubs: func(store *mockdb.MockQuerier, redis *mockdb.MockRedisQuerier) {
				store.EXPECT().
					GetURL(gomock.Any(), gomock.Any()).
					Times(0)

				redis.EXPECT().
					ExistBloom(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:     "Not found",
			shortUrl: url.ShortUrl,
			buildStubs: func(store *mockdb.MockQuerier, redis *mockdb.MockRedisQuerier) {
				store.EXPECT().
					GetURL(gomock.Any(), gomock.Any()).
					Times(0)

				redis.EXPECT().
					ExistBloom(gomock.Any(), gomock.Eq(url.ShortUrl)).
					Times(1).
					Return(false, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:     "InternalError",
			shortUrl: url.ShortUrl,
			buildStubs: func(store *mockdb.MockQuerier, redis *mockdb.MockRedisQuerier) {
				store.EXPECT().
					GetURL(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Url{}, sql.ErrConnDone)

				redis.EXPECT().
					ExistBloom(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		// go mock
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockQueries := mockdb.NewMockQuerier(ctrl)
		mockRedis := mockdb.NewMockRedisQuerier(ctrl)
		tc.buildStubs(mockQueries, mockRedis)

		// start test server and send request
		server := NewServer(mockQueries, nil)
		recorder := httptest.NewRecorder()

		apiUrl := fmt.Sprintf("/%s", tc.shortUrl)
		request, err := http.NewRequest(http.MethodGet, apiUrl, nil)
		require.NoError(t, err)

		server.router.ServeHTTP(recorder, request)
		tc.checkResponse(recorder)
	}
}
