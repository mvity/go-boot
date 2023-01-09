/*
 * Copyright © 2021 - 2023 vity <vityme@icloud.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 */

package rds

import (
	"github.com/go-redis/redis/v9"
	"github.com/mvity/go-box/x"
)

type near struct{}

// Near 附近地点
var Near near

// Set 设置附近地点
func (n *near) Set(tag string, key string, locs string) bool {
	rkey := RedisDataPrefix + "Near" + ":" + tag
	lng, lat := x.ParseLocation(locs)
	loc := &redis.GeoLocation{
		Name:      key,
		Longitude: lng,
		Latitude:  lat,
	}
	return Redis.GeoAdd(RedisContext, rkey, loc).Val() > 0
}

// Delete 删除附近地点
func (n *near) Delete(tag string, key string) {
	rkey := RedisDataPrefix + "Near" + ":" + tag
	Redis.ZRem(RedisContext, rkey, key)
}

// Query 查询附近地点
func (n *near) Query(tag string, locs string, meter int64, size int) []redis.GeoLocation {
	rkey := RedisDataPrefix + "Near" + ":" + tag
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
	return Redis.GeoRadius(RedisContext, rkey, lng, lat, query).Val()
}
