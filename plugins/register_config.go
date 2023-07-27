package plugins

import (
	"github.com/isyscore/isc-gobase/logger"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
	"reflect"
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
	PermissionHost  string `json:"permissionHost" yaml:"permissionHost"`
	GroupName       string `json:"groupName" yaml:"groupName"`
	ApiServiceCode  string `json:"apiServiceCode" yaml:"apiServiceCode"`
	ApiManagerHost  string `json:"apiManagerHost" yaml:"apiManagerHost"`
	Type            string `json:"type" yaml:"type"`
}

type AppRegisterConf struct {
	//Enable
	Enable         bool                  `json:"enable" yaml:"enable"`
	ConnectTimeout int                   `yaml:"connectTimeout"`
	ReadTimeout    int                   `yaml:"readTimeout"`
	AppId          string                `json:"appId" yaml:"appId"`
	AppCode        string                `json:"appCode" yaml:"appCode"`
	AppName        string                `json:"appName" yaml:"appName"`
	IsMainService  bool                  `json:"isMainService" yaml:"isMainService"`
	InMenu         int                   `json:"in-menu" yaml:"in-menu"`
	SecondType     int                   `json:"secondType" yaml:"secondType"`
	RegisterAppUrl string                `json:"registerAppUrl" yaml:"registerAppUrl"`
	RedirectUrl    string                `json:"redirectUrl" yaml:"redirectUrl"`
	AppVersion     string                `json:"app-version" yaml:"app-version"`
	Type           int                   `yaml:"type"`
	ServiceList    []ServiceRegisterInfo `json:"serviceList"`
	ServiceName    string                `json:"serviceName" yaml:"serviceName"`
	ServiceId      string                `json:"serviceId" yaml:"serviceId"`
	ServicePath    string                `json:"servicePath" yaml:"servicePath"`
	ServiceUrl     string                `json:"serviceUrl" yaml:"serviceUrl"`
	ExcludeUrl     string                `json:"excludeUrl" yaml:"excludeUrl"`
}

// 服务注册信息
type ServiceRegisterInfo struct {
	ServiceId     string `json:"serviceId" `
	ServicePath   string `json:"servicePath" `
	ServiceName   string `json:"serviceName" `
	ServiceUrl    string `json:"serviceUrl" `
	Protocol      string `json:"protocol" `
	ExcludeUrl    string `json:"excludeUrl"`
	SpecialUrl    string `json:"specialUrl"`
	ServiceEnable bool   `json:"serviceEnable"`
	Retryable     bool   `json:"retryable"`
	StripPrefix   bool   `json:"stripPrefix"`
}

func defaultApiConf() *ApiRegisterConf {
	return &ApiRegisterConf{
		Enable: false,
	}
}

func defaultAppConf() *AppRegisterConf {
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
	if NCNotEmpty[any](resultMap["apiRegister"]) {
		var apiConf ApiRegisterConf
		err = mapstructure.Decode(resultMap["apiRegister"], &apiConf)
		if err != nil {
			logger.Error("解析service-register配置信息错误", err.Error())
		}
		PluginConfig.ApiConf = &apiConf
	} else {
		PluginConfig.ApiConf = defaultApiConf()
	}

	if NCNotEmpty[any](resultMap["appRegister"]) {
		var appConf AppRegisterConf
		err = mapstructure.Decode(resultMap["appRegister"], &appConf)
		if err != nil {
			logger.Error("解析service-register配置信息错误", err.Error())
		}
		PluginConfig.AppConf = &appConf
	} else {
		PluginConfig.AppConf = defaultAppConf()
	}
}

func NCIsEmpty[T any](v T) bool {
	var zero T
	return reflect.DeepEqual(v, zero)
}

func NCNotEmpty[T any](v T) bool {
	var zero T
	return !reflect.DeepEqual(v, zero)
}
