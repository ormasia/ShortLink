package model

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// 初始化数据库
	if err := InitDB1("root:root@tcp(localhost:3306)/shortlink?charset=utf8mb4&parseTime=True&loc=Local"); err != nil {
		panic(err)
	}

	// 创建测试表
	if err := createTables(); err != nil {
		panic(err)
	}

	// 运行测试
	code := m.Run()

	// 清理测试数据
	if err := cleanup(); err != nil {
		panic(err)
	}

	os.Exit(code)
}

func cleanup() error {
	// 清空测试表
	_, err := db1.Exec("TRUNCATE TABLE id_generator")
	return err
}

func TestGenerateID(t *testing.T) {
	// 测试ID生成
	id, err := GenerateID()
	if err != nil {
		t.Fatalf("生成ID失败: %v", err)
	}

	if id <= 0 {
		t.Errorf("生成的ID应该大于0，实际得到: %d", id)
	}

	// 测试生成的ID是否递增
	nextID, err := GenerateID()
	if err != nil {
		t.Fatalf("生成ID失败: %v", err)
	}

	if nextID <= id {
		t.Errorf("第二个ID应该大于第一个ID，第一个: %d, 第二个: %d", id, nextID)
	}
}
