package database

import "jvm-vs-jsr.jtlapp.com/benchmark/config"

type AppKey struct {
	AppName    string
	AppVersion string
	AppConfig  config.AppConfig
}
