package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
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
		OriginUrl: util.RandomLongURL(),
	}

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockQuerier)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Success case",
			body: gin.H{
				"originUrl": url.OriginUrl,
			},
			buildStubs: func(store *mockdb.MockQuerier) {
				store.EXPECT().
					CreateURL(gomock.Any(), gomock.Any()).
					Times(1).
					Return(url, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Invalid request input",
			body: gin.H{
				"originUrl": "",
			},
			buildStubs: func(store *mockdb.MockQuerier) {
				store.EXPECT().
					CreateURL(gomock.Any(), gomock.Any()).
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
			buildStubs: func(store *mockdb.MockQuerier) {
				store.EXPECT().
					CreateURL(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Url{}, sql.ErrConnDone)
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
			tc.buildStubs(mockQueries)

			// start test server and send request
			server := NewServer(mockQueries)
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
