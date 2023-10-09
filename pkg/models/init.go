package models

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"time"

	"menu/config"

	"github.com/go-redis/redis/v8"
	"github.com/olivere/elastic/v7"
	"github.com/qiniu/qmgo"
	"github.com/spf13/viper"
)

var once sync.Once
var onceDBClientSync sync.Once

// var redisOnce sync.Once

func GetDB(cfg config.Provider) *qmgo.Database {
	var db *qmgo.Database
	once.Do(func() {
		db = initDB(cfg)
	})
	return db
}

func GetRedis(cfg config.Provider) *redis.Client {
	rdb := initRedis(cfg)
	// redisOnce.Do(func() {
	// 	rdb = initRedis(cfg)
	// })
	return rdb
}

func initDB(cfg config.Provider) *qmgo.Database {
	ctx := context.Background()
	username := cfg.GetString("mongo.username")
	password := cfg.GetString("mongo.password")
	host := cfg.GetString("mongo.host")
	port := cfg.GetString("mongo.port")
	dbname := cfg.GetString("mongo.dbname")

	URI := fmt.Sprintf("mongodb://%s:%s@%s:%s/?authSource=admin&directConnection=true", username, url.QueryEscape(password), host, port)
	timeout := cfg.GetInt64("mongo.conn_timeout")
	client, err := qmgo.NewClient(ctx, &qmgo.Config{
		Uri:              URI,
		ConnectTimeoutMS: &timeout,
	})
	if err != nil {
		panic(fmt.Sprintf("db conn err: %v", err))
	}

	err = client.Ping(5)
	if err != nil {
		panic(fmt.Sprintf("db ping err: %v", err))
	}

	return client.Database(dbname)

}

func initRedis(cfg config.Provider) *redis.Client {
	ctx := context.Background()

	host := fmt.Sprintf("%v:%v", cfg.GetString("redis.host"), cfg.GetString("redis.port"))
	password := cfg.GetString("redis.password")
	connTimeout := time.Duration(cfg.GetInt64("redis.conn_timeout"))
	dbIndex := cfg.GetInt("redis.db_index")

	rdb := redis.NewClient(&redis.Options{
		Addr:         host,
		Password:     password,
		DB:           dbIndex,
		PoolSize:     15,
		MinIdleConns: 1,
		DialTimeout:  time.Millisecond * connTimeout,
		ReadTimeout:  time.Millisecond * connTimeout,
		WriteTimeout: time.Millisecond * connTimeout,
	})

	err := rdb.Ping(ctx).Err()
	if err != nil {
		// panic(fmt.Sprintf("rdb ping err: %v", err))
		return nil
	}

	return rdb
}

// 初始化ES连接
func ElasticInit(urls []string, name string, pwd string, cfg *viper.Viper) *elastic.Client {
	var client *elastic.Client
	var err error
	if name != "" && pwd != "" {
		client, err = elastic.NewSimpleClient(elastic.SetURL(urls...), elastic.SetBasicAuth(name, pwd), elastic.SetSniff(true))
	} else {
		client, err = elastic.NewSimpleClient(elastic.SetURL(urls...), elastic.SetSniff(true))
	}
	if err != nil {
		fmt.Printf("Elastic client init ERROR: %v", err)
	}
	return client
}

func GetDBClient(cfg config.Provider) *qmgo.Client {
	var client *qmgo.Client

	onceDBClientSync.Do(func() {
		client = initDBClient(cfg)
	})

	return client
}

func initDBClient(cfg config.Provider) *qmgo.Client {
	ctx := context.Background()
	username := cfg.GetString("mongo.username")
	password := cfg.GetString("mongo.password")
	host := cfg.GetString("mongo.host")
	port := cfg.GetString("mongo.port")
	// dbname := cfg.GetString("mongo.dbname")

	// URI Reference: https://docs.mongodb.com/manual/reference/connection-string/
	URI := fmt.Sprintf("mongodb://%s:%s@%s:%s/?authSource=admin&directConnection=true", username, url.QueryEscape(password), host, port)
	timeout := cfg.GetInt64("mongo.conn_timeout")

	client, err := qmgo.NewClient(ctx, &qmgo.Config{
		Uri:              URI,
		ConnectTimeoutMS: &timeout,
	})
	if err != nil {
		panic(fmt.Sprintf("db conn err: %v", err))
	}

	err = client.Ping(5)
	if err != nil {
		panic(fmt.Sprintf("db ping err: %v", err))
	}

	return client

}
