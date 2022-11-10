package redis

import (
	"context"
	"encoding/json"
	"time"

	db "shortURL/db/sqlc"

	redis "github.com/go-redis/redis/v8"
)

func NewRedisQuery(client *redis.Client) *RedisQueries {
	return &RedisQueries{client: client}
}

type RedisQueries struct {
	client *redis.Client
}

// 放置進布隆過濾器， true 代表過濾器內無相同資料且成功放入資料
func (q *RedisQueries) SetBloom(ctx context.Context, shortUrl string) (bool, error) {
	ret, err := q.client.Do(ctx, "BF.ADD", "url-filter", shortUrl).Result()
	if err != nil {
		return false, err
	}

	retInt64 := ret.(int64)

	// 1 代表過濾器內無相同資料且成功放入資料
	if retInt64 == 1 {
		return true, nil
	}

	return false, nil
}

// 檢查過濾器內是否有資料
func (r *RedisQueries) ExistBloom(ctx context.Context, shortUrl string) (bool, error) {
	ret, err := r.client.Do(ctx, "bf.exists", "url-filter", shortUrl).Result()
	if err != nil {
		return false, err
	}

	retInt64 := ret.(int64)
	// 1 代表過濾器內可能有資料，0 代表過濾器內絕對不可能有資料
	if retInt64 == 1 {
		return true, nil
	}

	return false, nil
}

// 設置資料
func (r *RedisQueries) SetData(ctx context.Context, shortUrl string, url db.Url) error {
	urlByte, err := json.Marshal(url)
	if err != nil {
		return err
	}

	_, err = r.client.Set(ctx, shortUrl, urlByte, time.Hour).Result()
	if err != nil {
		return err
	}

	return nil
}

// 取出資料
func (r *RedisQueries) GetData(ctx context.Context, shortUrl string) (db.Url, bool, error) {
	ret, err := r.client.Get(ctx, shortUrl).Result()
	if err != nil {
		if err == redis.Nil {
			return db.Url{}, false, nil
		}
		return db.Url{}, false, err
	}

	url := db.Url{}

	err = json.Unmarshal([]byte(ret), &url)
	if err != nil {
		return db.Url{}, false, err
	}

	return url, true, nil
}
