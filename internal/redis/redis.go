package redis

import (
	"github.com/redis/go-redis/v9"
)

// localhost
var RedisClient = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

//docker
/* var redisClient = redis.NewClient(&redis.Options{
    Addr: "redis:6379",
}) */

//kubernetes
/* var redisClient = redis.NewClient(&redis.Options{
	Addr: "url-shortener-redis:6379",
}) */

var RedisNil = redis.Nil