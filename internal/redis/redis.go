// Package redis handles Redis client connection establishment and session management.
package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

var (
	ctx         = context.Background()
	redisClient *redis.Client
	DefaultTTL  = (60 * 60 * 24) * time.Second
)

// Initialize Redis Client Connection and return client
func InitRedisClient(host, port string) (*redis.Client, error) {
	// Connect to Redis
	fmt.Printf("Connecting to Redis at %s:%s\n", host, port)
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: "", // no password set
		// DB:       0,  // use default DB
	})

	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatal().Err(err).Msg("Error connecting to Redis")

		return nil, err
	}

	log.Info().Msgf("Successfully connected to Redis instance: %v", client.String())

	redisClient = client
	return client, nil
}

// Set a Key, Value pair in Redis
func SetCache(key string, value string, ttl time.Duration) error {
	if ttl == 0 {
		ttl = DefaultTTL
	}

	set, err := redisClient.SetNX(ctx, key, value, ttl).Result()
	if err != nil {
		log.Error().Err(err).Msg("Error setting key")
		return err
	}

	log.Info().Msgf("set: %v", set)

	return nil
}

// Get a Key, Value pair from Redis
func GetCache(key string) (string, error) {
	value, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		log.Error().Err(err).Msg("Error getting key")
		return "", nil
	}

	// log.Info().Msgf("value: %v", value)

	return value, nil
}

// Delete a Key, Value pair from Redis
func DeleteCache(key string) error {
	deleted, err := redisClient.Del(ctx, key).Result()
	if err != nil {
		log.Error().Err(err).Msg("Error deleting key")
		return err
	}

	log.Info().Msgf("deleted: %v", deleted)

	return nil
}
