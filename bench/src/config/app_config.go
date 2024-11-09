package config

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type AppInfo struct {
	AppName    string                 `json:"appName"`
	AppVersion string                 `json:"appVersion"`
	AppConfig  map[string]interface{} `json:"appConfig"`
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

	// TODO: Delete this?
	// for key, value := range appInfo.AppConfig {
	// 	switch v := value.(type) {
	// 	case string:
	// 		fmt.Printf("%s: (string) %s\n", key, v)
	// 	case float64:
	// 		fmt.Printf("%s: (number) %f\n", key, v)
	// 	case bool:
	// 		fmt.Printf("%s: (bool) %t\n", key, v)
	// 	case map[string]interface{}:
	// 		fmt.Printf("%s: (object) %v\n", key, v)
	// 	case []interface{}:
	// 		fmt.Printf("%s: (array) %v\n", key, v)
	// 	case nil:
	// 		fmt.Printf("%s: null\n", key)
	// 	default:
	// 		fmt.Printf("%s: (unknown type) %v\n", key, v)
	// 	}
	// }
}
