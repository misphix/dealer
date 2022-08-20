package main

import (
	"context"
	"dealer/internal/configmanager"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func newMySQL(config configmanager.DatabaseConfig) (*gorm.DB, error) {
	return gorm.Open(mysql.New(mysql.Config{
		DSN:                       config.DSN,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{})
}

func newCache(config configmanager.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: config.Address,
	})

	ctx, cancel := context.WithTimeout(context.Background(), config.DialTimeout)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}
