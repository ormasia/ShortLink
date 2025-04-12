package config

import (
	"strings"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
)

// Config 支付服务配置结构
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Nacos    NacosConfig    `mapstructure:"nacos"`
	Payment  PaymentConfig  `mapstructure:"payment"`
	Order    OrderConfig    `mapstructure:"order"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port int `mapstructure:"port"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// NacosConfig Nacos配置
type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      uint64 `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	Group     string `mapstructure:"group"`
	DataID    string `mapstructure:"dataid"`
}

// PaymentConfig 支付配置
type PaymentConfig struct {
	Alipay    AlipayConfig `mapstructure:"alipay"`
	WechatPay WechatConfig `mapstructure:"wechat"`
}

// OrderConfig 订单服务配置
type OrderConfig struct {
	OrderService struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"order_service"`
}

// AlipayConfig 支付宝配置
type AlipayConfig struct {
	AppID      string `mapstructure:"app_id"`
	PrivateKey string `mapstructure:"private_key"`
	PublicKey  string `mapstructure:"public_key"`
	NotifyURL  string `mapstructure:"notify_url"`
	ReturnURL  string `mapstructure:"return_url"`
	IsSandbox  bool   `mapstructure:"is_sandbox"`
}

// WechatConfig 微信支付配置
type WechatConfig struct {
	AppID     string `mapstructure:"app_id"`
	MchID     string `mapstructure:"mch_id"`
	Key       string `mapstructure:"key"`
	CertFile  string `mapstructure:"cert_file"`
	KeyFile   string `mapstructure:"key_file"`
	NotifyURL string `mapstructure:"notify_url"`
	IsSandbox bool   `mapstructure:"is_sandbox"`
}

// LoadConfig 加载配置，优先从本地文件加载，然后尝试从Nacos获取并合并配置
func LoadConfig() (*Config, error) {
	clientConfig := constant.ClientConfig{
		NamespaceId:         "public",
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "debug",
	}

	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: "localhost",
			Port:   8848,
		},
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		// 如果Nacos客户端创建失败，仅使用本地配置
		return loadLocalConfig()
	}

	// 从本地文件加载默认配置
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")
	v.AddConfigPath("../config")
	v.AddConfigPath("../../config")

	if readErr := v.ReadInConfig(); readErr != nil {
		// 如果找不到配置文件，使用默认配置
		v.SetDefault("server.port", 8083)
		v.SetDefault("database.host", "localhost")
		v.SetDefault("database.port", 3306)
		v.SetDefault("database.username", "root")
		v.SetDefault("database.password", "password")
		v.SetDefault("database.dbname", "shortlink")
		v.SetDefault("redis.host", "localhost")
		v.SetDefault("redis.port", 6379)
		v.SetDefault("redis.password", "")
		v.SetDefault("redis.db", 0)
	}

	// 尝试从Nacos获取配置
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: "payment-service",
		Group:  "DEFAULT_GROUP",
	})

	if err == nil && content != "" {
		// 解析Nacos配置
		v2 := viper.New()
		v2.SetConfigType("yaml")
		v2.ReadConfig(strings.NewReader(content))

		// 合并配置
		for _, key := range v2.AllKeys() {
			v.Set(key, v2.Get(key))
		}
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// loadLocalConfig 仅从本地加载配置，当Nacos不可用时使用
func loadLocalConfig() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")
	v.AddConfigPath("../config")
	v.AddConfigPath("../../config")

	if err := v.ReadInConfig(); err != nil {
		// 如果找不到配置文件，使用默认配置
		v.SetDefault("server.port", 8083)
		v.SetDefault("database.host", "localhost")
		v.SetDefault("database.port", 3306)
		v.SetDefault("database.username", "root")
		v.SetDefault("database.password", "password")
		v.SetDefault("database.dbname", "shortlink")
		v.SetDefault("redis.host", "localhost")
		v.SetDefault("redis.port", 6379)
		v.SetDefault("redis.password", "")
		v.SetDefault("redis.db", 0)
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
