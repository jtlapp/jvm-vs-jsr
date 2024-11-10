package config

import (
	"runtime"
)

const (
	maxReservedPorts = 4
)

type PlatformConfig struct {
	ClientConfig
	AppName          string
	AppVersion       string
	AppConfig        map[string]interface{}
	CPUsPerNode      uint
	MaxReservedPorts uint
}

func GetPlatformConfig(clientConfig ClientConfig) (*PlatformConfig, error) {
	appInfo, err := GetAppInfo(clientConfig.BaseAppUrl)
	if err != nil {
		return nil, err
	}

	return &PlatformConfig{
		ClientConfig:     clientConfig,
		AppName:          appInfo.AppName,
		AppVersion:       appInfo.AppVersion,
		AppConfig:        appInfo.AppConfig,
		CPUsPerNode:      uint(runtime.NumCPU()),
		MaxReservedPorts: maxReservedPorts,
	}, nil
}
