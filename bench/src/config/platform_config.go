package config

import (
	"fmt"
	"runtime"

	"jvm-vs-jsr.jtlapp.com/benchmark/util"
)

type PlatformConfig struct {
	ClientConfig
	AppName           string
	AppVersion        string
	AppConfig         map[string]interface{}
	CPUsPerNode       uint
	InitialPortsInUse uint
}

func GetPlatformConfig(clientConfig ClientConfig) (*PlatformConfig, error) {
	appInfo, err := GetAppInfo(clientConfig.BaseAppUrl)
	if err != nil {
		return nil, err
	}

	resources := util.NewResourceStatus()
	if resources.TimeWaitPortsCount != 0 {
		return nil, fmt.Errorf("%d ports are in TIME_WAIT state", resources.TimeWaitPortsCount)
	}

	return &PlatformConfig{
		ClientConfig:      clientConfig,
		AppName:           appInfo.AppName,
		AppVersion:        appInfo.AppVersion,
		AppConfig:         appInfo.AppConfig,
		CPUsPerNode:       uint(runtime.NumCPU()),
		InitialPortsInUse: resources.EstablishedPortsCount,
	}, nil
}
