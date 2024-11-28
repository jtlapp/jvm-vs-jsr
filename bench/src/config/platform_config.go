package config

import (
	"fmt"
	"os"
	"runtime"
)

const (
	baseAppUrlEnvVar = "BASE_APP_URL"
	maxReservedPorts = 4
)

type PlatformConfig struct {
	BaseAppUrl       string
	AppName          string
	AppVersion       string
	AppConfig        map[string]interface{}
	CPUsPerNode      uint
	MaxReservedPorts uint
}

func GetPlatformConfig() (*PlatformConfig, error) {

	baseAppUrl := os.Getenv(baseAppUrlEnvVar)
	if baseAppUrl == "" {
		err := fmt.Errorf("%s environment variable is required", baseAppUrlEnvVar)
		if err != nil {
			return nil, err
		}
	}

	appInfo, err := GetAppInfo(baseAppUrl)
	if err != nil {
		return nil, err
	}

	return &PlatformConfig{
		BaseAppUrl:       baseAppUrl,
		AppName:          appInfo.AppName,
		AppVersion:       appInfo.AppVersion,
		AppConfig:        appInfo.AppConfig,
		CPUsPerNode:      uint(runtime.NumCPU()),
		MaxReservedPorts: maxReservedPorts,
	}, nil
}

func (pc *PlatformConfig) Print() {
	fmt.Printf("BaseAppUrl: %s\n", pc.BaseAppUrl)
	fmt.Printf("AppName: %s\n", pc.AppName)
	fmt.Printf("AppVersion: %s\n", pc.AppVersion)
	fmt.Printf("CPUsPerNode: %d\n", pc.CPUsPerNode)
	fmt.Printf("MaxReservedPorts: %d\n", pc.MaxReservedPorts)
	fmt.Println("AppConfig:")
	for key, value := range pc.AppConfig {
		fmt.Printf("  %s: %v\n", key, value)
	}
}
