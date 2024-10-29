package util

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
	resp, err := http.Get(fmt.Sprintf("%s/api/info", baseAppUrl))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var appInfo AppInfo
	err = json.Unmarshal(body, &appInfo)
	if err != nil {
		return nil, err
	}

	return &appInfo, nil

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
