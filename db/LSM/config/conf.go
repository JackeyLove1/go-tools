package config

import "sync"

type Config struct {
    DataDir         string
    Level0MaxSum    int
    SStableNums     int
    KvNumsThreshold int
    CheckInterval   int
}

var (
    once   sync.Once
    config *Config
)

func InitConfig() {
    once.Do(func() {
        config = &Config{}
    })
}

func GetConfig() *Config {
    return config
}
