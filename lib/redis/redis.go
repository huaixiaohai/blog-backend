package redis

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
)

var redisPool *redis.Pool

func RegisterRedisPool(dsn string, pw string, maxIdle, db int) {
	pool := &redis.Pool{
		MaxIdle: maxIdle,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", dsn, redis.DialPassword(pw))
			if err != nil {
				return nil, err
			}
			if _, err = c.Do("SELECT", db); err != nil {
				return nil, err
			}
			return c, err
		},
		//TestOnBorrow: func(c redis.Conn, t time.Time) error {
		//	rel, err := c.Do("PING")
		//	return err
		//},
	}
	conn := pool.Get()
	//defer conn.Close()

	if conn.Err() != nil {
		fmt.Printf("failed to connect redis on %s/%d, max idle conn: %d, err: %s", dsn, db, maxIdle, conn.Err().Error())
		panic(conn.Err())
	}

	fmt.Println(fmt.Sprintf("connect redis on %s/%d, max idle conn: %d", dsn, db, maxIdle))

	if rel, err := conn.Do("PING"); err == nil {
		fmt.Println("Redis PING ", rel)
	}
	_ = conn.Close()
	redisPool = pool
}

func DEL(key interface{}) {
	r := redisPool.Get()
	//defer r.Close()
	_, _ = r.Do("DEL", key)
	_ = r.Close()
}

func EXPIRE(key interface{}, ex int) {
	r := redisPool.Get()
	//defer r.Close()
	_, _ = r.Do("EXPIRE", key, ex)
	_ = r.Close()
}

func EXISTS(key interface{}) bool {
	r := redisPool.Get()
	//defer r.Close()
	exists, _ := redis.Bool(r.Do("EXISTS", key))
	_ = r.Close()
	return exists
}

func HGET(name string, key interface{}) (interface{}, error) {
	r := redisPool.Get()
	//defer r.Close()
	res, err := r.Do("HGET", name, key)
	_ = r.Close()
	return res, err
}

func HSET(name string, key, value interface{}) bool {
	r := redisPool.Get()
	//defer r.Close()
	_, err := r.Do("HSET", name, key, value)
	_ = r.Close()
	return err == nil
}

func SADD(set string, item interface{}) bool {
	r := redisPool.Get()
	//defer r.Close()
	_, err := r.Do("SADD", set, item)
	_ = r.Close()
	return err == nil
}

func SISMEMBER(set string, item interface{}) bool {
	r := redisPool.Get()
	//defer r.Close()
	in, _ := redis.Bool(r.Do("SISMEMBER", set, item))
	_ = r.Close()
	return in
}

func SREM(set string, item interface{}) bool {
	r := redisPool.Get()
	//defer r.Close()
	_, err := r.Do("SREM", set, item)
	_ = r.Close()
	return err == nil
}

func GET(key string) (string, error) {
	r := redisPool.Get()
	//defer r.Close()
	res, err := redis.String(r.Do("GET", key))
	_ = r.Close()
	return res, err

}

func SET(key string, value ...interface{}) bool {
	r := redisPool.Get()
	//defer r.Close()
	_, err := r.Do("SET", key, value)
	_ = r.Close()
	return err == nil
}

func INCR(key string) bool {
	r := redisPool.Get()
	//defer r.Close()
	_, err := r.Do("INCR", key)
	_ = r.Close()
	return err == nil
}

func SETEX(ex int, key, value interface{}) bool {
	r := redisPool.Get()
	//defer r.Close()
	_, err := r.Do("SETEX", key, ex, value)
	_ = r.Close()
	return err == nil
}
