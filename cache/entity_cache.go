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
	path    string
}

type EntityCache interface {
	ClearCache()
	UpdateCache([]model.Entity)
	GetCache() []model.Entity
	SetPath(string)
	//DeleteEntity(model.Entity)
}

func GetEntityCache() (EntityCache, error) {
	if EntityCacheInstance != nil {
		return EntityCacheInstance, nil
	}

	redisAddress := fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: config.RedisPassword,
		DB:       0,
	})

	err := redisClient.Ping(context.Background()).Err()
	if err != nil {
		return EntityCacheInstance, fmt.Errorf("failed to connect to redis: %s", err)
	}

	EntityCacheInstance = &entityCache{connect: redisClient}

	closer.CloseFunctions = append(closer.CloseFunctions, EntityCacheInstance.CloseEntityCache())
	return EntityCacheInstance, nil
}

func (cache *entityCache) CloseEntityCache() func() {
	return func() {
		if err := cache.connect.Close(); err != nil {
			log.Printf("failed while closing redis connection: %s", err)
			return
		}
		log.Println("entityCache closed successfully")
	}
}

func (cache *entityCache) ClearCache() {
	status := cache.connect.FlushAll(context.Background())
	if status.Err() != nil {
		log.Printf("failed to clean cache: %s", status.Err())
		return
	}
	log.Println("cache cleaned")
}

func (cache *entityCache) UpdateCache(entitySlice []model.Entity) {
	ctx := context.Background()
	marshalledEntity, err := json.Marshal(entitySlice)
	if err != nil {
		log.Printf("Failed to marshal entity while saving to cache: %s", err)
	}
	key := fmt.Sprintf("entity_%s", cache.path)
	data := cache.connect.Set(ctx, key, marshalledEntity, time.Second*3600)
	if data.Err() != nil {
		log.Printf("failed to update cache: %s", err)
	}

	//сохраняем значения в кэш, обходя их в цикле
	/*for _, entity := range entitySlice {
		marshalledTest, err := json.Marshal(entity.Test)
		if err != nil {
			log.Printf("Failed to marshal entity while saving to cache: %s", err)
		}
		cache.connect.Set(ctx, entity.Id, marshalledTest, time.Second*3600)
	}*/

}

func (cache *entityCache) GetCache() []model.Entity {
	var entities []model.Entity
	ctx := context.Background()

	//получаем ключи, чтобы получить значения в цикле
	/*cacheKeysSlice := cache.connect.Keys(ctx, "*").Val()
	log.Println("here keys from cache: ", cacheKeysSlice*/

	//забираем значения из кэша
	key := fmt.Sprintf("entity_%s", cache.path)

	//получаем кэш
	data := cache.connect.Get(ctx, key)
	if err := data.Err(); err != nil && err.Error() != "redis: nil" {
		log.Printf("ошибка в получениии кэша: %s", err)
	}

	//если кэш пустой, возвращаем nil слайс
	if data.Val() == "" {
		return entities
	}

	//если в кэше что-то лежит, парсим json в слайс, возвращаем его
	err := json.Unmarshal([]byte(data.Val()), &entities)

	if err != nil {
		log.Println("Failed to unmarshal json in cache: ", err)
	}
	/*for _, key := range cacheKeysSlice {
		var test model.Test
		err := json.Unmarshal([]byte(cache.connect.Get(ctx, key).Val()), &test)
		if err != nil {
			log.Println("Failed to unmarshal json in cache: ", err)
		}
		entities = append(entities, model.Entity{Id: key, Test: test})
	}*/
	return entities
}

func (cache *entityCache) SetPath(path string) {
	cache.path = path
}

/*
func (cache *entityCache) DeleteEntity(entity model.Entity) {
	key := fmt.Sprintf("entity_%s", cache.path)
	marshalledEntity, err := json.Marshal(entity)
	if err != nil {
		log.Println("failed to marshal entity into json while deleting cache: ", err)
	}
	err = cache.connect.SRem(context.Background(), key, marshalledEntity).Err()
	if err != nil {
		log.Println("error while deleting value from cache: ", err)
	}
}
*/
