package common

import (
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"strings"
)

var Pool redis.Pool

//RedisSecret redis密钥
//const RedisSecret = "B9q1T0.`$Q83G;6h"

//RedisNetWork redis网路配置
const RedisNetWork = "tcp"

//RedisInit redis初始化配置
func RedisInit() {
	Pool = redis.Pool{
		MaxIdle:     16,
		MaxActive:   32,
		IdleTimeout: 120,
		Dial: func() (redis.Conn, error) {
			return redis.Dial(RedisNetWork, CONFIG.Redis.Address)
		},
	}
	if err := pingRedis(); err != nil {
		LogFatalf("Redis connection failed.", logrus.Fields{
			"err": err,
		})
		panic(fmt.Sprintf("Redis connection failed. Error is %s", err.Error()))
	}
}
func pingRedis() (err error) {
	conn := Pool.Get()
	var msg string
	msg, err = redis.String(conn.Do("ping"))
	if err != nil {
		LogErrorf("pingRedis Error may be redis information Error", logrus.Fields{"redis": CONFIG.Redis, "err": err})
	} else if strings.ToLower(msg) != "pong" {
		err = errors.New("ping redis get response error")
		LogErrorf("pingRedis Get wrong message", logrus.Fields{"redis": CONFIG.Redis, "msg": msg})
	} else {
		LogInfof("ping Redis Success", logrus.Fields{"msg": msg, "address": CONFIG.Redis.Address})
	}
	defer conn.Close()
	return
}
