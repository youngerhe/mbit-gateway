package nacos

import (
	"bytes"
	"fmt"
	"gateway/pkg/util"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	ENV_NACOS_ENDPOINTS = "ENV_NACOS_ENDPOINTS"
	ENV_APP             = "ENV_APP"
)

type Client struct {
	nameClient   naming_client.INamingClient
	configClient config_client.IConfigClient
	ServiceCache map[string][]model.Instance
	mu           sync.Mutex
	*ClientOptions
}

var NacosClient *Client

type ClientOptions struct {
	ServerAddr string
	Namespace  string
	DataId     string
	GroupName  string
}

func Init() {
	// 获取nacos节点
	endpoints, exist := os.LookupEnv(ENV_NACOS_ENDPOINTS)
	if !exist {
		panic("ENV_NACOS_ENDPOINTS not exist")
	}

	// 获取当前环境
	envApp, exist := os.LookupEnv(ENV_APP)
	if !exist {
		panic("ENV_APP not exist")
	}

	// 初始化Nacos配置
	client, err := NewNacosClient(&ClientOptions{
		ServerAddr: endpoints,
		Namespace:  envApp,
		DataId:     "gateway_config.yaml",
		GroupName:  "DEFAULT_GROUP",
	})

	if err != nil {
		panic(fmt.Sprintf("Nacos Init err:%v", err))
	}
	NacosClient = client
}

func RegisterInstance() {
	ip := util.LocalMulIPv4()
	err := NacosClient.RegisterInstance(Config.App.Name, ip[0], uint64(Config.App.Port))
	fmt.Println("RegisterInstance:", Config.App.Name, ip[0], uint64(Config.App.Port))
	if err != nil {
		panic(fmt.Sprintf("NacosClient.RegisterInstance err:%v", err))
		return
	}

	//go WatchService()
}

func DeregisterInstance() {
	ip := util.LocalMulIPv4()
	err := NacosClient.RegisterInstance(Config.App.Name, ip[0], uint64(Config.App.Port))
	if err != nil {
		panic(fmt.Sprintf("NacosClient.DeregisterInstance err:%v", err))
		return
	}
}

// WatchService 监听服务
func WatchService() {
	services := Config.App.Services
	for _, name := range services {

		err := NacosClient.WatchService(name, func(instances []model.Instance) {
			fmt.Printf("WatchService name:%s, data:%v", name, instances)
			NacosClient.mu.Lock()
			NacosClient.ServiceCache[name] = instances
			NacosClient.mu.Unlock()
		})
		if err != nil {
			fmt.Printf("WatchService name error:%s", name)
			log.Printf("WatchService name error:%s", name)
			continue
		}
	}
	ticker := time.NewTicker(2 * time.Second) // 设置定时器，每隔2秒触发一次
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Println("每隔2秒打印一次数据")
		}
	}
}

func NewNacosClient(co *ClientOptions) (*Client, error) {
	var serverConfigs []constant.ServerConfig
	values := strings.Split(co.ServerAddr, ",")
	for _, v := range values {
		vs := strings.Split(v, ":")
		if len(vs) != 2 {
			continue
		}
		port, _ := strconv.ParseUint(vs[1], 10, 64)
		serverConfigs = append(serverConfigs, *constant.NewServerConfig(vs[0], port))
	}
	clientConfig := *constant.NewClientConfig(
		constant.WithNamespaceId(co.Namespace),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir("/tmp/nacos/log"),
		constant.WithCacheDir("/tmp/nacos/cache"),
		constant.WithLogLevel("debug"),
	)
	param := vo.NacosClientParam{
		ClientConfig:  &clientConfig,
		ServerConfigs: serverConfigs,
	}

	nameClient, err := clients.NewNamingClient(param)
	if err != nil {
		return nil, err
	}
	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)

	if err != nil {
		return nil, err
	}
	return &Client{nameClient: nameClient, configClient: configClient, ServiceCache: map[string][]model.Instance{}, ClientOptions: co}, nil
}

func (nc *Client) GetConfig() (string, error) {

	// 获取配置
	content, err := nc.configClient.GetConfig(vo.ConfigParam{
		DataId: nc.DataId,
		Group:  nc.GroupName,
	})
	if err != nil {
		return "", err
	}
	go func() {
		err = nc.configClient.ListenConfig(vo.ConfigParam{
			DataId: nc.DataId,
			Group:  nc.GroupName,
			OnChange: func(namespace, group, dataId, data string) {
				dataByte := []byte(data)
				if err = viper.MergeConfig(bytes.NewBuffer(dataByte)); err != nil {
					fmt.Printf("viper MergeConfig err: %v", err)
				}
			},
		})
	}()
	return content, nil
}

func (nc *Client) RegisterInstance(serviceName, ip string, port uint64) error {
	param := vo.RegisterInstanceParam{
		Ip:          ip,
		Port:        port,
		ServiceName: serviceName,
		GroupName:   nc.GroupName,
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
	}
	success, err := nc.nameClient.RegisterInstance(param)
	if !success || err != nil {
		return err
	}
	return nil
}

func (nc *Client) DeregisterInstance(serviceName, ip string, port uint64) error {
	param := vo.DeregisterInstanceParam{
		Ip:          ip,
		Port:        port,
		ServiceName: serviceName,
		GroupName:   nc.GroupName,
		Ephemeral:   true,
	}
	success, err := nc.nameClient.DeregisterInstance(param)
	if success || err != nil {
		return err
	}
	return nil
}

func (nc *Client) GetAllInstances() (serviceList model.ServiceList, err error) {

	param := vo.GetAllServiceInfoParam{
		NameSpace: nc.Namespace,
		GroupName: nc.GroupName,
		PageNo:    1,
		PageSize:  100,
	}
	serviceList = model.ServiceList{}
	serviceList, err = nc.nameClient.GetAllServicesInfo(param)
	if err != nil {
		return serviceList, err
	}
	return serviceList, nil
}

func (nc *Client) WatchService(serviceName string, callback func(services []model.Instance)) error {
	err := nc.nameClient.Subscribe(&vo.SubscribeParam{
		ServiceName: serviceName,
		GroupName:   nc.GroupName,
		SubscribeCallback: func(services []model.Instance, err error) {
			fmt.Println("回调不了一点，不知道和redis集群那个问题是不是都是因为mac电脑的问题，日。。")
			callback(services)
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (nc *Client) GetOneHealthyInstance(serviceName string) (*model.Instance, error) {
	instance, err := nc.nameClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: serviceName,
		GroupName:   nc.GroupName,
	})
	if err != nil {
		return nil, err
	}
	return instance, err
}
