package config

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type AppInfo struct {
	AppName    string    `json:"appName"`
	AppVersion string    `json:"appVersion"`
	AppConfig  AppConfig `json:"appConfig"`
}

type AppConfig struct {
	Server   map[string]interface{} `json:"server"`
	DbClient map[string]interface{} `json:"dbclient"`
}

func GetAppInfo(baseAppUrl string) (*AppInfo, error) {
	url := fmt.Sprintf("%s/api/info", baseAppUrl)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error getting app info: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code %d getting app info from %s", resp.StatusCode, url)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading app info response: %v", err)
	}

	var appInfo AppInfo
	err = json.Unmarshal(body, &appInfo)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling app info: %v", err)
	}

	return &appInfo, nil
}

func (ac *AppConfig) ToJsonString() (string, error) {
	bytes, err := json.Marshal(ac)
	if err != nil {
		return "", fmt.Errorf("error marshalling app config: %v", err)
	}
	return string(bytes), nil
}
