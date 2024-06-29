package config

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"sync"
)

var (
	once          sync.Once
	mappingConfig *MappingConfig
)

type SftpConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type MappingConfig struct {
	Sftp struct {
		Dev       SftpConfig `mapstructure:"dev"`
		Uat       SftpConfig `mapstructure:"uat"`
		Pre       SftpConfig `mapstructure:"pre"`
		RemoteDir string     `mapstructure:"remoteDir"`
		LocalDir  string     `mapstructure:"localDir"`
		Export    string     `mapstructure:"export"`
	} `mapstructure:"sftp"`
}

func GetConfig() *MappingConfig {
	once.Do(func() {
		mappingConfig = loadConfig()
	})

	return mappingConfig
}

func GetEnvSiteData(envReq string) (config SftpConfig, err error) {
	switch envReq {
	case "dev":
		return mappingConfig.Sftp.Dev, nil
	case "uat":
		return mappingConfig.Sftp.Uat, nil
	case "pre":
		return mappingConfig.Sftp.Pre, nil
	default:
		return config, fmt.Errorf("invalid env value: %s", envReq)
	}
}

func loadConfig() *MappingConfig {
	viper.SetConfigFile("config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal().Err(err).Msg("Cannot load bineConfig")
	}

	var bineConfig MappingConfig
	if err := viper.Unmarshal(&bineConfig); err != nil {
		panic(fmt.Errorf("Error unmarshaling config: %s \n", err))
	}

	return &bineConfig
}
