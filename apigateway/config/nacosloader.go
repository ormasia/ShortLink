package config

import (
	"bytes"
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
)

func InitConfigFromNacos() error {
	serverConfigs := []constant.ServerConfig{
		*constant.NewServerConfig("nacos", 8848), // Nacos 地址
	}
	// // 打印连接信息
	// fmt.Printf("Nacos 服务器配置: %+v\n", serverConfigs)

	// 创建 Nacos 客户端配置
	clientConfig := *constant.NewClientConfig(
		constant.WithNamespaceId(""), // 空为 public
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),   // ✅ 启动不读缓存
		constant.WithUpdateCacheWhenEmpty(false), // ✅ 读不到时也不兜底
		constant.WithLogLevel("info"),
	)

	client, err := clients.NewConfigClient(vo.NacosClientParam{
		ClientConfig:  &clientConfig,
		ServerConfigs: serverConfigs,
	})
	if err != nil {
		return fmt.Errorf("❌ 创建 Nacos 客户端失败: %v", err)
	}
	fmt.Println("✅ Nacos 客户端创建成功")

	// 获取配置内容
	configParam := vo.ConfigParam{
		DataId: "docker", // 移除 .yaml 后缀
		Group:  "DEFAULT_GROUP",
	}
	fmt.Printf("正在获取配置: DataId=%s, Group=%s\n", configParam.DataId, configParam.Group)

	content, err := client.GetConfig(configParam)
	if err != nil {
		return fmt.Errorf("❌ 获取配置失败: %v", err)
	}
	fmt.Printf("获取到的配置内容:\n%s\n", content)

	// 用 viper 读取 YAML 内容
	viper.SetConfigType("yaml")
	err = viper.ReadConfig(bytes.NewBufferString(content))
	if err != nil {
		return fmt.Errorf("❌ viper 解析失败: %v", err)
	}

	// 绑定到结构体
	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		return fmt.Errorf("❌ 配置绑定结构体失败: %v", err)
	}
	// 热更新监听
	err = client.ListenConfig(vo.ConfigParam{
		DataId: "shortlink",
		Group:  "DEFAULT_GROUP",
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("✅ Nacos 配置更新！重新加载中...")
			_ = viper.ReadConfig(bytes.NewBufferString(data))
			_ = viper.Unmarshal(&GlobalConfig)
		},
	})
	if err != nil {
		return fmt.Errorf("❌ 监听配置失败: %v", err)
	}

	return nil
}
