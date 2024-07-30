package cache

import (
	"context"
	"customers_kuber/closer"
	"customers_kuber/config"
	"customers_kuber/model"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

var EntityCacheInstance *entityCache

type entityCache struct {
	connect *redis.Client
}

type EntityCache interface {
	ClearCache()
	UpdateCache([]model.Entity)
	GetCache() []model.Entity
}

func GetEntityCache() EntityCache {
	if EntityCacheInstance != nil {
		return EntityCacheInstance
	}

	redisAddress := fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: config.RedisPassword,
		DB:       0,
	})

	err := redisClient.Ping(context.Background()).Err()
	if err != nil {
		log.Println("failed to connect to redis: ", err)
	}

	EntityCacheInstance = &entityCache{redisClient}

	closer.CloseFunctions = append(closer.CloseFunctions, EntityCacheInstance.CloseEntityCache())
	return EntityCacheInstance
}

func (cache *entityCache) CloseEntityCache() func() {
	return func() {
		if err := cache.connect.Close(); err != nil {
			log.Println("failed while closing redis connection:", err)
			return
		}
		log.Println("entityCache closed successfully")
	}
}

func (cache *entityCache) ClearCache() {
	cache.connect.FlushAll(context.Background())
}
func (cache *entityCache) UpdateCache(entitySlice []model.Entity) {
	ctx := context.Background()

	//сохраняем значения в кэш, обходя их в цикле
	for _, entity := range entitySlice {
		marshalledTest, err := json.Marshal(entity.Test)
		if err != nil {
			log.Printf("Failed to marshal entity while saving to cache: %s", err)
		}
		cache.connect.Set(ctx, entity.Id, marshalledTest, time.Second*3600)
	}

}
func (cache *entityCache) GetCache() []model.Entity {
	var entities []model.Entity
	ctx := context.Background()

	//получаем ключи, чтобы получить значения в цикле
	cacheKeysSlice := cache.connect.Keys(ctx, "*").Val()

	//забираем значения из кэша
	for _, key := range cacheKeysSlice {
		var test model.Test
		err := json.Unmarshal([]byte(cache.connect.Get(ctx, key).Val()), &test)
		if err != nil {
			log.Println("Failed to unmarshal json in cache: ", err)
		}
		entities = append(entities, model.Entity{Id: key, Test: test})
	}
	return entities
}
