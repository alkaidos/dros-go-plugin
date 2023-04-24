package plugins

import (
	"github.com/isyscore/isc-gobase/logger"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
	"sync"
)

var PluginConfig = &PluginConf{}
var exist bool
var loadLock sync.Mutex

func init() {
	loadLock.Lock()
	defer loadLock.Unlock()
	if exist {
		return
	}
	readApplicationYaml()
	exist = true
}

type PluginConf struct {
	ApiConf *ApiRegisterConf
	AppConf *AppRegisterConf
}

type ApiRegisterConf struct {
	Enable          bool   `json:"enable" yaml:"enable"`
	RegisterHost    string `json:"registerHost" yaml:"registerHost"`
	SwaggerFilePath string `json:"swaggerFilePath" yaml:"swaggerFilePath"`
}

type AppRegisterConf struct {
	//Enable
	Enable         bool   `json:"enable" yaml:"enable"`
	ConnectTimeout int    `yaml:"connectTimeout"`
	ReadTimeout    int    `yaml:"readTimeout"`
	AppId          string `json:"appId" yaml:"appId"`
	AppCode        string `json:"appCode" yaml:"appCode"`
	AppName        string `json:"appName" yaml:"appName"`
	IsMainService  bool   `json:"isMainService" yaml:"isMainService"`
	InMenu         int    `json:"in-menu" yaml:"in-menu"`
	ServiceName    string `json:"serviceName" yaml:"serviceName"`
	ServiceId      string `json:"serviceId" yaml:"serviceId"`
	ServicePath    string `json:"servicePath" yaml:"servicePath"`
	ServerHost     string `json:"serverHost" yaml:"serverHost"`
	RegisterAppUrl string `json:"registerAppUrl" yaml:"registerAppUrl"`
	ServiceUrl     string `json:"serviceUrl" yaml:"serviceUrl"`
	ExcludeUrl     string `json:"excludeUrl" yaml:"excludeUrl"`
	AppVersion     string `json:"app-version" yaml:"app-version"`
	Type           int    `yaml:"type"`
}

func NewDefaultConf() *AppRegisterConf {
	return &AppRegisterConf{
		ConnectTimeout: 3000,
		Enable:         false,
		ReadTimeout:    15000,
		AppId:          "",
		AppCode:        "",
		AppName:        "",
		IsMainService:  false,
		InMenu:         2,
		ServiceName:    "",
		ServiceId:      "",
		ServicePath:    "",
		ServerHost:     "",
		RegisterAppUrl: "http://isc-permission-service:32100",
		ServiceUrl:     "",
		ExcludeUrl:     "/all",
		AppVersion:     "",
		Type:           3,
	}
}

func (conf *AppRegisterConf) ReadServiceRegisterYaml() {
	serverConf := &AppRegisterConf{}
	swaggerData, err := GetFileInfo("./service-register.yaml")
	if err != nil {
		logger.Error("读取service-register配置文件失败", err.Error())
		return
	}
	err = yaml.Unmarshal(swaggerData, serverConf)
	if err != nil {
		logger.Error("解析service-register配置信息错误", err.Error())
		return
	}
	*conf = *serverConf
}

func readApplicationYaml() {
	resultMap := make(map[string]any)
	swaggerData, err := GetFileInfo("./service-register.yaml")
	if err != nil {
		logger.Error("读取service-register配置文件失败", err.Error())
	}
	err = yaml.Unmarshal(swaggerData, &resultMap)
	var apiConf ApiRegisterConf
	err = mapstructure.Decode(resultMap["apiRegister"], &apiConf)
	if err != nil {
		logger.Error("解析service-register配置信息错误", err.Error())
	}
	var appConf AppRegisterConf
	err = mapstructure.Decode(resultMap["appRegister"], &appConf)
	if err != nil {
		logger.Error("解析service-register配置信息错误", err.Error())
	}
	PluginConfig.ApiConf = &apiConf
	PluginConfig.AppConf = &appConf
}
