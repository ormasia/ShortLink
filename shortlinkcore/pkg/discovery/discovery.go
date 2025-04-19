package discovery

import (
	"fmt"
	"shortLink/shortlinkcore/config"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

// NamingClient 定义服务发现客户端接口
type NamingClient interface {
	RegisterInstance(param vo.RegisterInstanceParam) (bool, error)
	DeregisterInstance(param vo.DeregisterInstanceParam) (bool, error)
	SelectInstances(param vo.SelectInstancesParam) ([]model.Instance, error)
}

// 实例化服务发现客户端接口
var namingClient NamingClient

// InitNamingClient 初始化服务发现客户端
func InitNamingClient() error {
	serverConfigs := []constant.ServerConfig{
		*constant.NewServerConfig("127.0.0.1", 8848),
	}

	// 创建 Nacos 客户端配置
	clientConfig := *constant.NewClientConfig(
		constant.WithNamespaceId(config.GlobalConfig.Nacos.Namespace),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogLevel("info"),
	)

	var err error
	namingClient, err = clients.NewNamingClient(vo.NacosClientParam{
		ClientConfig:  &clientConfig,
		ServerConfigs: serverConfigs,
	})
	if err != nil {
		return fmt.Errorf("创建服务发现客户端失败: %v", err)
	}

	return nil
}

// RegisterService 注册服务
func RegisterService() error {
	param := vo.RegisterInstanceParam{
		Ip:          config.GlobalConfig.Nacos.Ip,
		Port:        uint64(config.GlobalConfig.Nacos.Port),
		ServiceName: "shortlink-service",
		GroupName:   "DEFAULT_GROUP",
		Weight:      1,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
	}

	success, err := namingClient.RegisterInstance(param)
	if err != nil {
		return fmt.Errorf("注册服务失败: %v", err)
	}
	if !success {
		return fmt.Errorf("注册服务失败")
	}

	return nil
}

// DeregisterService 注销服务
func DeregisterService() error {
	param := vo.DeregisterInstanceParam{
		Ip:          config.GlobalConfig.Nacos.Ip,
		Port:        uint64(config.GlobalConfig.Nacos.Port),
		ServiceName: "shortlink-service",
		GroupName:   "DEFAULT_GROUP",
	}

	success, err := namingClient.DeregisterInstance(param)
	if err != nil {
		return fmt.Errorf("注销服务失败: %v", err)
	}
	if !success {
		return fmt.Errorf("注销服务失败")
	}

	return nil
}
