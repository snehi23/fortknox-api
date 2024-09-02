package main

import (
	"time"

	"github.com/patrickmn/go-cache"
)

func setUpCache() *cache.Cache {
	return cache.New(5*time.Second, 10*time.Second)
}
