// Package middleware contains custom middlewares for the application.
package middleware

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/mexirica/chi-template/internal/helpers"
	"github.com/mexirica/chi-template/internal/redis"
	"github.com/rs/zerolog/log"
)

// PrepareRouteKey generates a unique key for the request route, method, and query.
func PrepareRouteKey(r *http.Request) (string, error) {
	return r.Method + "." + r.URL.Path + "." + r.URL.RawQuery, nil
}

// PrepareCacheKey returns a base64 encoded string of the payload with the route and method prepended.
func PrepareCacheKey(payload any, routeKey string) (string, error) {
	stringPayload := ""
	byteData, err := io.ReadAll(payload.(io.Reader))
	if err != nil {
		log.Fatal().Err(err).Msg("Error reading payload")
		return "", err
	}
	stringPayload = string(byteData)
	encodedPayload := base64.StdEncoding.EncodeToString([]byte(stringPayload))
	encodedPayload = routeKey + "." + encodedPayload
	return encodedPayload, nil
}

// StringifyResponse returns a stringified version of the response object.
func StringifyResponse(response any) (string, error) {
	stringResponse, err := json.Marshal(response)
	if err != nil {
		log.Fatal().Err(err).Msg("Error marshalling response")
		return "", err
	}
	return string(stringResponse), nil
}

// SaveToCache saves the response to cache after stringifying it and returns the cache key.
func SaveToCache(r *http.Request, response any) (string, error) {
	routeKey, err := PrepareRouteKey(r)
	if err != nil {
		log.Error().Err(err).Msg("Error preparing route key")
		return "", err
	}
	cacheKey, err := PrepareCacheKey(r.Body, routeKey)
	if err != nil {
		log.Error().Err(err).Msg("Error preparing cache key")
		return "", err
	}
	stringResponse, err := StringifyResponse(response)
	if err != nil {
		log.Error().Err(err).Msg("Error stringifying response")
		return "", err
	}
	err = redis.SetCache(cacheKey, stringResponse, 0)
	if err != nil {
		log.Error().Err(err).Msg("Error saving to cache")
		return "", err
	}
	return cacheKey, nil
}

// CachedResponseToJSON retrieves the cached response and converts it to JSON.
func CachedResponseToJSON(cacheKey string) ([]map[string]any, error) {
	cachedResponse, err := redis.GetCache(cacheKey)
	if err != nil {
		log.Error().Err(err).Msg("Error getting from cache")
		return nil, nil
	}
	if cachedResponse == "" {
		log.Info().Msgf("No cached response found for %s", cacheKey)
		return nil, nil
	}
	var cachedResponseJSON []map[string]any
	err = json.Unmarshal([]byte(cachedResponse), &cachedResponseJSON)
	if err != nil {
		log.Error().Err(err).Msg("Error unmarshalling cached response")
		return nil, err
	}
	return cachedResponseJSON, nil
}

// CacheMiddleware is a Redis cache middleware for HTTP handlers with a configurable TTL.
func CacheMiddleware(ttl time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if ttl == 0 {
				ttl = redis.DefaultTTL
			}
			routeKey, err := PrepareRouteKey(r)
			if err != nil {
				log.Fatal().Err(err).Msg("Error preparing route key")
				return
			}
			cacheKey, err := PrepareCacheKey(r.Body, routeKey)
			if err != nil {
				log.Fatal().Err(err).Msg("Error preparing cache key")
				return
			}
			cachedResponseJSON, err := CachedResponseToJSON(cacheKey)
			if err != nil {
				log.Fatal().Err(err).Msg("Error getting cached response")
				return
			}
			if cachedResponseJSON == nil {
				next.ServeHTTP(w, r)
			}
			if len(cachedResponseJSON) > 0 {
				helpers.WriteJSON(w, http.StatusOK, cachedResponseJSON)
				return
			}
		})
	}
}
