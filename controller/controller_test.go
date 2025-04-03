package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/shorten", ShortenURL)
	return r
}

func TestShortenURL(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name       string
		body       interface{}
		wantStatus int
	}{
		{
			name: "有效的URL",
			body: map[string]string{
				"url": "https://www.example.com",
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "无效的JSON",
			body:       "invalid json",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "缺少URL字段",
			body:       map[string]string{},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest("POST", "/shorten", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("状态码 = %v, 期望 %v", w.Code, tt.wantStatus)
			}
		})
	}
}
