package main

import (
	"fmt"
	"log"

	"shortLink/previewservice/preview"
)

func main() {
	// 测试网页地址（推荐使用带有 Open Graph 标签的网页）
	targetURL := "https://cloud.tencent.com/developer/article/1776068"

	fmt.Println("🔍 开始提取网页信息:", targetURL)

	// 第一步：提取 OG 标签
	og, err := preview.FetchOGTagsA(targetURL)
	if err != nil {
		log.Fatalf("❌ 提取 OG 标签失败: %v", err)
	}
	content := fmt.Sprintf("content:%v", og)

	// 第二步：调用 LLM 摘要（可选，视情况启用）
	fmt.Println("\n🤖 使用 OpenAI 生成摘要...")
	summary, err := preview.GenerateSummary(content)
	if err != nil {
		log.Fatalf("❌ 摘要生成失败: %v", err)
	}
	fmt.Println("✅ 摘要生成:")
	fmt.Println(summary)
}

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"net/http"
// )

// func main() {
// 	url := "https://spark-api-open.xf-yun.com/v1/chat/completions"
// 	apiKey := "vlqaanDeQlNHXBvuKmYR:BDAQwEALbtjgORamizGp" // TODO：替换为你自己的 API Password

// 	// 构造请求体
// 	requestBody := map[string]interface{}{
// 		"model": "4.0Ultra",
// 		"messages": []map[string]string{
// 			{
// 				"role":    "user",
// 				"content": "你是谁",
// 			},
// 		},
// 	}

// 	jsonData, err := json.Marshal(requestBody)
// 	if err != nil {
// 		fmt.Println("❌ JSON 序列化失败:", err)
// 		return
// 	}

// 	// 创建请求
// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
// 	if err != nil {
// 		fmt.Println("❌ 请求创建失败:", err)
// 		return
// 	}
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Authorization", "Bearer "+apiKey)

// 	// 发送请求
// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		fmt.Println("❌ 请求发送失败:", err)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	// 读取响应
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		fmt.Println("❌ 响应读取失败:", err)
// 		return
// 	}

// 	fmt.Println("✅ 响应内容：")
// 	fmt.Println(string(body))
// }
