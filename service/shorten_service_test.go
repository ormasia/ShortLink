package service

import (
	"testing"
)

func TestShorten(t *testing.T) {
	// 测试用例表
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "有效的URL",
			url:     "https://www.example.com",
			wantErr: false,
		},
		{
			name:    "无效的URL",
			url:     "not-a-url",
			wantErr: true,
		},
		{
			name:    "空URL",
			url:     "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shortKey, err := Shorten(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("Shorten() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && shortKey == "" {
				t.Error("生成了空的短链接")
			}
		})
	}
}
