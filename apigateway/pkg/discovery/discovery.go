package discovery

import (
	"fmt"
	"shortLink/apigateway/config"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	// "go.uber.org/zap"
)

// NamingClient 定义服务发现客户端接口
type NamingClient interface {
	RegisterInstance(param vo.RegisterInstanceParam) (bool, error)
	DeregisterInstance(param vo.DeregisterInstanceParam) (bool, error)
	SelectInstances(param vo.SelectInstancesParam) ([]model.Instance, error)
	SelectAllInstances(param vo.SelectAllInstancesParam) ([]model.Instance, error)
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
		ServiceName: config.GlobalConfig.Nacos.ServiceName,
		GroupName:   config.GlobalConfig.Nacos.GroupName,
		Weight:      float64(config.GlobalConfig.Nacos.Weight),
		Metadata:    config.GlobalConfig.Nacos.Metadata,
	}

	success, err := namingClient.RegisterInstance(param)
	if err != nil {
		return fmt.Errorf("注册服务失败: %v", err)
	}
	if !success {
		return fmt.Errorf("注册服务失败")
	}

	// logger.Log.Info("服务注册成功",
	// 	zap.String("service", param.ServiceName),
	// 	zap.String("ip", param.Ip),
	// 	zap.Uint64("port", param.Port))

	return nil
}

// DeregisterService 注销服务
func DeregisterService() error {
	param := vo.DeregisterInstanceParam{
		Ip:          config.GlobalConfig.Nacos.Ip,
		Port:        uint64(config.GlobalConfig.Nacos.Port),
		ServiceName: config.GlobalConfig.Nacos.ServiceName,
		GroupName:   config.GlobalConfig.Nacos.GroupName,
	}

	success, err := namingClient.DeregisterInstance(param)
	if err != nil {
		return fmt.Errorf("注销服务失败: %v", err)
	}
	if !success {
		return fmt.Errorf("注销服务失败")
	}

	// logger.Log.Info("服务注销成功",
	// 	zap.String("service", param.ServiceName),
	// 	zap.String("ip", param.Ip),
	// 	zap.Uint64("port", param.Port))

	return nil
}

// GetServiceInstances 获取服务实例
func GetServiceInstances(serviceName string) ([]model.Instance, error) {
	param := vo.SelectAllInstancesParam{
		ServiceName: serviceName,
		GroupName:   "DEFAULT_GROUP",
	}

	instances, err := namingClient.SelectAllInstances(param)
	if err != nil {
		return nil, fmt.Errorf("获取服务实例失败: %v", err)
	}

	return instances, nil
}
