package rds

import (
	"github.com/go-redis/redis/v9"
	"github.com/mvity/go-box/x"
	redis2 "github.com/mvity/go-quickstart/internal/dao/redis"
)

type near struct{}

// Near 附近地点
var Near near

// Set 设置附近地点
func (n *near) Set(tag string, key string, locs string) bool {
	rkey := redis2.RedisDataPrefix + "Near" + ":" + tag
	lng, lat := x.ParseLocation(locs)
	loc := &redis.GeoLocation{
		Name:      key,
		Longitude: lng,
		Latitude:  lat,
	}
	return redis2.Redis.GeoAdd(redis2.RedisContext, rkey, loc).Val() > 0
}

// Delete 删除附近地点
func (n *near) Delete(tag string, key string) {
	rkey := redis2.RedisDataPrefix + "Near" + ":" + tag
	redis2.Redis.ZRem(redis2.RedisContext, rkey, key)
}

// Query 查询附近地点
func (n *near) Query(tag string, locs string, meter int64, size int) []redis.GeoLocation {
	rkey := redis2.RedisDataPrefix + "Near" + ":" + tag
	lng, lat := x.ParseLocation(locs)

	query := &redis.GeoRadiusQuery{
		Radius:      x.ToFloat64(meter),
		Unit:        "m",
		WithCoord:   true,
		WithDist:    true,
		WithGeoHash: false,
		Count:       size,
		Sort:        "ASC",
		Store:       "",
		StoreDist:   "",
	}
	return redis2.Redis.GeoRadius(redis2.RedisContext, rkey, lng, lat, query).Val()

}
