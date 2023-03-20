package config

import (
    "github.com/zeromicro/go-zero/rest"
    "github.com/zeromicro/go-zero/core/stores/cache"
    )


type Config struct {
    rest.RestConf
    Mysql struct{
        DataSource string
    }
    CacheRedis cache.CacheConf
    Auth      struct {
        AccessSecret string
        AccessExpire int64
    }
}