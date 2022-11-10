package main

import (
	"database/sql"
	"log"

	"shortURL/api"
	"shortURL/db/redis"
	db "shortURL/db/sqlc"
	"shortURL/util"

	goredis "github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	redisClient := goredis.NewClient(&goredis.Options{
		Addr: "localhost:6379",
	})

	store := db.NewQuery(conn)
	redisQuery := redis.NewRedisQuery(redisClient)

	server := api.NewServer(store, redisQuery)

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
