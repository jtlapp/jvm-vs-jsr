package config

import (
	"fmt"
	"os"
	"runtime"
	"sort"
)

const (
	dockerHostEnvVar = "DOCKER_HOST"
	appPortEnvVar    = "APP_PORT"
	maxReservedPorts = 4
)

type PlatformConfig struct {
	BaseAppUrl       string
	AppName          string
	AppVersion       string
	AppConfig        AppConfig
	CPUsPerNode      uint
	MaxReservedPorts uint
}

func GetPlatformConfig() (*PlatformConfig, error) {

	dockerHost := os.Getenv(dockerHostEnvVar)
	if dockerHost == "" {
		err := fmt.Errorf("%s environment variable is required", dockerHostEnvVar)
		if err != nil {
			return nil, err
		}
	}
	appPort := os.Getenv(appPortEnvVar)
	if appPort == "" {
		err := fmt.Errorf("%s environment variable is required", appPortEnvVar)
		if err != nil {
			return nil, err
		}
	}
	baseAppUrl := fmt.Sprintf("http://%s:%s", dockerHost, appPort)

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
	fmt.Println()

	fmt.Printf("  BaseAppUrl: %s\n", pc.BaseAppUrl)
	fmt.Printf("  AppName: %s\n", pc.AppName)
	fmt.Printf("  AppVersion: %s\n", pc.AppVersion)
	fmt.Printf("  CPUsPerNode: %d\n", pc.CPUsPerNode)
	fmt.Printf("  MaxReservedPorts: %d\n", pc.MaxReservedPorts)
	fmt.Println()

	fmt.Println("  AppConfig:")
	fmt.Println("    server:")
	printSortedStringMap(pc.AppConfig.Server)
	fmt.Println("    dbclient:")
	printSortedStringMap(pc.AppConfig.DbClient)
	fmt.Println()
}

func printSortedStringMap(stringMap map[string]interface{}) {
	keys := make([]string, 0, len(stringMap))
	for k := range stringMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("      %s: %v\n", k, stringMap[k])
	}
}
