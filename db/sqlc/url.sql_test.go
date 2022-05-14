package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"shortURL/util"

	"github.com/stretchr/testify/require"
)

func createRandomURL(t *testing.T) Url {
	randomLongURL := util.RandomLongURL()
	require.NotEmpty(t, randomLongURL)

	arg := CreateURLParams{
		OriginUrl: randomLongURL,
		ShortUrl:  util.RandomString(6),
	}

	url, err := testQueries.CreateURL(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, url)

	require.Equal(t, arg.OriginUrl, url.OriginUrl)
	require.Equal(t, arg.ShortUrl, url.ShortUrl)
	require.NotZero(t, url.ID)
	require.NotZero(t, url.CreatedAt)

	return url
}

func TestQueries_CreateURL(t *testing.T) {
	arg := CreateURLParams{
		OriginUrl: util.RandomLongURL(),
		ShortUrl:  util.RandomString(6),
	}

	res, err := testQueries.CreateURL(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, res)

	require.Equal(t, arg.OriginUrl, res.OriginUrl)
	require.Equal(t, arg.ShortUrl, res.ShortUrl)
	require.NotZero(t, res.ID)
	require.NotZero(t, res.CreatedAt)
}

func TestGetURL(t *testing.T) {
	url1 := createRandomURL(t)
	url2, err := testQueries.GetURL(context.Background(), url1.ShortUrl)
	require.NoError(t, err)
	require.NotEmpty(t, url2)

	require.Equal(t, url1.ID, url2.ID)
	require.Equal(t, url1.OriginUrl, url2.OriginUrl)
	require.Equal(t, url1.ShortUrl, url2.ShortUrl)
	require.WithinDuration(t, url1.CreatedAt, url2.CreatedAt, time.Second)
}

func TestUpdateURL(t *testing.T) {
	url1 := createRandomURL(t)

	arg := UpdateURLParams{
		ID:       url1.ID,
		ShortUrl: util.RandomString(6),
	}

	url2, err := testQueries.UpdateURL(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, url2)

	require.Equal(t, url1.ID, url2.ID)
	require.Equal(t, url1.OriginUrl, url2.OriginUrl)
	require.Equal(t, arg.ShortUrl, url2.ShortUrl)
	require.WithinDuration(t, url1.CreatedAt, url2.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	url1 := createRandomURL(t)
	err := testQueries.DeleteURL(context.Background(), url1.ID)
	require.NoError(t, err)

	url2, err := testQueries.GetURL(context.Background(), url1.ShortUrl)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, url2)
}
