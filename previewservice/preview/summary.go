package preview

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// )

// type SummaryRequest struct {
// 	Model    string `json:"model"`
// 	Messages []struct {
// 		Role    string `json:"role"`
// 		Content string `json:"content"`
// 	} `json:"messages"`
// 	Stream bool `json:"stream"`
// }

// type SummaryResponse struct {
// 	Choices []struct {
// 		Message struct {
// 			Content string `json:"content"`
// 		} `json:"message"`
// 	} `json:"choices"`
// }

// // GenerateSummary 调用 OpenAI API 为网页描述生成摘要
// func GenerateSummary(content string) (string, error) {
// 	apiKey := "YjczNzY4ZTVmY2JkMjRhZTk3MzY0YTk3"
// 	if apiKey == "" {
// 		return "", fmt.Errorf("missing OPENAI_API_KEY")
// 	}

// 	reqData := map[string]interface{}{
// 		"model": "4.0Ultra",
// 		"messages": []map[string]string{
// 			{"role": "system", "content": "你是一个网页摘要助手。请总结以下内容："},
// 			{"role": "user", "content": content},
// 		},
// 		"stream": false,
// 	}

// 	reqBody, _ := json.Marshal(reqData)
// 	req, _ := http.NewRequest("POST", "https://spark-api-open.xf-yun.com/v1/chat/completions", bytes.NewBuffer(reqBody))
// 	req.Header.Set("Authorization", "Bearer "+apiKey)
// 	req.Header.Set("Content-Type", "application/json")

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return "", err
// 	}
// 	defer resp.Body.Close()

// 	var res SummaryResponse
// 	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
// 		return "", err
// 	}
// 	if len(res.Choices) > 0 {
// 		return res.Choices[0].Message.Content, nil
// 	}

// 	return "", fmt.Errorf("no summary returned")
// }

// const apiURL = "https://spark-api-open.xf-yun.com/v1/chat/completions"
// const apiKey = "vlqaanDeQlNHXBvuKmYR:BDAQwEALbtjgORamizGp" // 替换为实际值

// type ChatRequest struct {
// 	Model string `json:"model"`
// 	// User     string  `json:"user"`
// 	Messages []Message `json:"message"`
// 	Stream   bool      `json:"stream"`
// 	// User      string    `json:"user,omitempty"`
// 	// MaxTokens int       `json:"max_tokens,omitempty"`
// }

// type Message struct {
// 	Role    string `json:"role"`
// 	Content string `json:"content"`
// }

// type ChatResponse struct {
// 	Code    int `json:"code"`
// 	Choices []struct {
// 		Message Message `json:"message"`
// 	} `json:"choices"`
// }

// // 调用讯飞大模型生成摘要
// func GenerateSummary(text string) (string, error) {
// 	messages := []Message{
// 		{
// 			Role:    "system",
// 			Content: "你是知识渊博的助理",
// 		},
// 	}
// 	body := ChatRequest{
// 		Model:    "generalv3.5",
// 		Messages: messages,
// 		Stream:   false,
// 	}
// 	data, _ := json.Marshal(body)

// 	req, _ := http.NewRequest("POST", apiURL, bytes.NewBuffer(data))
// 	req.Header.Set("Authorization", "Bearer "+apiKey)
// 	req.Header.Set("Content-Type", "application/json")

// 	client := http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return "", err
// 	}
// 	defer resp.Body.Close()

// 	bodyBytes, _ := io.ReadAll(resp.Body)
// 	var response ChatResponse
// 	if err := json.Unmarshal(bodyBytes, &response); err != nil {
// 		return "", err
// 	}
// 	if response.Code != 0 || len(response.Choices) == 0 {
// 		return "", fmt.Errorf("API 调用失败: %s", string(bodyBytes))
// 	}
// 	return response.Choices[0].Message.Content, nil
// }

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GenerateSummary(content string) (string, error) {
	url := "https://spark-api-open.xf-yun.com/v1/chat/completions"
	apiKey := "vlqaanDeQlNHXBvuKmYR:BDAQwEALbtjgORamizGp" // TODO：替换为你自己的 API Password

	// 构造请求体
	requestBody := map[string]interface{}{
		"model": "4.0Ultra",
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "请帮我描述这段内容是什么方面的,二十个字以内;能让人快速了解其中的内容",
			},
			{
				"role":    "user",
				"content": content,
			},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("❌ JSON 序列化失败:", err)
		return "", err
	}

	// 创建请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("❌ 请求创建失败:", err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("❌ 请求发送失败:", err)
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("❌ 响应读取失败:", err)
		return "", err
	}

	fmt.Println("✅ 响应内容：")
	fmt.Println(string(body))
	return string(body), nil
}
