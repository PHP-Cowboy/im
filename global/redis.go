package global

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"

	"im/utils/str_util"
	"im/utils/timeutil"
)

type RedisCli struct {
	Cli *redis.Client
}

func (r *RedisCli) GetIncrNumberByKey(key string, padLength int) (string, error) {
	dateNumber := time.Now().Format(timeutil.DateNumberFormat)
	//rds key
	redisKey := key + dateNumber

	val, err := r.Incr(redisKey)
	if err != nil {
		return "", err
	}

	//设置过期时间
	err = r.Expire(redisKey, 24*time.Hour)

	number := strconv.Itoa(int(val))

	No := dateNumber + str_util.StrPad(number, padLength, "0", 0)

	return No, nil
}

func (r *RedisCli) Incr(key string) (int64, error) {
	return r.Cli.Incr(context.Background(), key).Result()
}

func (r *RedisCli) TTL(key string) (time.Duration, error) {
	return r.Cli.TTL(context.Background(), key).Result()
}

func (r *RedisCli) Set(key, val string, second int) (string, error) {
	return r.Cli.Set(context.Background(), key, val, time.Duration(second)*time.Second).Result()
}

func (r *RedisCli) Get(key string) (val string, err error) {
	val, err = r.Cli.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return
}

func (r *RedisCli) SetNx(key string, val string) error {
	ok, err := r.Cli.SetNX(context.Background(), key, val, 0).Result()

	if !ok {
		return errors.New("设置失败")
	}

	return err
}

// DEL 命令用于删除已存在的键。不存在的 key 会被忽略。
// 返回值: 被删除 key 的数量。
func (r *RedisCli) Del(key string) (int64, error) {
	return r.Cli.Del(context.Background(), key).Result()
}

// 设置过期时间
func (r *RedisCli) Expire(key string, d time.Duration) error {
	ok, err := r.Cli.Expire(context.Background(), key, d).Result()

	if !ok {
		return errors.New("设置过期时间失败")
	}

	return err
}

func (r *RedisCli) ExpireAt(key string, second int) error {
	ok, err := r.Cli.ExpireAt(context.Background(), key, time.Now().Add(time.Duration(second))).Result()

	if !ok {
		return errors.New("设置过期时间失败")
	}

	return err
}

/*
Redis Sadd 命令将一个或多个成员元素加入到集合中，已经存在于集合的成员元素将被忽略。
假如集合 key 不存在，则创建一个只包含添加的元素作成员的集合。
当集合 key 不是集合类型时，返回一个错误。
注意：在 Redis2.4 版本以前， SADD 只接受单个成员值。
*/
func (r *RedisCli) SAdd(key string, members []string) (err error) {
	_, err = r.Cli.SAdd(context.Background(), key, members).Result()

	return
}

/*
Redis Srem 命令用于移除集合中的一个或多个成员元素，不存在的成员元素会被忽略。
当 key 不是集合类型，返回一个错误。
在 Redis 2.4 版本以前， SREM 只接受单个成员值。
*/
func (r *RedisCli) SRem(key string, members []string) (err error) {
	_, err = r.Cli.SRem(context.Background(), key, members).Result()

	return
}

/*
Redis Smembers 命令返回集合中的所有的成员。 不存在的集合 key 被视为空集合。
*/
func (r *RedisCli) SMembers(key string) (err error) {
	_, err = r.Cli.SMembers(context.Background(), key).Result()

	return
}

// redis防重复点击
func (r *RedisCli) AntiRepeatedClick(key string, d int) (err error) {
	var value string

	value, err = r.Get(key)
	if err != nil {
		return
	}

	if value != "" {
		err = errors.New("处理中，请稍后重试")
		return
	}

	_, err = r.Set(key, "1", d)

	if err != nil {
		return
	}

	return
}
