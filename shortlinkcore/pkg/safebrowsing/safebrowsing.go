package safebrowsing

import (
	"bytes"
	"encoding/json"
	"net/http"

	"shortLink/shortlinkcore/logger"

	"go.uber.org/zap"
)

const (
	apiKey = "AIzaSyDUzbJU3XMhMCpkzZkjv7kCxOKpCxbQHCg"
	apiURL = "https://safebrowsing.googleapis.com/v4/threatMatches:find"
)

type ThreatMatch struct {
	ThreatType      string `json:"threatType"`
	PlatformType    string `json:"platformType"`
	ThreatEntryType string `json:"threatEntryType"`
	Threat          struct {
		URL string `json:"url"`
	} `json:"threat"`
}

type RequestBody struct {
	Client struct {
		ClientID      string `json:"clientId"`
		ClientVersion string `json:"clientVersion"`
	} `json:"client"`
	ThreatInfo struct {
		ThreatTypes      []string `json:"threatTypes"`
		PlatformTypes    []string `json:"platformTypes"`
		ThreatEntryTypes []string `json:"threatEntryTypes"`
		ThreatEntries    []struct {
			URL string `json:"url"`
		} `json:"threatEntries"`
	} `json:"threatInfo"`
}

// CheckURL 检查URL是否安全
func CheckURL(url string) (bool, string, error) {
	reqBody := RequestBody{}
	reqBody.Client.ClientID = "shortlink"
	reqBody.Client.ClientVersion = "1.0.0"
	reqBody.ThreatInfo.ThreatTypes = []string{
		"MALWARE",
		"SOCIAL_ENGINEERING",
		"UNWANTED_SOFTWARE",
		"POTENTIALLY_HARMFUL_APPLICATION",
	}
	reqBody.ThreatInfo.PlatformTypes = []string{"ANY_PLATFORM"}
	reqBody.ThreatInfo.ThreatEntryTypes = []string{"URL"}
	reqBody.ThreatInfo.ThreatEntries = []struct {
		URL string `json:"url"`
	}{
		{URL: url},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		logger.Log.Error("序列化请求体失败", zap.String("error", err.Error()))
		return false, "", err
	}

	req, err := http.NewRequest("POST", apiURL+"?key="+apiKey, bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Log.Error("创建请求失败", zap.String("error", err.Error()))
		return false, "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Log.Error("请求失败", zap.String("error", err.Error()))
		return false, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		// 如果响应为空，说明URL是安全的
		if resp.ContentLength == 0 {
			return true, "", nil
		}

		var result struct {
			Matches []ThreatMatch `json:"matches"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			logger.Log.Error("解析响应失败", zap.String("error", err.Error()))
			return false, "", err
		}

		if len(result.Matches) > 0 {
			return false, result.Matches[0].ThreatType, nil
		}
	}

	return true, "", nil
}
